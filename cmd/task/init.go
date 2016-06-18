package task

import "github.com/spf13/cobra"

var (
	taskDefinition string
)

// taskCmd represents the task command
var TaskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage task definition and control running task on Amazon ECS",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		td, err := cmd.Flags().GetString("task-definition")
		if err != nil {
			return err
		}
		taskDefinition = td
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
