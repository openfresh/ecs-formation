package client

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stormcat24/ecs-formation/client/ecs"
	"github.com/stormcat24/ecs-formation/client/s3"
)

var (
    AWSCli AWSClient
)

type AWSClient struct {
	ECS *ecs.Client
	S3  *s3.Client
}

func Init(region string, isMock bool) {

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
	ses.Config.WithMaxRetries(aws.UseServiceDefaultRetries).WithRegion(region)

	ecsCli := ecs.NewClient(&ecs.Config{
        IsMock: isMock,
        Region: region
    })

    s3Cli := s3.NewClient(&s3.Config{
        IsMock: isMock,
        Region: region
    })

    AWSCli = AWSClient{
        ECS: ecsCli,
        S3: s3Cli,
    }
}
