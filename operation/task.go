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
		logger.Main.Infof("Task Definition '%s'", plan.Name)

		for _, add := range plan.NewContainers {
			util.PrintlnCyan(fmt.Sprintf("    (+) %s", add.Name))
			util.PrintlnCyan(fmt.Sprintf("      image: %s", add.Image))
			util.PrintlnCyan(fmt.Sprintf("      ports: %s", add.Ports))
			util.PrintlnCyan(fmt.Sprintf("      environment:\n%s", util.StringValueWithIndent(add.Environment, 4)))
			util.PrintlnCyan(fmt.Sprintf("      links: %s", add.Links))
			util.PrintlnCyan(fmt.Sprintf("      volumes: %s", add.Volumes))
		}

		util.Println()
	}

	return plans
}
