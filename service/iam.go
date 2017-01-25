package service

import (
	"github.com/stormcat24/ecs-formation/client"
	"github.com/stormcat24/ecs-formation/client/iam"
)

type IamService interface {
	SearchRoles() error
}

type ConcreteIamService struct {
	iamCli     iam.Client
	projectDir string
	targetRole string
	params     map[string]string
}

func NewIamService(projectDir string, targetRole string, params map[string]string) (IamService, error) {

	service := ConcreteIamService{
		iamCli:     client.AWSCli.IAM,
		projectDir: projectDir,
		targetRole: targetRole,
		params:     params,
	}

	return &service, nil
}

func (ss ConcreteIamService) SearchRoles() error {
	return nil
}
