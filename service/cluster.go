package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	awsecs "github.com/aws/aws-sdk-go/service/ecs"
	"github.com/fatih/color"
	"github.com/openfresh/ecs-formation/client"
	"github.com/openfresh/ecs-formation/client/applicationautoscaling"
	"github.com/openfresh/ecs-formation/client/ecs"
	"github.com/openfresh/ecs-formation/logger"
	"github.com/openfresh/ecs-formation/service/types"
	"github.com/openfresh/ecs-formation/util"
)

type ClusterService interface {
	SearchClusters() ([]types.Cluster, error)
	CreateServiceUpdatePlans() ([]*types.ServiceUpdatePlan, error)
	ApplyServicePlans(plans []*types.ServiceUpdatePlan) error
	ApplyServicePlan(plan *types.ServiceUpdatePlan) error
}

type ConcreteClusterService struct {
	ecsCli            ecs.Client
	appAutoscalingCli applicationautoscaling.Client
	projectDir        string
	clusters          []string
	targetService     string
	params            map[string]string
}

func NewClusterService(projectDir string, clusters []string, targetService string, params map[string]string) (ClusterService, error) {

	service := ConcreteClusterService{
		ecsCli:            client.AWSCli.ECS,
		appAutoscalingCli: client.AWSCli.ApplicationAutoscaling,
		projectDir:        projectDir,
		clusters:          clusters,
		targetService:     targetService,
		params:            params,
	}

	return &service, nil
}

func (s ConcreteClusterService) SearchClusters() ([]types.Cluster, error) {

	clusterDir := s.projectDir + "/service"
	clusters := []types.Cluster{}

	filePattern := regexp.MustCompile(`^.+\/(.+)\.yml$`)
	filepath.Walk(clusterDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || !strings.HasSuffix(path, ".yml") {
			return nil
		}

		flg := false
		for _, cluster := range s.clusters {
			if strings.HasSuffix(path, fmt.Sprintf("%s.yml", cluster)) {
				flg = true
			}
		}

		if !flg {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		merged := util.MergeYamlWithParameters(content, s.params)
		tokens := filePattern.FindStringSubmatch(path)

		clusterName := tokens[1]

		serviceMap, err := types.CreateServiceMap(merged)
		if err != nil {
			return err
		}
		cluster := types.Cluster{
			Name:     clusterName,
			Services: serviceMap,
		}

		clusters = append(clusters, cluster)

		return nil
	})

	return clusters, nil
}

func (s ConcreteClusterService) CreateServiceUpdatePlans() ([]*types.ServiceUpdatePlan, error) {

	plans := []*types.ServiceUpdatePlan{}

	clusters, err := s.SearchClusters()
	if err != nil {
		return make([]*types.ServiceUpdatePlan, 0), err
	}

	for _, cluster := range clusters {
		cl, err := s.createServiceUpdatePlan(cluster)
		if err != nil {
			return make([]*types.ServiceUpdatePlan, 0), err
		}

		if cl != nil {
			plans = append(plans, cl)
		}
	}

	return plans, nil
}

