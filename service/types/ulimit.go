package types

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func ToUlimits(entries map[string]Ulimit) []*ecs.Ulimit {

	values := []*ecs.Ulimit{}
	for name, limit := range entries {
		values = append(values, &ecs.Ulimit{
			Name:      aws.String(name),
			SoftLimit: aws.Int64(limit.Soft),
			HardLimit: aws.Int64(limit.Hard),
		})
	}

	return values
}
