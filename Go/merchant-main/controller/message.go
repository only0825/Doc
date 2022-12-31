package controller

import (
	g "github.com/doug-martin/goqu/v9"
	"github.com/valyala/fasthttp"
	"merchant/contrib/helper"
	"merchant/contrib/validator"
	"merchant/model"
	"strings"
)

type MessageController struct{}

// 站内信新增
func (that *MessageController) Insert(ctx *fasthttp.RequestCtx) {

	title := string(ctx.PostArgs().Peek("title"))        //标题
	content := string(ctx.PostArgs().Peek("content"))    //内容
	ty := ctx.PostArgs().GetUintOrZero("ty")             //1站内消息 2活动消息
	isTop := ctx.PostArgs().GetUintOrZero("is_top")      //0不置顶 1置顶
	isPush := ctx.PostArgs().GetUintOrZero("is_push")    //0不推送 1推送
	sendName := string(ctx.PostArgs().Peek("send_name")) //发送人名
	sendAt := string(ctx.PostArgs().Peek("send_at"))     //发送时间
	isVip := ctx.PostArgs().GetUintOrZero("is_vip")      //是否vip站内信 1 vip站内信 2 直属下级 3 全部下级
	level := string(ctx.PostArgs().Peek("level"))        //vip等级 0-10,多个逗号分割
	names := string(ctx.PostArgs().Peek("names"))        //会员名，多个用逗号分割

	if len(title) < 1 || len(title) > 255 ||
		//len(subTitle) < 1 || len(subTitle) > 255 ||
		len(content) == 0 ||
		len(sendName) < 1 || len(sendName) > 100 {

		helper.Print(ctx, false, helper.ContentLengthErr)
		return
	}

	if ty != 1 && ty != 2 || isVip > 3 {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	switch isVip {
	case 0:
		if names == "" {
			helper.Print(ctx, false, helper.UsernameErr)
			return
		}

		names = strings.ToLower(names)
		usernames := strings.Split(names, ",")
		for _, v := range usernames {
			if !validator.CheckUName(v, 5, 14) {
				helper.Print(ctx, false, helper.UsernameErr)
				return
			}
		}
	case 1:
		if level == "" {
			helper.Print(ctx, false, helper.ParamErr)
			return
		}

		lv := map[string]bool{
			"1":  true,
			"2":  true,
			"3":  true,
			"4":  true,
			"5":  true,
			"6":  true,
			"7":  true,
			"8":  true,
			"9":  true,
			"10": true,
		}
		for _, v := range strings.Split(level, ",") {
			if _, ok := lv[v]; !ok {
				helper.Print(ctx, false, helper.ParamErr)
				return
			}
		}
	case 2, 3:
		names = strings.ToLower(names)
		if !validator.CheckUName(names, 5, 14) {
			helper.Print(ctx, false, helper.UsernameErr)
			return
		}
	default:
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	admin, err := model.AdminToken(ctx)
	if err != nil {
		helper.Print(ctx, false, helper.AccessTokenExpires)
		return
	}

	record := g.Record{
		"id":          helper.GenId(),
		"title":       title,             //标题
		"content":     content,           //内容
		"is_top":      isTop,             //0不置顶 1置顶
		"is_push":     isPush,            //0不推送 1推送
		"is_vip":      isVip,             //是否是vip
		"level":       level,             //会员等级
		"ty":          ty,                //站内信类型
		"usernames":   names,             //会员名
		"state":       1,                 //1审核中 2审核通过 3审核拒绝 4已删除
		"send_state":  1,                 //1未发送 2已发送
		"send_name":   sendName,          //发送人名
		"apply_at":    ctx.Time().Unix(), //创建时间
		"apply_uid":   admin["id"],       //创建人uid
		"apply_name":  admin["name"],     //创建人名
		"review_at":   0,                 //审核时间
		"review_uid":  0,                 //审核人uid
		"review_name": "",                //审核人名
	}
	err = model.MessageInsert(record, sendAt)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// 站内信列表
func (that *MessageController) List(ctx *fasthttp.RequestCtx) {

	page := ctx.PostArgs().GetUintOrZero("page")
	pageSize := ctx.PostArgs().GetUintOrZero("page_size")
	flag := ctx.PostArgs().GetUintOrZero("flag")                        //1审核列表 2历史记录
	title := string(ctx.PostArgs().Peek("title"))                       //标题
	sendName := string(ctx.PostArgs().Peek("send_name"))                //发送人
	isVip := string(ctx.PostArgs().Peek("is_vip"))                      //是否vip站内信 1 vip站内信
	isPush := string(ctx.PostArgs().Peek("is_push"))                    //0不推送 1推送
	ty := ctx.PostArgs().GetUintOrZero("ty")                            //1站内消息 2活动消息
	sendStartTime := string(ctx.PostArgs().Peek("send_start_time"))     //发送开始时间
	sendEndTime := string(ctx.PostArgs().Peek("send_end_time"))         //发送结束时间
	startTime := string(ctx.PostArgs().Peek("start_time"))              //申请开始时间
	endTime := string(ctx.PostArgs().Peek("end_time"))                  //申请结束时间
	reviewStartTime := string(ctx.PostArgs().Peek("review_start_time")) //审核开始时间
	reviewEndTime := string(ctx.PostArgs().Peek("review_end_time"))     //审核结束时间

	ex := g.Ex{}
	if page == 0 {
		page = 1
	}
	if pageSize < 10 {
		pageSize = 10
	}

	if flag != 1 && flag != 2 {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	if flag == 1 {
		ex["state"] = 1
	} else {
		ex["state"] = []int{2, 3, 4}
	}

	if title != "" {
		ex["title"] = title
	}

	if sendName != "" {
		ex["send_name"] = sendName
	}

	if ty > 0 {
		if ty != 1 && ty != 2 {
			helper.Print(ctx, false, helper.ParamErr)
			return
		}

		ex["ty"] = ty
	}

	if isVip != "" {
		if isVip != "0" && isVip != "1" {
			helper.Print(ctx, false, helper.ParamErr)
			return
		}

		ex["is_vip"] = isVip
	}

	if isPush != "" {
		if isPush != "0" && isPush != "1" {
			helper.Print(ctx, false, helper.ParamErr)
			return
		}

		ex["is_push"] = isPush
	}

	data, err := model.MessageList(page, pageSize, sendStartTime, sendEndTime, startTime, endTime, reviewStartTime, reviewEndTime, ex)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}

// 站内信编辑
func (that *MessageController) Update(ctx *fasthttp.RequestCtx) {

	id := string(ctx.PostArgs().Peek("id"))
	title := string(ctx.PostArgs().Peek("title"))        //标题
	content := string(ctx.PostArgs().Peek("content"))    //内容
	isTop := string(ctx.PostArgs().Peek("is_top"))       //0不置顶 1置顶
	sendName := string(ctx.PostArgs().Peek("send_name")) //发送人名
	sendAt := string(ctx.PostArgs().Peek("send_at"))     //发送时间

	record := g.Record{}
	if !validator.CtypeDigit(id) {
		helper.Print(ctx, false, helper.IDErr)
		return
	}

	if title != "" {
		record["title"] = title
	}

	if content != "" {
		record["content"] = content
	}

	if isTop != "" {
		if isTop != "0" && isTop != "1" {
			helper.Print(ctx, false, helper.ParamErr)
			return
		}

		record["is_top"] = isTop
	}

	if sendName != "" {
		record["send_name"] = sendName
	}
	err := model.MessageUpdate(id, sendAt, record)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// 站内信编辑
func (that *MessageController) Review(ctx *fasthttp.RequestCtx) {

	id := string(ctx.PostArgs().Peek("id"))
	state := ctx.PostArgs().GetUintOrZero("state")
	flag := ctx.PostArgs().GetUintOrZero("flag") // 1 定时发送 2立即发送
	if !validator.CtypeDigit(id) {
		helper.Print(ctx, false, helper.IDErr)
		return
	}

	flags := map[int]bool{
		1: true,
		2: true,
	}
	if _, ok := flags[flag]; !ok {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	states := map[int]bool{
		2: true,
		3: true,
	}
	if _, ok := states[state]; !ok {
		helper.Print(ctx, false, helper.StateParamErr)
		return
	}

	admin, err := model.AdminToken(ctx)
	if err != nil {
		helper.Print(ctx, false, helper.AccessTokenExpires)
		return
	}

	err = model.MessageReview(id, state, flag, admin)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// 已发站内信详情
func (that *MessageController) Detail(ctx *fasthttp.RequestCtx) {

	page := ctx.QueryArgs().GetUintOrZero("page")
	pageSize := ctx.QueryArgs().GetUintOrZero("page_size")
	id := string(ctx.QueryArgs().Peek("id"))

	if !validator.CtypeDigit(id) {
		helper.Print(ctx, false, helper.IDErr)
		return
	}

	if page == 0 {
		page = 1
	}
	if pageSize < 10 {
		pageSize = 10
	}
	data, err := model.MessageDetail(id, page, pageSize)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}

// 已发系统站内信站内信列表
func (that *MessageController) System(ctx *fasthttp.RequestCtx) {

	page := ctx.PostArgs().GetUintOrZero("page")
	pageSize := ctx.PostArgs().GetUintOrZero("page_size")
	startTime := string(ctx.PostArgs().Peek("start_time"))
	endTime := string(ctx.PostArgs().Peek("end_time"))
	username := string(ctx.PostArgs().Peek("username"))
	title := string(ctx.PostArgs().Peek("title"))

	if page == 0 {
		page = 1
	}
	if pageSize < 10 {
		pageSize = 10
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
	if title != "" {
		ex["title"] = title
	}
	data, err := model.MessageSystemList(startTime, endTime, page, pageSize, ex)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}

// 站内信删除
func (that *MessageController) Delete(ctx *fasthttp.RequestCtx) {

	id := string(ctx.PostArgs().Peek("id"))
	tss := string(ctx.PostArgs().Peek("tss"))

	if id == "" && tss == "" {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	err := model.MessageDelete(id, tss)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}
