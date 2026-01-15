package logos

func (log *LogEntry) Error(message string, args ...any) {
	log.logger.Error(message, log.appLogger(args...)...)
}
