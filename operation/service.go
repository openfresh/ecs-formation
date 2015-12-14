package operation

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/stormcat24/ecs-formation/logger"
	"github.com/stormcat24/ecs-formation/service"
	"github.com/stormcat24/ecs-formation/util"
	"github.com/str1ngs/ansi/color"
	"os"
)

var commandService = cli.Command{
	Name:        "service",
	Usage:       "Manage ECS services on cluster",
	Description: "Manage services on ECS cluster.",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "json-output, jo",
			Usage: "Output json",
		},
		cli.StringSliceFlag{
			Name:  "params, p",
			Usage: "parameters",
		},
	},
	Action: doService,
}

func doService(c *cli.Context) {

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

	jsonOutput := c.Bool("json-output")
	clusterController, err := service.NewServiceController(awsManager, projectDir, operation.TargetResource, operation.Params)
	if err != nil {
		logger.Main.Error(color.Red(err.Error()))
		os.Exit(1)
	}

	plans, err := createClusterPlans(clusterController, projectDir, jsonOutput)

	if err != nil {
		logger.Main.Error(color.Red(err.Error()))
		os.Exit(1)
	}

	if operation.SubCommand == "apply" {
		clusterController.ApplyServicePlans(plans)
	}
}

func createClusterPlans(controller *service.ServiceController, projectDir string, jsonOutput bool) ([]*service.ServiceUpdatePlan, error) {

	if jsonOutput {
		util.Output = false
		defer func() {
			util.Output = true
		}()
	}

	util.Infoln("Checking services on clusters...")
	plans, err := controller.CreateServiceUpdatePlans()

	if err != nil {
		return []*service.ServiceUpdatePlan{}, err
	}

	for _, plan := range plans {

		util.PrintlnYellow(fmt.Sprintf("Current status of ECS Cluster '%s':", plan.Name))
		if len(plan.InstanceARNs) > 0 {
			util.PrintlnYellow("    Container Instances as follows:")
			for _, instance := range plan.InstanceARNs {
				util.PrintlnYellow(fmt.Sprintf("        %s", *instance))
			}
		}

		util.PrintlnYellow("    Services as follows:")
		if len(plan.CurrentServices) == 0 {
			util.PrintlnYellow("         No services are deployed.")
		}

		for _, cs := range plan.CurrentServices {
			util.PrintlnYellow(fmt.Sprintf("        ServiceName = %s", *cs.ServiceName))
			util.PrintlnYellow(fmt.Sprintf("        ServiceARN = %s", *cs.ServiceArn))
			util.PrintlnYellow(fmt.Sprintf("        TaskDefinition = %s", *cs.TaskDefinition))
			util.PrintlnYellow(fmt.Sprintf("        DesiredCount = %d", *cs.DesiredCount))
			util.PrintlnYellow(fmt.Sprintf("        PendingCount = %d", *cs.PendingCount))
			util.PrintlnYellow(fmt.Sprintf("        RunningCount = %d", *cs.RunningCount))
			for _, lb := range cs.LoadBalancers {
				util.PrintlnYellow(fmt.Sprintf("        ELB = %s:", *lb.LoadBalancerName))
				util.PrintlnYellow(fmt.Sprintf("            ContainerName = %s", *lb.ContainerName))
				util.PrintlnYellow(fmt.Sprintf("            ContainerName = %d", *lb.ContainerPort))
			}
			util.PrintlnYellow(fmt.Sprintf("        STATUS = %s", *cs.Status))
		}

		util.Println()
		util.PrintlnYellow(fmt.Sprintf("Service update plan '%s':", plan.Name))

		util.PrintlnYellow("    Services:")
		for _, add := range plan.NewServices {
			util.PrintlnYellow(fmt.Sprintf("        ----------[%s]----------", add.Name))
			util.PrintlnYellow(fmt.Sprintf("        TaskDefinition = %s", add.TaskDefinition))
			util.PrintlnYellow(fmt.Sprintf("        DesiredCount = %d", add.DesiredCount))
			util.PrintlnYellow(fmt.Sprintf("        KeepDesiredCount = %t", add.KeepDesiredCount))
			for _, lb := range add.LoadBalancers {
				util.PrintlnYellow(fmt.Sprintf("        ELB:%s", lb.Name))
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
