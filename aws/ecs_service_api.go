package aws

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
)

type EcsServiceApi struct {
	Credentials *credentials.Credentials
	Region      string
}