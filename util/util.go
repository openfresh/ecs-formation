package util

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"io"
	"regexp"
	"strings"
)

var (
	keyValuePattern = regexp.MustCompile(`^\s*(.+)\s*=\s*(.+)\s*$`)
	variablePattern = regexp.MustCompile(`\$\{([\w_-]+)\}`)
)

func ConvertPointerString(values []string) []*string {

	merged := make([]*string, len(values))

	for i := 0; i < len(values); i++ {
		merged[i] = &values[i]
	}

	return merged
}

func StringValueWithIndent(value interface{}, indent int) string {
	sr := strings.NewReader(awsutil.Prettify(value))

	reader := bufio.NewReader(sr)

	var buffer bytes.Buffer
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		for i := 0; i < indent; i++ {
			buffer.WriteString("  ")
		}
		buffer.Write(line)
		buffer.WriteString("\n")
	}

	return buffer.String()
}

func ToUpperCamel(s string) string {

	if len(s) == 0 {
		return ""
	}

	prefix := s[0:1]
	suffix := s[1:len(s)]
	return fmt.Sprintf("%s%s", strings.ToUpper(prefix), suffix)
}

func ParseKeyValues(slice []string) map[string]string {

	params := slice
	values := map[string]string{}
	for _, p := range params {
		tokens := keyValuePattern.FindStringSubmatch(p)
		if len(tokens) == 3 {
			key := strings.Trim(tokens[1], " ")
			values[key] = strings.Trim(tokens[2], " ")
		}
	}

	return values
}

func MergeYamlWithParameters(content []byte, params map[string]string) string {

	s := string(content)
	matched := variablePattern.FindAllStringSubmatch(s, -1)

	for _, tokens := range matched {
		key := tokens[1]

		if value, ok := params[key]; ok {
			s = strings.Replace(s, fmt.Sprintf("${%v}", key), value, -1)
		} else {
			s = strings.Replace(s, tokens[0], "", -1)
		}
	}
	return s
}
