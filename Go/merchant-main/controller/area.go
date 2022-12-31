package controller

import (
	"github.com/valyala/fasthttp"
	"merchant/contrib/helper"
	"merchant/model"
	"strings"
)

type AreaController struct{}

func (that *AreaController) View(ctx *fasthttp.RequestCtx) {
	str := string(ctx.PostArgs().Peek("ips"))

	if str == "" {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	ips := strings.SplitN(str, ",", 20)
	data, err := model.Area(ips)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.PrintJson(ctx, true, data)
}
