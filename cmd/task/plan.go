package task

import "github.com/spf13/cobra"

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Show plan to update task definiton",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
