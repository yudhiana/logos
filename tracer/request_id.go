package tracer

import "github.com/google/uuid"

func GenerateRequestID() string {
	return uuid.NewString()
}
