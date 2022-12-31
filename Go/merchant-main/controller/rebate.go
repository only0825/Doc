package controller

import (
	"fmt"
	"merchant/contrib/helper"
	"merchant/model"

	"github.com/valyala/fasthttp"
)

type RebateController struct{}

func (that *RebateController) Scale(ctx *fasthttp.RequestCtx) {

	vs := model.MemberRebateScale()
	s := fmt.Sprintf(
		`{"ty":"%s","zr":"%s","dj":"%s","qp":"%s","dz":"%s","cp":"%s","fc":"%s","by":"%s","cg_official_rebate":"%s","cg_high_rebate":"%s"}`,
		vs.TY.StringFixed(1),
		vs.ZR.StringFixed(1),
		vs.DJ.StringFixed(1),
		vs.QP.StringFixed(1),
		vs.DZ.StringFixed(1),
		vs.CP.StringFixed(1),
		vs.FC.StringFixed(1),
		vs.BY.StringFixed(1),
		vs.CGOfficialRebate.StringFixed(2),
		vs.CGHighRebate.StringFixed(2),
	)

	helper.PrintJson(ctx, true, s)
}

func (that *RebateController) EnableMod(ctx *fasthttp.RequestCtx) {

	enable := ctx.QueryArgs().GetBool("enable")
	err := model.MemberRebateEnableMod(enable)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

//RebatePersonal 会员管理-会员列表-详情-返水记录
func (that RebateController) RebatePersonal(ctx *fasthttp.RequestCtx) {

	page := ctx.QueryArgs().GetUintOrZero("page")
	pageSize := ctx.QueryArgs().GetUintOrZero("page_size")
	startTime := string(ctx.QueryArgs().Peek("start_time")) //开始时间
	endTime := string(ctx.QueryArgs().Peek("end_time"))     //结束时间
	username := string(ctx.QueryArgs().Peek("username"))    // 返水类型

	if page < 1 {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	if pageSize < 10 || pageSize > 100 {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	if len(username) <= 0 {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	data, err := model.RebatePersonalReport(username, startTime, endTime, page, pageSize)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}
