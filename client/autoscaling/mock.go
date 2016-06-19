package autoscaling

type MockClient struct {
}

func (c MockClient) DescribeAutoScalingGroups(groups []string) (map[string]*autoscaling.Group, error) {

	return nil, nil
}

func (c MockClient) DescribeLoadBalancerState(group string) (map[string]*autoscaling.LoadBalancerState, error) {

	return nil, nil
}

func (c MockClient) AttachLoadBalancers(group string, lb []string) error {

	return nil
}

func (c MockClient) DetachLoadBalancers(group string, lb []string) error {

	return nil
}
