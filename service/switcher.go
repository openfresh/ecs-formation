package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/fatih/color"
	"github.com/stormcat24/ecs-formation/client"
	"github.com/stormcat24/ecs-formation/logger"
	"github.com/stormcat24/ecs-formation/service/types"
)

type ELBSwitcher interface {
	Apply(clusterService ClusterService, bgplan *types.BlueGreenPlan, nodeploy bool) error
}

func NewELBSwitcher(awscli client.AWSClient, bgplan *types.BlueGreenPlan) ELBSwitcher {
	if bgplan.ElbV2 != nil && len(bgplan.ElbV2.TargetGroups) > 0 {
		return &ELBV2Switcher{
			awsCli: client.AWSCli,
		}
	} else {
		return &ELBV1Switcher{
			awsCli: client.AWSCli,
		}
	}
}

type ELBV1Switcher struct {
	awsCli client.AWSClient
}

type ELBV2Switcher struct {
	awsCli client.AWSClient
}

func (s ELBV1Switcher) Apply(clusterService ClusterService, bgplan *types.BlueGreenPlan, nodeploy bool) error {

	var targetGreen bool
	for _, lb := range bgplan.Blue.AutoScalingGroup.LoadBalancerNames {
		if *lb == bgplan.PrimaryElb {
			targetGreen = true
		}
	}

	var currentLabel string
	var nextLabel string
	var current *types.ServiceSet
	var next *types.ServiceSet
	primaryLb := bgplan.PrimaryElb
	standbyLb := bgplan.StandbyElb
	if targetGreen {
		current = bgplan.Blue
		next = bgplan.Green
		currentLabel = color.CyanString("blue")
		nextLabel = color.GreenString("green")
	} else {
		current = bgplan.Green
		next = bgplan.Blue
		currentLabel = color.GreenString("green")
		nextLabel = color.CyanString("blue")
	}

	primaryGroup := []string{primaryLb}
	standbyGroup := []string{standbyLb}
	for _, entry := range bgplan.ChainElb {
		primaryGroup = append(primaryGroup, entry.PrimaryElb)
		standbyGroup = append(standbyGroup, entry.StandbyElb)
	}

	logger.Main.Infof("Current status is '%s'", currentLabel)
	logger.Main.Infof("Start Blue-Green Deployment: %s to %s ...", currentLabel, nextLabel)
	if nodeploy {
		logger.Main.Infof("Without deployment. It only replaces load balancers.")
	} else {
		// deploy service
		logger.Main.Infof("Updating %s@%s service at %s ...", next.NewService.Service, next.NewService.Cluster, nextLabel)
		if err := clusterService.ApplyServicePlan(next.ClusterUpdatePlan); err != nil {
			return err
		}
	}

	// attach next group to primary lb
	if err := s.awsCli.Autoscaling.AttachLoadBalancers(*next.AutoScalingGroup.AutoScalingGroupName, primaryGroup); err != nil {
		return err
	}
	for _, e := range primaryGroup {
		logger.Main.Infof("Attached to attach %s group to %s(primary).", nextLabel, e)
	}

	if err := s.waitLoadBalancer(*next.AutoScalingGroup.AutoScalingGroupName, primaryLb); err != nil {
		return err
	}
	logger.Main.Infof("Added %s group to primary", nextLabel)

	// detach current group from primary lb
	if err := s.awsCli.Autoscaling.DetachLoadBalancers(*current.AutoScalingGroup.AutoScalingGroupName, primaryGroup); err != nil {
		return err
	}
	for _, e := range primaryGroup {
		logger.Main.Infof("Detached %s group from %s(primary).", currentLabel, e)
	}

	// detach next group from standby lb
	if err := s.awsCli.Autoscaling.DetachLoadBalancers(*next.AutoScalingGroup.AutoScalingGroupName, standbyGroup); err != nil {
		return err
	}
	for _, e := range standbyGroup {
		logger.Main.Infof("Detached %s group from %s(standby).", nextLabel, e)
	}

	// attach current group to standby lb
	if err := s.awsCli.Autoscaling.AttachLoadBalancers(*current.AutoScalingGroup.AutoScalingGroupName, standbyGroup); err != nil {
		return err
	}
	for _, e := range standbyGroup {
		logger.Main.Infof("Attached %s group to %s(standby).", currentLabel, e)
	}

	return nil
}

func (s ELBV1Switcher) waitLoadBalancer(group string, lb string) error {

	for {
		time.Sleep(5 * time.Second)

		result, err := s.awsCli.Autoscaling.DescribeLoadBalancerState(group)
		if err != nil {
			return err
		}

		if lbs, ok := result[lb]; ok {

			// *** LoadbalancerState
			// Adding - The instances in the group are being registered with the load balancer.
			// Added - All instances in the group are registered with the load balancer.
			// InService - At least one instance in the group passed an ELB health check.
			// Removing - The instances are being deregistered from the load balancer. If connection draining is enabled, Elastic Load Balancing waits for in-flight requests to complete before deregistering the instances.
			if *lbs.State == "Added" || *lbs.State == "InService" {
				return nil
			}

		} else {
			return fmt.Errorf("cannot get load balanracer '%s'", lb)
		}

	}

}

