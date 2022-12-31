package controller

import (
	g "github.com/doug-martin/goqu/v9"
	"github.com/valyala/fasthttp"
	"merchant/contrib/helper"
	"merchant/contrib/validator"
	"merchant/model"
	"strings"
)

type LinkController struct{}

func (that *LinkController) List(ctx *fasthttp.RequestCtx) {

	page := ctx.QueryArgs().GetUintOrZero("page")
	pageSize := ctx.QueryArgs().GetUintOrZero("page_size")
	username := string(ctx.QueryArgs().Peek("username"))
	shortURL := string(ctx.QueryArgs().Peek("short_url"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 15
	}
	ex := g.Ex{}
	if username != "" {
		username = strings.ToLower(username)
		if !validator.CheckUName(username, 5, 14) {
			helper.Print(ctx, false, helper.UsernameErr)
			return
		}

		ex["username"] = username
	}

	if shortURL != "" {
		if strings.Contains(shortURL, "/") {
			helper.Print(ctx, false, helper.ParamErr)
			return
		}
		ex["short_url"] = shortURL
	}

	data, err := model.LinkList(uint(page), uint(pageSize), ex)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}

// SetNoAd 设置是否显示广告页
func (that *LinkController) SetNoAd(ctx *fasthttp.RequestCtx) {

	shortCode := string(ctx.PostArgs().Peek("short_url"))
	noAd := string(ctx.PostArgs().Peek("no_ad"))

	if noAd != "0" && noAd != "1" {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	err := model.LinkSetNoAd(shortCode, noAd)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

func (that *LinkController) Delete(ctx *fasthttp.RequestCtx) {

	username := string(ctx.QueryArgs().Peek("username"))
	id := string(ctx.QueryArgs().Peek("id"))
	if !helper.CtypeDigit(id) {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	mb, err := model.MemberInfo(username)
	if err != nil {
		helper.Print(ctx, false, helper.UsernameErr)
		return
	}

	err = model.LinkDelete(mb.UID, id)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}
