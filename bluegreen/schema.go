package bluegreen

type BlueGreen struct {
	Blue       BlueGreenTarget    `yaml:"blue"`
	Green      BlueGreenTarget    `yaml:"green"`
	PrimaryElb string        `yaml:"primary_elb"`
	StandbyElb string        `yaml:"standby_elb"`
	ChainElb   []BlueGreenChainElb `yaml:"chain_elb"`
}

type BlueGreenChainElb struct {
	PrimaryElb string        `yaml:"primary_elb"`
	StandbyElb string        `yaml:"standby_elb"`
}

type BlueGreenTarget struct {
	Cluster          string    `yaml:"cluster"`
	Service          string    `yaml:"service"`
	AutoscalingGroup string    `yaml:"autoscaling_group"`
}
