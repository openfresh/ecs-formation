package task

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stormcat24/ecs-formation/client"
	cmdutil "github.com/stormcat24/ecs-formation/cmd/util"
	"github.com/stormcat24/ecs-formation/service"
	"github.com/stormcat24/ecs-formation/service/types"
	"github.com/stormcat24/ecs-formation/util"
)

var (
	projectDir     string
	taskDefinition string
	parameters     map[string]string
)

// taskCmd represents the task command
var TaskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage task definition and control running task on Amazon ECS",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		pd, err := cmdutil.GetProjectDir()
		if err != nil {
			return err
		}
		projectDir = pd

		region := viper.GetString("aws_region")
		client.Init(region, false)

		td, err := cmd.Flags().GetString("task-definition")
		if err != nil {
			return err
		}
		taskDefinition = td

		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			return err
		}

		if taskDefinition == "" && all == false {
			return errors.New("should specify '-t task_definition_name' or '-all' option")
		}

		paramTokens, err := cmd.Flags().GetStringSlice("parameter")
		if err != nil {
			return err
		}
		parameters = util.ParseKeyValues(paramTokens)

		return nil
	},
}

func init() {
	TaskCmd.AddCommand(planCmd)
	TaskCmd.AddCommand(applyCmd)
	TaskCmd.AddCommand(revisionCmd)
	TaskCmd.AddCommand(runCmd)

	TaskCmd.PersistentFlags().StringP("task-definition", "t", "", "Task Definition")
	TaskCmd.PersistentFlags().StringSliceP("parameter", "p", make([]string, 0), "parameter 'key=value'")

}

func createTaskPlans(srv service.TaskService) []*types.TaskUpdatePlan {

	taskDefs := srv.GetTaskDefinitions()
	plans := srv.CreateTaskUpdatePlans(taskDefs)

	for _, plan := range plans {
		for _, add := range plan.NewContainers {
			util.PrintlnCyan("    (+) %v", add.Name)
			util.PrintlnCyan("      image: %v", add.Image)
			util.PrintlnCyan("      ports: %v", add.Ports)
			util.PrintlnCyan("      environment:\n%v", util.StringValueWithIndent(add.Environment, 4))
			util.PrintlnCyan("      links: %v", add.Links)
			util.PrintlnCyan("      volumes: %v", add.Volumes)
			util.PrintlnCyan("      volumes_from: %v", add.VolumesFrom)
			if add.Memory != nil {
				util.PrintlnCyan("      memory: %v", *add.Memory)
			}
			if add.MemoryReservation != nil {
				util.PrintlnCyan("      memory_reservation: %#v", *add.MemoryReservation)
			}
			util.PrintlnCyan("      cpu_units: %v", add.CPUUnits)
			util.PrintlnCyan("      essential: %v", add.Essential)
			util.PrintlnCyan("      entry_point: %v", add.EntryPoint)
			util.PrintlnCyan("      command: %v", add.Command)
			util.PrintlnCyan("      disable_networking: %v", add.DisableNetworking)
			util.PrintlnCyan("      dns_search: %v", add.DNSSearchDomains)
			util.PrintlnCyan("      dns: %v", add.DNSServers)
			if len(add.DockerLabels) > 0 {
				util.PrintlnCyan("      labels: %v", util.StringValueWithIndent(add.DockerLabels, 4))
			}
			util.PrintlnCyan("      security_opt: %v", add.DockerSecurityOptions)
			util.PrintlnCyan("      extra_hosts: %v", add.ExtraHosts)
			util.PrintlnCyan("      hostname: %v", add.Hostname)
			util.PrintlnCyan("      log_driver: %v", add.LogDriver)
			if len(add.LogOpt) > 0 {
				util.PrintlnCyan("      log_opt: %v", util.StringValueWithIndent(add.LogOpt, 4))
			}
			util.PrintlnCyan("      privileged: %v", add.Privileged)
			util.PrintlnCyan("      read_only: %v", add.ReadonlyRootFilesystem)
			if len(add.Ulimits) > 0 {
				util.PrintlnCyan("      ulimits: %v", util.StringValueWithIndent(add.Ulimits, 4))
			}
			util.PrintlnCyan("      user: %v", add.User)
			util.PrintlnCyan("      working_dir: %v", add.WorkingDirectory)
		}

		util.Println()
	}

	return plans
}
