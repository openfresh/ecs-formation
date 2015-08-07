package aws

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/aws"
	"encoding/json"
)

type SnsApi struct {
	Credentials *credentials.Credentials
	Region      *string
}


func (self *SnsApi) Publish(topicArn string, message interface{}) (*sns.PublishOutput, error) {

	data, err := json.Marshal(&message)
	if err != nil {
		return &sns.PublishOutput{}, err
	}

	svc := sns.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &sns.PublishInput{
		TopicARN: aws.String(topicArn),
		Message: aws.String(string(data)),
		MessageStructure: aws.String("json"),
	}

	return svc.Publish(params)
}