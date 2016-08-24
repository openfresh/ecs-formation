package service

import (
	"github.com/spf13/cobra"
	"github.com/stormcat24/ecs-formation/service"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Update ecs service on target cluster",
	RunE: func(cmd *cobra.Command, args []string) error {

		srv, err := service.NewServiceService(projectDir, cluster, serviceName, parameters)
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
