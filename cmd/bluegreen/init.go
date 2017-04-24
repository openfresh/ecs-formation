package bluegreen

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/openfresh/ecs-formation/client"
	cmdutil "github.com/openfresh/ecs-formation/cmd/util"
	"github.com/openfresh/ecs-formation/service"
	"github.com/openfresh/ecs-formation/service/types"
	"github.com/openfresh/ecs-formation/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	projectDir    string
	bluegreenName string
	parameters    map[string]string
	jsonOutput    bool
	noDeploy      bool
)

type BlueGreenPlanJson struct {
	Blue       BlueGreenServiceJson
	Green      BlueGreenServiceJson
	Active     string
	PrimaryElb string
	StandbyElb string
}

type BlueGreenServiceJson struct {
	ClusterARN          string
	AutoScalingGroupARN string
	Instances           []*autoscaling.Instance
	TaskDefinition      string
	DesiredCount        int64
	PendingCount        int64
	RunningCount        int64
}

var BlueGreenCmd = &cobra.Command{
	Use:   "bluegreen",
	Short: "Manage Amazon ECS Service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		pd, err := cmdutil.GetProjectDir()
		if err != nil {
			return err
		}
		projectDir = pd

		region := viper.GetString("aws_region")
		client.Init(region, false)

		bg, err := cmd.Flags().GetString("group")
		if err != nil {
			return err
		}
		if bg == "" {
			return errors.New("-g (--group) is required")
		}
		bluegreenName = bg

		paramTokens, err := cmd.Flags().GetStringSlice("parameter")
		if err != nil {
			return err
		}
		parameters = util.ParseKeyValues(paramTokens)

		jo, err := cmd.Flags().GetBool("json-output")
		if err != nil {
			return err
		}
		jsonOutput = jo

		nd, err := cmd.Flags().GetBool("no-deploy")
		if err != nil {
			return err
		}
		noDeploy = nd

		return nil
	},
}

func init() {
	BlueGreenCmd.AddCommand(planCmd)
	BlueGreenCmd.AddCommand(applyCmd)

	BlueGreenCmd.PersistentFlags().StringP("group", "g", "", "BlueGreen group name")
	BlueGreenCmd.PersistentFlags().StringSliceP("parameter", "p", make([]string, 0), "parameter 'key=value'")
	BlueGreenCmd.PersistentFlags().BoolP("no-deploy", "", false, "Only change load balancer")
	BlueGreenCmd.PersistentFlags().BoolP("json-output", "j", false, "Print json format")
}

func createBlueGreenPlans(bgsrv service.BlueGreenService, csrv service.ClusterService) ([]*types.BlueGreenPlan, error) {
	if jsonOutput {
		util.Output = false
		defer func() {
			util.Output = true
		}()
	}

	bgmap := bgsrv.GetBlueGreenMap()

	cplans, err := csrv.CreateServiceUpdatePlans()
	if err != nil {
		return make([]*types.BlueGreenPlan, 0), err
	}

	bgplans, err := bgsrv.CreateBlueGreenPlans(bgmap, cplans)
	if err != nil {
		return bgplans, err
	}

	jsonItems := []BlueGreenPlanJson{}
	for _, bgplan := range bgplans {
		util.PrintlnCyan("    Blue:")
		util.PrintlnCyan(fmt.Sprintf("        Cluster = %s", bgplan.Blue.NewService.Cluster))
		util.PrintlnCyan(fmt.Sprintf("        AutoScalingGroupARN = %s", *bgplan.Blue.AutoScalingGroup.AutoScalingGroupARN))
		util.PrintlnCyan("        Current services as follows:")
		for _, bcss := range bgplan.Blue.ClusterUpdatePlan.CurrentServices {
			bcs := bcss.Service
			util.PrintlnCyan(fmt.Sprintf("            %s:", *bcs.ServiceName))
			util.PrintlnCyan(fmt.Sprintf("                ServiceARN = %s", *bcs.ServiceArn))
			util.PrintlnCyan(fmt.Sprintf("                TaskDefinition = %s", *bcs.TaskDefinition))
			util.PrintlnCyan(fmt.Sprintf("                DesiredCount = %d", *bcs.DesiredCount))
			util.PrintlnCyan(fmt.Sprintf("                PendingCount = %d", *bcs.PendingCount))
			util.PrintlnCyan(fmt.Sprintf("                RunningCount = %d", *bcs.RunningCount))
		}

		var active string
		if bgplan.IsBlueWithPrimaryElb() {
			active = "blue"
		} else {
			active = "green"
		}

		util.PrintlnGreen("    Green:")
		util.PrintlnGreen(fmt.Sprintf("        Cluster = %s", bgplan.Green.NewService.Cluster))
		util.PrintlnGreen(fmt.Sprintf("        AutoScalingGroupARN = %s", *bgplan.Green.AutoScalingGroup.AutoScalingGroupARN))
		util.PrintlnGreen("        Current services as follows:")
		for _, gcss := range bgplan.Green.ClusterUpdatePlan.CurrentServices {
			gcs := gcss.Service
			util.PrintlnGreen(fmt.Sprintf("            %s:", *gcs.ServiceName))
			util.PrintlnGreen(fmt.Sprintf("                ServiceARN = %s", *gcs.ServiceArn))
			util.PrintlnGreen(fmt.Sprintf("                TaskDefinition = %s", *gcs.TaskDefinition))
			util.PrintlnGreen(fmt.Sprintf("                DesiredCount = %d", *gcs.DesiredCount))
			util.PrintlnGreen(fmt.Sprintf("                PendingCount = %d", *gcs.PendingCount))
			util.PrintlnGreen(fmt.Sprintf("                RunningCount = %d", *gcs.RunningCount))
		}

		util.Println()

		jsonItems = append(jsonItems, BlueGreenPlanJson{

			Blue: BlueGreenServiceJson{
				ClusterARN:          *bgplan.Blue.CurrentService.ClusterArn,
				AutoScalingGroupARN: *bgplan.Blue.AutoScalingGroup.AutoScalingGroupARN,
				Instances:           bgplan.Blue.AutoScalingGroup.Instances,
				TaskDefinition:      *bgplan.Blue.CurrentService.TaskDefinition,
				DesiredCount:        *bgplan.Blue.CurrentService.DesiredCount,
				PendingCount:        *bgplan.Blue.CurrentService.PendingCount,
				RunningCount:        *bgplan.Blue.CurrentService.RunningCount,
			},
			Green: BlueGreenServiceJson{
				ClusterARN:          *bgplan.Green.CurrentService.ClusterArn,
				AutoScalingGroupARN: *bgplan.Green.AutoScalingGroup.AutoScalingGroupARN,
				Instances:           bgplan.Green.AutoScalingGroup.Instances,
				TaskDefinition:      *bgplan.Green.CurrentService.TaskDefinition,
				DesiredCount:        *bgplan.Green.CurrentService.DesiredCount,
				PendingCount:        *bgplan.Green.CurrentService.PendingCount,
				RunningCount:        *bgplan.Green.CurrentService.RunningCount,
			},
			PrimaryElb: bgplan.PrimaryElb,
			StandbyElb: bgplan.StandbyElb,
			Active:     active,
		})

	}

	if jsonOutput {
		bt, err := json.Marshal(&jsonItems)
		if err != nil {
			return bgplans, err
		}
		fmt.Println(string(bt))
	}

	return bgplans, nil
}
