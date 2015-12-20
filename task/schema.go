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
	Name                  string
	Image                 string            `yaml:"image"`
	Ports                 []string          `yaml:"ports"`
	Environment           map[string]string `yaml:"environment"`
	Links                 []string          `yaml:"links"`
	Volumes               []string          `yaml:"volumes"`
	VolumesFrom           []string          `yaml:"volumes_from"`
	Memory                int64             `yaml:"memory"`
	CpuUnits              int64             `yaml:"cpu_units"`
	Essential             bool              `yaml:"essential"`
	EntryPoint            string            `yaml:"entry_point"`
	Command               string            `yaml:"command"`
	DisableNetworking     bool              `yaml:"disable_networking"`
	DnsSearchDomains      []string          `yaml:"dns_search_domains"`
	DnsServers            []string          `yaml:"dns_servers"`
	DockerLabels          map[string]string `yaml:"docker_labels"`
	DockerSecurityOptions []string          `yaml:"docker_security_options"`
	ExtraHosts            []string          `yaml:"extra_hosts"`
	Hostname              string            `yaml:"hostname"`
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
