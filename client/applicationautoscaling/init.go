package applicationautoscaling

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
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
		service: applicationautoscaling.New(ses),
	}
}
