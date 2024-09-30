package middleware

import "github.com/kataras/iris/v12"

func ActiveRecord() iris.Handler {
	return func(ctx iris.Context) {
		ctx.Record()
		ctx.Next()
	}
}
