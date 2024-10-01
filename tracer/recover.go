package tracer

import (
	"context"
	"fmt"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/mataharibiz/sange"
	"github.com/mataharibiz/ward/tracer/models"
)

func Panic(data interface{}) {
	ctx := data.(context.Context)
	panicMessage := fmt.Sprintf("recover panic cause : %v", ctx.Value(PanicContextKey))
	irisCtx := ctx.Value(IrisContextKey).(iris.Context)

	apiPanicRequest := &models.APIRequest{
		RequestID: irisCtx.GetHeader("X-Request-Id"),
		Method:    irisCtx.Method(),
		URL:       irisCtx.Request().RequestURI,
		ClientIP:  irisCtx.RemoteAddr(),
		UserAgent: irisCtx.GetHeader("User-Agent"),
		AppOrigin: irisCtx.GetHeader("Dmp-Origin"),
		Headers:   irisCtx.Request().Header,
		ResponseBody: map[string]interface{}{
			"panic": panicMessage,
		},
		TimeStamp: time.Now().UTC(),
	}

	panicEvent := sange.EventData{
		EventType: "panic-observer",
		Data:      apiPanicRequest,
	}
	panicEvent.Publish(sange.GetEnv("OBSERVER_EVENT", "dmp_observer"))
}
