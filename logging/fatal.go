package logging

import "github.com/sirupsen/logrus"

func (log *LogEntry) Fatal(data interface{}, message string) {
	if data != nil {
		log.logger.WithFields(logrus.Fields{
			"data": data,
		}).Fatal(message)
		return
	}

	log.logger.Fatal(message)
}
