package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type EcsApi struct {
	service *ecs.ECS
}

// Cluster API
func (self *EcsApi) CreateCluster(clusterName string) (*ecs.CreateClusterOutput, error) {

	params := &ecs.CreateClusterInput{
		ClusterName: aws.String(clusterName),
	}

	result, err := self.service.CreateCluster(params)

	if isRateExceeded(err) {
		return self.CreateCluster(clusterName)
	}

	return result, err
}

func (self *EcsApi) DeleteCluster(clusterName string) (*ecs.DeleteClusterOutput, error) {

	params := &ecs.DeleteClusterInput{
		Cluster: aws.String(clusterName),
	}

	result, err := self.service.DeleteCluster(params)

	if isRateExceeded(err) {
		return self.DeleteCluster(clusterName)
	}

	return result, err
}

func (self *EcsApi) DescribeClusters(clusterNames []*string) (*ecs.DescribeClustersOutput, error) {

	params := &ecs.DescribeClustersInput{
		Clusters: clusterNames,
	}

	result, err := self.service.DescribeClusters(params)

	if isRateExceeded(err) {
		return self.DescribeClusters(clusterNames)
	}

	return result, err
}

func (self *EcsApi) ListClusters(maxResult int64) (*ecs.ListClustersOutput, error) {

	params := &ecs.ListClustersInput{
		MaxResults: &maxResult,
	}

	result, err := self.service.ListClusters(params)
	if isRateExceeded(err) {
		return self.ListClusters(maxResult)
	}

	return result, err
}

func (self *EcsApi) ListContainerInstances(cluster string) (*ecs.ListContainerInstancesOutput, error) {

	params := &ecs.ListContainerInstancesInput{
		Cluster: aws.String(cluster),
	}

	result, err := self.service.ListContainerInstances(params)
	if isRateExceeded(err) {
		return self.ListContainerInstances(cluster)
	}

	return result, err
}

// Service API
func (self *EcsApi) CreateService(cluster string, service string, desiredCount int64, lb []*ecs.LoadBalancer, taskDef string, role string) (*ecs.CreateServiceOutput, error) {

	params := &ecs.CreateServiceInput{
		ServiceName:    aws.String(service),
		Cluster:        aws.String(cluster),
		DesiredCount:   &desiredCount,
		LoadBalancers:  lb,
		TaskDefinition: aws.String(taskDef),
	}

	if role != "" {
		params.Role = aws.String(role)
	}

	result, err := self.service.CreateService(params)
	if isRateExceeded(err) {
		return self.CreateService(cluster, service, desiredCount, lb, taskDef, role)
	}

	return result, err
}

func (self *EcsApi) UpdateService(cluster string, service string, desiredCount int64, taskDef string) (*ecs.UpdateServiceOutput, error) {

	params := &ecs.UpdateServiceInput{
		Cluster:        aws.String(cluster),
		Service:        aws.String(service),
		DesiredCount:   &desiredCount,
		TaskDefinition: aws.String(taskDef),
	}

	result, err := self.service.UpdateService(params)
	if isRateExceeded(err) {
		return self.UpdateService(cluster, service, desiredCount, taskDef)
	}

	return result, err
}

func (self *EcsApi) DescribeService(cluster string, services []*string) (*ecs.DescribeServicesOutput, error) {

	params := &ecs.DescribeServicesInput{
		Cluster:  aws.String(cluster),
		Services: services,
	}

	result, err := self.service.DescribeServices(params)
	if isRateExceeded(err) {
		return self.DescribeService(cluster, services)
	}

	return result, err
}

func (self *EcsApi) DeleteService(cluster string, service string) (*ecs.DeleteServiceOutput, error) {

	params := &ecs.DeleteServiceInput{
		Cluster: aws.String(cluster),
		Service: aws.String(service),
	}

	result, err := self.service.DeleteService(params)
	if isRateExceeded(err) {
		return self.DeleteService(cluster, service)
	}
	return result, err
}

func (self *EcsApi) ListServices(cluster string) (*ecs.ListServicesOutput, error) {

	params := &ecs.ListServicesInput{
		Cluster: aws.String(cluster),
	}

	result, err := self.service.ListServices(params)
	if isRateExceeded(err) {
		return self.ListServices(cluster)
	}
	return result, err
}

// TASK API
func (self *EcsApi) DescribeTaskDefinition(defName string) (*ecs.DescribeTaskDefinitionOutput, error) {

	params := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(defName),
	}

	result, err := self.service.DescribeTaskDefinition(params)
	if isRateExceeded(err) {
		return self.DescribeTaskDefinition(defName)
	}
	return result, err
}

func (self *EcsApi) RegisterTaskDefinition(taskName string, conDefs []*ecs.ContainerDefinition, volumes []*ecs.Volume) (*ecs.RegisterTaskDefinitionOutput, error) {

	params := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: conDefs,
		Family:               aws.String(taskName),
		Volumes:              volumes,
	}

	result, err := self.service.RegisterTaskDefinition(params)
	if isRateExceeded(err) {
		return self.RegisterTaskDefinition(taskName, conDefs, volumes)
	}
	return result, err
}

func (self *EcsApi) DeregisterTaskDefinition(taskName string) (*ecs.DeregisterTaskDefinitionOutput, error) {

	params := &ecs.DeregisterTaskDefinitionInput{
		TaskDefinition: aws.String(taskName),
	}

	result, err := self.service.DeregisterTaskDefinition(params)
	if isRateExceeded(err) {
		return self.DeregisterTaskDefinition(taskName)
	}
	return result, err
}

func (self *EcsApi) ListTasks(cluster string, service string) (*ecs.ListTasksOutput, error) {

	params := &ecs.ListTasksInput{
		Cluster:     aws.String(cluster),
		ServiceName: aws.String(service),
	}

	result, err := self.service.ListTasks(params)
	if isRateExceeded(err) {
		return self.ListTasks(cluster, service)
	}
	return result, err
}

func (self *EcsApi) DescribeTasks(cluster string, tasks []*string) (*ecs.DescribeTasksOutput, error) {

	params := &ecs.DescribeTasksInput{
		Cluster: aws.String(cluster),
		Tasks:   tasks,
	}

	result, err := self.service.DescribeTasks(params)
	if isRateExceeded(err) {
		return self.DescribeTasks(cluster, tasks)
	}
	return result, err
}

func (self *EcsApi) StopTask(cluster string, task string) (*ecs.StopTaskOutput, error) {

	params := &ecs.StopTaskInput{
		Cluster: aws.String(cluster),
		Task:    aws.String(task),
	}

	result, err := self.service.StopTask(params)
	if isRateExceeded(err) {
		return self.StopTask(cluster, task)
	}
	return result, err
}
