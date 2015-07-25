package plan

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/stormcat24/ecs-formation/schema"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

type ServiceUpdatePlan struct {
	Name            string
	InstanceARNs    []*string
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
	PrimaryElb string
	StandbyElb string
}

func (self *BlueGreenPlan) IsBluwWithPrimaryElb() bool {

	for _, lb := range self.Blue.AutoScalingGroup.LoadBalancerNames {
		if *lb == self.PrimaryElb {
			return true
		}
	}

	return false
}

type ServiceSet struct {
	CurrentService *ecs.Service
	NewService *schema.BlueGreenTarget
	AutoScalingGroup *autoscaling.Group
	ClusterUpdatePlan *ServiceUpdatePlan
}