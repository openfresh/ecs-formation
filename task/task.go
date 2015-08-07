package task

import (
	"io/ioutil"
	"github.com/stormcat24/ecs-formation/schema"
	"strings"
	"regexp"
	"github.com/stormcat24/ecs-formation/aws"
	"github.com/stormcat24/ecs-formation/plan"
	"github.com/stormcat24/ecs-formation/logger"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type TaskDefinitionController struct {
	Ecs            *aws.AwsManager
	TargetResource string
	defmap         map[string]*schema.TaskDefinition
}

func NewTaskDefinitionController(ecs *aws.AwsManager, projectDir string, targetResource string) (*TaskDefinitionController, error) {

	con := &TaskDefinitionController{
		Ecs: ecs,
	}

	defmap, err := con.searchTaskDefinitions(projectDir)
	if err != nil {
		return con, err
	}
	con.defmap = defmap

	if targetResource != "" {
		con.TargetResource = targetResource
	}

	return con, nil
}

func (self *TaskDefinitionController) GetTaskDefinitionMap() map[string]*schema.TaskDefinition {
	return self.defmap
}

func (self *TaskDefinitionController) searchTaskDefinitions(projectDir string) (map[string]*schema.TaskDefinition, error) {

	taskDir := projectDir + "/task"
	files, err := ioutil.ReadDir(taskDir)

	taskDefMap := map[string]*schema.TaskDefinition{}

	if err != nil {
		return taskDefMap, err
	}

	filePattern := regexp.MustCompile("^(.+)\\.yml$")

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".yml") {
			content, _ := ioutil.ReadFile(taskDir + "/" + file.Name())

			tokens := filePattern.FindStringSubmatch(file.Name())
			taskDefName := tokens[1]

			taskDefinition, _ := schema.CreateTaskDefinition(taskDefName, content)

			taskDefMap[taskDefName] = taskDefinition
		}
	}

	return taskDefMap, nil
}

func (self *TaskDefinitionController) CreateTaskUpdatePlans(tasks map[string]*schema.TaskDefinition) []*plan.TaskUpdatePlan {

	plans := []*plan.TaskUpdatePlan{}
	for _, task := range tasks {
		if len(self.TargetResource) == 0 || self.TargetResource == task.Name {
			plans = append(plans, self.CreateTaskUpdatePlan(task))
		}
	}

	return plans
}

func (self *TaskDefinitionController) CreateTaskUpdatePlan(task *schema.TaskDefinition) *plan.TaskUpdatePlan {

	newContainers := map[string]*schema.ContainerDefinition{}

	for _, con := range task.ContainerDefinitions {
		newContainers[con.Name] = con
	}

	return &plan.TaskUpdatePlan{
		Name: task.Name,
		NewContainers: newContainers,
	}
}

func (self *TaskDefinitionController) ApplyTaskDefinitionPlans(plans []*plan.TaskUpdatePlan) ([]*ecs.RegisterTaskDefinitionOutput, error) {

	logger.Main.Info("Start apply Task definitions...")

	outputs := []*ecs.RegisterTaskDefinitionOutput{}
	for _, plan := range plans {

		result, err := self.ApplyTaskDefinitionPlan(plan)

		if err != nil {
			return []*ecs.RegisterTaskDefinitionOutput{}, err
		}

		outputs = append(outputs, result)
	}

	return outputs, nil
}

func (self *TaskDefinitionController) ApplyTaskDefinitionPlan(task *plan.TaskUpdatePlan) (*ecs.RegisterTaskDefinitionOutput, error) {

	containers := []*schema.ContainerDefinition{}
	for _, con := range task.NewContainers {
		containers = append(containers, con)
	}

	return self.Ecs.TaskApi().RegisterTaskDefinition(task.Name, containers)
}
