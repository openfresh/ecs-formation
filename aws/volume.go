package aws
import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"errors"
	"strings"
	"github.com/naoina/go-stringutil"
	"github.com/aws/aws-sdk-go/aws"
)


type VolumeInfo struct {
	Volume *ecs.Volume
	MountPoint *ecs.MountPoint
}

func CreateVolumeInfoItems(values []string) ([]*VolumeInfo, error) {

	volumes := []*VolumeInfo{}

	for _, value := range values {
		vi, err := CreateVolumeInfo(value)

		if err != nil {
			return []*VolumeInfo{}, err
		}

		volumes = append(volumes, vi)
	}

	return volumes, nil
}


func CreateVolumeInfo(value string) (*VolumeInfo, error) {

	if len(value) == 0 {
		return &VolumeInfo{}, errors.New("'volumes' element must not be empty.")
	}

	tokens := strings.Split(value, ":")
	length := len(tokens)

	if length == 0 {
		return &VolumeInfo{}, errors.New("'volumes' element must not be empty.")
	}

	var sourceVolume string
	var containerPath string
	var ro = false

	if length > 2 && tokens[2] == "ro" {
		ro = true
	}

	if length > 1 {
		containerPath = tokens[1]
	}

	if length == 1 {
		containerPath = tokens[0]
	}

	sourceVolume = tokens[0]
	volumeName, err := createVolumeName(&tokens[0])

	if err != nil {
		return &VolumeInfo{}, errors.New("'volumes' element must not be empty.")
	}

	return &VolumeInfo{
		Volume: &ecs.Volume{
			Name: aws.String(volumeName),
			Host: &ecs.HostVolumeProperties{
				SourcePath: aws.String(sourceVolume),
			},
		},
		MountPoint: &ecs.MountPoint{
			SourceVolume: aws.String(volumeName),
			ContainerPath: aws.String(containerPath),
			ReadOnly: &ro,
		},
	}, nil
}

func createVolumeName(path *string) (string, error) {

	if len(*path) == 0 {
		return "", errors.New("cannot create volume name from empty string.")
	}

	noDots := strings.Replace(*path, ".", "/", -1)
	tokens := strings.Split(noDots, "/")

	parts := []string{}
	for _, token := range tokens {
		s := stringutil.ToUpperCamelCase(token)
		parts = append(parts, s)
	}

	combined := strings.Join(parts, "")
	return combined, nil
}
