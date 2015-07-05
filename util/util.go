package util

import (
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"strings"
	"bufio"
	"io"
	"bytes"
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
