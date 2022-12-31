package controller

import (
	g "github.com/doug-martin/goqu/v9"
	"github.com/valyala/fasthttp"
	"merchant/contrib/helper"
	"merchant/model"
)

type PlatformController struct{}

//List 场馆列表
func (that *PlatformController) List(ctx *fasthttp.RequestCtx) {

	id := string(ctx.QueryArgs().Peek("id"))
	state := ctx.QueryArgs().GetUintOrZero("state")
	maintained := ctx.QueryArgs().GetUintOrZero("maintained")
	gameType := ctx.QueryArgs().GetUintOrZero("game_type")
	page := ctx.QueryArgs().GetUintOrZero("page")
	pageSize := ctx.QueryArgs().GetUintOrZero("page_size")

	if page < 1 {
		page = 1
	}
	if pageSize < 10 || pageSize > 100 {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	ex := g.Ex{}
	if id != "" {
		ex["id"] = id
	} else {
		if state > 0 {
			if state > 2 {
				helper.Print(ctx, false, helper.StateParamErr)
				return
			}

			ex["state"] = state
		}

		if maintained > 0 {
			if maintained > 2 {
				helper.Print(ctx, false, helper.StateParamErr)
				return
			}

			ex["maintained"] = maintained
		}

		if gameType > 0 {
			if gameType > 9 {
				helper.Print(ctx, false, helper.GameTypeErr)
				return
			}

			ex["game_type"] = gameType
		}
	}

	data, err := model.PlatformList(ex, page, pageSize)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}

//Update 场馆更新
func (that *PlatformController) Update(ctx *fasthttp.RequestCtx) {

	id := ctx.PostArgs().GetUintOrZero("id")
	state := ctx.PostArgs().GetUintOrZero("state")
	maintained := ctx.PostArgs().GetUintOrZero("maintained")
	//maintainedStart := string(ctx.PostArgs().Peek("maintained_start"))
	//maintainedEnd := string(ctx.PostArgs().Peek("maintained_end"))
	seq := ctx.PostArgs().GetUintOrZero("seq")

	if id == 0 {
		helper.Print(ctx, false, helper.PlatIDErr)
		return
	}

	ex := g.Ex{
		"id": id,
	}
	p, err := model.PlatformFindOne(ex)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	if state == 0 && maintained == 0 && seq == 0 ||
		p.State == state && p.Maintained == maintained && p.Seq == seq {
		helper.Print(ctx, false, helper.NoDataUpdate)
		return
	}

	record := g.Record{}
	if state > 0 {
		if state > 2 {
			helper.Print(ctx, false, helper.StateParamErr)
			return
		}

		record["state"] = state
	}

	if maintained > 0 {
		if maintained > 2 {
			helper.Print(ctx, false, helper.StateParamErr)
			return
		}

		if maintained == 2 {
			//if maintainedStart == "" || maintainedEnd == "" {
			//	helper.Print(ctx, false, helper.ParamErr)
			//	return
			//}
			//
			//record["maintained"] = maintained
			//record["maintained_start"] = maintainedStart
			record["maintained"] = maintained

		}

		if maintained == 1 {
			record["maintained"] = maintained
		}

		//record["maintained_end"] = maintainedEnd
	}

	if seq > 0 {
		if seq > 999 {
			helper.Print(ctx, false, helper.PlatSeqErr)
			return
		}

		record["seq"] = seq
	}

	err = model.PlatformUpdate(ex, record)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}
