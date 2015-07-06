package util

import (
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"strings"
	"bufio"
	"io"
	"bytes"
	"fmt"
)


func ConvertPointerString(values []string) []*string {

	merged := make([]*string, len(values))

	for i := 0; i < len(values); i++ {
		merged[i] = &values[i]
	}

	return merged
}

func StringValueWithIndent(value interface{}, indent int) string {
	sr := strings.NewReader(awsutil.StringValue(value))

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