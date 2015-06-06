package plan

import (
	"github.com/awslabs/aws-sdk-go/service/ecs"
	"github.com/stormcat24/ecs-formation/schema"
)

type ClusterUpdatePlan struct {
	Name            string
	CurrentServices map[string]*ecs.Service
	DeleteServices  map[string]*ecs.Service
	UpdateServices  map[string]*UpdateService
	NewServices     map[string]*schema.Service
}

type UpdateService struct {
	Before *ecs.Service
	After  *schema.Service
}

type TaskUpdatePlan struct {
	Name              string
	CurrentContainers map[string]*ecs.ContainerDefinition
	DeleteContainers  map[string]*ecs.ContainerDefinition
	NewContainers     map[string]*schema.ContainerDefinition
	UpdateContainers  map[string]*UpdateContainer
}

type UpdateContainer struct {
	Before *ecs.ContainerDefinition
	After  *schema.ContainerDefinition
}