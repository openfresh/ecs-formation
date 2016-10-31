package types

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func ToHostEntry(entries []string) ([]*ecs.HostEntry, error) {

	values := []*ecs.HostEntry{}
	for _, e := range entries {
		tokens := strings.Split(e, ":")
		if len(tokens) != 2 {
			return []*ecs.HostEntry{}, fmt.Errorf("'%v' is invalid extra_host definition.", e)
		}

		values = append(values, &ecs.HostEntry{
			Hostname:  aws.String(tokens[0]),
			IpAddress: aws.String(tokens[1]),
		})
	}

	return values, nil
}
