package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/stormcat24/ecs-formation/config"
)

type AwsManager struct {
	credentials *credentials.Credentials
	region      string
	retryCount  int
}

func NewAwsManager(accessKey string, secretKey string, region string) *AwsManager {

	cred := CreateAWSCredentials(accessKey, secretKey)

	manager := &AwsManager{
		credentials: cred,
		region:      region,
		retryCount:  config.AppConfig.RetryCount,
	}

	return manager
}

func (self *AwsManager) EcsApi() *EcsApi {
	return &EcsApi{
		service: ecs.New(&aws.Config{
			Credentials: self.credentials,
			Region:      &self.region,
			MaxRetries:  &self.retryCount,
		}),
	}
}

func (self *AwsManager) ElbApi() *ElbApi {
	return &ElbApi{
		service: elb.New(&aws.Config{
			Credentials: self.credentials,
			Region:      &self.region,
			MaxRetries:  &self.retryCount,
		}),
	}
}

func (self *AwsManager) AutoscalingApi() *AutoscalingApi {
	return &AutoscalingApi{
		service: autoscaling.New(&aws.Config{
			Credentials: self.credentials,
			Region:      &self.region,
			MaxRetries:  &self.retryCount,
		}),
	}
}

func (self *AwsManager) SnsApi() *SnsApi {
	return &SnsApi{
		service: sns.New(&aws.Config{
			Credentials: self.credentials,
			Region:      &self.region,
			MaxRetries:  &self.retryCount,
		}),
	}
}
