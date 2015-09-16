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

	result, err := self.service.DescribeLoadBalancers(params)
	if isRateExceeded(err) {
		return self.DescribeLoadBalancers(names)
	}

	return result, err
}

func (self *ElbApi) RegisterInstancesWithLoadBalancer(name string, instances []*elb.Instance) (*elb.RegisterInstancesWithLoadBalancerOutput, error) {

	params := &elb.RegisterInstancesWithLoadBalancerInput{
		LoadBalancerName: aws.String(name),
		Instances:        instances,
	}

	result, err := self.service.RegisterInstancesWithLoadBalancer(params)
	if isRateExceeded(err) {
		return self.RegisterInstancesWithLoadBalancer(name, instances)
	}

	return result, err
}

func (self *ElbApi) DeregisterInstancesFromLoadBalancer(lb string, instances []*elb.Instance) (*elb.DeregisterInstancesFromLoadBalancerOutput, error) {

	params := &elb.DeregisterInstancesFromLoadBalancerInput{
		LoadBalancerName: aws.String(lb),
		Instances:        instances,
	}

	result, err := self.service.DeregisterInstancesFromLoadBalancer(params)

	if isRateExceeded(err) {
		return self.DeregisterInstancesFromLoadBalancer(lb, instances)
	}

	return result, err
}
