package tracer

import (
	"context"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/mataharibiz/ward"
	"github.com/mataharibiz/ward/rmq"
	"github.com/mataharibiz/ward/tracer/models"
)

func AddSpans(ctx context.Context, data map[string]interface{}) {
	currentTime := time.Now().UTC()
	requestID := GenerateRequestID()

	irisCtx, ok := ctx.Value(IrisContextKey).(iris.Context)
	if ok {
		if xRequestID := irisCtx.GetHeader("X-Request-Id"); xRequestID != "" {
			requestID = xRequestID
		}
	}
	var span models.Span
	_ = ward.ParsePayloadData(data, &span)

	eventProcessing := models.APIProcessing{
		RequestID:    requestID,
		LastUpdateAt: currentTime,
		Span:         span,
	}

	spanEvent := rmq.EventData{
		EventType:   "spans",
		PublishDate: &currentTime,
		Data:        eventProcessing,
	}

	spanEvent.Publish(ward.GetEnv("OBSERVER_EVENT", "dmp_observer"))

}
