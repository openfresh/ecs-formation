package util

import (
	"os"
	"os/user"
	"strings"

	"github.com/spf13/viper"
)

func GetProjectDir() (string, error) {
	var projectDir string

	pd := viper.GetString("project_dir")
	if pd != "" {
		projectDir = pd
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		projectDir = wd
	}

	usr, _ := user.Current()
	projectDir = strings.Replace(projectDir, "~", usr.HomeDir, 1)

	return projectDir, nil
}
