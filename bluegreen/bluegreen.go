package bluegreen

import (
	"io/ioutil"
	"github.com/stormcat24/ecs-formation/aws"
	"github.com/stormcat24/ecs-formation/schema"
	"strings"
	"github.com/stormcat24/ecs-formation/plan"
	"time"
	"fmt"
	"errors"
	"github.com/stormcat24/ecs-formation/service"
	"github.com/stormcat24/ecs-formation/logger"
	"github.com/str1ngs/ansi/color"
	"regexp"
)

type BlueGreenController struct {
	Ecs *aws.ECSManager
	ClusterController *service.ServiceController
	blueGreenMap map[string]*schema.BlueGreen
	TargetResource string
}

func NewBlueGreenController(ecs *aws.ECSManager, projectDir string, targetResource string) (*BlueGreenController, error) {

	ccon, errcc := service.NewServiceController(ecs, projectDir, "")

	if errcc != nil {
		return nil, errcc
	}

	con := &BlueGreenController{
		Ecs: ecs,
		ClusterController: ccon,
		TargetResource: targetResource,
	}

	defs, errs := con.searchBlueGreen(projectDir)
	if errs != nil {
		return nil, errs
	}

	con.blueGreenMap = defs
	return con, nil
}

func (self *BlueGreenController) searchBlueGreen(projectDir string) (map[string]*schema.BlueGreen, error) {

	clusterDir := projectDir + "/bluegreen"
	files, err := ioutil.ReadDir(clusterDir)

	merged := map[string]*schema.BlueGreen{}

	if err != nil {
		return merged, err
	}

	filePattern := regexp.MustCompile("^(.+)\\.yml$")

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".yml") {
			content, _ := ioutil.ReadFile(clusterDir + "/" + file.Name())
			tokens := filePattern.FindStringSubmatch(file.Name())
			name := tokens[1]

			bg, err := schema.CreateBlueGreen(content)
			if err != nil {
				return merged, err
			}
			merged[name] = bg
		}
	}

	return merged, nil
}

func (self *BlueGreenController) GetBlueGreenMap() map[string]*schema.BlueGreen {
	return self.blueGreenMap
}

func (self *BlueGreenController) CreateBlueGreenPlans(bgmap map[string]*schema.BlueGreen, cplans []*plan.ServiceUpdatePlan) ([]*plan.BlueGreenPlan, error) {

	bgPlans := []*plan.BlueGreenPlan{}

	for name, bg := range bgmap {

		if len(self.TargetResource) == 0 || self.TargetResource == name {

			bgplan, err := self.CreateBlueGreenPlan(bg, cplans)
			if err != nil {
				return bgPlans, err
			}

			if bgplan.Blue.CurrentService == nil {
				return bgPlans, errors.New(fmt.Sprintf("Service '%s' is not found. ", bg.Blue.Service))
			}

			if bgplan.Green.CurrentService == nil {
				return bgPlans, errors.New(fmt.Sprintf("Service '%s' is not found. ", bg.Green.Service))
			}

			if bgplan.Blue.AutoScalingGroup == nil {
				return bgPlans, errors.New(fmt.Sprintf("AutoScaling Group '%s' is not found. ", bg.Blue.AutoscalingGroup))
			}

			if bgplan.Green.AutoScalingGroup == nil {
				return bgPlans, errors.New(fmt.Sprintf("AutoScaling Group '%s' is not found. ", bg.Green.AutoscalingGroup))
			}

			if bgplan.Blue.ClusterUpdatePlan == nil {
				return bgPlans, errors.New(fmt.Sprintf("ECS Cluster '%s' is not found. ", bg.Blue.Cluster))
			}

			if bgplan.Green.ClusterUpdatePlan == nil {
				return bgPlans, errors.New(fmt.Sprintf("ECS Cluster '%s' is not found. ", bg.Green.Cluster))
			}

			bgPlans = append(bgPlans, bgplan)
		}
	}

	return bgPlans, nil
}

func (self *BlueGreenController) CreateBlueGreenPlan(bluegreen *schema.BlueGreen, cplans []*plan.ServiceUpdatePlan) (*plan.BlueGreenPlan, error) {

	blue := bluegreen.Blue
	green := bluegreen.Green

	clusterMap := make(map[string]*plan.ServiceUpdatePlan, len(cplans))
	for _, cp := range cplans {
		clusterMap[cp.Name] = cp
	}

	bgPlan := plan.BlueGreenPlan{
		Blue: &plan.ServiceSet{
			ClusterUpdatePlan: clusterMap[blue.Cluster],
		},
		Green: &plan.ServiceSet{
			ClusterUpdatePlan: clusterMap[green.Cluster],
		},
		PrimaryElb: bluegreen.PrimaryElb,
		StandbyElb: bluegreen.StandbyElb,
	}

	// describe services
	bsrv, _ := self.Ecs.ServiceApi().DescribeService(blue.Cluster, []*string{
		&blue.Service,
	})

	bgPlan.Blue.NewService = &blue
	if len(bsrv.Services) > 0 {
		bgPlan.Blue.CurrentService = bsrv.Services[0]
	}

	gsrv, _ := self.Ecs.ServiceApi().DescribeService(green.Cluster, []*string{
		&green.Service,
	})

	bgPlan.Green.NewService = &green
	if len(gsrv.Services) > 0 {
		bgPlan.Green.CurrentService = gsrv.Services[0]
	}

	// describe autoscaling group
	asgmap, err := self.Ecs.AutoscalingApi().DescribeAutoScalingGroups([]string {
		blue.AutoscalingGroup,
		green.AutoscalingGroup,
	})

	if err != nil {
		return &bgPlan, err
	}

	if asgblue, ok := asgmap[blue.AutoscalingGroup]; ok {
		bgPlan.Blue.AutoScalingGroup = asgblue
	}

	if asggreen, ok := asgmap[green.AutoscalingGroup]; ok {
		bgPlan.Green.AutoScalingGroup = asggreen
	}

	return &bgPlan, nil
}


