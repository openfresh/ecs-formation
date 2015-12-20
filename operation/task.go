package operation

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/stormcat24/ecs-formation/logger"
	"github.com/stormcat24/ecs-formation/task"
	"github.com/stormcat24/ecs-formation/util"
	"github.com/str1ngs/ansi/color"
	"os"
)

var commandTask = cli.Command{
	Name:        "task",
	Usage:       "Manage ECS Task Definitions",
	Description: "Manage ECS Task Definitions.",
	Flags: []cli.Flag{
		cli.StringSliceFlag{
			Name:  "params, p",
			Usage: "parameters",
		},
	},
	Action: doTask,
}

func doTask(c *cli.Context) {

	awsManager, err := buildAwsManager()

	if err != nil {
		logger.Main.Error(color.Red(err.Error()))
		os.Exit(1)
	}

	operation, errSubCommand := createOperation(c)

	if errSubCommand != nil {
		logger.Main.Error(color.Red(errSubCommand.Error()))
		os.Exit(1)
	}

	projectDir, err := os.Getwd()
	if err != nil {
		logger.Main.Error(color.Red(err.Error()))
		os.Exit(1)
	}

	taskController, err := task.NewTaskDefinitionController(awsManager, projectDir, operation.TargetResource, operation.Params)
	if err != nil {
		logger.Main.Error(color.Red(err.Error()))
		os.Exit(1)
	}

	plans := createTaskPlans(taskController, projectDir)

	if operation.SubCommand == "apply" {
		results, errapp := taskController.ApplyTaskDefinitionPlans(plans)

		if errapp != nil {
			logger.Main.Error(color.Red(errapp.Error()))
			os.Exit(1)
		}

		for _, output := range results {
			logger.Main.Infof("Registered Task Definition '%s'", *output.TaskDefinition.Family)
			logger.Main.Info(color.Cyan(util.StringValueWithIndent(output.TaskDefinition, 1)))
		}
	}
}

func createTaskPlans(controller *task.TaskDefinitionController, projectDir string) []*task.TaskUpdatePlan {

	taskDefs := controller.GetTaskDefinitionMap()
	plans := controller.CreateTaskUpdatePlans(taskDefs)

	for _, plan := range plans {
		logger.Main.Infof("Task Definition '%v'", plan.Name)

		for _, add := range plan.NewContainers {
			util.PrintlnCyan(fmt.Sprintf("    (+) %v", add.Name))
			util.PrintlnCyan(fmt.Sprintf("      image: %v", add.Image))
			util.PrintlnCyan(fmt.Sprintf("      ports: %v", add.Ports))
			util.PrintlnCyan(fmt.Sprintf("      environment:\n%v", util.StringValueWithIndent(add.Environment, 4)))
			util.PrintlnCyan(fmt.Sprintf("      links: %v", add.Links))
			util.PrintlnCyan(fmt.Sprintf("      volumes: %v", add.Volumes))
			util.PrintlnCyan(fmt.Sprintf("      volumes_from: %v", add.VolumesFrom))
			util.PrintlnCyan(fmt.Sprintf("      memory: %v", add.Memory))
			util.PrintlnCyan(fmt.Sprintf("      cpu_units: %v", add.CpuUnits))
			util.PrintlnCyan(fmt.Sprintf("      essential: %v", add.Essential))
			util.PrintlnCyan(fmt.Sprintf("      entry_point: %v", add.EntryPoint))
			util.PrintlnCyan(fmt.Sprintf("      command: %v", add.Command))
			util.PrintlnCyan(fmt.Sprintf("      disable_networking: %v", add.DisableNetworking))
			util.PrintlnCyan(fmt.Sprintf("      dns_search: %v", add.DnsSearchDomains))
			util.PrintlnCyan(fmt.Sprintf("      dns: %v", add.DnsServers))
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
