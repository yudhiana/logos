package logging

import "github.com/sirupsen/logrus"

func setupLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true, // Enable colors in the output
	})
	return logger
}

// func main() {
// 	logger := setupLogger()

// 	// Now, you can use the logger to print colored output
// 	logger.WithFields(logrus.Fields{
// 		"animal": "walrus",
// 		"size":   10,
// 	}).Info("A group of walrus emerges from the ocean")
// }
