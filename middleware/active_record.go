package middleware

import "github.com/kataras/iris/v12"

func ActiveRecord(ctx iris.Context) {
	ctx.Record()
	ctx.Next()
}
