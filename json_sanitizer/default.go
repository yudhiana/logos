package jsonSanitizer

import "time"

const (
	DefaultMaxDepth         = 40
	DefaultRedactionMarker  = "[REDACTED]"
	DefaultSensitiveFields  = "password|token|secret|otp|auth|pin"
	DefaultSanitizeDuration = 1000 * time.Millisecond // it's mean 1 second
)
