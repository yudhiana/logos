package logging

import (
	"log/slog"
	"os"
)

type LogEntry struct {
	logger *slog.Logger
	app    map[string]interface{}
}

func NewLogger() *LogEntry {
	return &LogEntry{
		logger: ConfigureLogger(),
		app: map[string]interface{}{
			"source": os.Getenv("APP_NAME"),
		},
	}
}

func ConfigureLogger() *slog.Logger {
	logFormatter := slog.NewJSONHandler(os.Stdout, nil)
	return slog.New(logFormatter)
}

func (log *LogEntry) appLogger(args ...any) (results []any) {
	for k, v := range log.app {
		results = append(results, k, v)
	}

	countArgs := len(args)
	if countArgs > 0 {
		if countArgs%2 != 0 {
			log.logger.Error("args must be in pairs k/v", args...)
		} else {
			results = append(results, args...)
		}
	}

	return
}
