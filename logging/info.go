package logging

import "github.com/sirupsen/logrus"

func (log *LogEntry) Info(data interface{}, message string) {
	if data != nil {
		log.logger.WithFields(logrus.Fields{
			"data": data,
		}).Info(message)
		return
	}

	log.logger.Info(message)
}
