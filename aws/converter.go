package aws

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"strings"
	"strconv"
	"github.com/aws/aws-sdk-go/aws"
	"regexp"
	"errors"
	"fmt"
	"github.com/stormcat24/ecs-formation/schema"
)

var portPattern = regexp.MustCompile(`^(\d+)\/(tcp|udp)$`)

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

func toPortMappings(values []string) ([]*ecs.PortMapping, error) {

	mappings := []*ecs.PortMapping{}
	for _, value := range values {

		mp, err := toPortMapping(value)

		if err != nil {
			return []*ecs.PortMapping{}, err
		}

		mappings = append(mappings, mp);
	}

	return mappings, nil
}

func toPortMapping(value string) (*ecs.PortMapping, error) {

	tokens := strings.Split(value, ":")
	length := len(tokens)

	if length == 1 {

		if _, err := strconv.Atoi(tokens[0]); err != nil {

			return &ecs.PortMapping{}, errors.New(fmt.Sprintf("Invalid port mapping value '%s'", tokens[0]))
		}

		port, _ := strconv.ParseInt(tokens[0], 10, 64)

		return &ecs.PortMapping{
			HostPort: &port,
			ContainerPort: &port,
			Protocol: aws.String("tcp"),
		}, nil

	} else if length == 2 {

		prefixTokens := portPattern.FindStringSubmatch(tokens[0])
		suffixTokens := portPattern.FindStringSubmatch(tokens[1])

		var hPort int64
		var cPort int64
		var protocol string = "tcp"

		if len(prefixTokens) > 0 {
			hPort, _ = strconv.ParseInt(prefixTokens[1], 10, 64)
		} else {
			hPort, _ = strconv.ParseInt(tokens[0], 10, 64)
		}

		if len(suffixTokens) > 0 {
			cPort, _ = strconv.ParseInt(suffixTokens[1], 10, 64)
			protocol = suffixTokens[2]
		} else {
			cPort, _ = strconv.ParseInt(tokens[1], 10, 64)
		}

		return &ecs.PortMapping{
			HostPort: &hPort,
			ContainerPort: &cPort,
			Protocol: aws.String(protocol),
		}, nil

	} else {
		return &ecs.PortMapping{}, errors.New(fmt.Sprintf("Port mapping '%s' is invalid pattern.", value))
	}
}

func toLoadBalancers(values *[]schema.LoadBalancer) []*ecs.LoadBalancer {

	loadBalancers := []*ecs.LoadBalancer{}
	for _, lb := range *values {
		loadBalancers = append(loadBalancers, &ecs.LoadBalancer{
			LoadBalancerName: &lb.Name,
			ContainerName: &lb.ContainerName,
			ContainerPort: &lb.ContainerPort,
		})
	}

	return loadBalancers
}

func toVolumesFroms(values []string) ([]*ecs.VolumeFrom, error) {

	volumes := []*ecs.VolumeFrom{}
	for _, value := range values {

		vf, err := toVolumesFrom(value)

		if err != nil {
			return []*ecs.VolumeFrom{}, err
		}

		volumes = append(volumes, vf);
	}

	return volumes, nil
}

func toVolumesFrom(value string) (*ecs.VolumeFrom, error) {

	tokens := strings.Split(value, ":")
	length := len(tokens)

	var readOnly bool
	if length > 1 {
		ro := tokens[1]
		readOnly = (ro == "ro")

		return &ecs.VolumeFrom{
			SourceContainer: aws.String(tokens[0]),
			ReadOnly: &readOnly,
		}, nil
	} else if length == 1 {
		readOnly = false
		return &ecs.VolumeFrom{
			SourceContainer: aws.String(tokens[0]),
			ReadOnly: &readOnly,
		}, nil
	} else {
		return &ecs.VolumeFrom{}, errors.New(fmt.Sprintf("Invalid port mapping value '%s'", value))
	}
}