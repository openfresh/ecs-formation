package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

type EcsClusterApi struct {
	Credentials *credentials.Credentials
	Region      *string
}

func (self *EcsClusterApi) CreateCluster(clusterName string) (*ecs.CreateClusterOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.CreateClusterInput{
		ClusterName: aws.String(clusterName),
	}

	return svc.CreateCluster(params)
}

func (self *EcsClusterApi) DeleteCluster(clusterName string) (*ecs.DeleteClusterOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.DeleteClusterInput{
		Cluster: aws.String(clusterName),
	}

	return svc.DeleteCluster(params)
}

func (self *EcsClusterApi) DescribeClusters(clusterNames []*string) (*ecs.DescribeClustersOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.DescribeClustersInput{
		Clusters: clusterNames,
	}

	return svc.DescribeClusters(params)
}

func (self *EcsClusterApi) ListClusters(maxResult int64) (*ecs.ListClustersOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.ListClustersInput{
		MaxResults: &maxResult,
	}

	return svc.ListClusters(params)
}

func (self *EcsClusterApi) ListContainerInstances(cluster string) (*ecs.ListContainerInstancesOutput, error) {

	svc := ecs.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &ecs.ListContainerInstancesInput{
		Cluster: aws.String(cluster),
	}

	return svc.ListContainerInstances(params)
}