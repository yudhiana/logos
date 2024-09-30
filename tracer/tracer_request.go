package tracer

import (
	"context"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/mataharibiz/sange"
	"github.com/yudhiana99/ward"
	"github.com/yudhiana99/ward/tracer/models"
)

func TracerIncomingRequest(data interface{}) {
	defer ward.Recover()

	currentTime := time.Now().UTC()
	requestID := GenerateRequestID()

	ctx := data.(context.Context)
	irisCtx, ok := ctx.Value(IrisContextKey).(iris.Context)
	if ok {
		if xRequestID := irisCtx.GetHeader("X-Request-Id"); xRequestID != "" {
			requestID = xRequestID
		}

		apiRequest := models.APIRequest{
			RequestID:    requestID,
			LastUpdateAt: currentTime,
			Method:       irisCtx.Method(),
			URL:          irisCtx.Request().RequestURI,
			ClientIP:     irisCtx.RemoteAddr(),
			UserAgent:    irisCtx.GetHeader("User-Agent"),
			AppOrigin:    irisCtx.GetHeader("Dmp-Origin"),
			Headers:      irisCtx.Request().Header,
		}

		if body, _ := irisCtx.GetBody(); body != nil {
			apiRequest.RequestBody = body
		}

		sangeEvent := sange.EventData{
			EventType:   "api-requests",
			PublishDate: &currentTime,
			Data:        apiRequest,
		}
		sangeEvent.Publish(sange.GetEnv("OBSERVER_EVENT", "dmp_observer"))
	}
}
