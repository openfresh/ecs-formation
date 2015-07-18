package plan

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/stormcat24/ecs-formation/schema"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

type ClusterUpdatePlan struct {
	Name            string
	CurrentServices map[string]*ecs.Service
	NewServices     map[string]*schema.Service
}

type TaskUpdatePlan struct {
	Name              string
	NewContainers     map[string]*schema.ContainerDefinition
}

type UpdateContainer struct {
	Before *ecs.ContainerDefinition
	After  *schema.ContainerDefinition
}

type BlueGreenPlan struct {
	Blue *ServiceSet
	Green *ServiceSet
}

type ServiceSet struct {
	CurrentService *ecs.Service
	NewService *schema.BlueGreenTarget
	AutoScalingGroup *autoscaling.Group
	ClusterUpdatePlan *ClusterUpdatePlan
	LoadBalancer string
}


func (self *ServiceSet) HasOwnElb() bool {

	for _, lb := range self.AutoScalingGroup.LoadBalancerNames {
		if *lb == self.LoadBalancer {
			return true
		}
	}

	return false
}