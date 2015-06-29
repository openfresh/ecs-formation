package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/stormcat24/ecs-formation/schema"
	"strings"
	"github.com/stormcat24/ecs-formation/util"
	"time"
	"errors"
)


type ECSManager struct {
	Credentials *credentials.Credentials
	Region string
}

func NewECSManager(accessKey string, secretKey string, region string) *ECSManager {

	cred := CreateAWSCredentials(accessKey, secretKey)

	manager := &ECSManager{
		Credentials: cred,
		Region: region,
	}

	return manager
}

func (self *ECSManager) CreateCluster(clusterName string) (*ecs.CreateClusterOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.CreateClusterInput{
		ClusterName: aws.String(clusterName),
	}

	return svc.CreateCluster(params)
}

func (self *ECSManager) DescribeTaskDefinition(defName string) (*ecs.DescribeTaskDefinitionOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(defName),
	}

	return svc.DescribeTaskDefinition(params)
}

func (self *ECSManager) RegisterTaskDefinition(taskName string, containers []*schema.ContainerDefinition) (*ecs.RegisterTaskDefinitionOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	conDefs := []*ecs.ContainerDefinition{}

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

		conDef := &ecs.ContainerDefinition{
			CPU: aws.Long(con.CpuUnits),
			Command: commands,
			EntryPoint: entryPoints,
			Environment: toKeyValuePairs(con.Environment),
			Essential: aws.Boolean(con.Essential),
			Image: aws.String(con.Image),
			Links: util.ConvertPointerString(con.Links),
			Memory: aws.Long(con.Memory),
			// MountPoints
			Name: aws.String(con.Name),
			PortMappings: toPortMappings(con.Ports),
			// VolumesFrom
		}

		conDefs = append(conDefs, conDef)
	}

	params := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: conDefs,
		Family: aws.String(taskName),
		Volumes: []*ecs.Volume{
		},
	}

	return svc.RegisterTaskDefinition(params)
}

func (self *ECSManager) DeregisterTaskDefinition(taskName string) (*ecs.DeregisterTaskDefinitionOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.DeregisterTaskDefinitionInput{
		TaskDefinition: aws.String(taskName),
	}

	return svc.DeregisterTaskDefinition(params)
}

func (self *ECSManager) DescribeClusters(clusterNames []*string) (*ecs.DescribeClustersOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.DescribeClustersInput{
		Clusters: clusterNames,
	}

	return svc.DescribeClusters(params)
}

func (self *ECSManager) CreateService(cluster string, service schema.Service) (*ecs.CreateServiceOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.CreateServiceInput{
		ServiceName: aws.String(service.Name),
		Cluster: aws.String(cluster),
		DesiredCount: aws.Long(service.DesiredCount),
		LoadBalancers: []*ecs.LoadBalancer{
			// TODO
		},
		TaskDefinition: aws.String(service.TaskDefinition),
	}

	return svc.CreateService(params)
}

func (self *ECSManager) UpdateService(cluster string, service schema.Service) (*ecs.UpdateServiceOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.UpdateServiceInput{
		Cluster: aws.String(cluster),
		Service: aws.String(service.Name),
		DesiredCount: aws.Long(service.DesiredCount),
		TaskDefinition: aws.String(service.TaskDefinition),
	}

	return svc.UpdateService(params)
}

func (self *ECSManager) DescribeService(cluster string, services []*string) (*ecs.DescribeServicesOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.DescribeServicesInput{
		Cluster: aws.String(cluster),
		Services: services,
	}

	return svc.DescribeServices(params)
}

func (self *ECSManager) DeleteService(cluster string, service string) (*ecs.DeleteServiceOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.DeleteServiceInput{
		Cluster: aws.String(cluster),
		Service: aws.String(service),
	}

	return svc.DeleteService(params)
}

func (self *ECSManager) ListServices(cluster string) (*ecs.ListServicesOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.ListServicesInput{
		Cluster: aws.String(cluster),
	}

	return svc.ListServices(params)
}

func (self *ECSManager) WaitStoppingService(cluster string, service string) error {

	for {
		time.Sleep(5 * time.Second)

		result, err := self.DescribeService(cluster, []*string{&service})

		if err != nil {
			return err
		}

		if len(result.Services) == 0 {
			return errors.New("service not found")
		}

		target := result.Services[0]

		if *target.RunningCount == 0 {
			return nil
		}

		// TODO retry count restriction
	}
}
