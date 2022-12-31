package controller

import (
	g "github.com/doug-martin/goqu/v9"
	"github.com/olivere/elastic/v7"
	"github.com/valyala/fasthttp"
	"merchant/contrib/helper"
	"merchant/contrib/validator"
	"merchant/model"
)

type LogController struct{}

type adminLoginParam struct {
	StartTime string `name:"start_time" rule:"none" msg:"start_time error"`
	EndTime   string `name:"end_time" rule:"none" msg:"end_time error"`
	Name      string `name:"name" rule:"none" default:"" msg:"name error"`
	//Flag      int    `name:"flag" rule:"digit" default:"0" min:"0" max:"2" msg:"flag error"`
	Page     int `name:"page" rule:"digit" default:"1" msg:"page error"`
	PageSize int `name:"page_size" rule:"digit" default:"15" msg:"page_size error"`
}

type systemLogParam struct {
	StartTime string `name:"start_time" rule:"none" msg:"start_time error"`
	EndTime   string `name:"end_time" rule:"none" msg:"end_time error"`
	Name      string `name:"name" rule:"none" default:"" msg:"name error"`
	Title     string `name:"title" rule:"none" default:"" msg:"title error"`
	Page      int    `name:"page" rule:"digit" default:"1" msg:"page error"`
	PageSize  int    `name:"page_size" rule:"digit" default:"15" msg:"page_size error"`
}

// 登录日志
//func (that *LogController) AdminLoginLog(ctx *fasthttp.RequestCtx) {
//
//	params := adminLoginParam{}
//	err := validator.Bind(ctx, &params)
//	if err != nil {
//		helper.Print(ctx, false, helper.ParamErr)
//		return
//	}
//
//	query := elastic.NewBoolQuery()
//	if params.Flag != 0 {
//		query.Filter(elastic.NewTermQuery("flag", params.Flag))
//	}
//
//	if params.Name != "" {
//		query.Filter(elastic.NewTermQuery("name", params.Name))
//	}
//
//	data, err := model.AdminLoginLog(params.StartTime, params.EndTime, params.Page, params.PageSize, query)
//	if err != nil {
//		helper.Print(ctx, false, err.Error())
//		return
//	}
//
//	helper.Print(ctx, true, data)
//}

func (that *LogController) AdminLoginLog(ctx *fasthttp.RequestCtx) {

	params := adminLoginParam{}
	err := validator.Bind(ctx, &params)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	ex := g.Ex{}

	if params.Name != "" {
		ex["username"] = params.Name
	}

	if params.StartTime == "" || params.EndTime == "" {

	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 10 {
		params.PageSize = 10
	}

	data, err := model.AdminLoginLog(params.StartTime, params.EndTime, params.Page, params.PageSize, ex)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}

// 系统日志
func (that *LogController) SystemLog(ctx *fasthttp.RequestCtx) {

	params := systemLogParam{}
	err := validator.Bind(ctx, &params)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	query := elastic.NewBoolQuery()

	if params.Name != "" {
		query.Filter(elastic.NewTermQuery("name.keyword", params.Name))
	}

	if params.Title != "" {
		query.Filter(elastic.NewWildcardQuery("title.keyword", params.Title+"*"))
	}

	data, err := model.SystemLog(params.StartTime, params.EndTime, params.Page, params.PageSize, query)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}
