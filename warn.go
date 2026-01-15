package main

func (log *LogEntry) Warn(message string, args ...any) {
	log.logger.Warn(message, log.appLogger(args...)...)
}
