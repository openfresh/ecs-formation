package operation

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/stormcat24/ecs-formation/aws"

	"log"
	"os"
	"strings"
)

type Operation struct {
	SubCommand     string
	TargetResource string
}

var Commands = []cli.Command{
	commandService,
	commandTask,
	commandBluegreen,
}

func debug(v ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		log.Println(v...)
	}
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func buildAwsManager() (*aws.AwsManager, error) {

	region := strings.Trim(os.Getenv("AWS_REGION"), " ")

	return aws.NewAwsManager(region), nil
}

func createOperation(args cli.Args) (Operation, error) {

	if len(args) == 0 {
		return Operation{}, fmt.Errorf("subcommand is not specified.")
	}

	sub := args[0]
	if sub == "plan" || sub == "apply" {

		var targetResource string
		if len(args) > 1 {
			targetResource = args[1]
		}

		return Operation{
			SubCommand:     sub,
			TargetResource: targetResource,
		}, nil
	} else {
		return Operation{}, fmt.Errorf("'%s' is invalid subcommand.", sub)
	}
}
