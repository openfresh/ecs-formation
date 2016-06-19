package types

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var portPattern = regexp.MustCompile(`^(\d+)\/(tcp|udp)$`)

func ToPortMappings(values []string) ([]*ecs.PortMapping, error) {

	mappings := []*ecs.PortMapping{}
	for _, value := range values {

		mp, err := ToPortMapping(value)

		if err != nil {
			return []*ecs.PortMapping{}, err
		}

		mappings = append(mappings, mp)
	}

	return mappings, nil
}

func ToPortMapping(value string) (*ecs.PortMapping, error) {

	tokens := strings.Split(value, ":")
	length := len(tokens)

	if length == 1 {

		if _, err := strconv.Atoi(tokens[0]); err != nil {

			return &ecs.PortMapping{}, errors.New(fmt.Sprintf("Invalid port mapping value '%s'", tokens[0]))
		}

		port, _ := strconv.ParseInt(tokens[0], 10, 64)

		return &ecs.PortMapping{
			HostPort:      &port,
			ContainerPort: &port,
			Protocol:      aws.String("tcp"),
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
			HostPort:      &hPort,
			ContainerPort: &cPort,
			Protocol:      aws.String(protocol),
		}, nil

	} else {
		return &ecs.PortMapping{}, errors.New(fmt.Sprintf("Port mapping '%s' is invalid pattern.", value))
	}
}
