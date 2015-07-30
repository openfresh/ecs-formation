package aws
import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stormcat24/ecs-formation/util"
)


type AutoscalingApi struct {
	Credentials *credentials.Credentials
	Region      *string
}

func (self *AutoscalingApi) DescribeAutoScalingGroups(groups []string) (map[string]*autoscaling.Group, error) {

	svc := autoscaling.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: util.ConvertPointerString(groups),
	}

	asgmap := map[string]*autoscaling.Group{}
	result, err := svc.DescribeAutoScalingGroups(params)
	if err != nil {
		return asgmap, err
	}

	for _, asg := range result.AutoScalingGroups {
		asgmap[*asg.AutoScalingGroupName] = asg
	}

	return asgmap, nil
}

func (self *AutoscalingApi) DescribeLoadBalancerState(group string) (map[string]*autoscaling.LoadBalancerState, error) {

	svc := autoscaling.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &autoscaling.DescribeLoadBalancersInput{
		AutoScalingGroupName: aws.String(group),
	}

	lbmap := map[string]*autoscaling.LoadBalancerState{}
	result, err := svc.DescribeLoadBalancers(params)
	if err != nil {
		return lbmap, err
	}

	for _, lbs := range result.LoadBalancers {
		lbmap[*lbs.LoadBalancerName] = lbs
	}

	return lbmap, nil
}

func (self *AutoscalingApi) AttachLoadBalancers(group string, lb []string) (*autoscaling.AttachLoadBalancersOutput, error) {

	svc := autoscaling.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &autoscaling.AttachLoadBalancersInput{
		AutoScalingGroupName: aws.String(group),
		LoadBalancerNames: util.ConvertPointerString(lb),
	}

	return svc.AttachLoadBalancers(params)
}

func (self *AutoscalingApi) DetachLoadBalancers(group string, lb []string) (*autoscaling.DetachLoadBalancersOutput, error) {

	svc := autoscaling.New(&aws.Config{
		Region: self.Region,
		Credentials: self.Credentials,
	})

	params := &autoscaling.DetachLoadBalancersInput{
		AutoScalingGroupName: aws.String(group),
		LoadBalancerNames: util.ConvertPointerString(lb),
	}

	return svc.DetachLoadBalancers(params)
}