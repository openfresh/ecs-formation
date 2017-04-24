package types

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func ToPlacementConstraints(placementConstraints []PlacementConstraint) []*ecs.PlacementConstraint {

	slice := make([]*ecs.PlacementConstraint, len(placementConstraints))
	for i, pc := range placementConstraints {
		slice[i] = &ecs.PlacementConstraint{
			Expression: aws.String(pc.Expression),
			Type:       aws.String(pc.Type),
		}
	}
	return slice
}

func ToPlacementStrategy(placementStrategy []PlacementStrategy) []*ecs.PlacementStrategy {

	slice := make([]*ecs.PlacementStrategy, len(placementStrategy))
	for i, pc := range placementStrategy {
		slice[i] = &ecs.PlacementStrategy{
			Field: aws.String(pc.Field),
			Type:  aws.String(pc.Type),
		}
	}
	return slice
}
