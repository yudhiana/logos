package models

import (
	"net/http"
	"time"
)

type APIRequest struct {
	RequestID    string      `json:"request_id"`
	Status       int         `json:"status"`
	Method       string      `json:"method"`
	URL          string      `json:"url"`
	ClientIP     string      `json:"client_ip"`
	UserAgent    string      `json:"user_agent"`
	AppOrigin    string      `json:"app_origin"`
	Headers      http.Header `json:"headers"`
	RequestBody  interface{} `json:"request_body"`
	ResponseBody interface{} `json:"response_body"`
	Duration     int64       `json:"duration"`
	LastUpdateAt time.Time   `json:"last_update_at"`
}
