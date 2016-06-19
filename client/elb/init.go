package elb

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
)

type Config struct {
	IsMock bool
	Region string
}

func NewClient(ses *session.Session, conf *Config) Client {

	if conf.IsMock {
		return &MockClient{}
	}

	return &DefaultClient{
		service: elb.New(ses),
	}
}
