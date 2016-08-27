package elbv2

import "github.com/aws/aws-sdk-go/service/elbv2"

type MockClient struct {
}

func (c MockClient) DescribeLoadBalancers(names []string) (*elbv2.DescribeLoadBalancersOutput, error) {
	return nil, nil
}
