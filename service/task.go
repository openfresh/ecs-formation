package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/joho/godotenv"
	"github.com/stormcat24/ecs-formation/client"
	"github.com/stormcat24/ecs-formation/client/ecs"
	"github.com/stormcat24/ecs-formation/client/s3"
	"github.com/stormcat24/ecs-formation/logger"
	"github.com/stormcat24/ecs-formation/util"
)

type TaskService interface {
	SearchTaskDefinitions() (map[string]*TaskDefinition, error)
	CreateTaskPlans() []*TaskUpdatePlan
	CreateTaskUpdatePlans(tasks map[string]*TaskDefinition) []*TaskUpdatePlan
	CreateTaskUpdatePlan(task *TaskDefinition) *TaskUpdatePlan
	GetTaskDefinitions() map[string]*TaskDefinition
}

type ConcreteTaskService struct {
	ecsCli     ecs.Client
	s3Cli      s3.Client
	projectDir string
	target     string
	params     map[string]string
	taskDefs   map[string]*TaskDefinition
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

func (s ConcreteTaskService) SearchTaskDefinitions() (map[string]*TaskDefinition, error) {

	taskDir := s.projectDir + "/task"
	taskDefMap := map[string]*TaskDefinition{}
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

		taskDefinition, err := s.createTaskDefinition(taskDefName, merged, filepath.Dir(path))
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

func (s ConcreteTaskService) CreateTaskPlans() []*TaskUpdatePlan {

	plans := s.CreateTaskUpdatePlans(s.taskDefs)

	for _, plan := range plans {
		logger.Main.Infof("Task Definition '%v'", plan.Name)
	}

	return plans
}

func (s ConcreteTaskService) CreateTaskUpdatePlans(tasks map[string]*TaskDefinition) []*TaskUpdatePlan {
	plans := []*TaskUpdatePlan{}
	for _, task := range tasks {
		if len(s.target) == 0 || s.target == task.Name {
			plans = append(plans, s.CreateTaskUpdatePlan(task))
		}
	}

	return plans
}

func (s ConcreteTaskService) CreateTaskUpdatePlan(task *TaskDefinition) *TaskUpdatePlan {
	newContainers := map[string]*ContainerDefinition{}

	for _, con := range task.ContainerDefinitions {
		newContainers[con.Name] = con
	}

	return &TaskUpdatePlan{
		Name:          task.Name,
		NewContainers: newContainers,
	}
}

func (s ConcreteTaskService) GetTaskDefinitions() map[string]*TaskDefinition {
	return s.taskDefs
}

func (s ConcreteTaskService) createTaskDefinition(taskDefName string, data string, basedir string) (*TaskDefinition, error) {

	containerMap := map[string]ContainerDefinition{}
	if err := yaml.Unmarshal([]byte(data), &containerMap); err != nil {
		return nil, errors.New(fmt.Sprintf("%v\n\n%v", err.Error(), data))
	}

	containers := map[string]*ContainerDefinition{}
	for name, container := range containerMap {
		con := container
		con.Name = name

		environment := map[string]string{}
		if len(container.EnvFiles) > 0 {
			for _, envfile := range container.EnvFiles {
				var path string
				if envfile[0:10] == "https://s3" {
					_path, err := s.downloadS3(envfile)
					if err != nil {
						return nil, err
					}
					path = _path
					defer os.Remove(_path)
				} else if filepath.IsAbs(envfile) {
					path = envfile
				} else {
					path = fmt.Sprintf("%s/%s", basedir, envfile)
				}

				envmap, err := s.readEnvFile(path)
				if err != nil {
					return nil, err
				}

				for key, value := range envmap {
					environment[key] = value
				}
			}
		}

		for key, value := range container.Environment {
			environment[key] = value
		}

		con.Environment = environment
		containers[name] = &con
	}

	taskDef := TaskDefinition{
		Name:                 taskDefName,
		ContainerDefinitions: containers,
	}

	return &taskDef, nil
}

func (s ConcreteTaskService) downloadS3(path string) (string, error) {
	u, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	ps := strings.Split(u.Path, "/")
	bucket := ps[:2][1]
	key := strings.Join(ps[2:], "/")

	obj, err := s.s3Cli.GetObject(bucket, key)
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(obj.Body)
	if err != nil {
		return "", err
	}

	tempfile, err := ioutil.TempFile("", "ecs-formation")
	if err != nil {
		return "", err
	}
	defer tempfile.Close()
	tempfile.Write(b)
	return tempfile.Name(), nil
}

func (s ConcreteTaskService) readEnvFile(envpath string) (map[string]string, error) {

	envmap, err := godotenv.Read(envpath)
	if err != nil {
		return map[string]string{}, err
	}

	return envmap, nil
}
