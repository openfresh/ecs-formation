package main

import (
	"github.com/stormcat24/ecs-formation/operation"
	"github.com/codegangsta/cli"
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

	app.Run(os.Args)

}
