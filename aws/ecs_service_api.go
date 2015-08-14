package aws

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/aws"
)

type EcsServiceApi struct {
	Credentials *credentials.Credentials
	Region      *string
}

func (self *EcsServiceApi) CreateService(cluster string, service string, desiredCount int64, lb []*ecs.LoadBalancer, taskDef string, role string) (*ecs.CreateServiceOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.CreateServiceInput{
		ServiceName: aws.String(service),
		Cluster: aws.String(cluster),
		DesiredCount: &desiredCount,
		LoadBalancers: lb,
		TaskDefinition: aws.String(taskDef),
	}

	if role != "" {
		params.Role	= aws.String(role)
	}

	return svc.CreateService(params)
}

func (self *EcsServiceApi) UpdateService(cluster string, service string, desiredCount int64, taskDef string) (*ecs.UpdateServiceOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.UpdateServiceInput{
		Cluster: aws.String(cluster),
		Service: aws.String(service),
		DesiredCount: &desiredCount,
		TaskDefinition: aws.String(taskDef),
	}

	return svc.UpdateService(params)
}

func (self *EcsServiceApi) DescribeService(cluster string, services []*string) (*ecs.DescribeServicesOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.DescribeServicesInput{
		Cluster: aws.String(cluster),
		Services: services,
	}

	return svc.DescribeServices(params)
}

func (self *EcsServiceApi) DeleteService(cluster string, service string) (*ecs.DeleteServiceOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.DeleteServiceInput{
		Cluster: aws.String(cluster),
		Service: aws.String(service),
	}

	return svc.DeleteService(params)
}

func (self *EcsServiceApi) ListServices(cluster string) (*ecs.ListServicesOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.ListServicesInput{
		Cluster: aws.String(cluster),
	}

	return svc.ListServices(params)
}