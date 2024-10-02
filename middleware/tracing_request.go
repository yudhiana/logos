package middleware

import (
	"context"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/mataharibiz/ward/observer"
	"github.com/mataharibiz/ward/tracer"
	"github.com/mataharibiz/ward/tracer/models"
)

func TraceIncomingRequest(ctx iris.Context) {
	tracer.AuthenticateRequestId(ctx)

	observable := observer.NewObservable()
	goCtx := context.WithValue(context.TODO(), tracer.IrisContextKey, ctx)
	tracerMetaData := models.TracerCtx{Timestamp: time.Now().UTC()}
	tracerCtx := context.WithValue(goCtx, tracer.TracingRequestKey, tracerMetaData)
	observable.Register(observer.NewObserver("tracer request", tracer.TracingRequest))
	observable.TriggerEvent("tracer request", tracerCtx)

	ctx.Next()
}
