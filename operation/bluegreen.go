package operation

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/codegangsta/cli"
	"github.com/stormcat24/ecs-formation/bluegreen"
	"github.com/stormcat24/ecs-formation/logger"
	"github.com/stormcat24/ecs-formation/util"
	"github.com/str1ngs/ansi/color"
	"os"
)

type BlueGreenPlanJson struct {
	Blue       BlueGreenServiceJson
	Green      BlueGreenServiceJson
	Active     string
	PrimaryElb string
	StandbyElb string
}

type BlueGreenServiceJson struct {
	ClusterARN          string
	AutoScalingGroupARN string
	Instances           []*autoscaling.Instance
	TaskDefinition      string
	DesiredCount        int64
	PendingCount        int64
	RunningCount        int64
}

var commandBluegreen = cli.Command{
	Name:  "bluegreen",
	Usage: "Manage bluegreen deployment on ECS",
	Description: `
	Manage bluegreen deployment on ECS.
`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "nodeploy, nd",
			Usage: "without deployment, only swap load balancer",
		},
		cli.BoolFlag{
			Name:  "json-output, jo",
			Usage: "Output json",
		},
		cli.StringSliceFlag{
			Name:  "params, p",
			Usage: "parameters",
		},
	},
	Action: doBluegreen,
}

func doBluegreen(c *cli.Context) {

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

	bgController, err := bluegreen.NewBlueGreenController(awsManager, projectDir, operation.TargetResource, operation.Params)
	if err != nil {
		logger.Main.Error(color.Red(err.Error()))
		os.Exit(1)
	}

	jsonOutput := c.Bool("json-output")
	bgPlans, err := createBlueGreenPlans(bgController, jsonOutput)
	if err != nil {
		logger.Main.Error(color.Red(err.Error()))
		os.Exit(1)
	}

	// cluster check

	if operation.SubCommand == "apply" {

		nodeploy := c.Bool("nodeploy")

		if len(bgPlans) > 0 {
			errbg := bgController.ApplyBlueGreenDeploys(bgPlans, nodeploy)
			if errbg != nil {
				logger.Main.Error(color.Red(errbg.Error()))
				os.Exit(1)
			}
		} else {
			logger.Main.Infof("Not found Blue Green Definition")

			if len(operation.TargetResource) > 0 && !nodeploy {
				logger.Main.Infof("Try to update service '%s'", operation.TargetResource)
				doService(c)
			}

		}
	}
}

func createBlueGreenPlans(controller *bluegreen.BlueGreenController, jsonOutput bool) ([]*bluegreen.BlueGreenPlan, error) {

	if jsonOutput {
		util.Output = false
		defer func() {
			util.Output = true
		}()
	}

	bgmap := controller.GetBlueGreenMap()

	cplans, errcp := controller.ClusterController.CreateServiceUpdatePlans()
	if errcp != nil {
		return []*bluegreen.BlueGreenPlan{}, errcp
	}

	bgplans, errbgp := controller.CreateBlueGreenPlans(bgmap, cplans)
	if errbgp != nil {
		return bgplans, errbgp
	}

	jsonItems := []BlueGreenPlanJson{}
	for _, bgplan := range bgplans {
		util.PrintlnCyan("    Blue:")
		util.PrintlnCyan(fmt.Sprintf("        Cluster = %s", bgplan.Blue.NewService.Cluster))
		util.PrintlnCyan(fmt.Sprintf("        AutoScalingGroupARN = %s", *bgplan.Blue.AutoScalingGroup.AutoScalingGroupARN))
		util.PrintlnCyan("        Current services as follows:")
		for _, bcs := range bgplan.Blue.ClusterUpdatePlan.CurrentServices {
			util.PrintlnCyan(fmt.Sprintf("            %s:", *bcs.ServiceName))
			util.PrintlnCyan(fmt.Sprintf("                ServiceARN = %s", *bcs.ServiceArn))
			util.PrintlnCyan(fmt.Sprintf("                TaskDefinition = %s", *bcs.TaskDefinition))
			util.PrintlnCyan(fmt.Sprintf("                DesiredCount = %d", *bcs.DesiredCount))
			util.PrintlnCyan(fmt.Sprintf("                PendingCount = %d", *bcs.PendingCount))
			util.PrintlnCyan(fmt.Sprintf("                RunningCount = %d", *bcs.RunningCount))
		}

		var active string
		if bgplan.IsBlueWithPrimaryElb() {
			active = "blue"
		} else {
			active = "green"
		}

		util.PrintlnGreen("    Green:")
		util.PrintlnGreen(fmt.Sprintf("        Cluster = %s", bgplan.Green.NewService.Cluster))
		util.PrintlnGreen(fmt.Sprintf("        AutoScalingGroupARN = %s", *bgplan.Green.AutoScalingGroup.AutoScalingGroupARN))
		util.PrintlnGreen("        Current services as follows:")
		for _, gcs := range bgplan.Green.ClusterUpdatePlan.CurrentServices {
			util.PrintlnGreen(fmt.Sprintf("            %s:", *gcs.ServiceName))
			util.PrintlnGreen(fmt.Sprintf("                ServiceARN = %s", *gcs.ServiceArn))
			util.PrintlnGreen(fmt.Sprintf("                TaskDefinition = %s", *gcs.TaskDefinition))
			util.PrintlnGreen(fmt.Sprintf("                DesiredCount = %d", *gcs.DesiredCount))
			util.PrintlnGreen(fmt.Sprintf("                PendingCount = %d", *gcs.PendingCount))
			util.PrintlnGreen(fmt.Sprintf("                RunningCount = %d", *gcs.RunningCount))
		}

		util.Println()

		jsonItems = append(jsonItems, BlueGreenPlanJson{

			Blue: BlueGreenServiceJson{
				ClusterARN:          *bgplan.Blue.CurrentService.ClusterArn,
				AutoScalingGroupARN: *bgplan.Blue.AutoScalingGroup.AutoScalingGroupARN,
				Instances:           bgplan.Blue.AutoScalingGroup.Instances,
				TaskDefinition:      *bgplan.Blue.CurrentService.TaskDefinition,
				DesiredCount:        *bgplan.Blue.CurrentService.DesiredCount,
				PendingCount:        *bgplan.Blue.CurrentService.PendingCount,
				RunningCount:        *bgplan.Blue.CurrentService.RunningCount,
			},
			Green: BlueGreenServiceJson{
				ClusterARN:          *bgplan.Green.CurrentService.ClusterArn,
				AutoScalingGroupARN: *bgplan.Green.AutoScalingGroup.AutoScalingGroupARN,
				Instances:           bgplan.Green.AutoScalingGroup.Instances,
				TaskDefinition:      *bgplan.Green.CurrentService.TaskDefinition,
				DesiredCount:        *bgplan.Green.CurrentService.DesiredCount,
				PendingCount:        *bgplan.Green.CurrentService.PendingCount,
				RunningCount:        *bgplan.Green.CurrentService.RunningCount,
			},
			PrimaryElb: bgplan.PrimaryElb,
			StandbyElb: bgplan.StandbyElb,
			Active:     active,
		})
	}

	if jsonOutput {
		bt, err := json.Marshal(&jsonItems)
		if err != nil {
			return bgplans, err
		}
		fmt.Println(string(bt))
	}

	return bgplans, nil
}
