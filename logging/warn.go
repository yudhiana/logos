package logging

import "github.com/sirupsen/logrus"

func (log *LogEntry) Warn(data interface{}, message string) {
	if data != nil {
		log.logger.WithFields(logrus.Fields{
			"data": data,
		}).Warn(message)
		return
	}

	log.logger.Warn(message)
}
