package main

import (
	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
)

var (
	Version  = "1.0.0"
	Priority = 1
)

func main() {
	server.StartServer(New, Version, Priority)
}

type Config struct {
	RABBIT_HOST     string `json:"rabbit_host"`
	RABBIT_PORT     string `json:"rabbit_port"`
	RABBIT_USER     string `json:"rabbit_user"`
	RABBIT_PASSWORD string `json:"rabbit_password"`
}

func New() interface{} {
	return &Config{}
}

func (conf Config) Access(kong *pdk.PDK) {
	xRequestID, errGet := kong.Request.GetHeader("X-Request-Id")
	if errGet != nil {
		kong.Log.Err("Failed to get X-Request-Id", errGet)
		return
	}

	if xRequestID == "" {
		xRequestID = GenerateRequestID()
		kong.ServiceRequest.SetHeader("X-Request-Id", xRequestID)
		kong.Response.SetHeader("X-Request-Id", xRequestID)
	}

	NewTracer().Init(conf, kong.Log).BuildMessageRequest(kong).Publish("api-request")
}

func (conf Config) Response(kong *pdk.PDK) {
	NewTracer().Init(conf, kong.Log).BuildMessageResponse(kong).Publish("api-response")
}
