package aws

import (
	"testing"
	"github.com/aws/aws-sdk-go/aws/awsutil"
)

//func TestToKeyValuePairs(t *testing.T) {
//
//	input := map[string]string{
//		"key1": "value1",
//		"key2": "value2",
//	}
//
//	actual := toKeyValuePairs(input)
//
//	if len(input) != len(actual) {
//		t.Errorf("expect length = %d, but actual length = %d", len(input), len(actual))
//	}
//
//	if *actual[0].Name != "key1" || *actual[0].Value != input["key1"] {
//		t.Errorf("expect %s=%s, but actual %s=%s", "key1", input["key1"], *actual[0].Name, *actual[0].Value)
//	}
//
//	if *actual[1].Name != "key2" || *actual[1].Value != input["key2"] {
//		t.Errorf("expect %s=%s, but actual %s=%s", "key2", input["key2"], *actual[1].Name, *actual[1].Value)
//	}
//}

func TestToPortMappingsOnlyInt(t *testing.T) {

	input := "3000"

	actual, _ := toPortMapping(input)

	if *actual.HostPort != 3000 ||
	*actual.ContainerPort != 3000 ||
	*actual.Protocol != "tcp" {
		t.Errorf("Unexpected value. Actual = %s", awsutil.Prettify(actual))
	}
}

func TestToPortMappingsPairTcp(t *testing.T) {

	input := "3000:4000/tcp"

	actual, _ := toPortMapping(input)

	if *actual.HostPort != 3000 ||
	*actual.ContainerPort != 4000 ||
	*actual.Protocol != "tcp" {
		t.Errorf("Unexpected value. Actual = %s", awsutil.Prettify(actual))
	}
}

func TestToPortMappingsPairUdp(t *testing.T) {

	input := "3000:4000/udp"

	actual, _ := toPortMapping(input)

	if *actual.HostPort != 3000 ||
	*actual.ContainerPort != 4000 ||
	*actual.Protocol != "udp" {
		t.Errorf("Unexpected value. Actual = %s", awsutil.Prettify(actual))
	}
}

func TestToPortMappings(t *testing.T) {

	input := []string{
		"5000:5001",
		"6000:6001/tcp",
		"10000:10001/udp",
		"20000/udp:20001/udp",
	}

	actual, err := toPortMappings(input)
	if err != nil {
		t.Error(err)
	}

	if len(input) != len(actual) {
		t.Errorf("expect length = %d, but actual length = %d", len(input), len(actual))
	}

	if *actual[0].HostPort != 5000 ||
	*actual[0].ContainerPort != 5001 ||
	*actual[0].Protocol != "tcp" {
		t.Errorf("Unexpected value. Actual = %s", awsutil.Prettify(actual[0]))
	}

	if *actual[1].HostPort != 6000 ||
	*actual[1].ContainerPort != 6001 ||
	*actual[1].Protocol != "tcp" {
		t.Errorf("Unexpected value. Actual = %s", awsutil.Prettify(actual[1]))
	}

	if *actual[2].HostPort != 10000 ||
	*actual[2].ContainerPort != 10001 ||
	*actual[2].Protocol != "udp" {
		t.Errorf("Unexpected value. Actual = %s", awsutil.Prettify(actual[2]))
	}

	if *actual[3].HostPort != 20000 ||
	*actual[3].ContainerPort != 20001 ||
	*actual[3].Protocol != "udp" {
		t.Errorf("Unexpected value. Actual = %s", awsutil.Prettify(actual[3]))
	}
}

func TestToVolumesFrom(t *testing.T) {

	input := []string{
		"container1",
		"container2:ro",
	}

	actual, err := toVolumesFroms(input)

	if err != nil {
		t.Error(err)
	}

	if len(input) != len(actual) {
		t.Errorf("expect length = %d, but actual length = %d", len(input), len(actual))
	}

	if *actual[0].SourceContainer != "container1" ||
	*actual[0].ReadOnly != false {
		t.Errorf("Unexpected value. Actual = %s", awsutil.Prettify(actual[0]))
	}

	if *actual[1].SourceContainer != "container2" ||
	*actual[1].ReadOnly != true {
		t.Errorf("Unexpected value. Actual = %s", awsutil.Prettify(actual[0]))
	}
}

func TestToVolumesFromEmpty(t *testing.T) {

	input := []string{}

	actual, err := toVolumesFroms(input)

	if err != nil {
		t.Error(err)
	}

	if len(actual) != 0 {
		t.Errorf("expect length = %d, but actual length = %d", len(input), len(actual))
	}
}