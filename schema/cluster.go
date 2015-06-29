package schema

import "gopkg.in/yaml.v2"


type ClusterManager struct {
}

type Cluster struct {

	Name	string
	Services	map[string]Service
}

type Service struct {

	Name	string
	TaskDefinition	string	`yaml:"task_definition"`
	DesiredCount	int64	`yaml:"desired_count"`
}

func CreateServiceMap(data []byte) (map[string]Service, error) {

	servicesMap := map[string]Service{}
	err := yaml.Unmarshal(data, &servicesMap)

	for name, service := range servicesMap {
		service.Name = name
		servicesMap[name] = service
	}

	return servicesMap, err
}
