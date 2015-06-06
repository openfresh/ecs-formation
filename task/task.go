package task

import (
	"io/ioutil"
	"github.com/stormcat24/ecs-formation/schema"
	"fmt"
	"strings"
	"regexp"
	"github.com/stormcat24/ecs-formation/aws"
	"github.com/stormcat24/ecs-formation/plan"
	"github.com/awslabs/aws-sdk-go/service/ecs"
)

type TaskDefinitionController struct {
	Ecs *aws.ECSManager
	TargetResource string
}

func (self *TaskDefinitionController) SearchTaskDefinitions(projectDir string) map[string]*schema.TaskDefinition {

	taskDir := projectDir + "/task"
	files, err := ioutil.ReadDir(taskDir)

	if err != nil {
		panic(err)
	}

	filePattern := regexp.MustCompile("^(.+)\\.yml$")

	taskDefMap := map[string]*schema.TaskDefinition{}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".yml") {
			content, _ := ioutil.ReadFile(taskDir + "/" + file.Name())

			tokens := filePattern.FindStringSubmatch(file.Name())
			taskDefName := tokens[1]

			taskDefinition, _ := schema.CreateTaskDefinition(taskDefName, content)

			taskDefMap[taskDefName] = taskDefinition
		}
	}

	return taskDefMap
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

	result, _ := self.Ecs.DescribeTaskDefinition(task.Name)

	newContainers := map[string]*schema.ContainerDefinition{}
	currentContainers := map[string]*ecs.ContainerDefinition{}
	updateContainers := map[string]*plan.UpdateContainer{}
	deleteContainers := map[string]*ecs.ContainerDefinition{}

	if result.TaskDefinition != nil {

		for _, currentContainer := range result.TaskDefinition.ContainerDefinitions {
			currentContainers[*currentContainer.Name] = currentContainer
		}

		for name, currentContainer := range currentContainers {

			if newContainer, ok := task.ContainerDefinitions[name]; ok {
				// update
				updateContainers[name] = &plan.UpdateContainer{
					Before: currentContainer,
					After: newContainer,
				}
			} else {
				// delete
				deleteContainers[name] = currentContainer
			}
		}

		for name, con := range task.ContainerDefinitions {

			if _, ok := currentContainers[name]; !ok {
				newContainers[con.Name] = con
			}
		}

	} else {
		for _, con := range task.ContainerDefinitions {
			newContainers[con.Name] = con
		}
	}

	return &plan.TaskUpdatePlan{
		Name: task.Name,
		CurrentContainers: currentContainers,
		DeleteContainers: deleteContainers,
		NewContainers: newContainers,
		UpdateContainers: updateContainers,
	}
}

func (self *TaskDefinitionController) ApplyTaskDefinitionPlans(plans []*plan.TaskUpdatePlan) {

	fmt.Println("Start apply Task definitions...")

	for _, plan := range plans {
		self.ApplyTaskDefinitionPlan(plan)
	}
}

func (self *TaskDefinitionController) ApplyTaskDefinitionPlan(task *plan.TaskUpdatePlan) {

	containers := []*schema.ContainerDefinition{}
	for _, con := range task.NewContainers {
		containers = append(containers, con)
	}
	_, err := self.Ecs.RegisterTaskDefinition(task.Name, containers)

	if err != nil {
		panic(err)
	}

}
