package actions

import (
	"time"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/mataharibiz/sange"
)

func TracerRequest(data interface{}) {
	irisCtx, ok := data.(iris.Context)
	if ok {
		currentTime := time.Now().UTC()
		body, _ := irisCtx.GetBody()
		apiRequest := APIRequest{
			RequestID: uuid.NewString(),
			Timestamp: struct {
				Date time.Time `json:"date"`
			}{
				Date: currentTime,
			},
			Method:    irisCtx.Method(),
			URL:       irisCtx.Request().RequestURI,
			ClientIP:  irisCtx.RemoteAddr(),
			UserAgent: irisCtx.GetHeader("User-Agent"),
			RequestBody: struct {
				Body interface{} `json:"body"`
			}{
				Body: body,
			},
		}

		sangeEvent := sange.EventData{
			EventType: "api-requests",
			Data:      apiRequest,
		}
		sangeEvent.PublishDefault()
	}
}

type APIRequest struct {
	RequestID string `json:"request_id"`
	Timestamp struct {
		Date time.Time `json:"date"`
	} `json:"timestamp"`
	Method      string `json:"method"`
	URL         string `json:"url"`
	ClientIP    string `json:"client_ip"`
	UserAgent   string `json:"user_agent"`
	RequestBody struct {
		Body interface{} `json:"body"`
	} `json:"request_body"`
}
