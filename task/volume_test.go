package task

import (
	"testing"
	"github.com/aws/aws-sdk-go/aws/awsutil"
)


func TestCreateVolumeInfoEmpty(t *testing.T) {

	input := ""

	_, err := CreateVolumeInfo(input)

	if err == nil {
		t.Error("expect error, but success.")
	}
}

func TestCreateVolumeInfoOneOnly(t *testing.T) {

	input := "/var/log/hoge"

	actual, err := CreateVolumeInfo(input)

	if err != nil {
		t.Error(err)
	}

	if *actual.Volume.Name != "VarLogHoge" ||
		*actual.Volume.Host.SourcePath != input ||
		*actual.MountPoint.ContainerPath != input ||
		*actual.MountPoint.SourceVolume != "VarLogHoge" ||
		*actual.MountPoint.ReadOnly != false {
		t.Errorf("Unexpected value. Actual = %s", awsutil.Prettify(actual))
	}
}

func TestCreateVolumeInfoPair(t *testing.T) {

	input := "/var/log/container/nginx:/var/log/nginx"

	actual, err := CreateVolumeInfo(input)

	if err != nil {
		t.Error(err)
	}

	if *actual.Volume.Name != "VarLogContainerNginx" ||
		*actual.Volume.Host.SourcePath != "/var/log/container/nginx" ||
		*actual.MountPoint.ContainerPath != "/var/log/nginx" ||
		*actual.MountPoint.SourceVolume != "VarLogContainerNginx" ||
		*actual.MountPoint.ReadOnly != false {
		t.Errorf("Unexpected value. Actual = %s", awsutil.Prettify(actual))
	}
}

func TestCreateVolumeInfoPairWithReadOnly(t *testing.T) {

	input := "/var/log/container/nginx:/var/log/nginx:ro"

	actual, err := CreateVolumeInfo(input)

	if err != nil {
		t.Error(err)
	}

	if *actual.Volume.Name != "VarLogContainerNginx" ||
	*actual.Volume.Host.SourcePath != "/var/log/container/nginx" ||
	*actual.MountPoint.ContainerPath != "/var/log/nginx" ||
	*actual.MountPoint.SourceVolume != "VarLogContainerNginx" ||
	*actual.MountPoint.ReadOnly != true {
		t.Errorf("Unexpected value. Actual = %s", awsutil.Prettify(actual))
	}
}

func TestCreateVolumeNameEmpty(t *testing.T) {

	input := ""

	_, err := createVolumeName(&input)

	if err == nil {
		t.Error("expect error, but success.")
	}
}

func TestCreateVolumeNameSlash(t *testing.T) {

	expected := "VarLogHoge"
	input := "/var/log/hoge"

	actual, err := createVolumeName(&input)
	if err != nil {
		t.Error(err)
	}

	if actual != expected {
		t.Errorf("expect '%s', but actual '%s'.", expected, actual)
	}
}

func TestCreateVolumeNameDots(t *testing.T) {

	expected := "VarRunDockerSock"
	input := "/var/run/docker.sock"

	actual, err := createVolumeName(&input)
	if err != nil {
		t.Error(err)
	}

	if actual != expected {
		t.Errorf("expect '%s', but actual '%s'.", expected, actual)
	}
}
