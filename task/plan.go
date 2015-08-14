package task
import "github.com/aws/aws-sdk-go/service/ecs"

type TaskUpdatePlan struct {
	Name          string
	NewContainers map[string]*ContainerDefinition
}

type UpdateContainer struct {
	Before *ecs.ContainerDefinition
	After  *ContainerDefinition
}
