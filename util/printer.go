package util

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/stormcat24/ecs-formation/logger"
)

var (
	Output = true
)

func Println(a ...interface{}) (int, error) {
	if Output {
		return fmt.Println(a...)
	} else {
		return 0, nil
	}
}

func Print(a ...interface{}) (int, error) {
	if Output {
		return fmt.Print(a...)
	} else {
		return 0, nil
	}
}

func PrintlnCyan(format string, a ...interface{}) {
	if Output {
		color.Cyan(format, a...)
	}
}

func PrintlnGreen(format string, a ...interface{}) {
	if Output {
		color.Green(format, a...)
	}
}

func PrintlnYellow(format string, a ...interface{}) {
	if Output {
		color.Yellow(format, a...)
	}
}

func Infoln(a ...interface{}) {
	if Output {
		logger.Main.Infoln(a...)
	}
}
