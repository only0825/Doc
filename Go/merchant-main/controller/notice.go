package controller

import (
	"merchant/contrib/helper"
	"merchant/contrib/validator"
	"merchant/model"

	g "github.com/doug-martin/goqu/v9"
	"github.com/valyala/fasthttp"
)

type NoticeController struct{}

type noticeInsertParam struct {
	Title       string `name:"title" rule:"filter" min:"1" max:"255" msg:"title error[2-50]"`
	Content     string `name:"content" rule:"filter" min:"1" max:"1000" msg:"content error[1-1000]"`
	Redirect    int    `name:"redirect" rule:"digit" min:"1" max:"2" msg:"redirect error"` // 是否跳转 1 是 2否
	RedirectUrl string `name:"redirect_url" rule:"none"`                                   // 跳转链接
}

type noticeUpdateParam struct {
	ID          string `rule:"none" name:"id"`
	Title       string `name:"title" rule:"filter" required:"0" min:"1" max:"255" msg:"title error[2-50]"`
	Content     string `name:"content" rule:"filter" required:"0" min:"1" max:"1000" msg:"content error[1-1000]"`
	Redirect    int    `name:"redirect" rule:"digit" min:"1" max:"2" required:"0" default:"0"  msg:"redirect error"` // 是否跳转 1 是 2否
	RedirectUrl string `name:"redirect_url" rule:"none" required:"0"`                                                // 跳转链接
}

type noticeListParam struct {
	Title       string `name:"title" rule:"none" default:""`
	State       int    `name:"state" rule:"none" default:"0"`
	CreatedName string `name:"created_name" rule:"none" default:""`
	StartTime   string `name:"start_time" rule:"none"`
	EndTime     string `name:"end_time" rule:"none"`
	Page        uint   `name:"page" rule:"digit" default:"1" msg:"page error"`
	PageSize    uint   `name:"page_size" rule:"digit" default:"15" msg:"page_size error"`
}

func (that *NoticeController) Insert(ctx *fasthttp.RequestCtx) {

	param := noticeInsertParam{}
	err := validator.Bind(ctx, &param)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	admin, err := model.AdminToken(ctx)
	if err != nil {
		helper.Print(ctx, false, helper.AccessTokenExpires)
		return
	}

	data := model.Notice{
		ID:          helper.GenId(),
		Title:       param.Title,
		RedirectUrl: param.RedirectUrl,
		Redirect:    param.Redirect,
		Content:     param.Content,
		State:       1,
		CreatedAt:   ctx.Time().Unix(),
		CreatedUid:  admin["id"],
		CreatedName: admin["name"],
	}
	err = model.NoticeInsert(data)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// 系统公告列表
func (that *NoticeController) List(ctx *fasthttp.RequestCtx) {

	param := noticeListParam{}
	err := validator.Bind(ctx, &param)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	ex := g.Ex{}
	if param.Title != "" {
		ex["title"] = param.Title
	}

	if param.State != 0 {

		if param.State < 1 && param.State > 4 {
			helper.Print(ctx, false, helper.StateParamErr)
			return
		}

		ex["state"] = param.State
	}

	if param.CreatedName != "" {
		ex["created_name"] = param.CreatedName
	}

	data, err := model.NoticeList(param.Page, param.PageSize, param.StartTime, param.EndTime, ex)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}

// 系统公告编辑
func (that *NoticeController) Update(ctx *fasthttp.RequestCtx) {

	param := noticeUpdateParam{}
	err := validator.Bind(ctx, &param)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	ex := g.Ex{
		"id": param.ID,
	}
	record := g.Record{}
	if param.Title != "" {
		record["title"] = param.Title
	}

	if param.Content != "" {
		record["content"] = param.Content
	}

	if param.Redirect > 0 {
		record["redirect"] = param.Redirect
	}

	if param.RedirectUrl != "" {
		record["redirect_url"] = param.RedirectUrl
	}

	if param.Redirect == 2 {
		record["redirect_url"] = ""
	}

	if len(record) == 0 {
		helper.Print(ctx, false, helper.NoDataUpdate)
		return
	}

	err = model.NoticeUpdate(ex, record)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// 系统停用 启用 系统审核
func (that *NoticeController) UpdateState(ctx *fasthttp.RequestCtx) {

	id := string(ctx.PostArgs().Peek("id"))
	if id == "" || !validator.CheckStringDigit(id) {
		helper.Print(ctx, false, helper.IDErr)
		return
	}

	state := ctx.PostArgs().GetUintOrZero("state")
	s := map[int]bool{
		1: true, //停用
		2: true, //启用
	}
	if _, ok := s[state]; !ok {
		helper.Print(ctx, false, helper.StateParamErr)
		return
	}

	err := model.NoticeUpdateState(id, state)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// 系统公告删除
func (that *NoticeController) Delete(ctx *fasthttp.RequestCtx) {

	id := string(ctx.QueryArgs().Peek("id"))
	if id == "" || !validator.CheckStringDigit(id) {
		helper.Print(ctx, false, helper.IDErr)
		return
	}

	err := model.NoticeDelete(id)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}
