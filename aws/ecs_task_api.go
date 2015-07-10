package aws

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
)

type EcsTaskApi struct {
	Credentials *credentials.Credentials
	Region      string
}