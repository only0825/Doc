package controller

import (
	"github.com/go-redis/redis/v8"
	"github.com/valyala/fasthttp"
	"merchant/contrib/helper"
	"merchant/model"
)

type ShortURLController struct{}

func (that *ShortURLController) Get(ctx *fasthttp.RequestCtx) {

	res, err := model.ShortURLGet()
	if err != nil {
		if err == redis.Nil {
			helper.Print(ctx, true, nil)
			return
		}

		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, res)
}

func (that *ShortURLController) Set(ctx *fasthttp.RequestCtx) {

	uri := string(ctx.PostArgs().Peek("uri"))
	err := model.ShortURLSet(uri)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, "success")
}
