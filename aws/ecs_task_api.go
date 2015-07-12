package aws

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/stormcat24/ecs-formation/schema"
	"github.com/stormcat24/ecs-formation/util"
	"strings"
)

type EcsTaskApi struct {
	Credentials *credentials.Credentials
	Region      string
}

func (self *EcsTaskApi) DescribeTaskDefinition(defName string) (*ecs.DescribeTaskDefinitionOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(defName),
	}

	return svc.DescribeTaskDefinition(params)
}

func (self *EcsTaskApi) RegisterTaskDefinition(taskName string, containers []*schema.ContainerDefinition) (*ecs.RegisterTaskDefinitionOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	conDefs := []*ecs.ContainerDefinition{}
	volumes := []*ecs.Volume{}

	for _, con := range containers {

		var commands []*string
		if (len(con.Command) > 0) {
			for _, token := range strings.Split(con.Command, " ") {
				commands = append(commands, aws.String(token))
			}
		} else {
			commands = nil
		}

		var entryPoints []*string
		if (len(con.EntryPoint) > 0) {
			for _, token := range strings.Split(con.EntryPoint, " ") {
				entryPoints = append(entryPoints, aws.String(token))
			}
		} else {
			entryPoints = nil
		}

		portMappings, err := toPortMappings(con.Ports)
		if err != nil {
			return &ecs.RegisterTaskDefinitionOutput{}, err
		}

		volumeItems, err := CreateVolumeInfoItems(con.Volumes)
		if err != nil {
			return &ecs.RegisterTaskDefinitionOutput{}, err
		}

		mountPoints := []*ecs.MountPoint{}
		for _, vp := range volumeItems {
			volumes = append(volumes, vp.Volume)

			mountPoints = append(mountPoints, vp.MountPoint)
		}

		conDef := &ecs.ContainerDefinition{
			CPU: aws.Long(con.CpuUnits),
			Command: commands,
			EntryPoint: entryPoints,
			Environment: toKeyValuePairs(con.Environment),
			Essential: aws.Boolean(con.Essential),
			Image: aws.String(con.Image),
			Links: util.ConvertPointerString(con.Links),
			Memory: aws.Long(con.Memory),
			MountPoints: mountPoints,
			Name: aws.String(con.Name),
			PortMappings: portMappings,
			// VolumesFrom
		}

		conDefs = append(conDefs, conDef)
	}

	params := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: conDefs,
		Family: aws.String(taskName),
		Volumes: volumes,
	}

	return svc.RegisterTaskDefinition(params)
}

func (self *EcsTaskApi) DeregisterTaskDefinition(taskName string) (*ecs.DeregisterTaskDefinitionOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.DeregisterTaskDefinitionInput{
		TaskDefinition: aws.String(taskName),
	}

	return svc.DeregisterTaskDefinition(params)
}
