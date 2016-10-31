package types

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type BlueGreenPlan struct {
	Blue       *ServiceSet
	Green      *ServiceSet
	PrimaryElb string
	StandbyElb string
	ChainElb   []BlueGreenChainElb
	ElbV2      *BlueGreenElbV2
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
	ElbV2      *BlueGreenElbV2     `yaml:"elbv2"`
}

type BlueGreenChainElb struct {
	PrimaryElb string `yaml:"primary_elb"`
	StandbyElb string `yaml:"standby_elb"`
}

type BlueGreenElbV2 struct {
	TargetGroups []BlueGreenTargetGroupPair `yaml:"target_groups"`
}

type BlueGreenTargetGroupPair struct {
	PrimaryGroup string `yaml:"primary_group"`
	StandbyGroup string `yaml:"standby_group"`
}

type BlueGreenTarget struct {
	Cluster          string `yaml:"cluster"`
	Service          string `yaml:"service"`
	AutoscalingGroup string `yaml:"autoscaling_group"`
}

func (p *BlueGreenPlan) IsBlueWithPrimaryElb() bool {

	if p.ElbV2 != nil && len(p.ElbV2.TargetGroups) > 0 {
		for _, tg := range p.Blue.AutoScalingGroup.TargetGroupARNs {
			topTg := p.ElbV2.TargetGroups[0]
			// TODO name check
			if strings.Contains(*tg, topTg.PrimaryGroup) {
				return true
			}
		}
	} else {
		for _, lb := range p.Blue.AutoScalingGroup.LoadBalancerNames {
			if *lb == p.PrimaryElb {
				return true
			}
		}
	}

	return false
}