func (s ConcreteClusterService) createServiceUpdatePlan(cluster types.Cluster) (*types.ServiceUpdatePlan, error) {

	output, err := s.ecsCli.DescribeClusters([]*string{
		aws.String(cluster.Name),
	})

	if err != nil {
		return nil, err
	}

	if len(output.Failures) > 0 {
		return nil, fmt.Errorf("Cluster '%s' not found", cluster.Name)
	} else {
		logger.Main.Infof("Cluster '%v' is found.", cluster.Name)
	}

	lciResult, err := s.ecsCli.ListContainerInstances(cluster.Name)
	if err != nil {
		return nil, err
	}

	if len(lciResult.ContainerInstanceArns) == 0 {
		logger.Main.Warnf("ECS instances not found in cluster '%s' not found", cluster.Name)
		return nil, nil
	} else {
		target := output.Clusters[0]

		if *target.Status != "ACTIVE" {
			return nil, fmt.Errorf("Cluster '%s' is not ACTIVE.", cluster.Name)
		}

		lsResult, err := s.ecsCli.ListServices(cluster.Name)
		if err != nil {
			return nil, err
		}

		currentStacks := map[string]*types.ServiceStack{}
		if len(lsResult.ServiceArns) > 0 {

			resDescribeService, errds := s.ecsCli.DescribeService(cluster.Name, lsResult.ServiceArns)
			if errds != nil {
				return nil, errds
			}

			for _, service := range resDescribeService.Services {
				if s.targetService == "" || (s.targetService != "" && s.targetService == *service.ServiceName) {
					autoScaling, err := s.appAutoscalingCli.DescribeScalableTarget(cluster.Name, *service.ServiceName)
					if err != nil {
						return nil, err
					}

					currentStacks[*service.ServiceName] = &types.ServiceStack{
						Service:     service,
						AutoScaling: autoScaling,
					}
				}
			}
		}

		newServices := map[string]*types.Service{}
		for name, newService := range cluster.Services {
			if s.targetService == "" || (s.targetService != "" && s.targetService == newService.Name) {
				s := newService
				newServices[name] = &s
			}
		}

		return &types.ServiceUpdatePlan{
			Name:            cluster.Name,
			InstanceARNs:    lciResult.ContainerInstanceArns,
			CurrentServices: currentStacks,
			NewServices:     newServices,
		}, nil
	}
}

func (s ConcreteClusterService) ApplyServicePlans(plans []*types.ServiceUpdatePlan) error {
	logger.Main.Info("Start apply serivces...")

	for _, plan := range plans {
		if err := s.ApplyServicePlan(plan); err != nil {
			logger.Main.Error(color.RedString(err.Error()))
			return err
		}
	}
	return nil
}

