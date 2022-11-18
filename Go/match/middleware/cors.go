package middleware

import (
	"github.com/valyala/fasthttp"
	//"fmt"
	"strings"
)

func handlePreflight(ctx *fasthttp.RequestCtx) {
	originHeader := string(ctx.Request.Header.Peek("Origin"))
	if len(originHeader) == 0 {
		return
	}
	method := string(ctx.Request.Header.Peek("Access-Control-Request-Method"))

	headers := []string{}
	if len(ctx.Request.Header.Peek("Access-Control-Request-Headers")) > 0 {
		headers = strings.Split(string(ctx.Request.Header.Peek("Access-Control-Request-Headers")), ",")
	}

	//ctx.Response.Header.Set("Access-Control-Allow-Origin", originHeader)
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", method)
	if len(headers) > 0 {
		ctx.Response.Header.Set("Access-Control-Allow-Headers", strings.Join(headers, ", "))
	}
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
}

func handleActual(ctx *fasthttp.RequestCtx) {
	originHeader := string(ctx.Request.Header.Peek("Origin"))
	if len(originHeader) == 0 {
		return
	}
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	//ctx.Response.Header.Set("Access-Control-Allow-Origin", originHeader)
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
}

func CorsMiddleware(ctx *fasthttp.RequestCtx) error {

	if string(ctx.Method()) == "OPTIONS" {

		handlePreflight(ctx)
		ctx.SetStatusCode(200)
	} else {
		handleActual(ctx)
	}

	return nil
}
