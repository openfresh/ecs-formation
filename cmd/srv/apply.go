package srv

import "github.com/spf13/cobra"

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Update ecs service",
	RunE: func(cmd *cobra.Command, args []string) error {

		return nil
	},
}
