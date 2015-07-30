package aws

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
)


type ECSManager struct {
	credentials *credentials.Credentials
	region *string
}

func NewECSManager(accessKey string, secretKey string, region string) *ECSManager {

	cred := CreateAWSCredentials(accessKey, secretKey)

	manager := &ECSManager{
		credentials: cred,
		region: &region,
	}

	return manager
}

func (self *ECSManager) ClusterApi() *EcsClusterApi {
	return &EcsClusterApi{
		Credentials: self.credentials,
		Region: self.region,
	}
}

func (self *ECSManager) ServiceApi() *EcsServiceApi {
	return &EcsServiceApi{
		Credentials: self.credentials,
		Region: self.region,
	}
}

func (self *ECSManager) TaskApi() *EcsTaskApi {
	return &EcsTaskApi{
		Credentials: self.credentials,
		Region: self.region,
	}
}

func (self *ECSManager) ElbApi() *ElbApi {
	return &ElbApi{
		Credentials: self.credentials,
		Region: self.region,
	}
}

func (self *ECSManager) AutoscalingApi() *AutoscalingApi {
	return &AutoscalingApi{
		Credentials: self.credentials,
		Region: self.region,
	}
}