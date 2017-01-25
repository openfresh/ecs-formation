package iam

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type Config struct {
	IsMock bool
}

func NewClient(ses *session.Session, conf *Config) Client {

	if conf.IsMock {
		return &MockClient{}
	}

	return &DefaultClient{
		service: iam.New(ses),
	}
}
