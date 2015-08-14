package service

import "github.com/aws/aws-sdk-go/service/ecs"

type ServiceUpdatePlan struct {
	Name            string
	InstanceARNs    []*string
	CurrentServices map[string]*ecs.Service
	NewServices     map[string]*Service
}