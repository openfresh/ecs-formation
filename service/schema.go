package service

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
)

type Cluster struct {
	Name     string
	Services map[string]Service
}

type Service struct {
	Name             string
	TaskDefinition   string         `yaml:"task_definition"`
	DesiredCount     int64          `yaml:"desired_count"`
	KeepDesiredCount bool           `yaml:"keep_desired_count"`
	LoadBalancers    []LoadBalancer `yaml:"load_balancers"`
	Role             string         `yaml:"role"`
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
	Name          string `yaml:"name"`
	ContainerName string `yaml:"container_name"`
	ContainerPort int64  `yaml:"container_port"`
}

func CreateServiceMap(data string) (map[string]Service, error) {

	servicesMap := map[string]Service{}
	if err := yaml.Unmarshal([]byte(data), &servicesMap); err != nil {
		return nil, errors.New(fmt.Sprintf("%v\n\n%v", err.Error(), data))
	}

	for name, service := range servicesMap {
		service.Name = name
		servicesMap[name] = service
	}

	return servicesMap, nil
}
