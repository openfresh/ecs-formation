package util

import (
	"fmt"
	"github.com/str1ngs/ansi/color"
)

var (
	output = true
)

func Println(a ...interface{}) (int, error) {
	if output {
		return fmt.Println(a...)
	} else {
		return 0, nil
	}
}

func PrintlnCyan(a ...interface{}) (int, error) {
	return Println(color.Cyan(fmt.Sprint(a...)))
}

func PrintlnGreen(a ...interface{}) (int, error) {
	return Println(color.Green(fmt.Sprint(a...)))
}

func PrintlnYellow(a ...interface{}) (int, error) {
	return Println(color.Yellow(fmt.Sprint(a...)))
}