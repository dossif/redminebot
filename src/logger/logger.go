package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

func NewLogger(logLevel string) (*logrus.Logger, error) {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return &logrus.Logger{}, fmt.Errorf("failed to set log level %v: %v", logLevel, err)
	}
	logger := logrus.New()
	if level == logrus.DebugLevel {
		logger.SetFormatter(&logrus.TextFormatter{})
		logger.SetLevel(level)
	} else {
		logger.SetFormatter(&logrus.JSONFormatter{FieldMap: logrus.FieldMap{logrus.FieldKeyTime: "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message"}})
		logger.SetLevel(level)
	}
	logger.SetReportCaller(false)
	return logger, nil
}
