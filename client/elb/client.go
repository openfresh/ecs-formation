package elb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"

	"github.com/openfresh/ecs-formation/client/util"
)

type Client interface {
	DescribeLoadBalancers(names []string) (*elb.DescribeLoadBalancersOutput, error)
	RegisterInstancesWithLoadBalancer(name string, instances []*elb.Instance) ([]*elb.Instance, error)
	DeregisterInstancesFromLoadBalancer(lb string, instances []*elb.Instance) ([]*elb.Instance, error)
}

type DefaultClient struct {
	service *elb.ELB
}

func (c DefaultClient) DescribeLoadBalancers(names []string) (*elb.DescribeLoadBalancersOutput, error) {

	params := elb.DescribeLoadBalancersInput{
		LoadBalancerNames: aws.StringSlice(names),
	}

	result, err := c.service.DescribeLoadBalancers(&params)
	if util.IsRateExceeded(err) {
		return c.DescribeLoadBalancers(names)
	}

	return result, err
}

func (c DefaultClient) RegisterInstancesWithLoadBalancer(name string, instances []*elb.Instance) ([]*elb.Instance, error) {

	params := elb.RegisterInstancesWithLoadBalancerInput{
		LoadBalancerName: aws.String(name),
		Instances:        instances,
	}

	result, err := c.service.RegisterInstancesWithLoadBalancer(&params)
	if util.IsRateExceeded(err) {
		return c.RegisterInstancesWithLoadBalancer(name, instances)
	}

	return result.Instances, err
}

func (c DefaultClient) DeregisterInstancesFromLoadBalancer(lb string, instances []*elb.Instance) ([]*elb.Instance, error) {

	params := elb.DeregisterInstancesFromLoadBalancerInput{
		LoadBalancerName: aws.String(lb),
		Instances:        instances,
	}

	result, err := c.service.DeregisterInstancesFromLoadBalancer(&params)

	if util.IsRateExceeded(err) {
		return c.DeregisterInstancesFromLoadBalancer(lb, instances)
	}

	return result.Instances, err
}
