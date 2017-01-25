package iam

import "github.com/spf13/cobra"

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Show plan to update iam role",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO service
		// TODO createPlan
		return nil
	},
}
