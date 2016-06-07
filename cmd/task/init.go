package task

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	taskDefinition string
)

// taskCmd represents the task command
var TaskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage task definition and control running task on Amazon ECS",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Work your own magic here
		fmt.Println("task called")
		return nil
	},
}

func init() {
	TaskCmd.AddCommand(planCmd)
	TaskCmd.AddCommand(applyCmd)
	TaskCmd.AddCommand(revisionCmd)
	TaskCmd.AddCommand(runCmd)
	TaskCmd.PersistentFlags().StringVarP(&taskDefinition, "task-definition", "t", "", "Task Definition")
}
