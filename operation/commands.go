package operation

import (
	"log"
	"os"
	"github.com/codegangsta/cli"
	"fmt"
	"github.com/stormcat24/ecs-formation/aws"
	"github.com/stormcat24/ecs-formation/service"
	"github.com/stormcat24/ecs-formation/task"
	"strings"
	"github.com/str1ngs/ansi/color"
	"github.com/stormcat24/ecs-formation/util"
	"github.com/stormcat24/ecs-formation/bluegreen"
	"github.com/stormcat24/ecs-formation/logger"
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

var Commands = []cli.Command{
	commandService,
	commandTask,
	commandBluegreen,
}

var commandService = cli.Command{
	Name: "service",
	Usage: "Manage ECS services on cluster",
	Description: `
	Manage services on ECS cluster.
`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name: "json-output, jo",
			Usage: "Output json",
		},
	},
	Action: doService,
}

var commandTask = cli.Command{
	Name: "task",
	Usage: "Manage ECS Task Definitions",
	Description: `
	Manage ECS Task Definitions.
`,
	Action: doTask,
}

var commandBluegreen = cli.Command{
	Name: "bluegreen",
	Usage: "Manage bluegreen deployment on ECS",
	Description: `
	Manage bluegreen deployment on ECS.
`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name: "nodeploy, nd",
			Usage: "without deployment, only swap load balancer",
		},
		cli.BoolFlag{
			Name: "json-output, jo",
			Usage: "Output json",
		},
	},
	Action: doBluegreen,
}

func debug(v ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		log.Println(v...)
	}
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func doService(c *cli.Context) {

	awsManager, err := buildAwsManager()

	if err != nil {
		logger.Main.Error(color.Red(err.Error()))
		os.Exit(1)
	}

	operation, errSubCommand := createOperation(c.Args())

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
	clusterController, err := service.NewServiceController(awsManager, projectDir, operation.TargetResource)

	plans, err := createClusterPlans(clusterController, projectDir, jsonOutput)

	if err != nil {
		logger.Main.Error(color.Red(err.Error()))
		os.Exit(1)
	}

	if (operation.SubCommand == "apply") {
		clusterController.ApplyServicePlans(plans)
	}
}

type Hoge struct {
	Param1 string `json:param1`
}

