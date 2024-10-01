package models

import (
	"net/http"
	"time"
)

type APIRequest struct {
	RequestID    string      `json:"request_id"`
	Status       int         `json:"status"`
	Method       string      `json:"method"`
	Type         string      `json:"type"`
	URL          string      `json:"url"`
	ClientIP     string      `json:"client_ip"`
	UserAgent    string      `json:"user_agent"`
	AppOrigin    string      `json:"app_origin"`
	Headers      http.Header `json:"headers"`
	RequestBody  interface{} `json:"request_body"`
	LastUpdateAt time.Time   `json:"last_update_at"`
}
