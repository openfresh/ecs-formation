package elbv2

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"

	"github.com/stormcat24/ecs-formation/client/util"
)

type Client interface {
	DescribeLoadBalancers(names []string) (*elbv2.DescribeLoadBalancersOutput, error)
}

type DefaultClient struct {
	service *elbv2.ELBV2
}

func (c DefaultClient) DescribeLoadBalancers(names []string) (*elbv2.DescribeLoadBalancersOutput, error) {

	params := elbv2.DescribeLoadBalancersInput{
		LoadBalancerArns: aws.StringSlice(names),
	}

	result, err := c.service.DescribeLoadBalancers(&params)
	if util.IsRateExceeded(err) {
		return c.DescribeLoadBalancers(names)
	}

	return result, err
}
