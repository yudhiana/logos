package tracer

import (
	"context"
	"fmt"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/mataharibiz/sange"
	"github.com/mataharibiz/ward/rmq"
	"github.com/mataharibiz/ward/tracer/models"
)

func Panic(data interface{}) {
	ctx := data.(context.Context)
	panicMessage := fmt.Sprintf("recover panic cause : %v", ctx.Value(PanicContextKey))
	irisCtx := ctx.Value(IrisContextKey).(iris.Context)
	tracerCtx := ctx.Value(TracingRequestKey).(models.TracerCtx)
	apiPanicRequest := &models.APIRequest{
		RequestID:       tracerCtx.XRequestID,
		DurationAsMilis: time.Since(tracerCtx.Timestamp).Milliseconds(),
		Status:          500,
		Method:          irisCtx.Method(),
		URL:             irisCtx.Request().RequestURI,
		ClientIP:        irisCtx.RemoteAddr(),
		UserAgent:       irisCtx.GetHeader("User-Agent"),
		AppOrigin:       irisCtx.GetHeader("Dmp-Origin"),
		Headers:         irisCtx.Request().Header,
		ResponseBody: map[string]interface{}{
			"panic": panicMessage,
		},
		TimeStamp: time.Now().UTC(),
	}

	panicEvent := rmq.EventData{
		EventType: "api-request",
		Data:      apiPanicRequest,
	}
	panicEvent.Publish(sange.GetEnv("OBSERVER_EVENT", "dmp_observer"))
}
