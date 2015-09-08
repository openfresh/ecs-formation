package config

import (
	"github.com/codegangsta/cli"
)

var (
	AppConfig *ApplicationConfig
)

type ApplicationConfig struct {
	SnsTopic   string
	RetryCount int
}

func PrepareGlobalOptions(c *cli.Context) error {

	AppConfig = &ApplicationConfig{
		SnsTopic:   c.GlobalString("sns-topic"),
		RetryCount: c.GlobalInt("retry-count"),
	}

	return nil
}
