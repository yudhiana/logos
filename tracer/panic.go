package tracer

import (
	"context"

	"github.com/mataharibiz/ward/observer"
)

func Recover(ctx context.Context) {
	if r := recover(); r != nil {
		newCtx := context.WithValue(ctx, PanicContextKey, r)
		observer.NewObserver("panic", Panic).Notify(newCtx)
	}
}
