package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/yudhiana99/ward/actions"
	"github.com/yudhiana99/ward/observer"
)

func TracingRequest() iris.Handler {
	return func(ctx iris.Context) {
		go func(irisCtx iris.Context) {
			observable := observer.NewObservable()
			observable.Register(observer.NewObserver("tracer request", actions.TracerRequest))
			observable.TriggerEvent("tracer request", ctx)
		}(ctx)
		ctx.Next()
	}
}
