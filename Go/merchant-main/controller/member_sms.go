package controller

import (
	"merchant/contrib/helper"
	"merchant/contrib/validator"
	"merchant/model"
	"strings"

	"github.com/valyala/fasthttp"
)

type SmsRecordController struct{}

// List 验证码列表
func (that *SmsRecordController) List(ctx *fasthttp.RequestCtx) {

	page := ctx.QueryArgs().GetUintOrZero("page")
	pageSize := ctx.QueryArgs().GetUintOrZero("page_size")
	username := string(ctx.QueryArgs().Peek("username"))
	phone := string(ctx.QueryArgs().Peek("phone"))
	state := ctx.QueryArgs().GetUintOrZero("state")
	ty := ctx.QueryArgs().GetUintOrZero("ty")
	startTime := string(ctx.QueryArgs().Peek("start_time"))
	endTime := string(ctx.QueryArgs().Peek("end_time"))

	if page < 1 {
		page = 1
	}
	if pageSize < 10 {
		pageSize = 10
	}
	// 会员名校验
	if username != "" {
		username = strings.ToLower(username)
		if !validator.CheckUName(username, 5, 14) {
			helper.Print(ctx, false, helper.UsernameErr)
			return
		}
	}

	//// 手机号校验
	if phone != "" {
		if !validator.IsVietnamesePhone(phone) {
			helper.Print(ctx, false, helper.PhoneFMTErr)
			return
		}
	}

	if startTime == "" || endTime == "" {
		helper.Print(ctx, false, helper.DateTimeErr)
		return
	}

	data, err := model.SmsList(uint(page), uint(pageSize), startTime, endTime, username, phone, state, ty)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)

}
