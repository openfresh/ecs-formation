package aws

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"strings"
	"strconv"
	"github.com/aws/aws-sdk-go/aws"
)

func toKeyValuePairs(values map[string]string) []*ecs.KeyValuePair {

	pairs := []*ecs.KeyValuePair{}
	for key, value := range values {

		pair := ecs.KeyValuePair{
			Name: aws.String(key),
			Value: aws.String(value),
		}
		pairs = append(pairs, &pair)
	}

	return pairs
}

func toPortMappings(values []string) []*ecs.PortMapping {

	mappings := []*ecs.PortMapping{}
	for _, value := range values {
		tokens := strings.Split(value, ":")

		containerPort, _ := strconv.ParseInt(tokens[0], 10, 64)
		hostPort, _ := strconv.ParseInt(tokens[1], 10, 64)
		mapping := ecs.PortMapping{
			ContainerPort: aws.Long(containerPort),
			HostPort: aws.Long(hostPort),
		}

		mappings = append(mappings, &mapping);
	}

	return mappings
}
