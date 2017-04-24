package service

import (
	"github.com/openfresh/ecs-formation/service"
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Update ecs service on target cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, err := service.NewClusterService(projectDir, []string{cluster}, serviceName, parameters)
		if err != nil {
			return err
		}

		plans, err := createClusterPlans(srv)
		if err != nil {
			return err
		}

		return srv.ApplyServicePlans(plans)
	},
}
