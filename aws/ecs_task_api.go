package aws

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type EcsTaskApi struct {
	Credentials *credentials.Credentials
	Region      *string
}

func (self *EcsTaskApi) DescribeTaskDefinition(defName string) (*ecs.DescribeTaskDefinitionOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(defName),
	}

	return svc.DescribeTaskDefinition(params)
}

func (self *EcsTaskApi) RegisterTaskDefinition(taskName string, conDefs []*ecs.ContainerDefinition, volumes []*ecs.Volume) (*ecs.RegisterTaskDefinitionOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: conDefs,
		Family: aws.String(taskName),
		Volumes: volumes,
	}

	return svc.RegisterTaskDefinition(params)
}

func (self *EcsTaskApi) DeregisterTaskDefinition(taskName string) (*ecs.DeregisterTaskDefinitionOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.DeregisterTaskDefinitionInput{
		TaskDefinition: aws.String(taskName),
	}

	return svc.DeregisterTaskDefinition(params)
}

func (self *EcsTaskApi) ListTasks(cluster string, service string) (*ecs.ListTasksOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.ListTasksInput{
		Cluster: aws.String(cluster),
		ServiceName: aws.String(service),
	}

	return svc.ListTasks(params)
}

func (self *EcsTaskApi) DescribeTasks(cluster string, tasks []*string) (*ecs.DescribeTasksOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.DescribeTasksInput{
		Cluster: aws.String(cluster),
		Tasks: tasks,
	}

	return svc.DescribeTasks(params)
}