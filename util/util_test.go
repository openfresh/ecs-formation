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
		t.Errorf("expect length = %d, but actual length = %d", len(input), len(actual))
	}

	if &input[0] != actual[0] {
		t.Errorf("expect[0] address '%x', but actual[0] address = %x", &input[0], actual[0])
	}

	if &input[1] != actual[1] {
		t.Errorf("expect[1] address '%x', but actual[1] address = %x", &input[0], actual[0])
	}

}
