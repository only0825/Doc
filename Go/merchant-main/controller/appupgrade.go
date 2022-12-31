package controller

import (
	"github.com/valyala/fasthttp"
	"merchant/contrib/helper"
	"merchant/contrib/validator"
	"merchant/model"
)

type AppUpgradeController struct{}

var (
	apps = map[string]bool{
		"ios":     true,
		"android": true,
	}
)

// 修改
func (that *AppUpgradeController) Update(ctx *fasthttp.RequestCtx) {

	device := string(ctx.PostArgs().Peek("platform"))
	version := string(ctx.PostArgs().Peek("version"))
	content := string(ctx.PostArgs().Peek("content"))
	url := string(ctx.PostArgs().Peek("url"))
	isForce := string(ctx.PostArgs().Peek("is_force"))

	if _, ok := apps[device]; !ok {
		helper.Print(ctx, false, helper.DeviceErr)
		return
	}

	if len(content) > 255 {
		helper.Print(ctx, false, helper.ContentLengthErr)
		return
	}

	if !validator.CheckUrl(url) {
		helper.Print(ctx, false, helper.URLErr)
		return
	}

	if !validator.CheckIntScope(isForce, 0, 1) {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	admin, err := model.AdminToken(ctx)
	if err != nil {
		helper.Print(ctx, false, helper.AccessTokenExpires)
		return
	}

	err = model.AppUpgradeUpdate(device, version, content, url, isForce, admin)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// 版本管理列表
func (that *AppUpgradeController) List(ctx *fasthttp.RequestCtx) {

	data, err := model.AppUpgradeList()
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}
