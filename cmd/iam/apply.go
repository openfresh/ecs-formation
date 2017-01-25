package iam

import (
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Update IAM role",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO apply
		return nil
	},
}
