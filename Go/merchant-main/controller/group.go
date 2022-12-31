package controller

import (
	"merchant/contrib/helper"
	"merchant/contrib/validator"
	"merchant/model"

	"github.com/valyala/fasthttp"
)

type GroupController struct{}

/**
* @Description: 用户组信息更新
* @Author: carl
 */
func (that *GroupController) Update(ctx *fasthttp.RequestCtx) {

	data := model.Group{}
	err := validator.Bind(ctx, &data)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	id := string(ctx.PostArgs().Peek("id"))
	if !helper.CtypeDigit(id) {
		helper.Print(ctx, false, helper.IDErr)
		return
	}

	admin, err := model.AdminToken(ctx)
	if err != nil {
		helper.Print(ctx, false, helper.AccessTokenExpires)
		return
	}
	err = model.GroupUpdate(id, admin["group_id"], data)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, "succeed")
}

/**
* @Description: 用户组插入
* @Author: carl
 */
func (that *GroupController) Insert(ctx *fasthttp.RequestCtx) {

	group := model.Group{}
	err := validator.Bind(ctx, &group)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	admin, err := model.AdminToken(ctx)
	if err != nil {
		helper.Print(ctx, false, helper.AccessTokenExpires)
		return
	}
	// 新增权限信息
	group.CreateAt = int32(ctx.Time().Unix())
	err = model.GroupInsert(admin["group_id"], group)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, "succeed")
}

/**
* @Description: 用户组列表获取
* @Author: carl
 */
func (that *GroupController) List(ctx *fasthttp.RequestCtx) {

	gid := string(ctx.QueryArgs().Peek("gid"))
	if gid != "" {
		if !helper.CtypeDigit(gid) {
			helper.Print(ctx, false, helper.GroupIDErr)
			return
		}
	}

	admin, err := model.AdminToken(ctx)
	if err != nil {
		helper.Print(ctx, false, helper.AccessTokenExpires)
		return
	}

	// 获取权限列表
	data, err := model.GroupList(gid, admin["group_id"])
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.PrintJson(ctx, true, data)
}