func (s ConcreteClusterService) ApplyServicePlan(plan *types.ServiceUpdatePlan) error {

	// currentにあってnewにない（削除）
	for _, currentStack := range plan.CurrentServices {
		current := currentStack.Service
		if _, ok := plan.NewServices[*current.ServiceName]; !ok {
			logger.Main.Infof("Delating '%s' service on '%s' ...", *current.ServiceName, plan.Name)

			// set desired_count = 0
			params := awsecs.UpdateServiceInput{
				Cluster:        aws.String(plan.Name),
				Service:        current.ServiceName,
				DesiredCount:   aws.Int64(0),
				TaskDefinition: current.TaskDefinition,
			}
			if _, err := s.ecsCli.UpdateService(&params); err != nil {
				return err
			}
			logger.Main.Infof("Updated desired count = 0 of '%s' service on '%s' ...", *current.ServiceName, plan.Name)

			// wait to stop service
			logger.Main.Infof("Waiting to stop '%s' service on '%s' ...", *current.ServiceName, plan.Name)
			if err := s.waitStoppingService(plan.Name, *current.ServiceName); err != nil {
				return err
			}
			logger.Main.Infof("Stoped '%s' service on '%s'.", *current.ServiceName, plan.Name)

			// delete service
			dsrv, err := s.ecsCli.DeleteService(plan.Name, *current.ServiceArn)
			if err != nil {
				return err
			}

			if err := s.waitStoppingService(plan.Name, *current.ServiceName); err != nil {
				return err
			}

			logger.Main.Infof("Deleted service '%s' completely.", *dsrv.ServiceArn)
		}
	}
	// only new registration
	for _, add := range plan.NewServices {
		if s.targetService != "" && s.targetService != add.Name {
			continue
		}

		if _, ok := plan.CurrentServices[add.Name]; !ok {
			logger.Main.Infof("Creating '%s' service on '%s' ...", add.Name, plan.Name)

			p := awsecs.CreateServiceInput{
				Cluster:        aws.String(plan.Name),
				ServiceName:    aws.String(add.Name),
				DesiredCount:   aws.Int64(add.DesiredCount),
				LoadBalancers:  toLoadBalancersNew(add.LoadBalancers),
				Role:           aws.String(add.Role),
				TaskDefinition: aws.String(add.TaskDefinition),
			}
			if add.MinimumHealthyPercent.Valid && add.MaximumPercent.Valid {
				p.DeploymentConfiguration = &awsecs.DeploymentConfiguration{
					MinimumHealthyPercent: aws.Int64(add.MinimumHealthyPercent.Int64),
					MaximumPercent:        aws.Int64(add.MaximumPercent.Int64),
				}
			}

			csrv, err := s.ecsCli.CreateService(&p)
			if err != nil {
				return err
			}

			logger.Main.Infof("Created service '%s', task-definition is '%s'.", *csrv.ServiceArn, *csrv.TaskDefinition)
			if err := s.waitActiveService(plan.Name, add.Name); err != nil {
				return err
			}
			logger.Main.Infof("Started service '%s' completely.", *csrv.ServiceArn)
		}
	}

	// update
	for _, add := range plan.NewServices {
		currentStack, ok := plan.CurrentServices[add.Name]
		if ok {
			current := currentStack.Service
			logger.Main.Infof("Updating '%s' service on '%s' ...", add.Name, plan.Name)

			var nextDesiredCount int64
			if add.KeepDesiredCount {
				nextDesiredCount = *current.DesiredCount
				logger.Main.Infof("Keep DesiredCount = %d at '%s'", nextDesiredCount, add.Name)
			} else {
				nextDesiredCount = add.DesiredCount
				logger.Main.Infof("Next DesiredCount = %d at '%s'", nextDesiredCount, add.Name)
			}

			params := awsecs.UpdateServiceInput{
				Cluster:        aws.String(plan.Name),
				Service:        aws.String(add.Name),
				DesiredCount:   aws.Int64(nextDesiredCount),
				TaskDefinition: aws.String(add.TaskDefinition),
			}
			if add.MinimumHealthyPercent.Valid && add.MaximumPercent.Valid {
				params.DeploymentConfiguration = &awsecs.DeploymentConfiguration{
					MinimumHealthyPercent: aws.Int64(add.MinimumHealthyPercent.Int64),
					MaximumPercent:        aws.Int64(add.MaximumPercent.Int64),
				}
			}

			svc, err := s.ecsCli.UpdateService(&params)
			if err != nil {
				return err
			}
			logger.Main.Infof("Created service '%v', task-definition is '%v'.", *svc.ServiceArn, *svc.TaskDefinition)
			logger.Main.Infof("Launching task definition '%s' ...", *svc.TaskDefinition)

			if add.AutoScaling != nil {
				asgTarget := add.AutoScaling.Target
				if err := s.appAutoscalingCli.RegisterScalableTarget(plan.Name, add.Name, asgTarget.MinCapacity, asgTarget.MaxCapacity, asgTarget.Role); err != nil {
					return err
				}
				logger.Main.Infof("Update autoscaling MinCapacity:%v MaxCapacity:%v", asgTarget.MinCapacity, asgTarget.MaxCapacity)
			} else if currentStack.AutoScaling != nil {
				resourceID := *currentStack.AutoScaling.ResourceId
				if err := s.appAutoscalingCli.DeregisterScalableTarget(resourceID); err != nil {
					return err
				}
				logger.Main.Infof("Deregistered autoscaling ResourceID:%s", resourceID)
			}

			var targetServiceId string
			if len(svc.Deployments) > 1 {
				for _, dep := range svc.Deployments {
					if *dep.Status == "ACTIVE" {
						targetServiceId = *dep.Id
					}
				}
			} else {
				for _, dep := range svc.Deployments {
					targetServiceId = *dep.Id
				}
			}

			tasks, err := s.ecsCli.ListTasks(plan.Name, add.Name)
			if err != nil {
				return err
			}

			taskIds := []*string{}
			for _, tarn := range tasks.TaskArns {
				tokens := strings.Split(*tarn, "/")
				if len(tokens) == 2 {
					s := tokens[1]
					taskIds = append(taskIds, &s)
				}
			}

			if len(taskIds) > 0 {
				dts, err := s.ecsCli.DescribeTasks(plan.Name, taskIds)
				if err != nil {
					return err
				}

				for _, t := range dts.Tasks {
					if *t.StartedBy == targetServiceId {
						if _, err := s.ecsCli.StopTask(plan.Name, *t.TaskArn); err != nil {
							logger.Main.Warnf("Task '%s' is not found, so cannot stop.", *t.TaskArn)
						} else {
							logger.Main.Infof("Stopped Task '%s'", *t.TaskArn)
						}
					}
				}
			}

			if err := s.waitActiveService(plan.Name, add.Name); err != nil {
				return err
			}
			logger.Main.Infof("Started service '%s' completely.", *svc.ServiceArn)
		}

	}

	return nil
}

