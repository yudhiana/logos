package main

func (log *LogEntry) Info(message string, args ...any) {
	log.logger.Info(message, log.appLogger(args...)...)
}
