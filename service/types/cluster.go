package types

import (
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go/service/ecs"
	"gopkg.in/guregu/null.v3"
)

type TaskWatchStatus int

const (
	WatchContinue TaskWatchStatus = iota
	WatchFinish
	WatchTerminate
)

type Cluster struct {
	Name     string
	Services map[string]Service
}

type Service struct {
	Name                  string
	TaskDefinition        string                `yaml:"task_definition"`
	DesiredCount          int64                 `yaml:"desired_count"`
	KeepDesiredCount      bool                  `yaml:"keep_desired_count"`
	LoadBalancers         []LoadBalancer        `yaml:"load_balancers"`
	MinimumHealthyPercent null.Int              `yaml:"minimum_healthy_percent"`
	MaximumPercent        null.Int              `yaml:"maximum_percent"`
	Role                  string                `yaml:"role"`
	AutoScaling           *AutoScaling          `yaml:"autoscaling"`
	PlacementConstraints  []PlacementConstraint `yaml:"placement_constraints"`
	PlacementStrategy     []PlacementStrategy   `yaml:"placement_strategy"`
}

type LoadBalancer struct {
	Name           null.String `yaml:"name"`
	ContainerName  string      `yaml:"container_name"`
	ContainerPort  int64       `yaml:"container_port"`
	TargetGroupARN null.String `yaml:"target_group_arn"`
}

type ServiceStack struct {
	Service     *ecs.Service
	AutoScaling *applicationautoscaling.ScalableTarget
}

type ServiceUpdatePlan struct {
	Name            string
	InstanceARNs    []*string
	CurrentServices map[string]*ServiceStack
	NewServices     map[string]*Service
}

type AutoScaling struct {
	Target *ServiceScalableTarget `yaml:"target"`
}

type ServiceScalableTarget struct {
	MinCapacity uint   `yaml:"min_capacity"`
	MaxCapacity uint   `yaml:"max_capacity"`
	Role        string `yaml:"role"`
}

type PlacementConstraint struct {
	Expression string `yaml:"expression"`
	Type       string `yaml:"type"`
}

type PlacementStrategy struct {
	Field string `yaml:"field"`
	Type  string `yaml:"type"`
}
