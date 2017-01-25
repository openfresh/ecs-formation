package iam

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/stormcat24/ecs-formation/client/util"
)

type Client interface {
	CreateRole(name string) (*iam.Role, error)
	DeleteRole(name string) error
	CreateInstanceProfile(name string) (*iam.InstanceProfile, error)
	DeleteInstanceProfile(name string) error
	AddRoleToInstanceProfile(instanceProfileName string, roleName string) error
	CreatePolicy(name string, policy string) (*iam.Policy, error)
	DeletePolicy(name string) error
}

type DefaultClient struct {
	service *iam.IAM
}

const (
	defaultAssumePolicy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}`
)

func (c DefaultClient) CreateRole(name string) (*iam.Role, error) {

	params := iam.CreateRoleInput{
		RoleName:                 aws.String(name),
		AssumeRolePolicyDocument: aws.String(defaultAssumePolicy),
	}

	result, err := c.service.CreateRole(&params)
	if util.IsRateExceeded(err) {
		return c.CreateRole(name)
	}

	return result.Role, err
}

func (c DefaultClient) DeleteRole(name string) error {

	params := iam.DeleteRoleInput{
		RoleName: aws.String(name),
	}

	_, err := c.service.DeleteRole(&params)
	if util.IsRateExceeded(err) {
		return c.DeleteRole(name)
	}
	return err
}

func (c DefaultClient) CreateInstanceProfile(name string) (*iam.InstanceProfile, error) {

	params := iam.CreateInstanceProfileInput{
		InstanceProfileName: aws.String(name),
	}

	result, err := c.service.CreateInstanceProfile(&params)
	if util.IsRateExceeded(err) {
		return c.CreateInstanceProfile(name)
	}
	return result.InstanceProfile, err
}

func (c DefaultClient) DeleteInstanceProfile(name string) error {

	params := iam.DeleteInstanceProfileInput{
		InstanceProfileName: aws.String(name),
	}

	_, err := c.service.DeleteInstanceProfile(&params)
	if util.IsRateExceeded(err) {
		return c.DeleteInstanceProfile(name)
	}
	return err
}

func (c DefaultClient) AddRoleToInstanceProfile(instanceProfileName string, roleName string) error {

	params := iam.AddRoleToInstanceProfileInput{
		InstanceProfileName: aws.String(instanceProfileName),
		RoleName:            aws.String(roleName),
	}

	_, err := c.service.AddRoleToInstanceProfile(&params)
	if util.IsRateExceeded(err) {
		return c.AddRoleToInstanceProfile(instanceProfileName, roleName)
	}
	return err
}

func (c DefaultClient) CreatePolicy(name string, policy string) (*iam.Policy, error) {

	params := iam.CreatePolicyInput{
		PolicyName:     aws.String(name),
		PolicyDocument: aws.String(policy),
	}

	result, err := c.service.CreatePolicy(&params)
	if util.IsRateExceeded(err) {
		return c.CreatePolicy(name, policy)
	}

	return result.Policy, err
}

func (c DefaultClient) DeletePolicy(policyArn string) error {

	params := iam.DeletePolicyInput{
		PolicyArn: aws.String(policyArn),
	}

	_, err := c.service.DeletePolicy(&params)
	if util.IsRateExceeded(err) {
		return c.DeleteRole(policyArn)
	}

	return err
}
