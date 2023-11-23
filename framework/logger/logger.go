package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func Logger() *logrus.Logger {
	var logger = logrus.New()

	logger.Formatter = &logrus.JSONFormatter{}

	if os.Getenv("environment") == "development" {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	return logger
}
