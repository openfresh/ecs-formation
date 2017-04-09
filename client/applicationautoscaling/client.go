package applicationautoscaling

import (
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"

	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stormcat24/ecs-formation/client/util"
)

type Client interface {
	DeleteScalingPolicy(params *applicationautoscaling.DeleteScalingPolicyInput) error
	DeregisterScalableTarget(resourceID string) error
	DescribeScalableTarget(cluster, service string) (*applicationautoscaling.ScalableTarget, error)
	DescribeScalingActivities(params *applicationautoscaling.DescribeScalingActivitiesInput) ([]*applicationautoscaling.ScalingActivity, error)
	DescribeScalingPolicies(params *applicationautoscaling.DescribeScalingPoliciesInput) ([]*applicationautoscaling.ScalingPolicy, error)
	PutScalingPolicy(params *applicationautoscaling.PutScalingPolicyInput) (string, error)
	RegisterScalableTarget(cluster string, service string, min, max uint, role string) error
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

func (c DefaultClient) DescribeScalableTarget(cluster, service string) (*applicationautoscaling.ScalableTarget, error) {

	params := applicationautoscaling.DescribeScalableTargetsInput{
		ServiceNamespace: aws.String("ecs"),
		ResourceIds: aws.StringSlice([]string{
			fmt.Sprintf("service/%s/%s", cluster, service),
		}),
		ScalableDimension: aws.String("ecs:service:DesiredCount"),
	}

	result, err := c.service.DescribeScalableTargets(&params)
	if util.IsRateExceeded(err) {
		return c.DescribeScalableTarget(cluster, service)
	}

	if len(result.ScalableTargets) > 0 {
		return result.ScalableTargets[0], nil
	}

	return nil, err
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

func (c DefaultClient) DeregisterScalableTarget(resourceID string) error {

	input := applicationautoscaling.DeregisterScalableTargetInput{
		ResourceId:        aws.String(resourceID),
		ServiceNamespace:  aws.String("ecs"),
		ScalableDimension: aws.String("ecs:service:DesiredCount"),
	}

	_, err := c.service.DeregisterScalableTarget(&input)
	if util.IsRateExceeded(err) {
		return c.DeregisterScalableTarget(resourceID)
	}

	return err
}

func (c DefaultClient) RegisterScalableTarget(
	cluster string,
	service string,
	min, max uint,
	role string) error {

	input := applicationautoscaling.RegisterScalableTargetInput{
		ServiceNamespace:  aws.String("ecs"),
		ResourceId:        aws.String(fmt.Sprintf("service/%s/%s", cluster, service)),
		ScalableDimension: aws.String("ecs:service:DesiredCount"),
		MinCapacity:       aws.Int64(int64(min)),
		MaxCapacity:       aws.Int64(int64(max)),
		RoleARN:           aws.String(role),
	}

	_, err := c.service.RegisterScalableTarget(&input)
	if util.IsRateExceeded(err) {
		return c.RegisterScalableTarget(cluster, service, min, max, role)
	}

	return err
}
