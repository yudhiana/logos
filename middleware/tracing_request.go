package middleware

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/yudhiana99/ward/observer"
	"github.com/yudhiana99/ward/tracer"
)

func TraceIncomingRequest() iris.Handler {
	return func(ctx iris.Context) {
		go func(irisCtx iris.Context) {
			observable := observer.NewObservable()
			goCtx := context.WithValue(context.Background(), tracer.IrisContextKey, irisCtx)
			observable.Register(observer.NewObserver("tracer request", tracer.TracerIncomingRequest))
			observable.TriggerEvent("tracer request", goCtx)
		}(ctx)
		ctx.Next()
	}
}
