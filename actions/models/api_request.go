package models

import (
	"net/http"
	"time"
)

type APIRequest struct {
	RequestID    string      `json:"request_id"`
	Method       string      `json:"method"`
	URL          string      `json:"url"`
	ClientIP     string      `json:"client_ip"`
	UserAgent    string      `json:"user_agent"`
	Headers      http.Header `json:"headers"`
	RequestBody  interface{} `json:"request_body"`
	LastUpdateAt time.Time   `json:"last_update_at"`
}
