package middleware

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

type middleware_t func(ctx *fasthttp.RequestCtx) error

var MiddlewareList = []middleware_t{
	//CorsMiddleware,
	//CheckTokenMiddleware,

}

func Use(next fasthttp.RequestHandler) fasthttp.RequestHandler {

	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {

		for _, cb := range MiddlewareList {
			if err := cb(ctx); err != nil {
				fmt.Fprint(ctx, err)
				return
			}
		}

		next(ctx)
	})
}
