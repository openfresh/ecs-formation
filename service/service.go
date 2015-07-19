package service

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
	"github.com/stormcat24/ecs-formation/logger"
	"time"
	"errors"
)

type ServiceController struct {
	Ecs            *aws.ECSManager
	TargetResource string
	clusters       []schema.Cluster
}

func NewServiceController(ecs *aws.ECSManager, projectDir string, targetResource string) (*ServiceController, error) {

	con := &ServiceController{
		Ecs: ecs,
	}

	clusters, err := con.searchServices(projectDir)
	if err != nil {
		return nil, err
	}

	con.clusters = clusters

	if targetResource != "" {
		con.TargetResource = targetResource
	}

	return con, nil
}

func (self *ServiceController) searchServices(projectDir string) ([]schema.Cluster, error) {

	clusterDir := projectDir + "/service"
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

func (self *ServiceController) GetClusters() []schema.Cluster {
	return self.clusters
}

func (self *ServiceController) CreateServiceUpdatePlans() ([]*plan.ServiceUpdatePlan, error) {

	plans := []*plan.ServiceUpdatePlan{}
	for _, cluster := range self.GetClusters() {
		if len(self.TargetResource) == 0 || self.TargetResource == cluster.Name {

			cp, err := self.CreateServiceUpdatePlan(cluster)
			if err != nil {
				return plans, err
			}

			plans = append(plans, cp)
		}
	}
	return plans, nil
}

func (self *ServiceController) CreateServiceUpdatePlan(cluster schema.Cluster) (*plan.ServiceUpdatePlan, error) {

	output, errdc := self.Ecs.ClusterApi().DescribeClusters([]*string{&cluster.Name})

	if errdc != nil {
		return &plan.ServiceUpdatePlan{}, errdc
	}

	if len(output.Failures) > 0 {
		return &plan.ServiceUpdatePlan{}, errors.New(fmt.Sprintf("Cluster '%s' not found", cluster.Name))
	}

	target := output.Clusters[0]

	if *target.Status != "ACTIVE" {
		return &plan.ServiceUpdatePlan{}, errors.New(fmt.Sprintf("Cluster '%s' is not ACTIVE.", cluster.Name))
	}

	api := self.Ecs.ServiceApi()

	resListServices, errls := api.ListServices(cluster.Name)
	if errls != nil {
		return &plan.ServiceUpdatePlan{}, errls
	}

	resDescribeService, errds := api.DescribeService(cluster.Name, resListServices.ServiceARNs)
	if errds != nil {
		return &plan.ServiceUpdatePlan{}, errds
	}

	currentServices := map[string]*ecs.Service{}
	for _, service := range resDescribeService.Services {
		currentServices[*service.ServiceName] = service
	}

	newServices := map[string]*schema.Service{}
	for name, newService := range cluster.Services {
		newServices[name] = &newService
	}

	return &plan.ServiceUpdatePlan{
		Name: cluster.Name,
		CurrentServices: currentServices,
		NewServices: newServices,
	}, nil
}

func (self *ServiceController) ApplyServicePlans(plans []*plan.ServiceUpdatePlan) {

	logger.Main.Info("Start apply serivces...")

	for _, plan := range plans {
		if err := self.ApplyServicePlan(plan); err != nil {
			fmt.Fprintln(os.Stderr, color.Red(err.Error()))
			os.Exit(1)
		}
	}
}

func (self *ServiceController) ApplyServicePlan(plan *plan.ServiceUpdatePlan) error {

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
		logger.Main.Infof("Waiting to stop '%s' service on '%s' ...", *current.ServiceName, plan.Name)
		if err := self.waitStoppingService(plan.Name, *current.ServiceName); err != nil {
			return err
		}
		logger.Main.Infof("Stoped '%s' service on '%s'.", *current.ServiceName, plan.Name)


		// delete service
		result, err := api.DeleteService(plan.Name, *current.ServiceARN)
		if err != nil {
			return err
		}

		logger.Main.Infof("Waiting to delete '%s' service on '%s' ...", *current.ServiceName, plan.Name)
		if err := self.waitStoppingService(plan.Name, *current.ServiceName); err != nil {
			return err
		}

		logger.Main.Infof("Deleted service '%s'", *result.Service.ServiceARN)
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

		logger.Main.Infof("Created service '%s'", *result.Service.ServiceARN)
	}

	return nil
}

func (self *ServiceController) waitStoppingService(cluster string, service string) error {

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

		logger.Main.Infof("service '%s@%s' current status = %s", service, cluster, *target.Status)
		if *target.RunningCount == 0 && *target.Status != "DRAINING" {
			return nil
		}

	}
}

func (self *ServiceController) WaitActiveService(cluster string, service string) error {

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
			logger.Main.Infof("service '%s@%s' status = %s ...", service, cluster, *target.Status)
		}
	}
}
