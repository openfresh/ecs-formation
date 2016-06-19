package types

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"

	"github.com/stormcat24/ecs-formation/client/s3"
)

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
