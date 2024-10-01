package logging

import (
	"github.com/sirupsen/logrus"
)

func SetupLogger(formatter logrus.Formatter) *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(formatter)
	return log
}
