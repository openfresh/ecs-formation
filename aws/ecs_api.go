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

	return self.service.CreateCluster(params)
}

func (self *EcsApi) DeleteCluster(clusterName string) (*ecs.DeleteClusterOutput, error) {

	params := &ecs.DeleteClusterInput{
		Cluster: aws.String(clusterName),
	}

	return self.service.DeleteCluster(params)
}

func (self *EcsApi) DescribeClusters(clusterNames []*string) (*ecs.DescribeClustersOutput, error) {

	params := &ecs.DescribeClustersInput{
		Clusters: clusterNames,
	}

	return self.service.DescribeClusters(params)
}

func (self *EcsApi) ListClusters(maxResult int64) (*ecs.ListClustersOutput, error) {

	params := &ecs.ListClustersInput{
		MaxResults: &maxResult,
	}

	return self.service.ListClusters(params)
}

func (self *EcsApi) ListContainerInstances(cluster string) (*ecs.ListContainerInstancesOutput, error) {

	params := &ecs.ListContainerInstancesInput{
		Cluster: aws.String(cluster),
	}

	return self.service.ListContainerInstances(params)
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

	return self.service.CreateService(params)
}

func (self *EcsApi) UpdateService(cluster string, service string, desiredCount int64, taskDef string) (*ecs.UpdateServiceOutput, error) {

	params := &ecs.UpdateServiceInput{
		Cluster:        aws.String(cluster),
		Service:        aws.String(service),
		DesiredCount:   &desiredCount,
		TaskDefinition: aws.String(taskDef),
	}

	return self.service.UpdateService(params)
}

func (self *EcsApi) DescribeService(cluster string, services []*string) (*ecs.DescribeServicesOutput, error) {

	params := &ecs.DescribeServicesInput{
		Cluster:  aws.String(cluster),
		Services: services,
	}

	return self.service.DescribeServices(params)
}

func (self *EcsApi) DeleteService(cluster string, service string) (*ecs.DeleteServiceOutput, error) {

	params := &ecs.DeleteServiceInput{
		Cluster: aws.String(cluster),
		Service: aws.String(service),
	}

	return self.service.DeleteService(params)
}

func (self *EcsApi) ListServices(cluster string) (*ecs.ListServicesOutput, error) {

	params := &ecs.ListServicesInput{
		Cluster: aws.String(cluster),
	}

	return self.service.ListServices(params)
}

// TASK API
func (self *EcsApi) DescribeTaskDefinition(defName string) (*ecs.DescribeTaskDefinitionOutput, error) {

	params := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(defName),
	}

	return self.service.DescribeTaskDefinition(params)
}

func (self *EcsApi) RegisterTaskDefinition(taskName string, conDefs []*ecs.ContainerDefinition, volumes []*ecs.Volume) (*ecs.RegisterTaskDefinitionOutput, error) {

	params := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: conDefs,
		Family:               aws.String(taskName),
		Volumes:              volumes,
	}

	return self.service.RegisterTaskDefinition(params)
}

func (self *EcsApi) DeregisterTaskDefinition(taskName string) (*ecs.DeregisterTaskDefinitionOutput, error) {

	params := &ecs.DeregisterTaskDefinitionInput{
		TaskDefinition: aws.String(taskName),
	}

	return self.service.DeregisterTaskDefinition(params)
}

func (self *EcsApi) ListTasks(cluster string, service string) (*ecs.ListTasksOutput, error) {

	params := &ecs.ListTasksInput{
		Cluster:     aws.String(cluster),
		ServiceName: aws.String(service),
	}

	return self.service.ListTasks(params)
}

func (self *EcsApi) DescribeTasks(cluster string, tasks []*string) (*ecs.DescribeTasksOutput, error) {

	params := &ecs.DescribeTasksInput{
		Cluster: aws.String(cluster),
		Tasks:   tasks,
	}

	return self.service.DescribeTasks(params)
}
