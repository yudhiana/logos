package main

import (
	"time"

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
	}

	// responseStatus, errStatus := kong.Response.GetStatus()
	// if errStatus != nil {
	// 	kong.Log.Err("Failed to get response status", errStatus)
	// 	return
	// }

	httpMethod, errMethod := kong.Request.GetMethod()
	if errMethod != nil {
		kong.Log.Err("Failed to get http method", errMethod)
		return
	}

	appOrigin, errOrigin := kong.Request.GetHeader("Dmp-Origin")
	if errOrigin != nil {
		kong.Log.Err("Failed to get app origin", errOrigin)
		return
	}

	headers, errHeaders := kong.Request.GetHeaders(-1)
	if errHeaders != nil {
		kong.Log.Err("Failed to get headers", errHeaders)
		return
	}

	// body, errBody := kong.Request.GetRawBody()
	// if errBody != nil {
	// 	kong.Log.Err("Failed to get body", errBody)
	// 	return
	// }

	// clientIP, _ := kong.Client.GetIp()
	path, _ := kong.Request.GetPath()
	userAgent, _ := kong.Request.GetHeader("User-Agent")

	metadata := &APIRequest{
		RequestID: xRequestID,
		// Status:    responseStatus,
		Method: httpMethod,
		URL:    path,
		// ClientIP:  clientIP,
		UserAgent: userAgent,
		AppOrigin: appOrigin,
		Headers:   headers,
		TimeStamp: time.Now().UTC(),
	}

	// var mapBody map[string]interface{}
	// if errUnmarshal := json.Unmarshal(body, &mapBody); errUnmarshal != nil {
	// 	metadata.RequestBody = string(body)
	// } else {
	// 	metadata.RequestBody = mapBody
	// }

	// bodyResponse, _ := kong.ServiceResponse.GetRawBody()
	// var mapBodyResponse map[string]interface{}
	// if errUnmarshal := json.Unmarshal(bodyResponse, &mapBodyResponse); errUnmarshal != nil {
	// 	metadata.ResponseBody = string(bodyResponse)
	// } else {
	// 	metadata.ResponseBody = mapBodyResponse
	// }
	kong.Log.Info("Captured", "metadata", metadata)
	NewTracer().Captured(metadata)
}
