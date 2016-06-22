package srv

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stormcat24/ecs-formation/client"
)

// serviceCmd represents the service command
var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage service on Amazon ECS cluster",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		pd, err := cmd.PersistentFlags().GetString("project_dir")
		if err != nil {
			return err
		}

		jsonOutput, err := cmd.PersistentFlags().GetBool("json-output")
		if err != nil {
			return err
		}

		// TODO region
		client.Init("ap-northeast-1", false)
		// TODO: Work your own magic here
		fmt.Println("service called")
		return nil
	},
}

func init() {
	ServiceCmd.AddCommand(planCmd)
	ServiceCmd.AddCommand(applyCmd)

	ServiceCmd.PersistentFlags().BoolP("json-output", "j", false, "output json")
}
