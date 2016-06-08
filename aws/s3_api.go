package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Api struct {
	service *s3.S3
}

func (self *S3Api) GetObject(bucket string, key string) (*s3.GetObjectOutput, error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	result, err := self.service.GetObject(params)
	if err != nil {
		return nil, err
	}
	return result, nil
}
