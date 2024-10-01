package rmq

import (
	"fmt"
	"os"
)

func GetRabbitURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBIT_USER"),
		os.Getenv("RABBIT_PASSWORD"),
		os.Getenv("RABBIT_HOST"),
		os.Getenv("RABBIT_PORT"))
}
