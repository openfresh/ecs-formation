package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/stormcat24/ecs-formation/util"
)

type AutoscalingApi struct {
	service *autoscaling.AutoScaling
}

func (self *AutoscalingApi) DescribeAutoScalingGroups(groups []string) (map[string]*autoscaling.Group, error) {

	params := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: util.ConvertPointerString(groups),
	}

	asgmap := map[string]*autoscaling.Group{}
	result, err := self.service.DescribeAutoScalingGroups(params)
	if isRateExceeded(err) {
		return self.DescribeAutoScalingGroups(groups)
	}

	if err != nil {
		return asgmap, err
	}

	for _, asg := range result.AutoScalingGroups {
		asgmap[*asg.AutoScalingGroupName] = asg
	}

	return asgmap, nil
}

func (self *AutoscalingApi) DescribeLoadBalancerState(group string) (map[string]*autoscaling.LoadBalancerState, error) {

	params := &autoscaling.DescribeLoadBalancersInput{
		AutoScalingGroupName: aws.String(group),
	}

	lbmap := map[string]*autoscaling.LoadBalancerState{}
	result, err := self.service.DescribeLoadBalancers(params)
	if isRateExceeded(err) {
		return self.DescribeLoadBalancerState(group)
	}

	if err != nil {
		return lbmap, err
	}

	for _, lbs := range result.LoadBalancers {
		lbmap[*lbs.LoadBalancerName] = lbs
	}

	return lbmap, nil
}

func (self *AutoscalingApi) AttachLoadBalancers(group string, lb []string) (*autoscaling.AttachLoadBalancersOutput, error) {

	params := &autoscaling.AttachLoadBalancersInput{
		AutoScalingGroupName: aws.String(group),
		LoadBalancerNames:    util.ConvertPointerString(lb),
	}

	result, err := self.service.AttachLoadBalancers(params)
	if isRateExceeded(err) {
		return self.AttachLoadBalancers(group, lb)
	}
	return result, err
}

func (self *AutoscalingApi) DetachLoadBalancers(group string, lb []string) (*autoscaling.DetachLoadBalancersOutput, error) {

	params := &autoscaling.DetachLoadBalancersInput{
		AutoScalingGroupName: aws.String(group),
		LoadBalancerNames:    util.ConvertPointerString(lb),
	}

	result, err := self.service.DetachLoadBalancers(params)
	if isRateExceeded(err) {
		return self.DetachLoadBalancers(group, lb)
	}
	return result, err
}
