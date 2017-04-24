package service

import (
	"github.com/openfresh/ecs-formation/service"
	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Show plan to update ECS service",
	RunE: func(cmd *cobra.Command, args []string) error {

		srv, err := service.NewClusterService(projectDir, []string{cluster}, serviceName, parameters)
		if err != nil {
			return err
		}

		if _, err := createClusterPlans(srv); err != nil {
			return err
		}

		return nil
	},
}
