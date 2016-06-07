package task

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run task in specified ECS Cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Work your own magic here
		fmt.Println("run called")
		return nil
	},
}
