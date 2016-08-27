package ecs

import "github.com/aws/aws-sdk-go/service/ecs"

type MockClient struct {
}

func (c MockClient) CreateCluster(cluster string) (*ecs.Cluster, error) {
	return nil, nil
}

func (c MockClient) DeleteCluster(cluster string) (*ecs.Cluster, error) {
	return nil, nil
}

func (c MockClient) DescribeClusters(clusters []*string) (*ecs.DescribeClustersOutput, error) {
	return nil, nil
}

func (c MockClient) ListClusters(maxResult int) (*ecs.ListClustersOutput, error) {
	return nil, nil
}

func (c MockClient) ListContainerInstances(cluster string) (*ecs.ListContainerInstancesOutput, error) {
	return nil, nil
}

func (c MockClient) CreateService(params *ecs.CreateServiceInput) (*ecs.Service, error) {
	return nil, nil
}

func (c MockClient) UpdateService(params *ecs.UpdateServiceInput) (*ecs.Service, error) {
	return nil, nil
}

func (c MockClient) DescribeService(cluster string, services []*string) (*ecs.DescribeServicesOutput, error) {
	return nil, nil
}

func (c MockClient) DeleteService(cluster string, service string) (*ecs.Service, error) {
	return nil, nil
}

func (c MockClient) ListServices(cluster string) (*ecs.ListServicesOutput, error) {
	return nil, nil
}

func (c MockClient) DescribeTaskDefinition(td string) (*ecs.TaskDefinition, error) {
	return nil, nil
}

func (c MockClient) RegisterTaskDefinition(taskName string, containers []*ecs.ContainerDefinition, volumes []*ecs.Volume) (*ecs.TaskDefinition, error) {
	return nil, nil
}

func (c MockClient) DeregisterTaskDefinition(taskName string) (*ecs.TaskDefinition, error) {
	return nil, nil
}

func (c MockClient) ListTasks(cluster string, service string) (*ecs.ListTasksOutput, error) {
	return nil, nil
}

func (c MockClient) DescribeTasks(cluster string, tasks []*string) (*ecs.DescribeTasksOutput, error) {
	return nil, nil
}

func (c MockClient) StopTask(cluster string, task string) (*ecs.Task, error) {
	return nil, nil
}
