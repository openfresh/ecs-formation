package schema

type BlueGreen struct {
	Blue  BlueGreenTarget    `yaml:"blue"`
	Green BlueGreenTarget    `yaml:"green"`
}

type BlueGreenTarget struct {
	Cluster string    `yaml:"cluster"`
	Service string    `yaml:"service"`
	ElbName string    `yaml:"elb_name"`
	AutoscalingGroup string    `yaml:"autoscaling_group"`
}
