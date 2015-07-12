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

type ClusterControler struct {
	Ecs            *aws.ECSManager
	TargetResource string
}

func (self *ClusterControler) SearchClusters(projectDir string) []schema.Cluster {

	clusterDir := projectDir + "/cluster"
	files, err := ioutil.ReadDir(clusterDir)

	if err != nil {
		panic(err)
	}

	clusters := []schema.Cluster{}

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

	return clusters
}

func (self *ClusterControler) CreateClusterUpdatePlans(clusters []schema.Cluster) []*plan.ClusterUpdatePlan {

	plans := []*plan.ClusterUpdatePlan{}
	for _, cluster := range clusters {
		if len(self.TargetResource) == 0 || self.TargetResource == cluster.Name {
			plans = append(plans, self.CreateClusterUpdatePlan(cluster))
		}
	}
	return plans
}

func (self *ClusterControler) CreateClusterUpdatePlan(cluster schema.Cluster) *plan.ClusterUpdatePlan {

	output, err := self.Ecs.ClusterApi().DescribeClusters([]*string{&cluster.Name})

	if err != nil {
		fmt.Fprintln(os.Stderr, color.Red("[ERROR] discribe_cluster"))
		os.Exit(1)
	}

	if len(output.Failures) > 0 {
		fmt.Fprintln(os.Stderr, color.Red(fmt.Sprintf("[ERROR] Cluster '%s' not found", cluster.Name)))
		os.Exit(1)
	}

	target := output.Clusters[0]

	if *target.Status != "ACTIVE" {
		fmt.Fprintln(os.Stderr, color.Red(fmt.Sprintf("[ERROR] Cluster '%s' is not ACTIVE.", cluster.Name)))
		os.Exit(1)
	}

	api := self.Ecs.ServiceApi()
	resListServices, _ := api.ListServices(cluster.Name)

	resDescribeService, _ := api.DescribeService(cluster.Name, resListServices.ServiceARNs)

	currentServices := map[string]*ecs.Service{}
	for _, service := range resDescribeService.Services {
		currentServices[*service.ServiceName] = service
	}

	deleteServices := map[string]*ecs.Service{}
	updateServices := map[string]*plan.UpdateService{}
	for name, currentService := range currentServices {

		if newService, ok := cluster.Services[name]; ok {
			// update
			updateServices[name] = &plan.UpdateService{
				Before: currentService,
				After: &newService,
			}
		} else {
			// delete
			deleteServices[name] = currentService
		}
	}

	newServices := map[string]*schema.Service{}
	for name, newService := range cluster.Services {

		if _, ok := currentServices[name]; !ok {
			// create
			newServices[name] = &newService
		}
	}

	return &plan.ClusterUpdatePlan{
		Name: cluster.Name,
		CurrentServices: currentServices,
		DeleteServices: deleteServices,
		UpdateServices: updateServices,
		NewServices: newServices,
	}
}

func (self *ClusterControler) ApplyClusterPlans(plans []*plan.ClusterUpdatePlan) {

	fmt.Println("Start apply serivces...")

	for _, plan := range plans {
		self.ApplyClusterPlan(plan)
	}
}

func (self *ClusterControler) ApplyClusterPlan(plan *plan.ClusterUpdatePlan) {

	api := self.Ecs.ServiceApi()

	for _, delete := range plan.DeleteServices {

		// DesiredCount must be 0 to remove service.
		_, err := api.UpdateService(plan.Name, schema.Service{
			Name: *delete.ServiceName,
			DesiredCount: 0,
		})

		if err != nil {
			panic(err)
		}

		result, err := api.DeleteService(plan.Name, *delete.ServiceARN)
		if err != nil {
			panic(err)
		}

		fmt.Printf("[INFO] Removed service '%s'\n", *result.Service.ClusterARN)
	}

	for _, add := range plan.NewServices {

		result, err := api.CreateService(plan.Name, schema.Service{
			Name: add.Name,
			DesiredCount: add.DesiredCount,
			TaskDefinition: add.TaskDefinition,
		})

		if err != nil {
			panic(err)
		}

		fmt.Printf("[INFO] Created service '%s'\n", *result.Service.ServiceARN)
	}

	for _, update := range plan.UpdateServices {

		_, err1 := api.UpdateService(plan.Name, schema.Service{
			Name: *update.Before.ServiceName,
			DesiredCount: 0,
		})

		if err1 != nil {
			panic(err1)
		}

		fmt.Printf("[INFO] Waiting to stop '%s' service on '%s' ...\n", *update.Before.ServiceName, plan.Name)
		if err := self.waitStoppingService(plan.Name, *update.Before.ServiceName); err != nil {
			panic(err)
		}
		fmt.Printf("[INFO] Stoped '%s' service on '%s'.\n", *update.Before.ServiceName, plan.Name)

		result, err2 := api.UpdateService(plan.Name, schema.Service{
			Name: update.After.Name,
			DesiredCount: update.After.DesiredCount,
			TaskDefinition: update.After.TaskDefinition,
		})

		if err2 != nil {
			panic(err2)
		}

		fmt.Printf("[INFO] Updated service '%s'\n", *result.Service.ServiceARN)
	}

}

func (self *ClusterControler) waitStoppingService(cluster string, service string) error {

	api := self.Ecs.ServiceApi()

	for {
		time.Sleep(5 * time.Second)

		result, err := api.DescribeService(cluster, []*string{&service})

		if err != nil {
			return err
		}

		if len(result.Services) == 0 {
			return errors.New("service not found")
		}

		target := result.Services[0]

		if *target.RunningCount == 0 {
			return nil
		}

		// TODO retry count restriction
	}
}