package operation

import (
	"log"
	"os"
	"github.com/codegangsta/cli"
	"fmt"
	"github.com/stormcat24/ecs-formation/aws"
	"github.com/stormcat24/ecs-formation/cluster"
	"github.com/stormcat24/ecs-formation/task"
	"strings"
	"github.com/str1ngs/ansi/color"
	"github.com/stormcat24/ecs-formation/plan"
	"github.com/stormcat24/ecs-formation/util"
	"github.com/stormcat24/ecs-formation/bluegreen"
	"github.com/stormcat24/ecs-formation/logger"
	"errors"
)

var Commands = []cli.Command{
	commandCluster,
	commandTask,
	commandDeploy,
}

var commandCluster = cli.Command{
	Name: "cluster",
	Usage: "Manage ECS services on cluster",
	Description: `
	Manage ECS Clusters.
`,
	Action: doCluster,
}

var commandTask = cli.Command{
	Name: "task",
	Usage: "Manage ECS Task Definitions",
	Description: `
	Manage ECS Task Definitions.
`,
	Action: doTask,
}

var commandDeploy = cli.Command{
	Name: "deploy",
	Usage: "Manage bluegreen deployment on ECS",
	Description: `
	Manage bluegreen deployment on ECS.
`,
	Action: doDeploy,
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

func doCluster(c *cli.Context) {

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

	clusterController, err := cluster.NewClusterController(ecsManager, projectDir, operation.TargetResource)

	plans, err := createClusterPlans(clusterController, projectDir)

	if err != nil {
		logger.Main.Error(color.Red(err.Error()))
		os.Exit(1)
	}

	if (operation.SubCommand == "apply") {
		clusterController.ApplyClusterPlans(plans)
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

	// plan
	taskController := &task.TaskDefinitionController{
		Ecs: ecsManager,
		TargetResource: operation.TargetResource,
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

func doDeploy(c *cli.Context) {

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

	bgController, errbgc := bluegreen.NewBlueGreenController(ecsManager, projectDir)
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

		errbg := bgController.ApplyBlueGreenDeploys(bgPlans)
		if errbg != nil {
			logger.Main.Error(color.Red(errbg.Error()))
			os.Exit(1)
		}

	}
}

func createClusterPlans(controller *cluster.ClusterController, projectDir string) ([]*plan.ClusterUpdatePlan, error) {

	plans, err := controller.CreateClusterUpdatePlans()

	if err != nil {
		return []*plan.ClusterUpdatePlan{}, err
	}

	for _, plan := range plans {
		logger.Main.Infof("Cluster '%s'\n", plan.Name)

		logger.Main.Info(color.Cyan(fmt.Sprintf("\t[Add] num = %d", len(plan.NewServices))))
		for _, add := range plan.NewServices {
			logger.Main.Info(color.Cyan(fmt.Sprintf("\t\t (+) %s", add.Name)))

			for _, lb := range add.LoadBalancers {
				logger.Main.Info(color.Cyan(fmt.Sprintf("\t\t\t ELB:%s", lb.Name)))
			}
		}

		fmt.Println()
	}

	return plans, nil
}

func createTaskPlans(controller *task.TaskDefinitionController, projectDir string) []*plan.TaskUpdatePlan {

	taskDefs := controller.SearchTaskDefinitions(projectDir)
	plans := controller.CreateTaskUpdatePlans(taskDefs)

	for _, plan := range plans {
		logger.Main.Infof("Task Definition '%s'\n", plan.Name)

		logger.Main.Info(color.Cyan(fmt.Sprintf("  [Add] num = %d", len(plan.NewContainers))))
		for _, add := range plan.NewContainers {
			logger.Main.Info(color.Cyan(fmt.Sprintf("    (+) %s", add.Name)))
			logger.Main.Info(color.Cyan(fmt.Sprintf("      image: %s", add.Image)))
			logger.Main.Info(color.Cyan(fmt.Sprintf("      ports: %s", add.Ports)))
			logger.Main.Info(color.Cyan(fmt.Sprintf("      environment:\n%s", util.StringValueWithIndent(add.Environment, 4))))
			logger.Main.Info(color.Cyan(fmt.Sprintf("      links: %s", add.Links)))
			logger.Main.Info(color.Cyan(fmt.Sprintf("      volumes: %s", add.Volumes)))
		}

		fmt.Println()
	}

	return plans
}

func createBlueGreenPlans(controller *bluegreen.BlueGreenController) ([]*plan.BlueGreenPlan, error) {

	bgDefs := controller.GetBlueGreenDefs()
	bgPlans := []*plan.BlueGreenPlan{}

	cplans, errcp := controller.ClusterController.CreateClusterUpdatePlans()
	if errcp != nil {
		return bgPlans, errcp
	}

	for _, bg := range bgDefs {

		bgPlan, err := controller.CreateBlueGreenPlan(bg.Blue, bg.Green, cplans)
		if err != nil {
			return bgPlans, err
		}

		if bgPlan.Blue.CurrentService == nil {
			return bgPlans, errors.New(fmt.Sprintf("Service '%s' is not found. ", bg.Blue.Service))
		}

		if bgPlan.Green.CurrentService == nil {
			return bgPlans, errors.New(fmt.Sprintf("Service '%s' is not found. ", bg.Green.Service))
		}

		if bgPlan.Blue.AutoScalingGroup == nil {
			return bgPlans, errors.New(fmt.Sprintf("AutoScaling Group '%s' is not found. ", bg.Blue.AutoscalingGroup))
		}

		if bgPlan.Green.AutoScalingGroup == nil {
			return bgPlans, errors.New(fmt.Sprintf("AutoScaling Group '%s' is not found. ", bg.Green.AutoscalingGroup))
		}

		if bgPlan.Blue.ClusterUpdatePlan == nil {
			return bgPlans, errors.New(fmt.Sprintf("ECS Cluster '%s' is not found. ", bg.Blue.Cluster))
		}

		if bgPlan.Green.ClusterUpdatePlan == nil {
			return bgPlans, errors.New(fmt.Sprintf("ECS Cluster '%s' is not found. ", bg.Green.Cluster))
		}

		bgPlans = append(bgPlans, bgPlan)
	}

	return bgPlans, nil
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
