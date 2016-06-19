package autoscaling

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

type MockClient struct {
}

func (c MockClient) DescribeAutoScalingGroups(groups []string) (map[string]*autoscaling.Group, error) {

	return make(map[string]*autoscaling.Group, 0), nil
}

func (c MockClient) DescribeLoadBalancerState(group string) (map[string]*autoscaling.LoadBalancerState, error) {

	return make(map[string]*autoscaling.LoadBalancerState, 0), nil
}

func (c MockClient) AttachLoadBalancers(group string, lb []string) error {

	return nil
}

func (c MockClient) DetachLoadBalancers(group string, lb []string) error {

	return nil
}
