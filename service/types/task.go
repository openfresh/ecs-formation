package types

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/joho/godotenv"
	"github.com/openfresh/ecs-formation/client/s3"
	"gopkg.in/yaml.v2"
)

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
	Memory                 *int64            `yaml:"memory"`
	MemoryReservation      *int64            `yaml:"memory_reservation"`
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

type VolumeInfo struct {
	Volume     *ecs.Volume
	MountPoint *ecs.MountPoint
}

func CreateTaskDefinition(taskDefName string, data string, basedir string, s3Cli s3.Client) (*TaskDefinition, error) {

	containerMap := map[string]ContainerDefinition{}
	if err := yaml.Unmarshal([]byte(data), &containerMap); err != nil {
		return nil, errors.New(fmt.Sprintf("%v\n\n%v", err.Error(), data))
	}

	containers := map[string]*ContainerDefinition{}
	for name, container := range containerMap {
		con := container
		con.Name = name

		environment := map[string]string{}
		if len(container.EnvFiles) > 0 {
			for _, envfile := range container.EnvFiles {
				var path string
				if envfile[0:10] == "https://s3" {
					_path, err := downloadS3(envfile, s3Cli)
					if err != nil {
						return nil, err
					}
					path = _path
					defer os.Remove(_path)
				} else if filepath.IsAbs(envfile) {
					path = envfile
				} else {
					path = fmt.Sprintf("%s/%s", basedir, envfile)
				}

				envmap, err := readEnvFile(path)
				if err != nil {
					return nil, err
				}

				for key, value := range envmap {
					environment[key] = value
				}
			}
		}

		for key, value := range container.Environment {
			environment[key] = value
		}

		con.Environment = environment
		containers[name] = &con
	}

	taskDef := TaskDefinition{
		Name:                 taskDefName,
		ContainerDefinitions: containers,
	}

	return &taskDef, nil
}

func downloadS3(path string, s3Cli s3.Client) (string, error) {
	u, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	ps := strings.Split(u.Path, "/")
	bucket := ps[:2][1]
	key := strings.Join(ps[2:], "/")

	obj, err := s3Cli.GetObject(bucket, key)
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(obj.Body)
	if err != nil {
		return "", err
	}

	tempfile, err := ioutil.TempFile("", "ecs-formation")
	if err != nil {
		return "", err
	}
	defer tempfile.Close()
	tempfile.Write(b)
	return tempfile.Name(), nil
}

func readEnvFile(envpath string) (map[string]string, error) {

	envmap, err := godotenv.Read(envpath)
	if err != nil {
		return map[string]string{}, err
	}

	return envmap, nil
}