func (self *BlueGreenController) ApplyBlueGreenDeploys(plans []*plan.BlueGreenPlan, nodeploy bool) error {

	for _, plan := range plans {
		if err := self.ApplyBlueGreenDeploy(plan, nodeploy); err != nil {
			return err
		}
	}

	return nil
}

func (self *BlueGreenController) ApplyBlueGreenDeploy(bgplan *plan.BlueGreenPlan, nodeploy bool) error {

	apias := self.Ecs.AutoscalingApi()

	targetGreen := bgplan.IsBluwWithPrimaryElb()

	var currentLabel *color.Escape
	var nextLabel *color.Escape
	var current *plan.ServiceSet
	var next *plan.ServiceSet
	primaryLb := bgplan.PrimaryElb
	standbyLb := bgplan.StandbyElb
	if targetGreen {
		current = bgplan.Blue
		next = bgplan.Green
		currentLabel = color.Cyan("blue")
		nextLabel = color.Green("green")
	} else {
		current = bgplan.Green
		next = bgplan.Blue
		currentLabel = color.Green("green")
		nextLabel = color.Cyan("blue")
	}

	logger.Main.Infof("Current status is '%s'", currentLabel)
	logger.Main.Infof("Start Blue-Green Deployment: %s to %s ...", currentLabel, nextLabel)
	if nodeploy {
		logger.Main.Infof("Without deployment. It only replaces load balancers.")
	} else {
		// deploy service
		logger.Main.Infof("Updating %s@%s service at %s ...", next.NewService.Service, next.NewService.Cluster, nextLabel)
		if err := self.ClusterController.ApplyServicePlan(next.ClusterUpdatePlan); err != nil {
			return err
		}
	}

	// attach next group to primary lb
	_, erratt := apias.AttachLoadBalancers(*next.AutoScalingGroup.AutoScalingGroupName, []string{
		primaryLb,
	})
	if erratt != nil {
		return erratt
	}
	logger.Main.Infof("Attached to attach %s group to %s(primary).", nextLabel, primaryLb)

	errwlb := self.waitLoadBalancer(*next.AutoScalingGroup.AutoScalingGroupName, primaryLb)
	if errwlb != nil {
		return errwlb
	}
	logger.Main.Infof("Added %s group to primary", nextLabel)

	// detach current group from primary lb
	_, errelbb := apias.DetachLoadBalancers(*current.AutoScalingGroup.AutoScalingGroupName, []string{
		primaryLb,
	})
	if errelbb != nil {
		return errelbb
	}
	logger.Main.Infof("Detached %s group from %s(primary).", currentLabel, primaryLb)

	// detach next group from standby lb
	_, errelbg := apias.DetachLoadBalancers(*next.AutoScalingGroup.AutoScalingGroupName, []string{
		standbyLb,
	})
	if errelbg != nil {
		return errelbg
	}
	logger.Main.Infof("Detached %s group from %s(standby).", nextLabel, standbyLb)

	// attach current group to standby lb
	_, errelba := apias.AttachLoadBalancers(*current.AutoScalingGroup.AutoScalingGroupName, []string{
		standbyLb,
	})
	if errelba != nil {
		return errelba
	}
	logger.Main.Infof("Attached %s group to %s(standby).", currentLabel, standbyLb)

	return nil
}

func (self *BlueGreenController) waitLoadBalancer(group string, lb string) error {

	api := self.Ecs.AutoscalingApi()

	for {
		time.Sleep(5 * time.Second)

		result, err := api.DescribeLoadBalancerState(group)

		if err != nil {
			return err
		}

		if lbs, ok := result[lb]; ok {

			// *** LoadbalancerState
			// Adding - The instances in the group are being registered with the load balancer.
			// Added - All instances in the group are registered with the load balancer.
			// InService - At least one instance in the group passed an ELB health check.
			// Removing - The instances are being deregistered from the load balancer. If connection draining is enabled, Elastic Load Balancing waits for in-flight requests to complete before deregistering the instances.
			if *lbs.State == "Added" || *lbs.State == "InService" {
				return nil
			}

		} else {
			return errors.New(fmt.Sprintf("cannot get load balanracer '%s'", lb))
		}

	}

	return nil

}