package task

import (
	"testing"
)

func TestParseEntrypointWithSpace(t *testing.T) {
	input := `test1 test2 test3`
	result, err := parseEntrypoint(input)
	if err != nil {
		t.Error(err)
	}

	if len(result) != 3 {
		t.Error("len(result) expect = 3, but actual = %v", len(result))
	}

	if *result[0] != "test1" ||
		*result[1] != "test2" ||
		*result[2] != "test3" {
		t.Error("parse result is not expected")
	}
}

func TestParseEntrypointWithSpaceAndQuote(t *testing.T) {

	input := `nginx -g "daemon off;" -c /etc/nginx/nginx.conf`

	result, err := parseEntrypoint(input)
	if err != nil {
		t.Error(err)
	}

	if len(result) != 5 {
		t.Error("len(result) expect = 5, but actual = %v", len(result))
	}

	if *result[0] != "nginx" ||
		*result[1] != "-g" ||
		*result[2] != "daemon off;" ||
		*result[3] != "-c" ||
		*result[4] != "/etc/nginx/nginx.conf" {
		t.Error("parse result is not expected")
	}
}
