package main

import (
	"net/http"
	"os"
	"time"

	kLog "github.com/Kong/go-pdk/log"
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

type Tracer struct {
}

func NewTracer() *Tracer {
	return &Tracer{}
}

func (t *Tracer) Captured(msg *APIRequest) {
	event := EventData{
		EventType: "api-request",
		Data:      msg,
	}

	event.Publish(GetEnv("OBSERVER_EVENT", "dmp_observer"))
}

func (t *Tracer) Init(c Config, logger kLog.Log) *Tracer {
	t.setEnvString("RABBIT_HOST", &c.RABBIT_HOST, logger)
	t.setEnvString("RABBIT_PORT", &c.RABBIT_PORT, logger)
	t.setEnvString("RABBIT_USER", &c.RABBIT_USER, logger)
	t.setEnvString("RABBIT_PASS", &c.RABBIT_PASS, logger)
	return t
}

func (t *Tracer) setEnvString(env string, value *string, logger kLog.Log) {
	if value != nil {
		t.setEnv(env, *value, logger)
	}
}

func (t *Tracer) setEnv(env string, value string, logger kLog.Log) {
	_ = logger.Info("Setting ", env, " to ", value)
	err := os.Setenv(env, value)
	if err != nil {
		_ = logger.Err("Error setting environment ", env, " : ", err.Error())
		panic(err)
	}
}
