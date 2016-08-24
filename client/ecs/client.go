package ecs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/pkg/errors"
	"github.com/stormcat24/ecs-formation/client/util"
)

type Client interface {
	CreateCluster(cluster string) (*ecs.Cluster, error)
	DeleteCluster(cluster string) (*ecs.Cluster, error)
	DescribeClusters(clusters []*string) (*ecs.DescribeClustersOutput, error)
	ListClusters(maxResult int) (*ecs.ListClustersOutput, error)
	ListContainerInstances(cluster string) (*ecs.ListContainerInstancesOutput, error)
	CreateService(params *ecs.CreateServiceInput) (*ecs.Service, error)
	UpdateService(cluster string, service string, desiredCount int, taskDef string) (*ecs.Service, error)
	DescribeService(cluster string, services []*string) (*ecs.DescribeServicesOutput, error)
	DeleteService(cluster string, service string) (*ecs.Service, error)
	ListServices(cluster string) (*ecs.ListServicesOutput, error)
	DescribeTaskDefinition(td string) (*ecs.TaskDefinition, error)
	RegisterTaskDefinition(taskName string, containers []*ecs.ContainerDefinition, volumes []*ecs.Volume) (*ecs.TaskDefinition, error)
	DeregisterTaskDefinition(taskName string) (*ecs.TaskDefinition, error)
	ListTasks(cluster string, service string) (*ecs.ListTasksOutput, error)
	DescribeTasks(cluster string, tasks []*string) (*ecs.DescribeTasksOutput, error)
	StopTask(cluster string, task string) (*ecs.Task, error)
}

type DefaultClient struct {
	service *ecs.ECS
}

func (c DefaultClient) CreateCluster(cluster string) (*ecs.Cluster, error) {

	params := ecs.CreateClusterInput{
		ClusterName: aws.String(cluster),
	}

	result, err := c.service.CreateCluster(&params)
	if util.IsRateExceeded(err) {
		return c.CreateCluster(cluster)
	}
	return result.Cluster, err
}

func (c DefaultClient) DeleteCluster(cluster string) (*ecs.Cluster, error) {

	params := ecs.DeleteClusterInput{
		Cluster: aws.String(cluster),
	}

	result, err := c.service.DeleteCluster(&params)
	if util.IsRateExceeded(err) {
		return c.DeleteCluster(cluster)
	}

	return result.Cluster, err
}

func (c DefaultClient) DescribeClusters(clusters []*string) (*ecs.DescribeClustersOutput, error) {

	params := ecs.DescribeClustersInput{
		Clusters: clusters,
	}

	result, err := c.service.DescribeClusters(&params)
	if util.IsRateExceeded(err) {
		return c.DescribeClusters(clusters)
	}

	return result, err
}

func (c DefaultClient) ListClusters(maxResult int) (*ecs.ListClustersOutput, error) {

	params := ecs.ListClustersInput{
		MaxResults: aws.Int64(int64(maxResult)),
	}

	result, err := c.service.ListClusters(&params)
	if util.IsRateExceeded(err) {
		return c.ListClusters(maxResult)
	}

	return result, err
}

func (c DefaultClient) ListContainerInstances(cluster string) (*ecs.ListContainerInstancesOutput, error) {

	params := ecs.ListContainerInstancesInput{
		Cluster: aws.String(cluster),
	}

	result, err := c.service.ListContainerInstances(&params)
	if util.IsRateExceeded(err) {
		return c.ListContainerInstances(cluster)
	}

	return result, err
}

func (c DefaultClient) CreateService(params *ecs.CreateServiceInput) (*ecs.Service, error) {

	result, err := c.service.CreateService(params)
	if util.IsRateExceeded(err) {
		return c.CreateService(params)
	}

	return result.Service, err
}

