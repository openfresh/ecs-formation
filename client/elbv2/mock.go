package elbv2

import "github.com/aws/aws-sdk-go/service/elbv2"

type MockClient struct {
}

func (c MockClient) DescribeLoadBalancers(names []string) (*elbv2.DescribeLoadBalancersOutput, error) {
	return nil, nil
}

func (c MockClient) CreateRule(params *elbv2.CreateRuleInput) ([]*elbv2.Rule, error) {
	return []*elbv2.Rule{}, nil
}

func (c MockClient) DeleteRule(ruleArn string) error {

	return nil
}

func (c MockClient) DescribeRule(params *elbv2.DescribeRulesInput) ([]*elbv2.Rule, error) {

	return []*elbv2.Rule{}, nil
}

func (c MockClient) ModifyRule(params *elbv2.ModifyRuleInput) ([]*elbv2.Rule, error) {

	return []*elbv2.Rule{}, nil
}

func (c MockClient) CreateTargetGroup(params *elbv2.CreateTargetGroupInput) ([]*elbv2.TargetGroup, error) {

	return []*elbv2.TargetGroup{}, nil
}

func (c MockClient) DeleteTargetGroup(targetGroupArn string) error {

	return nil
}

func (c MockClient) DescribeTargetGroup(groupNames []string) (map[string]*elbv2.TargetGroup, error) {

	return map[string]*elbv2.TargetGroup{}, nil
}

func (c MockClient) ModifyTargetGroup(params *elbv2.ModifyTargetGroupInput) ([]*elbv2.TargetGroup, error) {

	return []*elbv2.TargetGroup{}, nil
}
