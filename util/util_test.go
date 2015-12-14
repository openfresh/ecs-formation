package util

import (
	"testing"
)

func TestConvertPointerString(t *testing.T) {

	input := []string{
		"value1",
		"value2",
	}

	actual := ConvertPointerString(input)

	if len(input) != len(actual) {
		t.Fatalf("expect length = %d, but actual length = %d", len(input), len(actual))
	}

	if &input[0] != actual[0] {
		t.Errorf("expect[0] address '%x', but actual[0] address = %x", &input[0], actual[0])
	}

	if &input[1] != actual[1] {
		t.Errorf("expect[1] address '%x', but actual[1] address = %x", &input[0], actual[0])
	}

}

func TestParseKeyValues(t *testing.T) {

	input := []string{
		"key1=value1",
		"key2= value2 ",
		"key3 = value3",
		" key4 = value4",
	}

	actual := ParseKeyValues(input)

	if len(input) != len(actual) {
		t.Fatalf("expect length = %d, but actual length = %d", len(input), len(actual))
	}

	if value, ok := actual["key1"]; ok {
		if value != "value1" {
			t.Errorf("expect[%v] is '%v', but actual is '%v'", "key1", "value1", value)
		}
	} else {
		t.Errorf("actual do not include '%v'.", "key1")
	}

	if value, ok := actual["key2"]; ok {
		if value != "value2" {
			t.Errorf("expect[%v] is '%v', but actual is '%v'", "key2", "value2", value)
		}
	} else {
		t.Errorf("actual do not include '%v'.", "key2")
	}

	if value, ok := actual["key3"]; ok {
		if value != "value3" {
			t.Errorf("expect[%v] is '%v', but actual is '%v'", "key3", "value3", value)
		}
	} else {
		t.Errorf("actual do not include '%v'.", "key3")
	}

	if value, ok := actual["key4"]; ok {
		if value != "value4" {
			t.Errorf("expect[%v] is '%v', but actual is '%v'", "key4", "value4", value)
		}
	} else {
		t.Errorf("actual do not include '%v'.", "key4")
	}

}

func TestMergeYamlWithParameters(t *testing.T) {

	yaml := `
	nginx:
		image: stormcat24/nginx:${NGINX_VERSION}
		ports:
			- 80:${NGINX_PORT}
		environment:
			PARAM: "${PARAM}"
	`

	expect := `
	nginx:
		image: stormcat24/nginx:latest
		ports:
			- 80:80
		environment:
			PARAM: ""
	`

	params := map[string]string{
		"NGINX_VERSION": "latest",
		"NGINX_PORT":    "80",
	}

	actual := MergeYamlWithParameters([]byte(yaml), params)

	if expect != actual {
		t.Errorf("actula merged string is %v", actual)
	}

}
