package actions

import (
	"time"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/mataharibiz/sange"
	"github.com/yudhiana99/ward/actions/models"
)

func TracerRequest(data interface{}) {
	currentTime := time.Now().UTC()
	irisCtx, ok := data.(iris.Context)
	if ok {
		requestID := uuid.NewString()

		if xRequestID := irisCtx.GetHeader("X-Request-Id"); xRequestID != "" {
			requestID = xRequestID
		}

		apiRequest := models.APIRequest{
			RequestID: requestID,
			Timestamp: currentTime,
			Method:    irisCtx.Method(),
			URL:       irisCtx.Request().RequestURI,
			ClientIP:  irisCtx.RemoteAddr(),
			UserAgent: irisCtx.GetHeader("User-Agent"),
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
