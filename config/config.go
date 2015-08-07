package config

import (
	"github.com/codegangsta/cli"
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

	return nil
}