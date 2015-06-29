package plan

import (
	"github.com/aws/aws-sdk-go/service/ecs"
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
	NewContainers     map[string]*schema.ContainerDefinition
}

type UpdateContainer struct {
	Before *ecs.ContainerDefinition
	After  *schema.ContainerDefinition
}
