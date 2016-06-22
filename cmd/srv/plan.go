package srv

import "github.com/spf13/cobra"

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Show plan to update ecs service",
	RunE: func(cmd *cobra.Command, args []string) error {

		return nil
	},
}
