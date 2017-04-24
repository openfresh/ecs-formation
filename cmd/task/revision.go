package task

import (
	"github.com/fatih/color"
	"github.com/openfresh/ecs-formation/service"
	"github.com/spf13/cobra"
)

var revisionCmd = &cobra.Command{
	Use:   "revision",
	Short: "Show current revision of task definition",
	RunE: func(cmd *cobra.Command, args []string) error {

		ts, err := service.NewTaskService(projectDir, taskDefinition, parameters)
		if err != nil {
			return err
		}

		revision, err := ts.GetCurrentRevision(taskDefinition)
		if err != nil {
			return err
		}
		color.Cyan("Current revision %s:%v", taskDefinition, revision)

		return nil
	},
}
