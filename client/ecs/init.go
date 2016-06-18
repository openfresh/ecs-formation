package ecs

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type Config struct {
	IsMock bool
	Region string
}

func NewClient(conf *Config) Client {

	if conf.IsMock {
		return &MockClient{}
	}

	ses := session.New()
	cred := credentials.NewChainCredentials([]credentials.Provider{
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{Filename: "", Profile: ""},
		&ec2rolecreds.EC2RoleProvider{
			Client:       ec2metadata.New(ses),
			ExpiryWindow: 5 * time.Minute,
		},
	})

	ses.Config.Credentials = cred
	ses.Config.WithMaxRetries(aws.UseServiceDefaultRetries).WithRegion(conf.Region)

	return &DefaultClient{
		service: ecs.New(ses),
	}
}
