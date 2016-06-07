package task

import (
	"fmt"

	"github.com/spf13/cobra"
)

var revisionCmd = &cobra.Command{
	Use:   "revision",
	Short: "Show current revision of task definition",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Work your own magic here
		fmt.Println("revision called")
		return nil
	},
}
