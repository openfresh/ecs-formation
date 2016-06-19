package s3

import (
	"github.com/aws/aws-sdk-go/service/s3"
)

type MockClient struct {
}

func (c MockClient) GetObject(bucket string, key string) (*s3.GetObjectOutput, error) {

	return nil, nil
}
