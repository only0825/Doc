package controller

import (
	g "github.com/doug-martin/goqu/v9"
	"merchant/contrib/helper"
	"merchant/contrib/validator"
	"merchant/model"

	"github.com/valyala/fasthttp"
)

type TagsController struct{}

type tagParam struct {
	Name        string `rule:"filter" min:"1" max:"20" msg:"name error" name:"name"`
	Description string `rule:"filter" min:"1" max:"100" msg:"description error" name:"description"`
	Sort        string `rule:"digit" min:"0" max:"999" msg:"sort error" name:"sort"`
	Flags       string `rule:"digit" min:"1" max:"6" msg:"flags error" name:"flags"`
}
type tagUpdateParam struct {
	ID          int64  `rule:"digit" default:"0" msg:"id error" name:"id"`
	Name        string `rule:"none" min:"1" max:"20" msg:"name error" name:"name"`
	Description string `rule:"none" min:"1" max:"50" msg:"description error" name:"description"`
	Sort        string `rule:"none" min:"0" max:"999" msg:"sort error" name:"sort"`
	Flags       string `rule:"none" min:"1" max:"6" msg:"flags error" name:"flags"`
}

// List 会员配置-标签管理（列表）
func (that *TagsController) List(ctx *fasthttp.RequestCtx) {

	name := string(ctx.PostArgs().Peek("name"))           // 标签名称
	flags := ctx.PostArgs().GetUintOrZero("flags")        // 标签类型 1:官方代理、2:普通代理、3:风控、4:财务、5:电维、6:其他
	all := ctx.PostArgs().GetUintOrZero("all")            // 是否取全部标签（用户列表需要全部的标签）
	page := ctx.PostArgs().GetUintOrZero("page")          //页数
	pageSize := ctx.PostArgs().GetUintOrZero("page_size") //页大小
	if name != "" && !validator.CheckStringLength(name, 1, 20) {
		helper.Print(ctx, false, helper.UserTagErr)
		return
	}
	if page < 1 {
		page = 1
	}

	if all == 0 && (pageSize < 10 || pageSize > 200) || flags > 6 {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	data, err := model.TagList(name, flags, all, page, pageSize)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}

// Insert 会员配置-添加标签
func (that *TagsController) Insert(ctx *fasthttp.RequestCtx) {

	params := tagParam{}
	err := validator.Bind(ctx, &params)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	// 校验name和flags是否存在
	exist, _ := model.TagByNameAndFlag(params.Name, params.Flags)
	if exist.ID != 0 {
		helper.Print(ctx, false, helper.RecordExistErr)
		return
	}

	ts := ctx.Time().Unix()
	fields := map[string]interface{}{
		"name":        params.Name,
		"description": params.Description,
		"sort":        params.Sort,
		"flags":       params.Flags,
		"created_at":  ts,
		"updated_at":  ts,
	}
	err = model.TagInsert(fields)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// Update 会员配置-修改标签
func (that *TagsController) Update(ctx *fasthttp.RequestCtx) {

	params := tagUpdateParam{}
	err := validator.Bind(ctx, &params)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	if params.ID < 1 {
		helper.Print(ctx, false, helper.IDErr)
		return
	}

	if len(params.Name) > 0 && len(params.Flags) > 0 {
		// 校验name和flags是否存在
		exist, _ := model.TagByNameAndFlag(params.Name, params.Flags)
		if exist.ID != 0 {
			helper.Print(ctx, false, helper.RecordExistErr)
			return
		}
	}

	record := g.Record{}
	if len(params.Name) > 0 {
		if !validator.CheckStringLength(params.Name, 1, 20) {
			helper.Print(ctx, false, helper.ParamErr)
			return
		}

		record["name"] = params.Name
	}

	if len(params.Description) > 0 {
		if !validator.CheckStringLength(params.Name, 1, 50) {
			helper.Print(ctx, false, helper.ParamErr)
			return
		}

		record["description"] = params.Description
	}

	if len(params.Sort) > 0 {
		if !validator.CheckIntScope(params.Sort, 1, 999) {
			helper.Print(ctx, false, helper.ParamErr)
			return
		}

		record["sort"] = params.Sort
	}

	if len(params.Flags) > 0 {
		if !validator.CheckIntScope(params.Flags, 1, 6) {
			helper.Print(ctx, false, helper.ParamErr)
			return
		}

		record["flags"] = params.Flags
	}

	if len(record) == 0 {
		helper.Print(ctx, false, helper.NoDataUpdate)
		return
	}

	record["updated_at"] = ctx.Time().Unix()

	ex := g.Ex{
		"id": params.ID,
	}
	err = model.TagUpdate(ex, record)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// Delete 会员配置-删除标签
func (that *TagsController) Delete(ctx *fasthttp.RequestCtx) {

	id := string(ctx.QueryArgs().Peek("id"))
	if !validator.CheckStringDigit(id) {
		helper.Print(ctx, false, helper.IDErr)
		return
	}

	err := model.TagDelete(id)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}
