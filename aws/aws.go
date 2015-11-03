package aws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/sns"
)

type AwsManager struct {
	session *session.Session
}

func NewAwsManager(region string) *AwsManager {

	cred := credentials.NewChainCredentials([]credentials.Provider{
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{Filename: "", Profile: ""},
		&ec2rolecreds.EC2RoleProvider{ExpiryWindow: 5 * time.Minute},
	})

	conf := aws.NewConfig().WithCredentials(cred).WithMaxRetries(aws.UseServiceDefaultRetries).WithRegion(region)

	return &AwsManager{
		session: session.New(conf),
	}
}

func (self *AwsManager) EcsApi() *EcsApi {
	return &EcsApi{
		service: ecs.New(self.session),
	}
}

func (self *AwsManager) ElbApi() *ElbApi {
	return &ElbApi{
		service: elb.New(self.session),
	}
}

func (self *AwsManager) AutoscalingApi() *AutoscalingApi {
	return &AutoscalingApi{
		service: autoscaling.New(self.session),
	}
}

func (self *AwsManager) SnsApi() *SnsApi {
	return &SnsApi{
		service: sns.New(self.session),
	}
}
