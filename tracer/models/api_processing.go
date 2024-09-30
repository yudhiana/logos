package models

import "time"

type APIProcessing struct {
	RequestID    string    `json:"request_id"`
	LastUpdateAt time.Time `json:"last_update_at"`
	Span
}

type Span struct {
	SpanName     string `json:"span_name"`
	SpanID       string `json:"span_id"`
	ParentSpanID string `json:"parent_span_id"`
	Duration     int64  `json:"duration"`
}
