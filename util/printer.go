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

func PrintlnCyan(a ...interface{}) {
	color.Cyan(fmt.Sprint(a...))
}

func PrintlnGreen(a ...interface{}) {
	color.Green(fmt.Sprint(a...))
}

func PrintlnYellow(a ...interface{}) {
	color.Yellow(fmt.Sprint(a...))
}

func Infoln(a ...interface{}) {
	if Output {
		logger.Main.Infoln(a...)
	}
}
