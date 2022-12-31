package controller

import (
	g "github.com/doug-martin/goqu/v9"
	"github.com/valyala/fasthttp"
	"merchant/contrib/helper"
	"merchant/contrib/validator"
	"merchant/model"
)

// SMSChannelController 会员端接口
type SMSChannelController struct{}

// List 短信通道列表及按 渠道名称，创建人 筛选
func (*SMSChannelController) List(ctx *fasthttp.RequestCtx) {

	channelName := string(ctx.PostArgs().Peek("name"))
	createdName := string(ctx.PostArgs().Peek("created_name"))

	ex := g.Ex{}

	if channelName != "" {
		if len(channelName) < 5 || len(channelName) >= 30 {
			helper.Print(ctx, false, helper.ParamErr)
			return
		}

		ex["name"] = channelName
	}

	if createdName != "" {
		if !validator.CheckAName(createdName, 5, 20) {
			helper.Print(ctx, false, helper.AdminNameErr)
			return
		}

		ex["created_name"] = createdName
	}

	list, err := model.SMSChannelList(ex)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, list)
}

func (*SMSChannelController) UpdateState(ctx *fasthttp.RequestCtx) {

	id := string(ctx.PostArgs().Peek("id"))             // 短信通道ID
	txtState := ctx.PostArgs().GetUintOrZero("txt")     // 短信通道文字状态
	voiceState := ctx.PostArgs().GetUintOrZero("voice") // 短信通道语音状态

	if !validator.CtypeDigit(id) {
		helper.Print(ctx, false, helper.DBErr)
		return
	}

	if txtState != 0 {
		if txtState != 1 && txtState != 2 {
			helper.Print(ctx, false, helper.StateParamErr)
			return
		}
	}

	if voiceState != 0 {
		if voiceState != 1 && voiceState != 2 {
			helper.Print(ctx, false, helper.StateParamErr)
			return
		}
	}

	err := model.SMSChannelUpdateState(id, txtState, voiceState)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, "success")
}

func (*SMSChannelController) Update(ctx *fasthttp.RequestCtx) {

	id := string(ctx.PostArgs().Peek("id"))             // 短信通道ID
	channelName := string(ctx.PostArgs().Peek("name"))  // 短信通道名称
	remark := string(ctx.PostArgs().Peek("remark"))     // 短信通道备注
	alias := string(ctx.PostArgs().Peek("alias"))       // 短信通道别名
	txtState := ctx.PostArgs().GetUintOrZero("txt")     // 短信通道文字状态
	voiceState := ctx.PostArgs().GetUintOrZero("voice") // 短信通道语音状态

	rc := g.Record{}

	if !validator.CtypeDigit(id) {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	if channelName != "" {
		if len(channelName) < 5 || len(channelName) >= 30 {
			helper.Print(ctx, false, helper.ParamErr)
			return
		}
		rc["name"] = channelName
	}

	if alias != "" {
		if !validator.CheckStringAlpha(alias) || len(alias) < 3 {
			helper.Print(ctx, false, helper.ParamErr)
			return
		}
		rc["alias"] = alias
	}

	if remark != "" {
		rc["remark"] = remark
	}

	tm := map[int]int{
		0: 1,
		1: 2,
		2: 3,
	}

	if txtState < 0 && txtState > 2 {
		helper.Print(ctx, false, helper.StateParamErr)
		return
	}
	rc["txt"] = tm[txtState]

	if voiceState < 0 && voiceState > 2 {
		helper.Print(ctx, false, helper.StateParamErr)
		return
	}
	rc["voice"] = tm[voiceState]

	err := model.SMSChannelUpdate(id, rc)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, "success")
}

func (*SMSChannelController) Insert(ctx *fasthttp.RequestCtx) {

	channelName := string(ctx.PostArgs().Peek("name"))  // 短信通道名称
	remark := string(ctx.PostArgs().Peek("remark"))     // 短信通道备注
	alias := string(ctx.PostArgs().Peek("alias"))       // 短信通道别名
	txtState := ctx.PostArgs().GetUintOrZero("txt")     // 短信通道文字状态
	voiceState := ctx.PostArgs().GetUintOrZero("voice") // 短信通道语音状态

	admin, err := model.AdminToken(ctx)
	if err != nil {
		helper.Print(ctx, false, helper.AccessTokenExpires)
		return
	}

	if channelName == "" || alias == "" || (txtState == 0 && voiceState == 0) {
		helper.Print(ctx, false, helper.ParamNull)
		return
	}

	data := model.SMSChannel{}

	if channelName != "" {
		if len(channelName) < 5 || len(channelName) >= 30 {
			helper.Print(ctx, false, helper.ParamErr)
			return
		}
		data.Name = channelName
	}

	if alias != "" {
		if !validator.CheckStringAlpha(alias) || len(alias) < 4 {
			helper.Print(ctx, false, helper.ParamErr)
			return
		}
		data.Alias = alias
	}

	if remark != "" {
		data.Remark = remark
	}

	//tm := map[int]string{
	//	0: "1", // 没有
	//	1: "2", // 开启
	//	2: "3", // 关闭
	//}
	data.Txt = "0"
	data.Voice = "0"

	if txtState != 0 {
		//data.Txt = tm[2]
		data.Txt = "2"
	}

	if voiceState != 0 {
		//data.Voice = tm[2]
		data.Voice = "2"
	}

	data.CreatedUid = admin["id"]
	data.CreatedName = admin["name"]

	err = model.SMSChannelInsert(&data)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, "success")
}
