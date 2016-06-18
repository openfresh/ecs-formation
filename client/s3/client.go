package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Client interface {
	GetObject(bucket string, key string) (*s3.GetObjectOutput, error)
}

type DefaultClient struct {
	service *s3.S3
}

func (c DefaultClient) GetObject(bucket string, key string) (*s3.GetObjectOutput, error) {

	params := s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	return c.service.GetObject(&params)
}