func (s ELBV2Switcher) Apply(clusterService ClusterService, bgplan *types.BlueGreenPlan, nodeploy bool) error {

	var targetGreen bool
	for _, tg := range bgplan.Blue.AutoScalingGroup.TargetGroupARNs {
		topTg := bgplan.ElbV2.TargetGroups[0]
		// TODO name check
		if strings.Contains(*tg, topTg.PrimaryGroup) {
			targetGreen = true
		}
	}

	var currentLabel string
	var nextLabel string
	var current *types.ServiceSet
	var next *types.ServiceSet

	topGroup := bgplan.ElbV2.TargetGroups[0]

	primaryTg := topGroup.PrimaryGroup
	if targetGreen {
		current = bgplan.Blue
		next = bgplan.Green
		currentLabel = color.CyanString("blue")
		nextLabel = color.GreenString("green")
	} else {
		current = bgplan.Green
		next = bgplan.Blue
		currentLabel = color.GreenString("green")
		nextLabel = color.CyanString("blue")
	}

	allGroup := []string{}
	primaryGroup := []string{}
	standbyGroup := []string{}
	for _, tg := range bgplan.ElbV2.TargetGroups {
		primaryGroup = append(primaryGroup, tg.PrimaryGroup)
		standbyGroup = append(standbyGroup, tg.StandbyGroup)
		allGroup = append(allGroup, tg.PrimaryGroup)
		allGroup = append(allGroup, tg.StandbyGroup)
	}

	tgmap, err := s.awsCli.ELBV2.DescribeTargetGroup(allGroup)
	if err != nil {
		return err
	}

	primaryGroupARNs := []string{}
	standbyGroupARNs := []string{}
	for name, tg := range tgmap {
		for _, pg := range primaryGroup {
			if pg == name {
				primaryGroupARNs = append(primaryGroupARNs, *tg.TargetGroupArn)
			}
		}

		for _, sg := range standbyGroup {
			if sg == name {
				standbyGroupARNs = append(standbyGroupARNs, *tg.TargetGroupArn)
			}
		}
	}

	logger.Main.Infof("Current status is '%s'", currentLabel)
	logger.Main.Infof("Start Blue-Green Deployment: %s to %s ...", currentLabel, nextLabel)
	if nodeploy {
		logger.Main.Infof("Without deployment. It only replaces load balancers.")
	} else {
		// deploy service
		logger.Main.Infof("Updating %s@%s service at %s ...", next.NewService.Service, next.NewService.Cluster, nextLabel)
		if err := clusterService.ApplyServicePlan(next.ClusterUpdatePlan); err != nil {
			return err
		}
	}

	// attach next group to primary target group
	if err := s.awsCli.Autoscaling.AttachLoadBalancerTargetGroups(*next.AutoScalingGroup.AutoScalingGroupName, aws.StringSlice(primaryGroupARNs)); err != nil {
		return err
	}
	for _, e := range primaryGroup {
		logger.Main.Infof("Attached to attach %s group to %s(primary).", nextLabel, e)
	}

	if err := s.waitTargetGroup(*next.AutoScalingGroup.AutoScalingGroupName, primaryTg); err != nil {
		return err
	}
	logger.Main.Infof("Added %s group to primary", nextLabel)

	time.Sleep(5 * time.Second)

	// detach current group from primary target group
	if err := s.awsCli.Autoscaling.DetachLoadBalancerTargetGroups(*current.AutoScalingGroup.AutoScalingGroupName, aws.StringSlice(primaryGroupARNs)); err != nil {
		return err
	}
	for _, e := range primaryGroup {
		logger.Main.Infof("Detached %s group from %s(primary).", currentLabel, e)
	}

	// detach next group from standby lb
	if err := s.awsCli.Autoscaling.DetachLoadBalancerTargetGroups(*next.AutoScalingGroup.AutoScalingGroupName, aws.StringSlice(standbyGroupARNs)); err != nil {
		return err
	}
	for _, e := range standbyGroup {
		logger.Main.Infof("Detached %s group from %s(standby).", nextLabel, e)
	}

	// attach current group to standby lb
	if err := s.awsCli.Autoscaling.AttachLoadBalancerTargetGroups(*current.AutoScalingGroup.AutoScalingGroupName, aws.StringSlice(standbyGroupARNs)); err != nil {
		return err
	}
	for _, e := range standbyGroup {
		logger.Main.Infof("Attached %s group to %s(standby).", currentLabel, e)
	}

	return nil
}

func (s ELBV2Switcher) waitTargetGroup(group string, tg string) error {
	for {
		time.Sleep(5 * time.Second)
		targetGroups, err := s.awsCli.Autoscaling.DescribeLoadBalancerTargetGroups(group)
		if err != nil {
			return err
		}

		// *** LoadBalancerTargetGroupState
		//    Adding - The Auto Scaling instances are being registered with the target group.
		//    Added - All Auto Scaling instances are registered with the target group.
		//    InService - At least one Auto Scaling instance passed an ELB health check.
		//    Removing - The Auto Scaling instances are being deregistered from the
		// target group. If connection draining is enabled, Elastic Load Balancing waits
		// for in-flight requests to complete before deregistering the instances.
		//    Removed - All Auto Scaling instances are deregistered from the target group.

		if len(targetGroups) == 0 {
			return fmt.Errorf("cannot get target group %s at %s", tg, group)
		}

		for _, targetGroup := range targetGroups {
			if strings.Contains(*targetGroup.LoadBalancerTargetGroupARN, tg) {
				state := *targetGroup.State
				if state == "Added" || state == "InService" {
					return nil
				}
			}
		}

		return nil
	}
}
