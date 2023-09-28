package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	Logger *logrus.Logger
}

func New(level int) *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.Level(level))
	logger.SetFormatter(&logrus.TextFormatter{})
	return logger
}
