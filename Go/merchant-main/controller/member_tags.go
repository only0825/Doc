package controller

import (
	"merchant/contrib/helper"
	"merchant/contrib/validator"
	"merchant/model"
	"strings"

	"github.com/valyala/fasthttp"
)

// Tags 用户标签
func (that *MemberController) Tags(ctx *fasthttp.RequestCtx) {

	uid := string(ctx.QueryArgs().Peek("uid"))
	if !validator.CheckStringDigit(uid) {
		helper.Print(ctx, false, helper.UIDErr)
		return
	}

	data, err := model.MemberTagsList(uid)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)

}

// CancelTags 取消用户标签
// 取消单个用户标签时uid为用户id, 批量取消多个用户标签的时候uid用`,`分割
func (that *MemberController) CancelTags(ctx *fasthttp.RequestCtx) {

	params := setTagParam{}
	err := validator.Bind(ctx, &params)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	err = model.MemberTagsCancel(params.uid, params.tags)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// SetTags 设置用户标签
// 为单个用户设置标签时batch=0,uid为用户id, 为多个用户批量设置标签的时候batch=1,uid用`,`分割
func (that *MemberController) SetTags(ctx *fasthttp.RequestCtx) {

	params := setTagParam{}
	err := validator.Bind(ctx, &params)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	// 校验tags，并拆解组装成slice
	var tags []string
	for _, v := range strings.Split(params.tags, ",") {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}

		tags = append(tags, v)
	}
	if len(tags) == 0 {
		helper.Print(ctx, false, helper.UserTagErr)
		return
	}

	// 校验uid，并拆解组装成slice
	uids := strings.Split(params.uid, ",")
	var ids []string
	for _, v := range uids {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}

		ids = append(ids, v)
	}

	if len(ids) == 0 {
		helper.Print(ctx, false, helper.UIDErr)
		return
	}

	admin, err := model.AdminToken(ctx)
	if err != nil {
		helper.Print(ctx, false, helper.AccessTokenExpires)
		return
	}

	err = model.MemberTagsSet(params.Batch, admin["id"], ids, tags, ctx.Time().Unix())
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}
