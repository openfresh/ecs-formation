package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/openfresh/ecs-formation/client"
	"github.com/openfresh/ecs-formation/service/types"
	"github.com/openfresh/ecs-formation/util"
)

type BlueGreenService interface {
	GetBlueGreenMap() map[string]*types.BlueGreen
	CreateBlueGreenPlans(bgmap map[string]*types.BlueGreen, cplans []*types.ServiceUpdatePlan) ([]*types.BlueGreenPlan, error)
	CreateClusterService() (ClusterService, error)
	ApplyBlueGreenDeploys(clusterService ClusterService, plans []*types.BlueGreenPlan, nodeploy bool) error
}

type ConcreteBlueGreenService struct {
	awsCli        client.AWSClient
	projectDir    string
	blueGreenName string
	blueGreenMap  map[string]*types.BlueGreen
	params        map[string]string
}

func NewBlueGreenService(projectDir string, blueGreenName string, params map[string]string) (BlueGreenService, error) {

	defs, err := searchBlueGreen(projectDir, blueGreenName, params)
	if err != nil {
		return nil, err
	}

	return &ConcreteBlueGreenService{
		awsCli:        client.AWSCli,
		projectDir:    projectDir,
		blueGreenName: blueGreenName,
		blueGreenMap:  defs,
		params:        params,
	}, nil
}

func searchBlueGreen(projectDir string, blueGreenName string, params map[string]string) (map[string]*types.BlueGreen, error) {
	clusterDir := projectDir + "/bluegreen"
	bgmap := map[string]*types.BlueGreen{}

	filePattern := regexp.MustCompile(`^.+\/(.+)\.yml$`)

	filepath.Walk(clusterDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || !strings.HasSuffix(path, ".yml") {
			return nil
		}

		if blueGreenName != "" && !strings.HasSuffix(path, fmt.Sprintf("%s.yml", blueGreenName)) {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		merged := util.MergeYamlWithParameters(content, params)
		tokens := filePattern.FindStringSubmatch(path)
		name := tokens[1]

		bg, err := createBlueGreen(merged)
		if err != nil {
			return err
		}
		bgmap[name] = bg

		return nil
	})

	return bgmap, nil
}

func createBlueGreen(data string) (*types.BlueGreen, error) {

	bg := types.BlueGreen{}
	if err := yaml.Unmarshal([]byte(data), &bg); err != nil {
		return nil, errors.New(fmt.Sprintf("%v\n\n%v", err.Error(), data))
	}

	return &bg, nil
}

func (s ConcreteBlueGreenService) GetBlueGreenMap() map[string]*types.BlueGreen {
	return s.blueGreenMap
}

func (s ConcreteBlueGreenService) CreateBlueGreenPlans(bgmap map[string]*types.BlueGreen, cplans []*types.ServiceUpdatePlan) ([]*types.BlueGreenPlan, error) {
	bgPlans := []*types.BlueGreenPlan{}

	for name, bg := range bgmap {
		if s.blueGreenName == "" || s.blueGreenName == name {
			bgplan, err := s.createBlueGreenPlan(bg, cplans)
			if err != nil {
				return bgPlans, err
			}

			if bgplan.Blue.CurrentService == nil {
				return bgPlans, fmt.Errorf("Service '%s' is not found. ", bg.Blue.Service)
			}

			if bgplan.Green.CurrentService == nil {
				return bgPlans, fmt.Errorf("Service '%s' is not found. ", bg.Green.Service)
			}

			if bgplan.Blue.AutoScalingGroup == nil {
				return bgPlans, fmt.Errorf("AutoScaling Group '%s' is not found. ", bg.Blue.AutoscalingGroup)
			}

			if bgplan.Green.AutoScalingGroup == nil {
				return bgPlans, fmt.Errorf("AutoScaling Group '%s' is not found. ", bg.Green.AutoscalingGroup)
			}

			if bgplan.Blue.ClusterUpdatePlan == nil {
				return bgPlans, fmt.Errorf("ECS Cluster '%s' is not found. ", bg.Blue.Cluster)
			}

			if bgplan.Green.ClusterUpdatePlan == nil {
				return bgPlans, fmt.Errorf("ECS Cluster '%s' is not found. ", bg.Green.Cluster)
			}

			bgPlans = append(bgPlans, bgplan)
		}
	}

	return bgPlans, nil
}

func (s ConcreteBlueGreenService) createBlueGreenPlan(bluegreen *types.BlueGreen, cplans []*types.ServiceUpdatePlan) (*types.BlueGreenPlan, error) {

	blue := bluegreen.Blue
	green := bluegreen.Green

	clusterMap := make(map[string]*types.ServiceUpdatePlan, len(cplans))
	for _, cp := range cplans {
		clusterMap[cp.Name] = cp
	}

	bgPlan := types.BlueGreenPlan{
		Blue: &types.ServiceSet{
			ClusterUpdatePlan: clusterMap[blue.Cluster],
		},
		Green: &types.ServiceSet{
			ClusterUpdatePlan: clusterMap[green.Cluster],
		},
		PrimaryElb: bluegreen.PrimaryElb,
		StandbyElb: bluegreen.StandbyElb,
		ChainElb:   bluegreen.ChainElb,
		ElbV2:      bluegreen.ElbV2,
	}

	// describe services
	bsrv, err := s.awsCli.ECS.DescribeService(blue.Cluster, []*string{aws.String(blue.Service)})
	if err != nil {
		return nil, err
	}

	bgPlan.Blue.NewService = &blue
	if len(bsrv.Services) > 0 {
		bgPlan.Blue.CurrentService = bsrv.Services[0]
	}

	gsrv, err := s.awsCli.ECS.DescribeService(green.Cluster, []*string{aws.String(green.Service)})
	if err != nil {
		return nil, err
	}

	bgPlan.Green.NewService = &green
	if len(gsrv.Services) > 0 {
		bgPlan.Green.CurrentService = gsrv.Services[0]
	}

	// describe autoscaling group
	asgmap, err := s.awsCli.Autoscaling.DescribeAutoScalingGroups([]string{
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

func (s ConcreteBlueGreenService) CreateClusterService() (ClusterService, error) {

	bg, _ := s.blueGreenMap[s.blueGreenName]
	if bg == nil {
		return nil, fmt.Errorf("load bluegreen data is failed. %s", s.blueGreenName)
	}

	clusters := []string{
		bg.Blue.Cluster,
		bg.Green.Cluster,
	}
	// TODO service名渡すか？ bluegreenコマンドでもserviceを渡せるようにする必要あり
	return NewClusterService(s.projectDir, clusters, "", s.params)
}

func (s ConcreteBlueGreenService) ApplyBlueGreenDeploys(clusterService ClusterService, plans []*types.BlueGreenPlan, nodeploy bool) error {

	for _, plan := range plans {
		if err := s.applyBlueGreenDeploy(clusterService, plan, nodeploy); err != nil {
			return err
		}
	}

	return nil
}

func (s ConcreteBlueGreenService) applyBlueGreenDeploy(clusterService ClusterService, bgplan *types.BlueGreenPlan, nodeploy bool) error {

	switcher := NewELBSwitcher(s.awsCli, bgplan)
	return switcher.Apply(clusterService, bgplan, nodeploy)
}
