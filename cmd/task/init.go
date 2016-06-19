package task

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stormcat24/ecs-formation/service"
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
		wd, err := cmd.Flags().GetString("working-dir")
		if err != nil {
			return err
		}
		if wd == "" {
			pd, err := os.Getwd()
			if err != nil {
				return err
			}
			projectDir = pd
		} else {
			projectDir = wd
		}

		td, err := cmd.Flags().GetString("task-definition")
		if err != nil {
			return err
		}
		taskDefinition = td

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

	TaskCmd.PersistentFlags().StringP("working-dir", "d", "", "working directory")
	TaskCmd.PersistentFlags().StringP("task-definition", "t", "", "Task Definition")
	TaskCmd.PersistentFlags().StringSliceP("parameter", "p", make([]string, 0), "parameter 'key=value'")

}

func createTaskPlans(srv service.TaskService) []*service.TaskUpdatePlan {

	taskDefs := srv.GetTaskDefinitions()
	plans := srv.CreateTaskUpdatePlans(taskDefs)

	for _, plan := range plans {
		for _, add := range plan.NewContainers {
			util.PrintlnCyan(fmt.Sprintf("    (+) %v", add.Name))
			util.PrintlnCyan(fmt.Sprintf("      image: %v", add.Image))
			util.PrintlnCyan(fmt.Sprintf("      ports: %v", add.Ports))
			util.PrintlnCyan(fmt.Sprintf("      environment:\n%v", util.StringValueWithIndent(add.Environment, 4)))
			util.PrintlnCyan(fmt.Sprintf("      links: %v", add.Links))
			util.PrintlnCyan(fmt.Sprintf("      volumes: %v", add.Volumes))
			util.PrintlnCyan(fmt.Sprintf("      volumes_from: %v", add.VolumesFrom))
			util.PrintlnCyan(fmt.Sprintf("      memory: %v", add.Memory))
			util.PrintlnCyan(fmt.Sprintf("      cpu_units: %v", add.CPUUnits))
			util.PrintlnCyan(fmt.Sprintf("      essential: %v", add.Essential))
			util.PrintlnCyan(fmt.Sprintf("      entry_point: %v", add.EntryPoint))
			util.PrintlnCyan(fmt.Sprintf("      command: %v", add.Command))
			util.PrintlnCyan(fmt.Sprintf("      disable_networking: %v", add.DisableNetworking))
			util.PrintlnCyan(fmt.Sprintf("      dns_search: %v", add.DNSSearchDomains))
			util.PrintlnCyan(fmt.Sprintf("      dns: %v", add.DNSServers))
			if len(add.DockerLabels) > 0 {
				util.PrintlnCyan(fmt.Sprintf("      labels: %v", util.StringValueWithIndent(add.DockerLabels, 4)))
			}
			util.PrintlnCyan(fmt.Sprintf("      security_opt: %v", add.DockerSecurityOptions))
			util.PrintlnCyan(fmt.Sprintf("      extra_hosts: %v", add.ExtraHosts))
			util.PrintlnCyan(fmt.Sprintf("      hostname: %v", add.Hostname))
			util.PrintlnCyan(fmt.Sprintf("      log_driver: %v", add.LogDriver))
			if len(add.LogOpt) > 0 {
				util.PrintlnCyan(fmt.Sprintf("      log_opt: %v", util.StringValueWithIndent(add.LogOpt, 4)))
			}
			util.PrintlnCyan(fmt.Sprintf("      privileged: %v", add.Privileged))
			util.PrintlnCyan(fmt.Sprintf("      read_only: %v", add.ReadonlyRootFilesystem))
			if len(add.Ulimits) > 0 {
				util.PrintlnCyan(fmt.Sprintf("      ulimits: %v", util.StringValueWithIndent(add.Ulimits, 4)))
			}
			util.PrintlnCyan(fmt.Sprintf("      user: %v", add.User))
			util.PrintlnCyan(fmt.Sprintf("      working_dir: %v", add.WorkingDirectory))
		}

		util.Println()
	}

	return plans
}
