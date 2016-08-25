package types

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type BlueGreenPlan struct {
	Blue       *ServiceSet
	Green      *ServiceSet
	PrimaryElb string
	StandbyElb string
	ChainElb   []BlueGreenChainElb
}

type ServiceSet struct {
	CurrentService    *ecs.Service
	NewService        *BlueGreenTarget
	AutoScalingGroup  *autoscaling.Group
	ClusterUpdatePlan *ServiceUpdatePlan
}

type BlueGreen struct {
	Blue       BlueGreenTarget     `yaml:"blue"`
	Green      BlueGreenTarget     `yaml:"green"`
	PrimaryElb string              `yaml:"primary_elb"`
	StandbyElb string              `yaml:"standby_elb"`
	ChainElb   []BlueGreenChainElb `yaml:"chain_elb"`
}

type BlueGreenChainElb struct {
	PrimaryElb string `yaml:"primary_elb"`
	StandbyElb string `yaml:"standby_elb"`
}

type BlueGreenTarget struct {
	Cluster          string `yaml:"cluster"`
	Service          string `yaml:"service"`
	AutoscalingGroup string `yaml:"autoscaling_group"`
}

func (p *BlueGreenPlan) IsBlueWithPrimaryElb() bool {

	for _, lb := range p.Blue.AutoScalingGroup.LoadBalancerNames {
		if *lb == p.PrimaryElb {
			return true
		}
	}

	return false
}
