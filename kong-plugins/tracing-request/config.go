package main

import (
	"fmt"
	"os"

	"github.com/google/uuid"
)

func GetEnv(key string, fallback string) string {
	env := os.Getenv(key)

	if len(env) == 0 {
		env = fallback
	}

	return env
}

func GetRabbitURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBIT_USER"),
		os.Getenv("RABBIT_PASSWORD"),
		os.Getenv("RABBIT_HOST"),
		os.Getenv("RABBIT_PORT"))
}

func GenerateRequestID() string {
	return uuid.NewString()
}
