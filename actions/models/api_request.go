package models

import "time"

type APIRequest struct {
	RequestID   string      `json:"request_id"`
	Timestamp   time.Time   `json:"timestamp"`
	Method      string      `json:"method"`
	URL         string      `json:"url"`
	ClientIP    string      `json:"client_ip"`
	UserAgent   string      `json:"user_agent"`
	RequestBody interface{} `json:"request_body"`
}
