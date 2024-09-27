package actions

import (
	"time"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/mataharibiz/sange"
	"github.com/yudhiana99/ward/actions/models"
)

func TracerRequest(data interface{}) {
	irisCtx, ok := data.(iris.Context)
	if ok {
		currentTime := time.Now().UTC()
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
			EventType: "api-requests",
			Data:      apiRequest,
		}
		sangeEvent.PublishDefault()
	}
}
