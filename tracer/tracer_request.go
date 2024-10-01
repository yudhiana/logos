package tracer

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/mataharibiz/sange"
	"github.com/yudhiana99/ward/observer"
	"github.com/yudhiana99/ward/rmq"
	"github.com/yudhiana99/ward/tracer/models"
)

func TracerIncomingRequest(data interface{}) {
	defer observer.Recover()
	currentTime := time.Now().UTC()

	ctx := data.(context.Context)
	irisCtx, ok := ctx.Value(IrisContextKey).(iris.Context)
	requestID := irisCtx.GetHeader("X-Request-Id")
	if ok {
		apiRequest := models.APIRequest{
			RequestID:    requestID,
			LastUpdateAt: currentTime,
			Method:       irisCtx.Method(),
			Type:         "request",
			Status:       irisCtx.GetStatusCode(),
			URL:          irisCtx.Request().RequestURI,
			ClientIP:     irisCtx.RemoteAddr(),
			UserAgent:    irisCtx.GetHeader("User-Agent"),
			AppOrigin:    irisCtx.GetHeader("Dmp-Origin"),
			Headers:      irisCtx.Request().Header,
		}

		if body := GetRequestBody(irisCtx); body != nil {
			var request map[string]interface{}
			_ = json.Unmarshal(body, &request)
			apiRequest.RequestBody = request
		}

		event := rmq.EventData{
			EventType:   "api-requests",
			PublishDate: &currentTime,
			Data:        apiRequest,
		}
		event.Publish(sange.GetEnv("OBSERVER_EVENT", "dmp_observer"))
	}
}

func AuthenticateRequestId(ctx iris.Context) {
	requestID := GenerateRequestID()
	if xRequestID := ctx.GetHeader("X-Request-Id"); xRequestID == "" {
		ctx.Request().Header.Set("X-Request-Id", requestID)
	}
}

func TracerOutgoingRequest(data interface{}) {
	defer observer.Recover()
	currentTime := time.Now().UTC()

	ctx := data.(context.Context)
	irisCtx, ok := ctx.Value(IrisContextKey).(iris.Context)
	requestID := irisCtx.GetHeader("X-Request-Id")
	if ok {
		apiRequest := models.APIResponse{
			RequestID:    requestID,
			Status:       irisCtx.GetStatusCode(),
			LastUpdateAt: currentTime,
			Method:       irisCtx.Method(),
			Type:         "response",
			URL:          irisCtx.Request().RequestURI,
			ClientIP:     irisCtx.RemoteAddr(),
			UserAgent:    irisCtx.GetHeader("User-Agent"),
			AppOrigin:    irisCtx.GetHeader("Dmp-Origin"),
			Headers:      irisCtx.Request().Header,
		}

		if f, fok := irisCtx.IsRecording(); fok {
			body := f.Body()
			var response map[string]interface{}
			_ = json.Unmarshal(body, &response)
			if body != nil {
				apiRequest.ResponseBody = response
			}
			f.FlushResponse()
			f.ResetBody()
		}

		event := rmq.EventData{
			EventType:   "api-responses",
			PublishDate: &currentTime,
			Data:        apiRequest,
		}
		event.Publish(sange.GetEnv("OBSERVER_EVENT", "dmp_observer"))
	}
}

func GetRequestBody(ctx iris.Context) (bodyRequest []byte) {
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return
	}
	defer ctx.Request().Body.Close()

	bodyRequest = body
	ctx.Request().Body = io.NopCloser(bytes.NewBuffer(bodyRequest))

	return
}
