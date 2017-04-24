package bluegreen

import (
	"github.com/openfresh/ecs-formation/logger"
	"github.com/openfresh/ecs-formation/service"
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply bluegreen deployment",
	RunE: func(cmd *cobra.Command, args []string) error {

		bgsrv, err := service.NewBlueGreenService(projectDir, bluegreenName, parameters)
		if err != nil {
			return err
		}

		csrv, err := bgsrv.CreateClusterService()
		if err != nil {
			return err
		}

		plans, err := createBlueGreenPlans(bgsrv, csrv)
		if err != nil {
			return err
		}

		if len(plans) > 0 {
			if err := bgsrv.ApplyBlueGreenDeploys(csrv, plans, noDeploy); err != nil {
				return err
			}
		} else {
			if bluegreenName != "" && !noDeploy {
				logger.Main.Infof("Try to update service '%s'", bluegreenName)
				// TODO doService(c)
			}
		}

		return nil
	},
}
