package bluegreen

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/stormcat24/ecs-formation/service"
)

type BlueGreenPlan struct {
	Blue       *ServiceSet
	Green      *ServiceSet
	PrimaryElb string
	StandbyElb string
	ChainElb   []BlueGreenChainElb
}

func (self *BlueGreenPlan) IsBlueWithPrimaryElb() bool {

	for _, lb := range self.Blue.AutoScalingGroup.LoadBalancerNames {
		if *lb == self.PrimaryElb {
			return true
		}
	}

	return false
}

type ServiceSet struct {
	CurrentService    *ecs.Service
	NewService        *BlueGreenTarget
	AutoScalingGroup  *autoscaling.Group
	ClusterUpdatePlan *service.ServiceUpdatePlan
}