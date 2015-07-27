package config
import (
	"github.com/codegangsta/cli"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awsutil"
)

var (
	AppConfig *ApplicationConfig
)

type ApplicationConfig struct {
	SnsTopic string
}

func PrepareGlobalOptions(c *cli.Context) error {

	AppConfig = &ApplicationConfig{
		SnsTopic: c.GlobalString("sns-topic"),
	}

	fmt.Println(awsutil.StringValue(AppConfig))

	return nil
}