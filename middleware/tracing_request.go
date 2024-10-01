package middleware

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/mataharibiz/ward/observer"
	"github.com/mataharibiz/ward/tracer"
)

func TraceIncomingRequest(ctx iris.Context) {
	tracer.AuthenticateRequestId(ctx)

	observable := observer.NewObservable()
	goCtx := context.WithValue(context.Background(), tracer.IrisContextKey, ctx)
	observable.Register(observer.NewObserver("tracer request", tracer.TracingRequest))
	observable.TriggerEvent("tracer request", goCtx)

	ctx.Next()
}
