package logging

import "github.com/sirupsen/logrus"

func (log *LogEntry) Panic(data interface{}, message string) {
	if data != nil {
		log.logger.WithFields(logrus.Fields{
			"data": data,
		}).Panic(message)
		return
	}

	log.logger.Panic(message)
}
