package autoscaling

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"

	"github.com/stormcat24/ecs-formation/client/util"
)

type Client interface {
	DescribeAutoScalingGroups(groups []string) (map[string]*autoscaling.Group, error)
	DescribeLoadBalancerState(group string) (map[string]*autoscaling.LoadBalancerState, error)
	AttachLoadBalancers(group string, lb []string) error
	DetachLoadBalancers(group string, lb []string) error
}

type DefaultClient struct {
	service *autoscaling.AutoScaling
}

func (c DefaultClient) DescribeAutoScalingGroups(groups []string) (map[string]*autoscaling.Group, error) {

	params := autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: aws.StringSlice(groups),
	}

	asgmap := make(map[string]*autoscaling.Group, 0)
	result, err := c.service.DescribeAutoScalingGroups(&params)
	if util.IsRateExceeded(err) {
		return c.DescribeAutoScalingGroups(groups)
	}

	if err != nil {
		return nil, err
	}

	for _, asg := range result.AutoScalingGroups {
		asgmap[*asg.AutoScalingGroupName] = asg
	}

	return asgmap, nil
}

func (c DefaultClient) DescribeLoadBalancerState(group string) (map[string]*autoscaling.LoadBalancerState, error) {

	params := autoscaling.DescribeLoadBalancersInput{
		AutoScalingGroupName: aws.String(group),
	}

	lbmap := map[string]*autoscaling.LoadBalancerState{}
	result, err := c.service.DescribeLoadBalancers(&params)
	if util.IsRateExceeded(err) {
		return c.DescribeLoadBalancerState(group)
	}

	if err != nil {
		return lbmap, err
	}

	for _, lbs := range result.LoadBalancers {
		lbmap[*lbs.LoadBalancerName] = lbs
	}

	return lbmap, nil
}

func (c DefaultClient) AttachLoadBalancers(group string, lb []string) error {

	params := autoscaling.AttachLoadBalancersInput{
		AutoScalingGroupName: aws.String(group),
		LoadBalancerNames:    aws.StringSlice(lb),
	}

	_, err := c.service.AttachLoadBalancers(&params)
	if util.IsRateExceeded(err) {
		return c.AttachLoadBalancers(group, lb)
	}

	return err
}

func (c DefaultClient) DetachLoadBalancers(group string, lb []string) error {

	params := autoscaling.DetachLoadBalancersInput{
		AutoScalingGroupName: aws.String(group),
		LoadBalancerNames:    aws.StringSlice(lb),
	}

	_, err := c.service.DetachLoadBalancers(&params)
	if util.IsRateExceeded(err) {
		return c.DetachLoadBalancers(group, lb)
	}

	return err
}
