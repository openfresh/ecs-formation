package aws

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
)


type AwsManager struct {
	credentials *credentials.Credentials
	region *string
}

func NewAwsManager(accessKey string, secretKey string, region string) *AwsManager {

	cred := CreateAWSCredentials(accessKey, secretKey)

	manager := &AwsManager{
		credentials: cred,
		region: &region,
	}

	return manager
}

func (self *AwsManager) ClusterApi() *EcsClusterApi {
	return &EcsClusterApi{
		Credentials: self.credentials,
		Region: self.region,
	}
}

func (self *AwsManager) ServiceApi() *EcsServiceApi {
	return &EcsServiceApi{
		Credentials: self.credentials,
		Region: self.region,
	}
}

func (self *AwsManager) TaskApi() *EcsTaskApi {
	return &EcsTaskApi{
		Credentials: self.credentials,
		Region: self.region,
	}
}

func (self *AwsManager) ElbApi() *ElbApi {
	return &ElbApi{
		Credentials: self.credentials,
		Region: self.region,
	}
}

func (self *AwsManager) AutoscalingApi() *AutoscalingApi {
	return &AutoscalingApi{
		Credentials: self.credentials,
		Region: self.region,
	}
}

func (self *AwsManager) SnsApi() *SnsApi {
	return &SnsApi{
		Credentials: self.credentials,
		Region: self.region,
	}
}