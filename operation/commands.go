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
	"github.com/stormcat24/ecs-formation/plan"
	"github.com/stormcat24/ecs-formation/util"
	"github.com/stormcat24/ecs-formation/bluegreen"
	"github.com/stormcat24/ecs-formation/logger"
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
			Usage: "bbb",
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

	ecsManager, err := buildECSManager()

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

	clusterController, err := service.NewServiceController(ecsManager, projectDir, operation.TargetResource)

	plans, err := createClusterPlans(clusterController, projectDir)

	if err != nil {
		logger.Main.Error(color.Red(err.Error()))
		os.Exit(1)
	}

	if (operation.SubCommand == "apply") {
		clusterController.ApplyServicePlans(plans)
	}
}

func doTask(c *cli.Context) {

	ecsManager, err := buildECSManager()

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

	taskController, err := task.NewTaskDefinitionController(ecsManager, projectDir, operation.TargetResource)
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

	ecsManager, err := buildECSManager()

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

	bgController, errbgc := bluegreen.NewBlueGreenController(ecsManager, projectDir, operation.TargetResource)
	if errbgc != nil {
		logger.Main.Error(color.Red(errbgc.Error()))
		os.Exit(1)
	}

	bgPlans, err := createBlueGreenPlans(bgController)

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

func createClusterPlans(controller *service.ServiceController, projectDir string) ([]*plan.ServiceUpdatePlan, error) {

	logger.Main.Infoln("Checking services on clusters...")
	plans, err := controller.CreateServiceUpdatePlans()

	if err != nil {
		return []*plan.ServiceUpdatePlan{}, err
	}

	for _, plan := range plans {

		fmt.Println(color.Yellow(fmt.Sprintf("Current status of ECS Cluster '%s':", plan.Name)))

		if len(plan.CurrentServices) > 0 {
			fmt.Println(color.Yellow("    Services as follows:"))
		} else {
			fmt.Println(color.Yellow("    No services are deployed."))
		}

		for _, cs := range plan.CurrentServices {
			fmt.Println(color.Yellow(fmt.Sprintf("        ServiceName = %s", *cs.ServiceName)))
			fmt.Println(color.Yellow(fmt.Sprintf("        ServiceARN = %s", *cs.ServiceARN)))
			fmt.Println(color.Yellow(fmt.Sprintf("        TaskDefinition = %s", *cs.TaskDefinition)))
			fmt.Println(color.Yellow(fmt.Sprintf("        DesiredCount = %d", *cs.DesiredCount)))
			fmt.Println(color.Yellow(fmt.Sprintf("        PendingCount = %d", *cs.PendingCount)))
			fmt.Println(color.Yellow(fmt.Sprintf("        RunningCount = %d", *cs.RunningCount)))
			for _, lb := range cs.LoadBalancers {
				fmt.Println(color.Yellow(fmt.Sprintf("        ELB = %s:", *lb.LoadBalancerName)))
				fmt.Println(color.Yellow(fmt.Sprintf("            ContainerName = %s", *lb.ContainerName)))
				fmt.Println(color.Yellow(fmt.Sprintf("            ContainerName = %d", *lb.ContainerPort)))
			}
			fmt.Println(color.Yellow(fmt.Sprintf("        STATUS = %s", *cs.Status)))
		}

		for _, add := range plan.NewServices {
			for _, lb := range add.LoadBalancers {
				logger.Main.Info(color.Cyan(fmt.Sprintf("            ELB:%s", lb.Name)))
			}
		}

		fmt.Println()
	}

	return plans, nil
}

func createTaskPlans(controller *task.TaskDefinitionController, projectDir string) []*plan.TaskUpdatePlan {

	taskDefs := controller.GetTaskDefinitionMap()
	plans := controller.CreateTaskUpdatePlans(taskDefs)

	for _, plan := range plans {
		logger.Main.Infof("Task Definition '%s'", plan.Name)

		for _, add := range plan.NewContainers {
			fmt.Println(color.Cyan(fmt.Sprintf("    (+) %s", add.Name)))
			fmt.Println(color.Cyan(fmt.Sprintf("      image: %s", add.Image)))
			fmt.Println(color.Cyan(fmt.Sprintf("      ports: %s", add.Ports)))
			fmt.Println(color.Cyan(fmt.Sprintf("      environment:\n%s", util.StringValueWithIndent(add.Environment, 4))))
			fmt.Println(color.Cyan(fmt.Sprintf("      links: %s", add.Links)))
			fmt.Println(color.Cyan(fmt.Sprintf("      volumes: %s", add.Volumes)))
		}

		fmt.Println()
	}

	return plans
}

func createBlueGreenPlans(controller *bluegreen.BlueGreenController) ([]*plan.BlueGreenPlan, error) {

	bgmap := controller.GetBlueGreenMap()

	cplans, errcp := controller.ClusterController.CreateServiceUpdatePlans()
	if errcp != nil {
		return []*plan.BlueGreenPlan{}, errcp
	}

	bgplans, errbgp := controller.CreateBlueGreenPlans(bgmap, cplans)
	if errbgp != nil {
		return bgplans, errbgp
	}

	for _, bgplan := range bgplans {
		fmt.Println(color.Cyan("    Blue:"))
		fmt.Println(color.Cyan(fmt.Sprintf("        Cluster = %s", bgplan.Blue.NewService.Cluster)))
		fmt.Println(color.Cyan(fmt.Sprintf("        AutoScalingGroupARN = %s", *bgplan.Blue.AutoScalingGroup.AutoScalingGroupARN)))
		fmt.Println(color.Cyan("        Current services as follows:"))
		for _, bcs := range bgplan.Blue.ClusterUpdatePlan.CurrentServices {
			fmt.Println(color.Cyan(fmt.Sprintf("            %s:", *bcs.ServiceName)))
			fmt.Println(color.Cyan(fmt.Sprintf("                ServiceARN = %s", *bcs.ServiceARN)))
			fmt.Println(color.Cyan(fmt.Sprintf("                TaskDefinition = %s", *bcs.TaskDefinition)))
			fmt.Println(color.Cyan(fmt.Sprintf("                DesiredCount = %d", *bcs.DesiredCount)))
			fmt.Println(color.Cyan(fmt.Sprintf("                PendingCount = %d", *bcs.PendingCount)))
			fmt.Println(color.Cyan(fmt.Sprintf("                RunningCount = %d", *bcs.RunningCount)))
		}

		fmt.Println(color.Green("    Green:"))
		fmt.Println(color.Green(fmt.Sprintf("        Cluster = %s", bgplan.Green.NewService.Cluster)))
		fmt.Println(color.Green(fmt.Sprintf("        AutoScalingGroupARN = %s", *bgplan.Green.AutoScalingGroup.AutoScalingGroupARN)))
		fmt.Println(color.Green("        Current services as follows:"))
		for _, gcs := range bgplan.Green.ClusterUpdatePlan.CurrentServices {
			fmt.Println(color.Green(fmt.Sprintf("            %s:", *gcs.ServiceName)))
			fmt.Println(color.Green(fmt.Sprintf("                ServiceARN = %s", *gcs.ServiceARN)))
			fmt.Println(color.Green(fmt.Sprintf("                TaskDefinition = %s", *gcs.TaskDefinition)))
			fmt.Println(color.Green(fmt.Sprintf("                DesiredCount = %d", *gcs.DesiredCount)))
			fmt.Println(color.Green(fmt.Sprintf("                PendingCount = %d", *gcs.PendingCount)))
			fmt.Println(color.Green(fmt.Sprintf("                RunningCount = %d", *gcs.RunningCount)))
		}

		fmt.Println()
	}

	return bgplans, nil
}

func buildECSManager() (*aws.ECSManager, error) {

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

	return aws.NewECSManager(accessKey, accessSecretKey, region), nil
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
