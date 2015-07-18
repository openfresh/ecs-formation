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
	"github.com/stormcat24/ecs-formation/cluster"
	"github.com/str1ngs/ansi/color"
)

type BlueGreenController struct {
	Ecs *aws.ECSManager
	ClusterController *cluster.ClusterController
	blueGreenDef []schema.BlueGreen
}

func NewBlueGreenController(ecs *aws.ECSManager, projectDir string) (*BlueGreenController, error) {

	ccon, errcc := cluster.NewClusterController(ecs, projectDir, "")

	if errcc != nil {
		return nil, errcc
	}

	con := &BlueGreenController{
		Ecs: ecs,
		ClusterController: ccon,
	}

	defs, errs := con.searchBlueGreen(projectDir)
	if errs != nil {
		return nil, errs
	}

	con.blueGreenDef = defs
	return con, nil
}

func (self *BlueGreenController) searchBlueGreen(projectDir string) ([]schema.BlueGreen, error) {

	clusterDir := projectDir + "/bluegreen"
	files, err := ioutil.ReadDir(clusterDir)

	bluegreenItems := []schema.BlueGreen{}

	if err != nil {
		return bluegreenItems, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".yml") {
			content, _ := ioutil.ReadFile(clusterDir + "/" + file.Name())

			bgMap, _ := schema.CreateBlueGreenMap(content)

			for _, bg := range bgMap {
				bluegreenItems = append(bluegreenItems, bg)
			}
		}
	}

	return bluegreenItems, nil
}

func (self *BlueGreenController) GetBlueGreenDefs() []schema.BlueGreen {
	return self.blueGreenDef
}

func (self *BlueGreenController) CreateBlueGreenPlan(blue schema.BlueGreenTarget, green schema.BlueGreenTarget,
	cplans []*plan.ClusterUpdatePlan) (*plan.BlueGreenPlan, error) {

	clusterMap := make(map[string]*plan.ClusterUpdatePlan, len(cplans))
	for _, cp := range cplans {
		clusterMap[cp.Name] = cp
	}

	bgPlan := plan.BlueGreenPlan{
		Blue: &plan.ServiceSet{
			LoadBalancer: blue.ElbName,
			ClusterUpdatePlan: clusterMap[blue.Cluster],
		},
		Green: &plan.ServiceSet{
			LoadBalancer: green.ElbName,
			ClusterUpdatePlan: clusterMap[green.Cluster],
		},
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


func (self *BlueGreenController) ApplyBlueGreenDeploys(plans []*plan.BlueGreenPlan) error {

	for _, plan := range plans {
		if err := self.ApplyBlueGreenDeploy(plan); err != nil {
			return err
		}
	}

	return nil
}

func (self *BlueGreenController) ApplyBlueGreenDeploy(bgplan *plan.BlueGreenPlan) error {

	apias := self.Ecs.AutoscalingApi()

	targetGreen := bgplan.Blue.HasOwnElb()

	var currentLabel *color.Escape
	var nextLabel *color.Escape
	var current *plan.ServiceSet
	var next *plan.ServiceSet
	primaryLb := bgplan.Blue.LoadBalancer
	standbyLb := bgplan.Green.LoadBalancer
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

	fmt.Printf("[INFO] Current status is '%s'\n", currentLabel)
	fmt.Printf("[INFO] Start Blue-Green Deployment: %s to %s ...\n", currentLabel, nextLabel)

	// deploy service
	fmt.Printf("[INFO] Updating %s@%s service at %s ...\n", next.NewService.Service, next.NewService.Cluster, nextLabel)
	if err := self.ClusterController.ApplyClusterPlan(next.ClusterUpdatePlan); err != nil {
		return err
	}

	fmt.Println("[INFO] Start to check whether service is running ...")
	self.ClusterController.WaitActiveService(next.NewService.Cluster, next.NewService.Service)
	fmt.Printf("[INFO] Service '%s' is running\n", next.NewService.Service)

	// attach next group to primary lb
	_, erratt := apias.AttachLoadBalancers(*next.AutoScalingGroup.AutoScalingGroupName, []string{
		primaryLb,
	})
	if erratt != nil {
		return erratt
	}
	fmt.Printf("[INFO] Attached to attach %s group to %s(primary).\n", nextLabel, primaryLb)

	errwlb := self.waitLoadBalancer(*next.AutoScalingGroup.AutoScalingGroupName, current.LoadBalancer)
	if errwlb != nil {
		return errwlb
	}
	fmt.Printf("[INFO] Added %s group to primary\n", nextLabel)

	// detach current group from primary lb
	_, errelbb := apias.DetachLoadBalancers(*current.AutoScalingGroup.AutoScalingGroupName, []string{
		primaryLb,
	})
	if errelbb != nil {
		return errelbb
	}
	fmt.Printf("[INFO] Detached %s group from %s(primary).\n", currentLabel, primaryLb)

	// detach next group from standby lb
	_, errelbg := apias.DetachLoadBalancers(*next.AutoScalingGroup.AutoScalingGroupName, []string{
		standbyLb,
	})
	if errelbg != nil {
		return errelbg
	}
	fmt.Printf("[INFO] Detached %s group from %s(standby).\n", nextLabel, standbyLb)

	// attach current group to standby lb
	_, errelba := apias.AttachLoadBalancers(*current.AutoScalingGroup.AutoScalingGroupName, []string{
		standbyLb,
	})
	if errelba != nil {
		return errelba
	}
	fmt.Printf("[INFO] Attached %s group to %s(standby).\n", currentLabel, standbyLb)

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