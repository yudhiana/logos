package models

import (
	"net/http"
	"time"
)

type APIResponse struct {
	RequestID    string      `json:"request_id"`
	Method       string      `json:"method"`
	URL          string      `json:"url"`
	ClientIP     string      `json:"client_ip"`
	UserAgent    string      `json:"user_agent"`
	AppOrigin    string      `json:"app_origin"`
	Headers      http.Header `json:"headers"`
	ResponseBody interface{} `json:"response_body"`
	LastUpdateAt time.Time   `json:"last_update_at"`
}
