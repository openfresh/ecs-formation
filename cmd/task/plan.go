package task

import (
	"github.com/spf13/cobra"

	"github.com/stormcat24/ecs-formation/service"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Show plan to update task definiton",
	RunE: func(cmd *cobra.Command, args []string) error {

		ts, err := service.NewTaskService(projectDir, taskDefinition, parameters)
		if err != nil {
			return err
		}

		createTaskPlans(ts)
		return nil
	},
}
