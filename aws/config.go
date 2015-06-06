package aws

import (
	"github.com/awslabs/aws-sdk-go/aws/credentials"
)


func CreateAWSCredentials(accessKey string, secretKey string) *credentials.Credentials {

	return credentials.NewStaticCredentials(accessKey, secretKey, "")
}