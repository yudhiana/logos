package logos

func (log *LogEntry) Debug(message string, args ...any) {
	log.logger.Debug(message, log.appLogger(args...)...)
}
