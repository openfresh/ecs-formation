package task

import (
	"github.com/spf13/cobra"
	"github.com/stormcat24/ecs-formation/logger"
	"github.com/stormcat24/ecs-formation/service"
	"github.com/stormcat24/ecs-formation/util"
	"github.com/str1ngs/ansi/color"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Update task definiton",
	RunE: func(cmd *cobra.Command, args []string) error {
		ts, err := service.NewTaskService(projectDir, taskDefinition, parameters)
		if err != nil {
			return err
		}

		plans := createTaskPlans(ts)
		result, err := ts.ApplyTaskDefinitionPlans(plans)
		if err != nil {
			logger.Main.Error(color.Red(err.Error()))
			return err
		}

		for _, output := range result {
			logger.Main.Infof("Registered Task Definition '%s'", *output.Family)
			logger.Main.Info(color.Cyan(util.StringValueWithIndent(output, 1)))
		}

		return nil
	},
}
