package logging

import "github.com/sirupsen/logrus"

func (log *LogEntry) Error(data interface{}, message string) {
	if data != nil {
		log.logger.WithFields(logrus.Fields{
			"data": data,
		}).Error(message)
		return
	}

	log.logger.Error(message)
}
