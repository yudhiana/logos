package logging

import (
	"time"

	"github.com/sirupsen/logrus"
)

type LogEntry struct {
	logger *logrus.Logger
}

func NewLogEntry(formatter ...logrus.Formatter) *LogEntry {
	var logFormatter logrus.Formatter

	// set default to JSON formatter
	logFormatter = &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	}

	if len(formatter) > 0 {
		logFormatter = formatter[0]
	}

	return &LogEntry{
		logger: SetupLogger(logFormatter),
	}
}
