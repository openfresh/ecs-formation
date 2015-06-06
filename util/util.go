package util

import (
	"github.com/awslabs/aws-sdk-go/aws/awsutil"
	"strings"
	"bufio"
	"io"
	"bytes"
)


func ConvertPointerString(values []string) []*string {

	merged := []*string{}
	for _, value := range values {
		merged = append(merged, &value)
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
			buffer.WriteString("\t")
		}
		buffer.Write(line)
	}

	return buffer.String()
}