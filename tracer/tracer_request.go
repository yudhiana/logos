package tracer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/mataharibiz/sange"
	"github.com/yudhiana99/ward"
	"github.com/yudhiana99/ward/tracer/models"
)

func TracerIncomingRequest(data interface{}) {
	defer ward.Recover()
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

		if body, _ := irisCtx.GetBody(); body != nil {
			apiRequest.RequestBody = string(body)
		}

		sangeEvent := sange.EventData{
			EventType:   "api-requests",
			PublishDate: &currentTime,
			Data:        apiRequest,
		}
		sangeEvent.Publish(sange.GetEnv("OBSERVER_EVENT", "dmp_observer"))
	}
}

func AuthenticateRequestId(ctx iris.Context) {
	requestID := GenerateRequestID()
	if xRequestID := ctx.GetHeader("X-Request-Id"); xRequestID == "" {
		ctx.Request().Header.Set("X-Request-Id", requestID)
	}
}

func TracerOutgoingRequest(data interface{}) {
	defer ward.Recover()
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

		sangeEvent := sange.EventData{
			EventType:   "api-responses",
			PublishDate: &currentTime,
			Data:        apiRequest,
		}
		sangeEvent.Publish(sange.GetEnv("OBSERVER_EVENT", "dmp_observer"))
	}
}
