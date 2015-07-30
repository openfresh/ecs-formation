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
	"github.com/stormcat24/ecs-formation/util"
)

type TaskWatchStatus int

const (
	WatchContinue TaskWatchStatus = iota
	WatchFinish
	WatchTerminate
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

	clusterApi := self.Ecs.ClusterApi()
	output, errdc := clusterApi.DescribeClusters([]*string{&cluster.Name})

	if errdc != nil {
		return &plan.ServiceUpdatePlan{}, errdc
	}

	if len(output.Failures) > 0 {
		return &plan.ServiceUpdatePlan{}, errors.New(fmt.Sprintf("Cluster '%s' not found", cluster.Name))
	}

	rlci, errlci := clusterApi.ListContainerInstances(cluster.Name)
	if errlci != nil {
		return &plan.ServiceUpdatePlan{}, errlci
	}

	if len(rlci.ContainerInstanceARNs) == 0 {
		return &plan.ServiceUpdatePlan{}, errors.New(fmt.Sprintf("ECS instances not found in cluster '%s' not found", cluster.Name))
	}

	target := output.Clusters[0]

	if *target.Status != "ACTIVE" {
		return &plan.ServiceUpdatePlan{}, errors.New(fmt.Sprintf("Cluster '%s' is not ACTIVE.", cluster.Name))
	}

	serviceApi := self.Ecs.ServiceApi()

	resListServices, errls := serviceApi.ListServices(cluster.Name)
	if errls != nil {
		return &plan.ServiceUpdatePlan{}, errls
	}

	currentServices := map[string]*ecs.Service{}
	if len(resListServices.ServiceARNs) > 0 {
		resDescribeService, errds := serviceApi.DescribeService(cluster.Name, resListServices.ServiceARNs)
		if errds != nil {
			return &plan.ServiceUpdatePlan{}, errds
		}

		for _, service := range resDescribeService.Services {
			currentServices[*service.ServiceName] = service
		}
	}

	newServices := map[string]*schema.Service{}
	for name, newService := range cluster.Services {
		s := newService
		newServices[name] = &s
	}

	return &plan.ServiceUpdatePlan{
		Name: cluster.Name,
		InstanceARNs: rlci.ContainerInstanceARNs,
		CurrentServices: currentServices,
		NewServices: newServices,
	}, nil
}

func (self *ServiceController) ApplyServicePlans(plans []*plan.ServiceUpdatePlan) {

	logger.Main.Info("Start apply serivces...")

	for _, plan := range plans {
		if err := self.ApplyServicePlan(plan); err != nil {
			logger.Main.Error(color.Red(err.Error()))
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
		errwas := self.WaitActiveService(plan.Name, add.Name)
		if errwas != nil {
			return errwas
		}
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
	taskApi := self.Ecs.TaskApi()

	var flag = false
	var taskARNs []*string

	for {
		time.Sleep(5 * time.Second)

		result, err := api.DescribeService(cluster, []*string{&service})

		if err != nil {
			return err
		}

		if len(result.Services) == 0 {
			continue
		}

		target := result.Services[0]

		// The status of the service. The valid values are ACTIVE, DRAINING, or INACTIVE.
		logger.Main.Infof("service '%s@%s' status = %s ...", service, cluster, *target.Status)

		if *target.Status == "ACTIVE" {

			if len(target.Events) > 0 && strings.Contains(*target.Events[0].Message, "was unable to place a task"){
				return errors.New(*target.Events[0].Message)
			}

			if !flag {
				reslt, errlt := taskApi.ListTasks(cluster, service)
				if errlt != nil {
					return errlt
				}

				if len(reslt.TaskARNs) == 0 {
					continue
				} else {
					taskARNs = reslt.TaskARNs
					flag = true
				}
			}

			resdt, errdt := taskApi.DescribeTasks(cluster, taskARNs)
			if errdt != nil {
				return errdt
			}


			watchStatus := self.checkRunningTask(resdt)
			if watchStatus == WatchFinish {
				logger.Main.Info("At least one of task has started successfully.")
				return nil
			} else if watchStatus == WatchTerminate {
				logger.Main.Error("Stopped watching task, because task has stopped.")
				return errors.New("Task has been stopped for some reason.")
			}

		}
	}
}

func (self *ServiceController) checkRunningTask(dto *ecs.DescribeTasksOutput) TaskWatchStatus {

	logger.Main.Info("Current task conditions as follows:")

	status := []string{}
	for _, task := range dto.Tasks {
		util.Println(fmt.Sprintf("    %s:", *task.TaskARN))
		util.Println(fmt.Sprintf("        LastStatus:%s", self.RoundColorStatus(*task.LastStatus)))
		util.Println("        Containers:")

		for _, con := range task.Containers {
			util.Println(fmt.Sprintf("            ----------[%s]----------", *con.Name))
			util.Println(fmt.Sprintf("            ContainerARN:%s", *con.ContainerARN))
			util.Println(fmt.Sprintf("            Status:%s", self.RoundColorStatus(*con.LastStatus)))
			util.Println()
		}

		status = append(status, *task.LastStatus)
	}

	// if RUNNING at least one, ecs-formation deals with ok.
	for _, s := range status {
		if s == "RUNNING" {
			return WatchFinish
		} else if s == "STOPPED" {
			return WatchTerminate
		}
	}

	return WatchContinue
}

func (self *ServiceController) RoundColorStatus(status string) *color.Escape {

	if status == "RUNNING" {
		return color.Green(status)
	} else if status == "PENDING" {
		return color.Yellow(status)
	} else if status == "STOPPED" {
		return color.Red(status)
	} else {
		return color.Magenta(status)
	}
}