func (c DefaultClient) UpdateService(cluster string, service string, desiredCount int, taskDef string) (*ecs.Service, error) {

	params := ecs.UpdateServiceInput{
		Cluster:        aws.String(cluster),
		Service:        aws.String(service),
		DesiredCount:   aws.Int64(int64(desiredCount)),
		TaskDefinition: aws.String(taskDef),
	}

	result, err := c.service.UpdateService(&params)
	if util.IsRateExceeded(err) {
		return c.UpdateService(cluster, service, desiredCount, taskDef)
	}

	return result.Service, err
}

func (c DefaultClient) DescribeService(cluster string, services []*string) (*ecs.DescribeServicesOutput, error) {

	params := ecs.DescribeServicesInput{
		Cluster:  aws.String(cluster),
		Services: services,
	}

	result, err := c.service.DescribeServices(&params)
	if util.IsRateExceeded(err) {
		return c.DescribeService(cluster, services)
	}

	return result, err
}

func (c DefaultClient) DeleteService(cluster string, service string) (*ecs.Service, error) {

	params := ecs.DeleteServiceInput{
		Cluster: aws.String(cluster),
		Service: aws.String(service),
	}

	result, err := c.service.DeleteService(&params)
	if util.IsRateExceeded(err) {
		return c.DeleteService(cluster, service)
	}

	return result.Service, err
}

func (c DefaultClient) ListServices(cluster string) (*ecs.ListServicesOutput, error) {

	params := ecs.ListServicesInput{
		Cluster: aws.String(cluster),
	}

	result, err := c.service.ListServices(&params)
	if util.IsRateExceeded(err) {
		return c.ListServices(cluster)
	}

	return result, err
}

func (c DefaultClient) DescribeTaskDefinition(td string) (*ecs.TaskDefinition, error) {

	params := ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(td),
	}

	result, err := c.service.DescribeTaskDefinition(&params)
	if util.IsRateExceeded(err) {
		return c.DescribeTaskDefinition(td)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "Describe ECS Cluster '%s' is failed", td)
	}

	return result.TaskDefinition, nil
}

func (c DefaultClient) RegisterTaskDefinition(taskName string, containers []*ecs.ContainerDefinition, volumes []*ecs.Volume) (*ecs.TaskDefinition, error) {
	params := ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: containers,
		Family:               aws.String(taskName),
		Volumes:              volumes,
	}

	result, err := c.service.RegisterTaskDefinition(&params)
	if util.IsRateExceeded(err) {
		return c.RegisterTaskDefinition(taskName, containers, volumes)
	}
	return result.TaskDefinition, err
}

func (c DefaultClient) DeregisterTaskDefinition(taskName string) (*ecs.TaskDefinition, error) {

	params := ecs.DeregisterTaskDefinitionInput{
		TaskDefinition: aws.String(taskName),
	}

	result, err := c.service.DeregisterTaskDefinition(&params)
	if util.IsRateExceeded(err) {
		return c.DeregisterTaskDefinition(taskName)
	}
	return result.TaskDefinition, err
}

func (c DefaultClient) ListTasks(cluster string, service string) (*ecs.ListTasksOutput, error) {

	params := ecs.ListTasksInput{
		Cluster:     aws.String(cluster),
		ServiceName: aws.String(service),
	}

	result, err := c.service.ListTasks(&params)
	if util.IsRateExceeded(err) {
		return c.ListTasks(cluster, service)
	}

	return result, err
}

func (c DefaultClient) DescribeTasks(cluster string, tasks []*string) (*ecs.DescribeTasksOutput, error) {

	params := ecs.DescribeTasksInput{
		Cluster: aws.String(cluster),
		Tasks:   tasks,
	}

	result, err := c.service.DescribeTasks(&params)
	if util.IsRateExceeded(err) {
		return c.DescribeTasks(cluster, tasks)
	}

	return result, err
}

func (c DefaultClient) StopTask(cluster string, task string) (*ecs.Task, error) {

	params := ecs.StopTaskInput{
		Cluster: aws.String(cluster),
		Task:    aws.String(task),
	}

	result, err := c.service.StopTask(&params)
	if util.IsRateExceeded(err) {
		return c.StopTask(cluster, task)
	}
	return result.Task, err
}
