package elb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/stormcat24/ecs-formation/client"
)

type MockClient struct {
}

func (c MockClient) DescribeLoadBalancers(names []string) (*elb.DescribeLoadBalancersOutput, error) {

	params := elb.DescribeLoadBalancersInput{
		LoadBalancerNames: aws.StringSlice(names),
	}

	result, err := c.service.DescribeLoadBalancers(&names)
	if client.IsRateExceeded(err) {
		return c.DescribeLoadBalancers(names)
	}

	return result, err
}

func (c MockClient) RegisterInstancesWithLoadBalancer(name string, instances []*elb.Instance) ([]*elb.Instance, error) {

	params := elb.RegisterInstancesWithLoadBalancerInput{
		LoadBalancerName: aws.String(name),
		Instances:        instances,
	}

	result, err := c.service.RegisterInstancesWithLoadBalancer(&params)
	if client.IsRateExceeded(err) {
		return c.RegisterInstancesWithLoadBalancer(name, instances)
	}

	return result.Instances, err
}

func (c MockClient) DeregisterInstancesFromLoadBalancer(lb string, instances []*elb.Instance) ([]*elb.Instance, error) {

	params := elb.DeregisterInstancesFromLoadBalancerInput{
		LoadBalancerName: aws.String(lb),
		Instances:        instances,
	}

	result, err := c.service.DeregisterInstancesFromLoadBalancer(&params)

	if client.IsRateExceeded(err) {
		return c.DeregisterInstancesFromLoadBalancer(lb, instances)
	}

	return result, err
}
