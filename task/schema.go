package task

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
)

type TaskDefinition struct {
	Name                 string
	ContainerDefinitions map[string]*ContainerDefinition
}

type ContainerDefinition struct {
	Name                   string
	Image                  string            `yaml:"image"`
	Ports                  []string          `yaml:"ports"`
	Environment            map[string]string `yaml:"environment"`
	Links                  []string          `yaml:"links"`
	Volumes                []string          `yaml:"volumes"`
	VolumesFrom            []string          `yaml:"volumes_from"`
	Memory                 int64             `yaml:"memory"`
	CpuUnits               int64             `yaml:"cpu_units"`
	Essential              bool              `yaml:"essential"`
	EntryPoint             string            `yaml:"entry_point"`
	Command                string            `yaml:"command"`
	DisableNetworking      bool              `yaml:"disable_networking"`
	DnsSearchDomains       []string          `yaml:"dns_search"`
	DnsServers             []string          `yaml:"dns"`
	DockerLabels           map[string]string `yaml:"labels"`
	DockerSecurityOptions  []string          `yaml:"security_opt"`
	ExtraHosts             []string          `yaml:"extra_hosts"`
	Hostname               string            `yaml:"hostname"`
	LogDriver              string            `yaml:"log_driver"`
	LogOpt                 map[string]string `yaml:"log_opt"`
	Privileged             bool              `yaml:"privileged"`
	ReadonlyRootFilesystem bool              `yaml:"read_only"`
	Ulimits                map[string]Ulimit `yaml:"ulimits"`
	User                   string            `yaml:"user"`
	WorkingDirectory       string            `yaml:"working_dir"`
}

type Ulimit struct {
	Soft int64 `yaml:"soft"`
	Hard int64 `yaml:"hard"`
}

func CreateTaskDefinition(taskDefName string, data string) (*TaskDefinition, error) {

	containerMap := map[string]ContainerDefinition{}
	if err := yaml.Unmarshal([]byte(data), &containerMap); err != nil {
		return nil, errors.New(fmt.Sprintf("%v\n\n%v", err.Error(), data))
	}

	containers := map[string]*ContainerDefinition{}
	for name, container := range containerMap {
		con := container
		con.Name = name
		containers[name] = &con
	}

	taskDef := TaskDefinition{
		Name:                 taskDefName,
		ContainerDefinitions: containers,
	}

	return &taskDef, nil
}
