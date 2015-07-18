package cluster

import (
	"io/ioutil"
	"github.com/stormcat24/ecs-formation/schema"
	"fmt"
	"strings"
	"regexp"
	"github.com/stormcat24/ecs-formation/aws"
	"os"
	"github.com/str1ngs/ansi/color"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/stormcat24/ecs-formation/plan"
	"time"
	"errors"
)

type ClusterController struct {
	Ecs            *aws.ECSManager
	TargetResource string
	clusters       []schema.Cluster
}

func NewClusterController(ecs *aws.ECSManager, projectDir string, targetResource string) (*ClusterController, error) {

	con := &ClusterController{
		Ecs: ecs,
	}

	clusters, err := con.searchClusters(projectDir)
	if err != nil {
		return nil, err
	}

	con.clusters = clusters

	if targetResource != "" {
		con.TargetResource = targetResource
	}

	return con, nil
}

func (self *ClusterController) searchClusters(projectDir string) ([]schema.Cluster, error) {

	clusterDir := projectDir + "/cluster"
	files, err := ioutil.ReadDir(clusterDir)

	clusters := []schema.Cluster{}

	if err != nil {
		return clusters, err
	}

	filePattern := regexp.MustCompile("^(.+)\\.yml$")

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".yml") {
			content, _ := ioutil.ReadFile(clusterDir + "/" + file.Name())

			tokens := filePattern.FindStringSubmatch(file.Name())
			clusterName := tokens[1]

			serviceMap, _ := schema.CreateServiceMap(content)
			cluster := schema.Cluster{
				Name: clusterName,
				Services: serviceMap,
			}

			clusters = append(clusters, cluster)
		}
	}

	return clusters, nil
}

func (self *ClusterController) GetClusters() []schema.Cluster {
	return self.clusters
}

func (self *ClusterController) CreateClusterUpdatePlans() ([]*plan.ClusterUpdatePlan, error) {

	plans := []*plan.ClusterUpdatePlan{}
	for _, cluster := range self.GetClusters() {
		if len(self.TargetResource) == 0 || self.TargetResource == cluster.Name {

			cp, err := self.CreateClusterUpdatePlan(cluster)
			if err != nil {
				return plans, err
			}

			plans = append(plans, cp)
		}
	}
	return plans, nil
}

func (self *ClusterController) CreateClusterUpdatePlan(cluster schema.Cluster) (*plan.ClusterUpdatePlan, error) {

	output, errdc := self.Ecs.ClusterApi().DescribeClusters([]*string{&cluster.Name})

	if errdc != nil {
		return &plan.ClusterUpdatePlan{}, errdc
	}

	if len(output.Failures) > 0 {
		return &plan.ClusterUpdatePlan{}, errors.New(fmt.Sprintf("[ERROR] Cluster '%s' not found", cluster.Name))
	}

	target := output.Clusters[0]

	if *target.Status != "ACTIVE" {
		return &plan.ClusterUpdatePlan{}, errors.New(fmt.Sprintf("[ERROR] Cluster '%s' is not ACTIVE.", cluster.Name))
	}

	api := self.Ecs.ServiceApi()

	resListServices, errls := api.ListServices(cluster.Name)
	if errls != nil {
		return &plan.ClusterUpdatePlan{}, errls
	}

	resDescribeService, errds := api.DescribeService(cluster.Name, resListServices.ServiceARNs)
	if errds != nil {
		return &plan.ClusterUpdatePlan{}, errds
	}

	currentServices := map[string]*ecs.Service{}
	for _, service := range resDescribeService.Services {
		currentServices[*service.ServiceName] = service
	}

	newServices := map[string]*schema.Service{}
	for name, newService := range cluster.Services {
		newServices[name] = &newService
	}

	return &plan.ClusterUpdatePlan{
		Name: cluster.Name,
		CurrentServices: currentServices,
		NewServices: newServices,
	}, nil
}

func (self *ClusterController) ApplyClusterPlans(plans []*plan.ClusterUpdatePlan) {

	fmt.Println("Start apply serivces...")

	for _, plan := range plans {
		if err := self.ApplyClusterPlan(plan); err != nil {
			fmt.Fprintln(os.Stderr, color.Red(err.Error()))
			os.Exit(1)
		}
	}
}

func (self *ClusterController) ApplyClusterPlan(plan *plan.ClusterUpdatePlan) error {

	api := self.Ecs.ServiceApi()

	for _, current := range plan.CurrentServices {

		// set desired_count = 0
		if _, err := api.UpdateService(plan.Name, schema.Service{
				Name: *current.ServiceName,
				DesiredCount: 0,
			}); err != nil {
			return err
		}

		// wait to stop service
		fmt.Printf("[INFO] Waiting to stop '%s' service on '%s' ...\n", *current.ServiceName, plan.Name)
		if err := self.waitStoppingService(plan.Name, *current.ServiceName); err != nil {
			return err
		}
		fmt.Printf("[INFO] Stoped '%s' service on '%s'.\n", *current.ServiceName, plan.Name)


		// delete service
		result, err := api.DeleteService(plan.Name, *current.ServiceARN)
		if err != nil {
			return err
		}

		fmt.Printf("[INFO] Waiting to delete '%s' service on '%s' ...\n", *current.ServiceName, plan.Name)
		if err := self.waitStoppingService(plan.Name, *current.ServiceName); err != nil {
			return err
		}

		fmt.Printf("[INFO] Deleted service '%s'\n", *result.Service.ServiceARN)
	}

	for _, add := range plan.NewServices {

		result, err := api.CreateService(plan.Name, schema.Service{
			Name: add.Name,
			DesiredCount: add.DesiredCount,
			LoadBalancers: add.LoadBalancers,
			TaskDefinition: add.TaskDefinition,
			Role: add.Role,
		})

		if err != nil {
			return err
		}

		fmt.Printf("[INFO] Created service '%s'\n", *result.Service.ServiceARN)
	}

	return nil
}

func (self *ClusterController) waitStoppingService(cluster string, service string) error {

	api := self.Ecs.ServiceApi()

	for {
		time.Sleep(5 * time.Second)

		result, err := api.DescribeService(cluster, []*string{&service})

		if err != nil {
			return err
		}

		if len(result.Services) == 0 {
			return nil
		}

		target := result.Services[0]

		fmt.Printf("[INFO] service '%s@%s' current status = %s \n", service, cluster, *target.Status)
		if *target.RunningCount == 0 && *target.Status != "DRAINING" {
			return nil
		}

	}
}

func (self *ClusterController) WaitActiveService(cluster string, service string) error {

	api := self.Ecs.ServiceApi()

	for {
		time.Sleep(5 * time.Second)

		result, err := api.DescribeService(cluster, []*string{&service})

		if err != nil {
			return err
		}

		if len(result.Services) == 0 {
			return nil
		}

		target := result.Services[0]

		// The status of the service. The valid values are ACTIVE, DRAINING, or INACTIVE.
		if *target.Status == "ACTIVE" {

			if len(target.Events) > 0 && strings.Contains(*target.Events[0].Message, "has started") {
				return nil
			}
		} else {
			fmt.Printf("[INFO] service '%s@%s' status = %s ...\n", service, cluster, *target.Status)
		}
	}
}
