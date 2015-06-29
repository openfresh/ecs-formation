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
)

var Commands = []cli.Command{
	commandCluster,
	commandTask,
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
		fmt.Fprintf(os.Stderr, "[ERROR]%s\n", color.Red(err.Error()))
		os.Exit(1)
	}

	operation, errSubCommand := createOperation(c.Args())

	if errSubCommand != nil {
		fmt.Fprintf(os.Stderr, "[ERROR]%s\n", color.Red(errSubCommand.Error()))
		os.Exit(1)
	}

	projectDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	clusterController := cluster.ClusterControler{
		Ecs: ecsManager,
		TargetResource: operation.TargetResource,
	}

	plans := createClusterPlans(&clusterController, projectDir)

	if (operation.SubCommand == "apply") {
		clusterController.ApplyClusterPlans(plans)
	}
}

func doTask(c *cli.Context) {

	ecsManager, err := buildECSManager()

	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR]%s\n", color.Red(err.Error()))
		os.Exit(1)
	}

	operation, errSubCommand := createOperation(c.Args())

	if errSubCommand != nil {
		fmt.Fprintf(os.Stderr, "[ERROR]%s\n", color.Red(errSubCommand.Error()))
		os.Exit(1)
	}

	projectDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// plan
	taskController := &task.TaskDefinitionController{
		Ecs: ecsManager,
		TargetResource: operation.TargetResource,
	}

	plans := createTaskPlans(taskController, projectDir)

	if (operation.SubCommand == "apply") {
		results := taskController.ApplyTaskDefinitionPlans(plans)

		for _, output := range results {
			fmt.Printf("Registered Task Definition '%s'", *output.TaskDefinition.Family)
			fmt.Print(color.Cyan(util.StringValueWithIndent(output.TaskDefinition, 1)))
		}
	}
}

func createClusterPlans(controller *cluster.ClusterControler, projectDir string) []*plan.ClusterUpdatePlan {

	clusters := controller.SearchClusters(projectDir)
	plans := controller.CreateClusterUpdatePlans(clusters)

	for _, plan := range plans {
		fmt.Printf("Cluster '%s'\n", plan.Name)

		fmt.Println(color.Cyan(fmt.Sprintf("\t[Add] num = %d", len(plan.NewServices))))
		for _, add := range plan.NewServices {
			fmt.Println(color.Cyan(fmt.Sprintf("\t\t (+) %s", add.Name)))
		}

		fmt.Println(color.Green(fmt.Sprintf("\t[Update] num = %d", len(plan.UpdateServices))))
		for _, update := range plan.UpdateServices {
			fmt.Println(color.Green(fmt.Sprintf("\t\t (+) %s(%s)", *update.Before.ServiceName, *update.Before.ClusterARN)))
		}

		fmt.Println(color.Red(fmt.Sprintf("\t[Remove] num = %d", len(plan.DeleteServices))))
		for _, delete := range plan.DeleteServices {
			fmt.Println(color.Red(fmt.Sprintf("\t\t (-) %s(%s)", *delete.ServiceName, *delete.ClusterARN)))
		}
		fmt.Println()
	}

	return plans
}

func createTaskPlans(controller *task.TaskDefinitionController, projectDir string) []*plan.TaskUpdatePlan {

	taskDefs := controller.SearchTaskDefinitions(projectDir)
	plans := controller.CreateTaskUpdatePlans(taskDefs)

	for _, plan := range plans {
		fmt.Printf("Task Definition '%s'\n", plan.Name)

		fmt.Println(color.Cyan(fmt.Sprintf("  [Add] num = %d", len(plan.NewContainers))))
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
	SubCommand	string
	TargetResource	string
}
