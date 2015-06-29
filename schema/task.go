package schema

import (
	"gopkg.in/yaml.v2"
)

type TaskDefinition struct {

	Name string
	ContainerDefinitions map[string]*ContainerDefinition
}

type ContainerDefinition struct {
	Name	string
	Image	string	`yaml:"image"`
	Ports	[]string	`yaml:"ports"`
	Environment	map[string]string	`yaml:"environment"`
	Links	[]string	`yaml:"links"`
	Volumes	[]string	`yaml:"volumes"`
//	volumes:
//- /var/lib/mysql
//- cache/:/tmp/cache
//- ~/configs:/etc/configs/:ro
	Memory	int64	`yaml:"memory"`
	CpuUnits	int64	`yaml:"cpu_units"`
	Essential	bool	`yaml:"essential"`
	EntryPoint	string	`yaml:"entry_point"`
	Command	string	`yaml:"command"`
}

func CreateTaskDefinition(taskDefName string, data []byte) (*TaskDefinition, error) {

	containerMap := map[string]ContainerDefinition{}
	err := yaml.Unmarshal(data, &containerMap)

	containers := map[string]*ContainerDefinition{}
	for name, container := range containerMap {
		con := container
		con.Name = name
		containers[name] = &con
	}

	taskDef := TaskDefinition{
		Name: taskDefName,
		ContainerDefinitions: containers,
	}

	return &taskDef, err
}
