package applicationautoscaling

import (
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
)

type MockClient struct {
}

func (c MockClient) DeleteScalingPolicy(params *applicationautoscaling.DeleteScalingPolicyInput) error {

	return nil
}

func (c MockClient) PutScalingPolicy(params *applicationautoscaling.PutScalingPolicyInput) (string, error) {

	return "", nil
}

func (c MockClient) DescribeScalableTargets(params *applicationautoscaling.DescribeScalableTargetsInput) ([]*applicationautoscaling.ScalableTarget, error) {

	return []*applicationautoscaling.ScalableTarget{}, nil
}

func (c MockClient) DescribeScalingActivities(params *applicationautoscaling.DescribeScalingActivitiesInput) ([]*applicationautoscaling.ScalingActivity, error) {

	return []*applicationautoscaling.ScalingActivity{}, nil
}

func (c MockClient) DescribeScalingPolicies(params *applicationautoscaling.DescribeScalingPoliciesInput) ([]*applicationautoscaling.ScalingPolicy, error) {

	return []*applicationautoscaling.ScalingPolicy{}, nil
}

func (c MockClient) DeregisterScalableTarget(params *applicationautoscaling.DeregisterScalableTargetInput) error {

	return nil
}

func (c MockClient) RegisterScalableTarget(params *applicationautoscaling.RegisterScalableTargetInput) error {

	return nil
}
