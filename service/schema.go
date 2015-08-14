package service

import (
	"gopkg.in/yaml.v2"
)

type Cluster struct {

	Name	string
	Services	map[string]Service
}

type Service struct {

	Name	string
	TaskDefinition	string	`yaml:"task_definition"`
	DesiredCount	int64	`yaml:"desired_count"`
	LoadBalancers	[]LoadBalancer	`yaml:"load_balancers"`
	Role	string	`yaml:"role"`
}

func (self *Service) FindLoadBalancerByContainer(conname string, port int64) *LoadBalancer {

	for _, lb := range self.LoadBalancers {
		if lb.ContainerName == conname &&
		lb.ContainerPort == port {
			return &lb
		}
	}
	return nil
}

func (self *Service) FindLoadBalancerByName(name string) *LoadBalancer {

	for _, lb := range self.LoadBalancers {
		if lb.Name == name {
			return &lb
		}
	}

	return nil
}

type LoadBalancer struct {

	Name	string	`yaml:"name"`
	ContainerName	string	`yaml:"container_name"`
	ContainerPort	int64	`yaml:"container_port"`
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