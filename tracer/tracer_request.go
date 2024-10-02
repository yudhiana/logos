package tracer

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/mataharibiz/sange"
	"github.com/mataharibiz/ward/rmq"
	"github.com/mataharibiz/ward/tracer/models"
)

func TracingRequest(data interface{}) {
	defer Recover(data.(context.Context))

	ctx := data.(context.Context)
	tracerCtx := ctx.Value(TracingRequestKey).(models.TracerCtx)
	irisCtx, ok := ctx.Value(IrisContextKey).(iris.Context)
	if ok {
		apiRequest := &models.APIRequest{
			RequestID: tracerCtx.XRequestID,
			Method:    irisCtx.Method(),
			URL:       irisCtx.Request().RequestURI,
			ClientIP:  irisCtx.RemoteAddr(),
			UserAgent: irisCtx.GetHeader("User-Agent"),
			AppOrigin: irisCtx.GetHeader("Dmp-Origin"),
			Headers:   irisCtx.Request().Header,
		}

		if body := GetRequestBody(irisCtx); body != nil {
			var request map[string]interface{}
			if e := json.Unmarshal(body, &request); e != nil {
				apiRequest.RequestBody = string(body)
			}
			apiRequest.RequestBody = request
		}

		irisCtx.Record()
		irisCtx.Next()

		// waiting to callback handlers
		if f, fok := irisCtx.IsRecording(); fok {
			apiRequest.Status = f.StatusCode()
			body := f.Body()
			var response map[string]interface{}
			if e := json.Unmarshal(body, &response); e != nil {
				apiRequest.ResponseBody = string(body)
			}
			apiRequest.ResponseBody = response

			endTime := time.Now().UTC()
			latency := endTime.Sub(tracerCtx.Timestamp)
			apiRequest.DurationAsMilis = latency.Milliseconds()
			f.FlushResponse()
			f.ResetBody()
		}
		apiRequest.TimeStamp = time.Now().UTC()

		eventResponse := rmq.EventData{
			EventType: "api-request",
			Data:      apiRequest,
		}
		eventResponse.Publish(sange.GetEnv("OBSERVER_EVENT", "dmp_observer"))
	}
}

func AuthenticateRequestId(ctx iris.Context) string {
	xRequestID := ctx.GetHeader("X-Request-Id")
	if xRequestID == "" {
		xRequestID = GenerateRequestID()
		ctx.Request().Header.Set("X-Request-Id", xRequestID)
	}

	return xRequestID
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
