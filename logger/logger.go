package logger

import "github.com/Sirupsen/logrus"

var (
	Main *logrus.Logger
)

func init() {

	logger := logrus.New()
	logger.Level = logrus.InfoLevel

	Main = logger
}