package service

import (
	"github.com/spf13/cobra"
	"github.com/stormcat24/ecs-formation/service"
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
