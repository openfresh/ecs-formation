package service

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	awsecs "github.com/aws/aws-sdk-go/service/ecs"
	"github.com/fatih/color"
	"github.com/openfresh/ecs-formation/client"
	"github.com/openfresh/ecs-formation/client/ecs"
	"github.com/openfresh/ecs-formation/client/s3"
	"github.com/openfresh/ecs-formation/logger"
	"github.com/openfresh/ecs-formation/service/types"
	"github.com/openfresh/ecs-formation/util"
)

type TaskService interface {
	SearchTaskDefinitions() (map[string]*types.TaskDefinition, error)
	CreateTaskPlans() []*types.TaskUpdatePlan
	CreateTaskUpdatePlans(tasks map[string]*types.TaskDefinition) []*types.TaskUpdatePlan
	CreateTaskUpdatePlan(task *types.TaskDefinition) *types.TaskUpdatePlan
	GetTaskDefinitions() map[string]*types.TaskDefinition
	ApplyTaskDefinitionPlans(plans []*types.TaskUpdatePlan) ([]*awsecs.TaskDefinition, error)
	ApplyTaskDefinitionPlan(task *types.TaskUpdatePlan) (*awsecs.TaskDefinition, error)
	GetCurrentRevision(td string) (int64, error)
}

type ConcreteTaskService struct {
	ecsCli     ecs.Client
	s3Cli      s3.Client
	projectDir string
	target     string
	params     map[string]string
	taskDefs   map[string]*types.TaskDefinition
}

func NewTaskService(projectDir string, target string, params map[string]string) (TaskService, error) {
	service := ConcreteTaskService{
		ecsCli:     client.AWSCli.ECS,
		s3Cli:      client.AWSCli.S3,
		projectDir: projectDir,
		target:     target,
		params:     params,
	}

	defs, err := service.SearchTaskDefinitions()
	if err != nil {
		return nil, err
	}

	service.taskDefs = defs

	return &service, nil
}

func (s ConcreteTaskService) SearchTaskDefinitions() (map[string]*types.TaskDefinition, error) {

	taskDir := s.projectDir + "/task"
	taskDefMap := map[string]*types.TaskDefinition{}
	filePattern := regexp.MustCompile(`^.+\/(.+)\.yml$`)

	err := filepath.Walk(taskDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || !strings.HasSuffix(path, ".yml") {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		merged := util.MergeYamlWithParameters(content, s.params)
		tokens := filePattern.FindStringSubmatch(path)
		taskDefName := tokens[1]

		taskDefinition, err := types.CreateTaskDefinition(taskDefName, merged, filepath.Dir(path), s.s3Cli)
		if err != nil {
			return err
		}

		taskDefMap[taskDefName] = taskDefinition

		return nil
	})

	if err != nil {
		return taskDefMap, err
	}

	return taskDefMap, nil
}

func (s ConcreteTaskService) CreateTaskPlans() []*types.TaskUpdatePlan {

	plans := s.CreateTaskUpdatePlans(s.taskDefs)

	for _, plan := range plans {
		logger.Main.Infof("Task Definition '%v'", plan.Name)
	}

	return plans
}

func (s ConcreteTaskService) CreateTaskUpdatePlans(tasks map[string]*types.TaskDefinition) []*types.TaskUpdatePlan {
	plans := []*types.TaskUpdatePlan{}
	for _, task := range tasks {
		if len(s.target) == 0 || s.target == task.Name {
			plans = append(plans, s.CreateTaskUpdatePlan(task))
		}
	}

	return plans
}

func (s ConcreteTaskService) CreateTaskUpdatePlan(task *types.TaskDefinition) *types.TaskUpdatePlan {
	newContainers := map[string]*types.ContainerDefinition{}

	for _, con := range task.ContainerDefinitions {
		newContainers[con.Name] = con
	}

	return &types.TaskUpdatePlan{
		Name:          task.Name,
		NewContainers: newContainers,
	}
}

func (s ConcreteTaskService) GetTaskDefinitions() map[string]*types.TaskDefinition {
	return s.taskDefs
}

func (s ConcreteTaskService) ApplyTaskDefinitionPlans(plans []*types.TaskUpdatePlan) ([]*awsecs.TaskDefinition, error) {

	logger.Main.Info("Start apply Task definitions...")

	outputs := []*awsecs.TaskDefinition{}
	for _, plan := range plans {

		result, err := s.ApplyTaskDefinitionPlan(plan)

		if err != nil {
			logger.Main.Errorf("Register Task Definition '%s' is error.", plan.Name)
			return []*awsecs.TaskDefinition{}, err
		}
		logger.Main.Infof("Register Task Definition '%s' is success.", color.CyanString(plan.Name))
		time.Sleep(1 * time.Second)
		outputs = append(outputs, result)
	}

	return outputs, nil
}

func (s ConcreteTaskService) ApplyTaskDefinitionPlan(task *types.TaskUpdatePlan) (*awsecs.TaskDefinition, error) {

	containers := []*types.ContainerDefinition{}
	for _, con := range task.NewContainers {
		containers = append(containers, con)
	}

	conDefs := []*awsecs.ContainerDefinition{}
	volumes := []*awsecs.Volume{}

	for _, con := range containers {
		conDef, volumeItems, err := types.CreateContainerDefinition(con)
		if err != nil {
			return nil, err
		}
		conDefs = append(conDefs, conDef)

		for _, v := range volumeItems {
			volumes = append(volumes, v)
		}
	}

	return s.ecsCli.RegisterTaskDefinition(task.Name, conDefs, volumes)
}

func (s ConcreteTaskService) GetCurrentRevision(td string) (int64, error) {

	result, err := s.ecsCli.DescribeTaskDefinition(td)
	if err != nil {
		return 0, err
	}

	return *result.Revision, nil
}
