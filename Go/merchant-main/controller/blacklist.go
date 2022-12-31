package controller

import (
	"merchant/contrib/helper"
	"merchant/contrib/validator"
	"merchant/model"
	"net"
	"strconv"
	"strings"

	g "github.com/doug-martin/goqu/v9"
	"github.com/valyala/fasthttp"
)

type BlacklistController struct{}

func (that *BlacklistController) LogList(ctx *fasthttp.RequestCtx) {

	username := string(ctx.QueryArgs().Peek("username"))
	parentName := string(ctx.QueryArgs().Peek("parent_name"))
	topName := string(ctx.QueryArgs().Peek("top_name"))
	deviceNo := string(ctx.QueryArgs().Peek("device_no"))
	device := string(ctx.QueryArgs().Peek("device"))
	page := ctx.QueryArgs().GetUintOrZero("page")
	pageSize := ctx.QueryArgs().GetUintOrZero("page_size")
	startTime := string(ctx.QueryArgs().Peek("start_time"))
	endTime := string(ctx.QueryArgs().Peek("end_time"))

	if page < 1 {
		page = 1
	}
	if pageSize < 10 {
		pageSize = 10
	}

	ex := g.Ex{}
	if len(username) > 0 {
		username = strings.ToLower(username)
		if !validator.CheckUName(username, 5, 14) {
			helper.Print(ctx, false, helper.UsernameErr)
			return
		}

		ex["username"] = username
	}

	if topName != "" {
		if !validator.CheckUName(topName, 5, 14) {
			helper.Print(ctx, false, helper.AgentNameErr)
			return
		}

		ex["top_name"] = topName
	}

	if parentName != "" {
		parentName = strings.ToLower(parentName)
		if !validator.CheckUName(parentName, 5, 14) {
			helper.Print(ctx, false, helper.AgentNameErr)
			return
		}

		ex["parent_name"] = parentName
	}

	if len(deviceNo) > 0 {
		ex["device_no"] = deviceNo
	}

	ip := string(ctx.QueryArgs().Peek("ip"))
	if len(ip) > 0 {
		ex["ip"] = ip
	}

	if len(device) > 0 {
		i, err := strconv.Atoi(device)
		if err != nil {
			helper.Print(ctx, false, helper.DeviceTypeErr)
			return
		}

		if _, ok := model.DeviceMap[i]; !ok {
			helper.Print(ctx, false, helper.DeviceTypeErr)
			return
		}

		ex["device"] = device
	}

	data, err := model.MemberLoginLogList(startTime, endTime, page, pageSize, ex)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}

//AssociateList 条件下 会员列表信息
func (that *BlacklistController) AssociateList(ctx *fasthttp.RequestCtx) {

	page := ctx.QueryArgs().GetUintOrZero("page")
	pageSize := ctx.QueryArgs().GetUintOrZero("page_size")
	tys := string(ctx.QueryArgs().Peek("ty"))
	value := string(ctx.QueryArgs().Peek("value"))

	if !validator.CheckIntScope(tys, model.TyDevice, model.TyIP) {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	if !validator.CheckStringLength(value, 1, 60) {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	ty, _ := strconv.Atoi(tys)

	ex := g.Ex{}
	if ty == model.TyDevice {
		ex["device_no"] = value
	} else if ty == model.TyIP {
		ipAddr := net.ParseIP(value)
		if ipAddr == nil {
			helper.Print(ctx, false, helper.IPErr)
			return
		}
		ex["ip"] = ipAddr.String()
	}

	if page < 1 {
		page = 1
	}

	if pageSize < 10 {
		pageSize = 10
	}

	data, err := model.MemberAccessList(page, pageSize, ex)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}

func (that *BlacklistController) List(ctx *fasthttp.RequestCtx) {

	startTime := string(ctx.QueryArgs().Peek("start_time"))
	endTime := string(ctx.QueryArgs().Peek("end_time"))
	page := ctx.QueryArgs().GetUintOrZero("page")
	pageSize := ctx.QueryArgs().GetUintOrZero("page_size")
	ty := ctx.QueryArgs().GetUintOrZero("ty")

	if _, ok := model.BlackTy[ty]; !ok {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	ex := g.Ex{
		"ty": ty,
	}
	value := string(ctx.QueryArgs().Peek("value"))
	if len(value) > 0 {
		ex["value"] = value
	}
	data, err := model.BlacklistList(uint(page), uint(pageSize), startTime, endTime, ty, ex)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}

func (that *BlacklistController) Insert(ctx *fasthttp.RequestCtx) {

	ty := ctx.PostArgs().GetUintOrZero("ty")
	if _, ok := model.BlackTy[ty]; !ok {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	value := string(ctx.PostArgs().Peek("value"))
	switch ty {
	case model.TyBankcard:
		if !validator.CheckStringLength(value, 6, 20) || !validator.CheckStringDigit(value) {
			helper.Print(ctx, false, helper.ParamErr)
			return
		}
	case model.TyRebate, model.TyCGRebate, model.TyPromoteLink:
		value = strings.ToLower(value)
		if !validator.CheckUName(value, 5, 14) {
			helper.Print(ctx, false, helper.UsernameErr)
			return
		}

		if !model.MemberExist(value) {
			helper.Print(ctx, false, helper.UserNotExist)
			return
		}

	default:
		if !validator.CheckStringLength(value, 1, 60) {
			helper.Print(ctx, false, helper.ParamErr)
			return
		}
	}

	remark := string(ctx.PostArgs().Peek("remark"))
	if !validator.CheckStringLength(remark, 1, 1000) {
		helper.Print(ctx, false, helper.RemarkFMTErr)
		return
	}

	record := g.Record{
		"id":     helper.GenId(),
		"ty":     ty,
		"value":  value,
		"remark": remark,
	}
	err := model.BlacklistInsert(ctx, ty, value, record)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// 只能更新remark
func (that *BlacklistController) Update(ctx *fasthttp.RequestCtx) {

	id := string(ctx.PostArgs().Peek("id"))
	if !validator.CheckStringDigit(id) {
		helper.Print(ctx, false, helper.IDErr)
		return
	}

	remark := string(ctx.PostArgs().Peek("remark"))
	if !validator.CheckStringLength(remark, 1, 1000) {
		helper.Print(ctx, false, helper.RemarkFMTErr)
		return
	}

	remark = validator.FilterInjection(remark)
	err := model.BlacklistUpdate(id, remark)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

func (that *BlacklistController) Delete(ctx *fasthttp.RequestCtx) {

	id := string(ctx.QueryArgs().Peek("id"))
	if !validator.CheckStringDigit(id) {
		helper.Print(ctx, false, helper.IDErr)
		return
	}

	/// 从数据库 和 redis删除黑名单
	err := model.BlacklistDelete(id)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

func (that *BlacklistController) ClearPhone(ctx *fasthttp.RequestCtx) {

	phone := string(ctx.PostArgs().Peek("phone"))
	if !validator.IsVietnamesePhone(phone) {
		helper.Print(ctx, false, helper.PhoneFMTErr)
		return
	}

	err := model.BlacklistClearPhone(phone)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, "success")
}
