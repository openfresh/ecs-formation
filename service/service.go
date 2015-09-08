package service

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/stormcat24/ecs-formation/aws"
	"github.com/stormcat24/ecs-formation/logger"
	"github.com/stormcat24/ecs-formation/util"
	"github.com/str1ngs/ansi/color"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

type TaskWatchStatus int

const (
	WatchContinue TaskWatchStatus = iota
	WatchFinish
	WatchTerminate
)

type ServiceController struct {
	manager        *aws.AwsManager
	TargetResource string
	clusters       []Cluster
}

func NewServiceController(manager *aws.AwsManager, projectDir string, targetResource string) (*ServiceController, error) {

	con := &ServiceController{
		manager: manager,
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

func (self *ServiceController) searchServices(projectDir string) ([]Cluster, error) {

	clusterDir := projectDir + "/service"
	files, err := ioutil.ReadDir(clusterDir)

	clusters := []Cluster{}

	if err != nil {
		return clusters, err
	}

	filePattern := regexp.MustCompile("^(.+)\\.yml$")

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".yml") {
			content, _ := ioutil.ReadFile(clusterDir + "/" + file.Name())

			tokens := filePattern.FindStringSubmatch(file.Name())
			clusterName := tokens[1]

			serviceMap, _ := CreateServiceMap(content)
			cluster := Cluster{
				Name:     clusterName,
				Services: serviceMap,
			}

			clusters = append(clusters, cluster)
		}
	}

	return clusters, nil
}

func (self *ServiceController) GetClusters() []Cluster {
	return self.clusters
}

func (self *ServiceController) CreateServiceUpdatePlans() ([]*ServiceUpdatePlan, error) {

	plans := []*ServiceUpdatePlan{}
	for _, cluster := range self.GetClusters() {
		if len(self.TargetResource) == 0 || self.TargetResource == cluster.Name {
			cp, err := self.CreateServiceUpdatePlan(cluster)
			if err != nil {
				return plans, err
			}

			if cp != nil {
				plans = append(plans, cp)
			}
		}
	}
	return plans, nil
}

func (self *ServiceController) CreateServiceUpdatePlan(cluster Cluster) (*ServiceUpdatePlan, error) {

	api := self.manager.EcsApi()
	output, errdc := api.DescribeClusters([]*string{&cluster.Name})

	if errdc != nil {
		return nil, errdc
	}

	if len(output.Failures) > 0 {
		return nil, errors.New(fmt.Sprintf("Cluster '%s' not found", cluster.Name))
	}

	rlci, errlci := api.ListContainerInstances(cluster.Name)
	if errlci != nil {
		return nil, errlci
	}

	if len(rlci.ContainerInstanceArns) == 0 {
		logger.Main.Warnf("ECS instances not found in cluster '%s' not found", cluster.Name)
		return nil, nil

	} else {
		target := output.Clusters[0]

		if *target.Status != "ACTIVE" {
			return &ServiceUpdatePlan{}, errors.New(fmt.Sprintf("Cluster '%s' is not ACTIVE.", cluster.Name))
		}

		resListServices, errls := api.ListServices(cluster.Name)
		if errls != nil {
			return &ServiceUpdatePlan{}, errls
		}

		currentServices := map[string]*ecs.Service{}
		if len(resListServices.ServiceArns) > 0 {
			resDescribeService, errds := api.DescribeService(cluster.Name, resListServices.ServiceArns)
			if errds != nil {
				return &ServiceUpdatePlan{}, errds
			}

			for _, service := range resDescribeService.Services {
				currentServices[*service.ServiceName] = service
			}
		}

		newServices := map[string]*Service{}
		for name, newService := range cluster.Services {
			s := newService
			newServices[name] = &s
		}

		return &ServiceUpdatePlan{
			Name:            cluster.Name,
			InstanceARNs:    rlci.ContainerInstanceArns,
			CurrentServices: currentServices,
			NewServices:     newServices,
		}, nil
	}
}

func (self *ServiceController) ApplyServicePlans(plans []*ServiceUpdatePlan) {

	logger.Main.Info("Start apply serivces...")

	for _, plan := range plans {
		if err := self.ApplyServicePlan(plan); err != nil {
			logger.Main.Error(color.Red(err.Error()))
			os.Exit(1)
		}
	}
}

func (self *ServiceController) ApplyServicePlan(plan *ServiceUpdatePlan) error {

	api := self.manager.EcsApi()

	for _, current := range plan.CurrentServices {

		// set desired_count = 0
		if _, err := api.UpdateService(plan.Name, *current.ServiceName, 0, *current.TaskDefinition); err != nil {
			return err
		}

		// wait to stop service
		logger.Main.Infof("Waiting to stop '%s' service on '%s' ...", *current.ServiceName, plan.Name)
		if err := self.waitStoppingService(plan.Name, *current.ServiceName); err != nil {
			return err
		}
		logger.Main.Infof("Stoped '%s' service on '%s'.", *current.ServiceName, plan.Name)

		// delete service
		result, err := api.DeleteService(plan.Name, *current.ServiceArn)
		if err != nil {
			return err
		}

		if err := self.waitStoppingService(plan.Name, *current.ServiceName); err != nil {
			return err
		}

		logger.Main.Infof("Deleted service '%s'", *result.Service.ServiceArn)
	}

	for _, add := range plan.NewServices {

		result, err := api.CreateService(plan.Name, add.Name, add.DesiredCount, toLoadBalancers(&add.LoadBalancers), add.TaskDefinition, add.Role)
		if err != nil {
			return err
		}

		logger.Main.Infof("Created service '%s'", *result.Service.ServiceArn)
		errwas := self.WaitActiveService(plan.Name, add.Name)
		if errwas != nil {
			return errwas
		}
	}

	return nil
}

func toLoadBalancers(values *[]LoadBalancer) []*ecs.LoadBalancer {

	loadBalancers := []*ecs.LoadBalancer{}
	for _, lb := range *values {
		loadBalancers = append(loadBalancers, &ecs.LoadBalancer{
			LoadBalancerName: &lb.Name,
			ContainerName:    &lb.ContainerName,
			ContainerPort:    &lb.ContainerPort,
		})
	}

	return loadBalancers
}

func (self *ServiceController) waitStoppingService(cluster string, service string) error {

	api := self.manager.EcsApi()

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

	api := self.manager.EcsApi()

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

			if len(target.Events) > 0 && strings.Contains(*target.Events[0].Message, "was unable to place a task") {
				return errors.New(*target.Events[0].Message)
			}

			if !flag {
				reslt, errlt := api.ListTasks(cluster, service)
				if errlt != nil {
					return errlt
				}

				if len(reslt.TaskArns) == 0 {
					continue
				} else {
					taskARNs = reslt.TaskArns
					flag = true
				}
			}

			resdt, errdt := api.DescribeTasks(cluster, taskARNs)
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
		util.Println(fmt.Sprintf("    %s:", *task.TaskArn))
		util.Println(fmt.Sprintf("        LastStatus:%s", self.RoundColorStatus(*task.LastStatus)))
		util.Println("        Containers:")

		for _, con := range task.Containers {
			util.Println(fmt.Sprintf("            ----------[%s]----------", *con.Name))
			util.Println(fmt.Sprintf("            ContainerARN:%s", *con.ContainerArn))
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
