package aws

import (
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stormcat24/ecs-formation/util"
)

type ElbApi struct {
	Credentials *credentials.Credentials
	Region      *string
}

func (self *ElbApi) DescribeLoadBalancers(names []string) (*elb.DescribeLoadBalancersOutput, error) {

	svc := elb.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: util.ConvertPointerString(names),
	}

	return svc.DescribeLoadBalancers(params)
}

func (self *ElbApi) RegisterInstancesWithLoadBalancer(name string, instances []*elb.Instance) (*elb.RegisterInstancesWithLoadBalancerOutput, error) {

	svc := elb.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &elb.RegisterInstancesWithLoadBalancerInput{
		LoadBalancerName: aws.String(name),
		Instances: instances,
	}

	return svc.RegisterInstancesWithLoadBalancer(params)
}

func (self *ElbApi) DeregisterInstancesFromLoadBalancer(lb string, instances []*elb.Instance) (*elb.DeregisterInstancesFromLoadBalancerOutput, error) {

	svc := elb.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &elb.DeregisterInstancesFromLoadBalancerInput{
		LoadBalancerName: aws.String(lb),
		Instances: instances,
	}

	return svc.DeregisterInstancesFromLoadBalancer(params)
}