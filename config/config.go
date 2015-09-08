package config

import (
	"github.com/codegangsta/cli"
	"github.com/stormcat24/ecs-formation/logger"
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

	logger.Main.Infof("retry-count=%d", AppConfig.RetryCount)

	return nil
}
