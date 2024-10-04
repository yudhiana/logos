package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/Kong/go-pdk"
	kLog "github.com/Kong/go-pdk/log"
)

type APIRequest struct {
	RequestID   string      `json:"request_id,omitempty"`
	Method      string      `json:"method,omitempty"`
	URL         string      `json:"url,omitempty"`
	ClientIP    string      `json:"client_ip,omitempty"`
	UserAgent   string      `json:"user_agent,omitempty"`
	AppOrigin   string      `json:"app_origin,omitempty"`
	Headers     http.Header `json:"headers,omitempty"`
	RequestBody interface{} `json:"request_body,omitempty"`
	Timestamp   time.Time   `json:"timestamp,omitempty"`
}

type APIResponse struct {
	RequestID    string      `json:"request_id,omitempty"`
	Status       int         `json:"status,omitempty"`
	ResponseBody interface{} `json:"response_body,omitempty"`
	Timestamp    time.Time   `json:"timestamp,omitempty"`
}

type Tracer struct {
	Request  *APIRequest
	Response *APIResponse
}

func NewTracer() *Tracer {
	return &Tracer{}
}
func (t *Tracer) BuildMessageRequest(kong *pdk.PDK) *Tracer {
	xRequestID, errGetRequestId := kong.Request.GetHeader("X-Request-Id")
	if errGetRequestId != nil {
		kong.Log.Err("Failed to get X-Request-Id", errGetRequestId)
	}

	httpMethod, errMethod := kong.Request.GetMethod()
	if errMethod != nil {
		kong.Log.Err("Failed to get http method", errMethod)
	}

	appOrigin, errOrigin := kong.Request.GetHeader("Dmp-Origin")
	if errOrigin != nil {
		kong.Log.Err("Failed to get app origin", errOrigin)
	}

	headers, errHeaders := kong.Request.GetHeaders(-1)
	if errHeaders != nil {
		kong.Log.Err("Failed to get headers", errHeaders)
	}

	clientIP, errIp := kong.Client.GetIp()
	if errIp != nil {
		kong.Log.Err("Failed to get client ip", errIp)

	}

	path, errPath := kong.Request.GetPath()
	if errPath != nil {
		kong.Log.Err("Failed to get path", errPath)
	}

	userAgent, errGetAgent := kong.Request.GetHeader("User-Agent")
	if errGetAgent != nil {
		kong.Log.Err("Failed to get user agent", errGetAgent)
	}

	metadata := &APIRequest{
		RequestID: xRequestID,
		Method:    httpMethod,
		URL:       path,
		ClientIP:  clientIP,
		UserAgent: userAgent,
		AppOrigin: appOrigin,
		Headers:   headers,
		Timestamp: time.Now().UTC(),
	}

	body, errBody := kong.Request.GetRawBody()
	if errBody != nil {
		kong.Log.Err("Failed to get body", errBody)
	}

	var mapBody map[string]interface{}
	if errUnmarshal := json.Unmarshal(body, &mapBody); errUnmarshal != nil {
		metadata.RequestBody = string(body)
	} else {
		metadata.RequestBody = mapBody
	}

	t.Request = metadata

	return t
}

func (t *Tracer) BuildMessageResponse(kong *pdk.PDK) *Tracer {
	xRequestID, errGetRequestId := kong.Request.GetHeader("X-Request-Id")
	if errGetRequestId != nil {
		kong.Log.Err("Failed to get X-Request-Id", errGetRequestId)
	}

	httpStatus, eStatus := kong.Response.GetStatus()
	if eStatus != nil {
		kong.Log.Err("Failed to get response status", eStatus)
	}

	metadata := &APIResponse{
		RequestID: xRequestID,
		Status:    httpStatus,
		Timestamp: time.Now().UTC(),
	}

	bodyResponse, errGetBody := kong.ServiceResponse.GetRawBody()
	if errGetBody != nil {
		kong.Log.Err("Failed to get response body", errGetBody)
	}

	if bodyResponse != nil {
		var mapBodyResponse map[string]interface{}
		if errUnmarshal := json.Unmarshal(bodyResponse, &mapBodyResponse); errUnmarshal != nil {
			metadata.ResponseBody = string(bodyResponse)
		} else {
			metadata.ResponseBody = mapBodyResponse
		}
	}

	t.Response = metadata

	return t
}

func (t *Tracer) Publish(eventType string) {
	event := &EventData{
		EventType: eventType,
	}

	switch eventType {
	case "api-request":
		event.Data = t.Request
	case "api-response":
		event.Data = t.Response
	}

	event.Publish(GetEnv("OBSERVER_EVENT", "dmp_observer"))
}

func (t *Tracer) Init(c Config, logger kLog.Log) *Tracer {
	t.configuredEnv(c, logger)
	return t
}

func (t *Tracer) configuredEnv(c Config, logger kLog.Log) *Tracer {
	t.setEnvString("RABBIT_HOST", &c.RABBIT_HOST, logger)
	t.setEnvString("RABBIT_PORT", &c.RABBIT_PORT, logger)
	t.setEnvString("RABBIT_USER", &c.RABBIT_USER, logger)
	t.setEnvString("RABBIT_PASSWORD", &c.RABBIT_PASSWORD, logger)
	return t
}

func (t *Tracer) setEnvString(env string, value *string, logger kLog.Log) {
	if value != nil {
		t.setEnv(env, *value, logger)
	}
}

func (t *Tracer) setEnv(env string, value string, logger kLog.Log) {
	err := os.Setenv(env, value)
	if err != nil {
		_ = logger.Err("Error setting environment ", env, " : ", err.Error())
		return
	}
}
