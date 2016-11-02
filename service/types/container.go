package types

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func CreateContainerDefinition(con *ContainerDefinition) (*ecs.ContainerDefinition, []*ecs.Volume, error) {

	var commands []*string
	if len(con.Command) > 0 {
		for _, token := range strings.Split(con.Command, " ") {
			commands = append(commands, aws.String(token))
		}
	} else {
		commands = nil
	}

	var entryPoints []*string
	if len(con.EntryPoint) > 0 {
		ep, err := ParseEntrypoint(con.EntryPoint)
		if err != nil {
			return nil, []*ecs.Volume{}, err
		}
		entryPoints = ep
	} else {
		entryPoints = nil
	}

	portMappings, err := ToPortMappings(con.Ports)
	if err != nil {
		return nil, []*ecs.Volume{}, err
	}

	volumeItems, err := CreateVolumeInfoItems(con.Volumes)
	if err != nil {
		return nil, []*ecs.Volume{}, err
	}

	mountPoints := []*ecs.MountPoint{}
	volumes := []*ecs.Volume{}
	for _, vp := range volumeItems {
		volumes = append(volumes, vp.Volume)

		mountPoints = append(mountPoints, vp.MountPoint)
	}

	volumesFrom, err := ToVolumesFroms(con.VolumesFrom)
	if err != nil {
		return nil, []*ecs.Volume{}, err
	}

	extraHosts, err := ToHostEntry(con.ExtraHosts)
	if err != nil {
		return nil, []*ecs.Volume{}, err
	}

	cd := &ecs.ContainerDefinition{
		Cpu:                    aws.Int64(con.CPUUnits),
		Command:                commands,
		EntryPoint:             entryPoints,
		Environment:            ToKeyValuePairs(con.Environment),
		Essential:              aws.Bool(con.Essential),
		Image:                  aws.String(con.Image),
		Links:                  aws.StringSlice(con.Links),
		Memory:                 con.Memory,
		MemoryReservation:      con.MemoryReservation,
		MountPoints:            mountPoints,
		Name:                   aws.String(con.Name),
		PortMappings:           portMappings,
		VolumesFrom:            volumesFrom,
		DisableNetworking:      aws.Bool(con.DisableNetworking),
		DnsSearchDomains:       aws.StringSlice(con.DNSSearchDomains),
		DnsServers:             aws.StringSlice(con.DNSServers),
		DockerLabels:           aws.StringMap(con.DockerLabels),
		DockerSecurityOptions:  aws.StringSlice(con.DockerSecurityOptions),
		ExtraHosts:             extraHosts,
		Privileged:             aws.Bool(con.Privileged),
		ReadonlyRootFilesystem: aws.Bool(con.ReadonlyRootFilesystem),
		Ulimits:                ToUlimits(con.Ulimits),
	}

	if con.Hostname != "" {
		cd.Hostname = aws.String(con.Hostname)
	}
	if con.LogDriver != "" {
		cd.LogConfiguration = &ecs.LogConfiguration{
			LogDriver: aws.String(con.LogDriver),
			Options:   aws.StringMap(con.LogOpt),
		}
	}
	if con.User != "" {
		cd.User = aws.String(con.User)
	}
	if con.WorkingDirectory != "" {
		cd.WorkingDirectory = aws.String(con.WorkingDirectory)
	}

	return cd, volumes, nil
}
