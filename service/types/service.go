package types

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

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