func (s ConcreteClusterService) waitStoppingService(cluster string, service string) error {

	for {
		time.Sleep(10 * time.Second)

		result, err := s.ecsCli.DescribeService(cluster, []*string{&service})

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

func (s ConcreteClusterService) waitActiveService(cluster string, service string) error {

	var flag = false
	var taskARNs []*string

	for {
		time.Sleep(10 * time.Second)

		result, err := s.ecsCli.DescribeService(cluster, []*string{aws.String(service)})
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
				reslt, errlt := s.ecsCli.ListTasks(cluster, service)
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

			resdt, errdt := s.ecsCli.DescribeTasks(cluster, taskARNs)
			if errdt != nil {
				return errdt
			}

			watchStatus := s.checkRunningTask(resdt)
			if watchStatus == types.WatchFinish {
				logger.Main.Info("At least one of task has started successfully.")
				return nil
			} else if watchStatus == types.WatchTerminate {
				logger.Main.Error("Stopped watching task, because task has stopped.")
				return errors.New("Task has been stopped for some reason.")
			}

		}
	}
}

func (s ConcreteClusterService) checkRunningTask(dto *awsecs.DescribeTasksOutput) types.TaskWatchStatus {

	logger.Main.Info("Current task conditions as follows:")

	status := []string{}
	for _, task := range dto.Tasks {
		util.Println(fmt.Sprintf("    %s:", *task.TaskArn))
		util.Println(fmt.Sprintf("        LastStatus:%s", s.roundColorStatus(*task.LastStatus)))
		util.Println("        Containers:")

		for _, con := range task.Containers {
			util.Println(fmt.Sprintf("            ----------[%s]----------", *con.Name))
			util.Println(fmt.Sprintf("            ContainerARN:%s", *con.ContainerArn))
			util.Println(fmt.Sprintf("            Status:%s", s.roundColorStatus(*con.LastStatus)))
			util.Println()
		}

		status = append(status, *task.LastStatus)
	}

	// if RUNNING at least one, ecs-formation deals with ok.
	for _, s := range status {
		if s == "RUNNING" {
			return types.WatchFinish
		} else if s == "STOPPED" {
			return types.WatchTerminate
		}
	}

	return types.WatchContinue
}

func (s ConcreteClusterService) roundColorStatus(status string) string {

	if status == "RUNNING" {
		return color.GreenString(status)
	} else if status == "PENDING" {
		return color.YellowString(status)
	} else if status == "STOPPED" {
		return color.RedString(status)
	} else {
		return color.MagentaString(status)
	}
}

func toLoadBalancersNew(values []types.LoadBalancer) []*awsecs.LoadBalancer {

	loadBalancers := []*awsecs.LoadBalancer{}
	for _, lb := range values {
		addElb := awsecs.LoadBalancer{
			ContainerName: &lb.ContainerName,
			ContainerPort: &lb.ContainerPort,
		}
		if lb.Name.Valid {
			addElb.LoadBalancerName = aws.String(lb.Name.String)
		}
		if lb.TargetGroupARN.Valid {
			addElb.TargetGroupArn = aws.String(lb.TargetGroupARN.String)
		}

		loadBalancers = append(loadBalancers, &addElb)
	}
	fmt.Println(loadBalancers)

	return loadBalancers
}