func doTask(c *cli.Context) {

	awsManager, err := buildAwsManager()

	if err != nil {
		logger.Main.Error(color.Red(err.Error()))
		os.Exit(1)
	}

	operation, errSubCommand := createOperation(c.Args())

	if errSubCommand != nil {
		logger.Main.Error(color.Red(errSubCommand.Error()))
		os.Exit(1)
	}

	projectDir, err := os.Getwd()
	if err != nil {
		logger.Main.Error(color.Red(err.Error()))
		os.Exit(1)
	}

	taskController, err := task.NewTaskDefinitionController(awsManager, projectDir, operation.TargetResource)
	if err != nil {
		logger.Main.Error(color.Red(err.Error()))
		os.Exit(1)
	}

	plans := createTaskPlans(taskController, projectDir)

	if (operation.SubCommand == "apply") {
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

func doBluegreen(c *cli.Context) {

	awsManager, err := buildAwsManager()

	if err != nil {
		logger.Main.Error(color.Red(err.Error()))
		os.Exit(1)
	}

	operation, errSubCommand := createOperation(c.Args())

	if errSubCommand != nil {
		logger.Main.Error(color.Red(errSubCommand.Error()))
		os.Exit(1)
	}

	projectDir, err := os.Getwd()
	if err != nil {
		logger.Main.Error(color.Red(err.Error()))
		os.Exit(1)
	}

	bgController, errbgc := bluegreen.NewBlueGreenController(awsManager, projectDir, operation.TargetResource)
	if errbgc != nil {
		logger.Main.Error(color.Red(errbgc.Error()))
		os.Exit(1)
	}

	jsonOutput := c.Bool("json-output")
	bgPlans, err := createBlueGreenPlans(bgController, jsonOutput)

	if err != nil {
		logger.Main.Error(color.Red(err.Error()))
		os.Exit(1)
	}

	// cluster check

	if (operation.SubCommand == "apply") {

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
			util.PrintlnYellow(fmt.Sprintf("        ServiceARN = %s", *cs.ServiceARN))
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
			util.PrintlnCyan(fmt.Sprintf("                ServiceARN = %s", *bcs.ServiceARN))
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
			util.PrintlnGreen(fmt.Sprintf("                ServiceARN = %s", *gcs.ServiceARN))
			util.PrintlnGreen(fmt.Sprintf("                TaskDefinition = %s", *gcs.TaskDefinition))
			util.PrintlnGreen(fmt.Sprintf("                DesiredCount = %d", *gcs.DesiredCount))
			util.PrintlnGreen(fmt.Sprintf("                PendingCount = %d", *gcs.PendingCount))
			util.PrintlnGreen(fmt.Sprintf("                RunningCount = %d", *gcs.RunningCount))
		}

		util.Println()

		jsonItems = append(jsonItems, BlueGreenPlanJson{

			Blue: BlueGreenServiceJson{
				ClusterARN: *bgplan.Blue.CurrentService.ClusterARN,
				AutoScalingGroupARN: *bgplan.Blue.AutoScalingGroup.AutoScalingGroupARN,
				Instances: bgplan.Blue.AutoScalingGroup.Instances,
				TaskDefinition: *bgplan.Blue.CurrentService.TaskDefinition,
				DesiredCount: *bgplan.Blue.CurrentService.DesiredCount,
				PendingCount: *bgplan.Blue.CurrentService.PendingCount,
				RunningCount: *bgplan.Blue.CurrentService.RunningCount,
			},
			Green: BlueGreenServiceJson{
				ClusterARN: *bgplan.Green.CurrentService.ClusterARN,
				AutoScalingGroupARN: *bgplan.Green.AutoScalingGroup.AutoScalingGroupARN,
				Instances: bgplan.Green.AutoScalingGroup.Instances,
				TaskDefinition: *bgplan.Green.CurrentService.TaskDefinition,
				DesiredCount: *bgplan.Green.CurrentService.DesiredCount,
				PendingCount: *bgplan.Green.CurrentService.PendingCount,
				RunningCount: *bgplan.Green.CurrentService.RunningCount,
			},
			PrimaryElb: bgplan.PrimaryElb,
			StandbyElb: bgplan.StandbyElb,
			Active: active,
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

func buildAwsManager() (*aws.AwsManager, error) {

	accessKey := strings.Trim(os.Getenv("AWS_ACCESS_KEY"), " ")
	accessSecretKey := strings.Trim(os.Getenv("AWS_SECRET_ACCESS_KEY"), " ")
	region := strings.Trim(os.Getenv("AWS_REGION"), " ")

	if len(accessKey) == 0 {
		return nil, fmt.Errorf("'AWS_ACCESS_KEY' is not specified.")
	}

	if len(accessSecretKey) == 0 {
		return nil, fmt.Errorf("'AWS_SECRET_ACCESS_KEY' is not specified.")
	}

	if len(region) == 0 {
		return nil, fmt.Errorf("'AWS_REGION' is not specified.")
	}

	return aws.NewAwsManager(accessKey, accessSecretKey, region), nil
}

func createOperation(args cli.Args) (Operation, error) {

	if len(args) == 0 {
		return Operation{}, fmt.Errorf("subcommand is not specified.")
	}

	sub := args[0]
	if sub == "plan" || sub == "apply" {

		var targetResource string
		if len(args) > 1 {
			targetResource = args[1]
		}

		return Operation{
			SubCommand: sub,
			TargetResource: targetResource,
		}, nil
	} else {
		return Operation{}, fmt.Errorf("'%s' is invalid subcommand.", sub)
	}
}

type Operation struct {
	SubCommand     string
	TargetResource string
}
