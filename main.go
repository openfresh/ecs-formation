package main

import (
	"github.com/codegangsta/cli"
	"github.com/stormcat24/ecs-formation/config"
	"github.com/stormcat24/ecs-formation/operation"
	"os"
)

func main() {

	app := cli.NewApp()
	app.Name = "ecs-formation"
	app.Version = operation.Version
	app.Usage = "Manage EC2 Container Service(ECS)"
	app.Author = "Akinori Yamada(@stormcat24)"
	app.Email = "a.yamada24@gmail.com"
	app.Commands = operation.Commands
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "sns-topic, s",
			Usage:  "AWS SNS Topic Name",
			EnvVar: "ECSF_SNS_TOPIC",
		},
		cli.IntFlag{
			Name:   "retry-count, r",
			Usage:  "AWS API Max Retry Count",
			EnvVar: "ECSF_RETRY_COUNT",
		},
	}
	app.Before = config.PrepareGlobalOptions

	app.Run(os.Args)

}
