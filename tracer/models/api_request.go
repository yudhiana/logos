package models

import (
	"net/http"
	"time"
)

type APIRequest struct {
	RequestID       string      `json:"request_id,omitempty"`
	Status          int         `json:"status,omitempty"`
	Method          string      `json:"method,omitempty"`
	URL             string      `json:"url,omitempty"`
	ClientIP        string      `json:"client_ip,omitempty"`
	UserAgent       string      `json:"user_agent,omitempty"`
	AppOrigin       string      `json:"app_origin,omitempty"`
	Headers         http.Header `json:"headers,omitempty"`
	RequestBody     interface{} `json:"request_body,omitempty"`
	ResponseBody    interface{} `json:"response_body,omitempty"`
	DurationAsMilis int64       `json:"duration_as_milis,omitempty"`
	TimeStamp       time.Time   `json:"timestamp,omitempty"`
}

type TracerCtx struct {
	Timestamp  time.Time
	XRequestID string
}
