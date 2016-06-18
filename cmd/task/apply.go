package task

import (
	"fmt"

	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Update task definiton",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Work your own magic here

		fmt.Println("apply called")
		return nil
	},
}
