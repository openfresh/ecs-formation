package service

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stormcat24/ecs-formation/client"
	cmdutil "github.com/stormcat24/ecs-formation/cmd/util"
	"github.com/stormcat24/ecs-formation/service"
	"github.com/stormcat24/ecs-formation/service/types"
	"github.com/stormcat24/ecs-formation/util"
)

var (
	projectDir  string
	cluster     string
	serviceName string
	parameters  map[string]string
	jsonOutput  bool
)

var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage Amazon ECS Service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		pd, err := cmdutil.GetProjectDir()
		if err != nil {
			return err
		}
		projectDir = pd

		region := viper.GetString("aws_region")
		client.Init(region, false)

		cl, err := cmd.Flags().GetString("cluster")
		if err != nil {
			return err
		}
		cluster = cl

		sv, err := cmd.Flags().GetString("service")
		if err != nil {
			return err
		}
		serviceName = sv

		paramTokens, err := cmd.Flags().GetStringSlice("parameter")
		if err != nil {
			return err
		}
		parameters = util.ParseKeyValues(paramTokens)

		jo, err := cmd.Flags().GetBool("json-output")
		if err != nil {
			return err
		}
		jsonOutput = jo

		return nil
	},
}

func init() {
	ServiceCmd.AddCommand(planCmd)
	ServiceCmd.AddCommand(applyCmd)

	ServiceCmd.PersistentFlags().StringP("cluster", "c", "", "ECS Cluster")
	ServiceCmd.PersistentFlags().StringP("service", "s", "", "ECS Service")
	ServiceCmd.PersistentFlags().StringSliceP("parameter", "p", make([]string, 0), "parameter 'key=value'")
	ServiceCmd.PersistentFlags().BoolP("json-output", "j", false, "Print json format")
}

func createClusterPlans(srv service.ServiceService) ([]*types.ServiceUpdatePlan, error) {

	if jsonOutput {
		util.Output = false
		defer func() {
			util.Output = true
		}()
	}

	util.PrintlnYellow("Checking services on clusters...")
	plans, err := srv.CreateServiceUpdatePlans()
	if err != nil {
		return make([]*types.ServiceUpdatePlan, 0), err
	}

	for _, plan := range plans {
		util.PrintlnYellow("Current status of ECS Cluster '%s':", plan.Name)
		if len(plan.InstanceARNs) > 0 {
			util.PrintlnYellow("    Container Instances as follows:")
			for _, instance := range plan.InstanceARNs {
				util.PrintlnYellow("        %s:", *instance)
			}
		}

		util.Println()
		util.PrintlnYellow("    Services as follows:")
		if len(plan.CurrentServices) == 0 {
			util.PrintlnYellow("         No services are deployed.")
		}

		for _, cs := range plan.CurrentServices {
			util.PrintlnYellow("        ####[%s]####\n", *cs.ServiceName)
			util.PrintlnYellow("        ServiceARN = %s", *cs.ServiceArn)
			util.PrintlnYellow("        TaskDefinition = %s", *cs.TaskDefinition)
			util.PrintlnYellow("        DesiredCount = %d", *cs.DesiredCount)
			util.PrintlnYellow("        PendingCount = %d", *cs.PendingCount)
			util.PrintlnYellow("        RunningCount = %d", *cs.RunningCount)
			for _, lb := range cs.LoadBalancers {
				util.PrintlnYellow("        ELB = %s:", *lb.LoadBalancerName)
				util.PrintlnYellow("            ContainerName = %s", *lb.ContainerName)
				util.PrintlnYellow("            ContainerName = %d", *lb.ContainerPort)
			}
			util.PrintlnYellow("        STATUS = %s", *cs.Status)
			util.Println()
		}

		util.Println()
		util.PrintlnYellow("Service update plan '%s':", plan.Name)

		util.PrintlnYellow("    Services:")
		for _, add := range plan.NewServices {
			util.PrintlnYellow("        ####[%s]####\n", add.Name)
			util.PrintlnYellow("        TaskDefinition = %s", add.TaskDefinition)
			util.PrintlnYellow("        DesiredCount = %d", add.DesiredCount)
			util.PrintlnYellow("        KeepDesiredCount = %t", add.KeepDesiredCount)
			for _, lb := range add.LoadBalancers {
				util.PrintlnYellow("        ELB:%s", lb.Name)
			}
			util.Println()
		}

		util.Println()
	}

	if jsonOutput {
		bt, err := json.Marshal(&plans)
		if err != nil {
			return plans, err
		}
		fmt.Println(string(bt))
	}

	return plans, nil
}
