package types

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func ToKeyValuePairs(values map[string]string) []*ecs.KeyValuePair {

	pairs := []*ecs.KeyValuePair{}
	for key, value := range values {

		pair := ecs.KeyValuePair{
			Name:  aws.String(key),
			Value: aws.String(value),
		}
		pairs = append(pairs, &pair)
	}

	return pairs
}
