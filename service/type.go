package service

type TaskDefinition struct {
	Name                 string
	ContainerDefinitions map[string]*ContainerDefinition
}

type ContainerDefinition struct {
	Name                   string
	Image                  string            `yaml:"image"`
	Ports                  []string          `yaml:"ports"`
	Environment            map[string]string `yaml:"environment"`
	EnvFiles               []string          `yaml:"env_file"`
	Links                  []string          `yaml:"links"`
	Volumes                []string          `yaml:"volumes"`
	VolumesFrom            []string          `yaml:"volumes_from"`
	Memory                 int64             `yaml:"memory"`
	CPUUnits               int64             `yaml:"cpu_units"`
	Essential              bool              `yaml:"essential"`
	EntryPoint             string            `yaml:"entry_point"`
	Command                string            `yaml:"command"`
	DisableNetworking      bool              `yaml:"disable_networking"`
	DNSSearchDomains       []string          `yaml:"dns_search"`
	DNSServers             []string          `yaml:"dns"`
	DockerLabels           map[string]string `yaml:"labels"`
	DockerSecurityOptions  []string          `yaml:"security_opt"`
	ExtraHosts             []string          `yaml:"extra_hosts"`
	Hostname               string            `yaml:"hostname"`
	LogDriver              string            `yaml:"log_driver"`
	LogOpt                 map[string]string `yaml:"log_opt"`
	Privileged             bool              `yaml:"privileged"`
	ReadonlyRootFilesystem bool              `yaml:"read_only"`
	Ulimits                map[string]Ulimit `yaml:"ulimits"`
	User                   string            `yaml:"user"`
	WorkingDirectory       string            `yaml:"working_dir"`
}

type Ulimit struct {
	Soft int64 `yaml:"soft"`
	Hard int64 `yaml:"hard"`
}

type TaskUpdatePlan struct {
	Name          string
	NewContainers map[string]*ContainerDefinition
}
