package aws

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

type SnsApi struct {
	service *sns.SNS
}

func (self *SnsApi) Publish(topicArn string, message interface{}) (*sns.PublishOutput, error) {

	data, err := json.Marshal(&message)
	if err != nil {
		return &sns.PublishOutput{}, err
	}

	params := &sns.PublishInput{
		TopicArn:         aws.String(topicArn),
		Message:          aws.String(string(data)),
		MessageStructure: aws.String("json"),
	}

	return self.service.Publish(params)
}
