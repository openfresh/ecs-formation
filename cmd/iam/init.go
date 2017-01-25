package iam

import (
	"github.com/spf13/cobra"
	"github.com/stormcat24/ecs-formation/client"
	cmdutil "github.com/stormcat24/ecs-formation/cmd/util"
	"github.com/stormcat24/ecs-formation/util"
)

var (
	projectDir string
	role       string
	parameters map[string]string
)

var IamCmd = &cobra.Command{
	Use:   "iam",
	Short: "Manage IAM role for ECS task",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		pd, err := cmdutil.GetProjectDir()
		if err != nil {
			return err
		}
		projectDir = pd

		// IAM is global service
		client.Init("", false)

		r, err := cmd.Flags().GetString("role")
		if err != nil {
			return err
		}
		role = r

		paramTokens, err := cmd.Flags().GetStringSlice("parameter")
		if err != nil {
			return err
		}
		parameters = util.ParseKeyValues(paramTokens)

		return nil
	},
}

func init() {
	IamCmd.AddCommand(planCmd)
	IamCmd.AddCommand(applyCmd)

	IamCmd.PersistentFlags().StringP("role", "r", "", "IAM Role")
	IamCmd.PersistentFlags().StringSliceP("parameter", "p", make([]string, 0), "parameter 'key=value'")
}
