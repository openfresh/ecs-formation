package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/stormcat24/ecs-formation/util"
)

type ElbApi struct {
	service *elb.ELB
}

func (self *ElbApi) DescribeLoadBalancers(names []string) (*elb.DescribeLoadBalancersOutput, error) {

	params := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: util.ConvertPointerString(names),
	}

	return self.service.DescribeLoadBalancers(params)
}

func (self *ElbApi) RegisterInstancesWithLoadBalancer(name string, instances []*elb.Instance) (*elb.RegisterInstancesWithLoadBalancerOutput, error) {

	params := &elb.RegisterInstancesWithLoadBalancerInput{
		LoadBalancerName: aws.String(name),
		Instances:        instances,
	}

	return self.service.RegisterInstancesWithLoadBalancer(params)
}

func (self *ElbApi) DeregisterInstancesFromLoadBalancer(lb string, instances []*elb.Instance) (*elb.DeregisterInstancesFromLoadBalancerOutput, error) {

	params := &elb.DeregisterInstancesFromLoadBalancerInput{
		LoadBalancerName: aws.String(lb),
		Instances:        instances,
	}

	return self.service.DeregisterInstancesFromLoadBalancer(params)
}
