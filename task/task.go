package task

import (
	"io/ioutil"
	"strings"
	"regexp"
	efaws "github.com/stormcat24/ecs-formation/aws"
	"github.com/stormcat24/ecs-formation/logger"
	"github.com/aws/aws-sdk-go/service/ecs"
	"time"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stormcat24/ecs-formation/util"
)

type TaskDefinitionController struct {
	Ecs            *efaws.AwsManager
	TargetResource string
	defmap         map[string]*TaskDefinition
}

func NewTaskDefinitionController(ecs *efaws.AwsManager, projectDir string, targetResource string) (*TaskDefinitionController, error) {

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

func (self *TaskDefinitionController) GetTaskDefinitionMap() map[string]*TaskDefinition {
	return self.defmap
}

func (self *TaskDefinitionController) searchTaskDefinitions(projectDir string) (map[string]*TaskDefinition, error) {

	taskDir := projectDir + "/task"
	files, err := ioutil.ReadDir(taskDir)

	taskDefMap := map[string]*TaskDefinition{}

	if err != nil {
		return taskDefMap, err
	}

	filePattern := regexp.MustCompile("^(.+)\\.yml$")

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".yml") {
			content, _ := ioutil.ReadFile(taskDir + "/" + file.Name())

			tokens := filePattern.FindStringSubmatch(file.Name())
			taskDefName := tokens[1]

			taskDefinition, _ := CreateTaskDefinition(taskDefName, content)

			taskDefMap[taskDefName] = taskDefinition
		}
	}

	return taskDefMap, nil
}

func (self *TaskDefinitionController) CreateTaskUpdatePlans(tasks map[string]*TaskDefinition) []*TaskUpdatePlan {

	plans := []*TaskUpdatePlan{}
	for _, task := range tasks {
		if len(self.TargetResource) == 0 || self.TargetResource == task.Name {
			plans = append(plans, self.CreateTaskUpdatePlan(task))
		}
	}

	return plans
}

func (self *TaskDefinitionController) CreateTaskUpdatePlan(task *TaskDefinition) *TaskUpdatePlan {

	newContainers := map[string]*ContainerDefinition{}

	for _, con := range task.ContainerDefinitions {
		newContainers[con.Name] = con
	}

	return &TaskUpdatePlan{
		Name: task.Name,
		NewContainers: newContainers,
	}
}

func (self *TaskDefinitionController) ApplyTaskDefinitionPlans(plans []*TaskUpdatePlan) ([]*ecs.RegisterTaskDefinitionOutput, error) {

	logger.Main.Info("Start apply Task definitions...")

	outputs := []*ecs.RegisterTaskDefinitionOutput{}
	for _, plan := range plans {

		result, err := self.ApplyTaskDefinitionPlan(plan)

		if err != nil {
			return []*ecs.RegisterTaskDefinitionOutput{}, err
		}
		time.Sleep(1 * time.Second)

		outputs = append(outputs, result)
	}

	return outputs, nil
}

func (self *TaskDefinitionController) ApplyTaskDefinitionPlan(task *TaskUpdatePlan) (*ecs.RegisterTaskDefinitionOutput, error) {

	containers := []*ContainerDefinition{}
	for _, con := range task.NewContainers {
		containers = append(containers, con)
	}

	conDefs := []*ecs.ContainerDefinition{}
	volumes := []*ecs.Volume{}

	for _, con := range containers {

		var commands []*string
		if (len(con.Command) > 0) {
			for _, token := range strings.Split(con.Command, " ") {
				commands = append(commands, aws.String(token))
			}
		} else {
			commands = nil
		}

		var entryPoints []*string
		if (len(con.EntryPoint) > 0) {
			for _, token := range strings.Split(con.EntryPoint, " ") {
				entryPoints = append(entryPoints, aws.String(token))
			}
		} else {
			entryPoints = nil
		}

		portMappings, err := toPortMappings(con.Ports)
		if err != nil {
			return &ecs.RegisterTaskDefinitionOutput{}, err
		}

		volumeItems, err := CreateVolumeInfoItems(con.Volumes)
		if err != nil {
			return &ecs.RegisterTaskDefinitionOutput{}, err
		}

		mountPoints := []*ecs.MountPoint{}
		for _, vp := range volumeItems {
			volumes = append(volumes, vp.Volume)

			mountPoints = append(mountPoints, vp.MountPoint)
		}

		volumesFrom, err := toVolumesFroms(con.VolumesFrom)
		if err != nil {
			return &ecs.RegisterTaskDefinitionOutput{}, err
		}

		conDef := &ecs.ContainerDefinition{
			CPU: &con.CpuUnits,
			Command: commands,
			EntryPoint: entryPoints,
			Environment: toKeyValuePairs(con.Environment),
			Essential: &con.Essential,
			Image: aws.String(con.Image),
			Links: util.ConvertPointerString(con.Links),
			Memory: &con.Memory,
			MountPoints: mountPoints,
			Name: aws.String(con.Name),
			PortMappings: portMappings,
			VolumesFrom: volumesFrom,
		}

		conDefs = append(conDefs, conDef)
	}

	return self.Ecs.TaskApi().RegisterTaskDefinition(task.Name, conDefs, volumes)
}
