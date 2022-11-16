package controller

import (
	"github.com/valyala/fasthttp"
	"goblog/common/helper"
	"goblog/model"
)

type UserController struct{}

// Reg 用户注册
func (that *UserController) Reg(ctx *fasthttp.RequestCtx) {

	username := string(ctx.PostArgs().Peek("username"))
	password := string(ctx.PostArgs().Peek("password"))

	if username == "" || password == "" {
		helper.Print(ctx, false, helper.ParamErr)
	}

	isSet := model.UserExist(username)
	if isSet != 0 {
		helper.Print(ctx, false, helper.UserExist)
	}

	uid, err := model.UserReg(username, password)
	if err != nil {
		helper.Print(ctx, false, err.Error())
	}

	helper.Print(ctx, true, uid)
}
