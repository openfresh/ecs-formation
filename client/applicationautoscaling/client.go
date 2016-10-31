package applicationautoscaling

import (
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"

	"github.com/stormcat24/ecs-formation/client/util"
)

type Client interface {
	DeleteScalingPolicy(params *applicationautoscaling.DeleteScalingPolicyInput) error
	DeregisterScalableTarget(params *applicationautoscaling.DeregisterScalableTargetInput) error
	DescribeScalableTargets(params *applicationautoscaling.DescribeScalableTargetsInput) ([]*applicationautoscaling.ScalableTarget, error)
	DescribeScalingActivities(params *applicationautoscaling.DescribeScalingActivitiesInput) ([]*applicationautoscaling.ScalingActivity, error)
	DescribeScalingPolicies(params *applicationautoscaling.DescribeScalingPoliciesInput) ([]*applicationautoscaling.ScalingPolicy, error)
	PutScalingPolicy(params *applicationautoscaling.PutScalingPolicyInput) (string, error)
	RegisterScalableTarget(params *applicationautoscaling.RegisterScalableTargetInput) error
}

type DefaultClient struct {
	service *applicationautoscaling.ApplicationAutoScaling
}

func (c DefaultClient) DeleteScalingPolicy(params *applicationautoscaling.DeleteScalingPolicyInput) error {

	_, err := c.service.DeleteScalingPolicy(params)
	if util.IsRateExceeded(err) {
		return c.DeleteScalingPolicy(params)
	}

	return err
}

func (c DefaultClient) PutScalingPolicy(params *applicationautoscaling.PutScalingPolicyInput) (string, error) {

	result, err := c.service.PutScalingPolicy(params)
	if util.IsRateExceeded(err) {
		return c.PutScalingPolicy(params)
	}

	return *result.PolicyARN, err
}

func (c DefaultClient) DescribeScalableTargets(params *applicationautoscaling.DescribeScalableTargetsInput) ([]*applicationautoscaling.ScalableTarget, error) {

	result, err := c.service.DescribeScalableTargets(params)
	if util.IsRateExceeded(err) {
		return c.DescribeScalableTargets(params)
	}

	return result.ScalableTargets, err
}

func (c DefaultClient) DescribeScalingActivities(params *applicationautoscaling.DescribeScalingActivitiesInput) ([]*applicationautoscaling.ScalingActivity, error) {

	result, err := c.service.DescribeScalingActivities(params)
	if util.IsRateExceeded(err) {
		return c.DescribeScalingActivities(params)
	}

	return result.ScalingActivities, err
}

func (c DefaultClient) DescribeScalingPolicies(params *applicationautoscaling.DescribeScalingPoliciesInput) ([]*applicationautoscaling.ScalingPolicy, error) {

	result, err := c.service.DescribeScalingPolicies(params)
	if util.IsRateExceeded(err) {
		return c.DescribeScalingPolicies(params)
	}

	return result.ScalingPolicies, err
}

func (c DefaultClient) DeregisterScalableTarget(params *applicationautoscaling.DeregisterScalableTargetInput) error {

	_, err := c.service.DeregisterScalableTarget(params)
	if util.IsRateExceeded(err) {
		return c.DeregisterScalableTarget(params)
	}

	return err
}

func (c DefaultClient) RegisterScalableTarget(params *applicationautoscaling.RegisterScalableTargetInput) error {

	_, err := c.service.RegisterScalableTarget(params)
	if util.IsRateExceeded(err) {
		return c.RegisterScalableTarget(params)
	}

	return err
}
