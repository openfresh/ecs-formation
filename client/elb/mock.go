package elb

import "github.com/aws/aws-sdk-go/service/elb"

type MockClient struct {
}

func (c MockClient) DescribeLoadBalancers(names []string) (*elb.DescribeLoadBalancersOutput, error) {

	return nil, nil
}

func (c MockClient) RegisterInstancesWithLoadBalancer(name string, instances []*elb.Instance) ([]*elb.Instance, error) {

	return make([]*elb.Instance, 0), nil
}

func (c MockClient) DeregisterInstancesFromLoadBalancer(lb string, instances []*elb.Instance) ([]*elb.Instance, error) {

	return make([]*elb.Instance, 0), nil
}
