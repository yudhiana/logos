package logging

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kataras/pio"
)

// LogEntry represents the structure of the log entry in JSON format
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}

func LogAsJSON(printer *pio.Printer, level string, message string) {
	// Create a log entry object
	entry := LogEntry{
		Level:     level,
		Message:   message,
		Timestamp: time.Now().UTC(),
	}

	// Marshal the log entry to JSON format
	logData, err := json.Marshal(entry)
	if err != nil {
		printer.Println(pio.Rich(fmt.Sprintf("Failed to marshal log entry: %s", err.Error()), pio.Red))
		return
	}

	// Log the JSON formatted log entry
	printer.Println(string(logData))
	// log.Println("[LOG]", string(logData))
}
