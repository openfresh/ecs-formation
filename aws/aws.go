package aws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/sns"
)

type AwsManager struct {
	conf *aws.Config
}

func NewAwsManager(region string) *AwsManager {

	cred := credentials.NewChainCredentials([]credentials.Provider{
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{Filename: "", Profile: ""},
		&ec2rolecreds.EC2RoleProvider{ExpiryWindow: 5 * time.Minute},
	})

	return &AwsManager{
		conf: aws.NewConfig().WithCredentials(cred).WithMaxRetries(aws.DefaultRetries).WithRegion(region),
	}
}

func (self *AwsManager) EcsApi() *EcsApi {
	return &EcsApi{
		service: ecs.New(self.conf),
	}
}

func (self *AwsManager) ElbApi() *ElbApi {
	return &ElbApi{
		service: elb.New(self.conf),
	}
}

func (self *AwsManager) AutoscalingApi() *AutoscalingApi {
	return &AutoscalingApi{
		service: autoscaling.New(self.conf),
	}
}

func (self *AwsManager) SnsApi() *SnsApi {
	return &SnsApi{
		service: sns.New(self.conf),
	}
}
