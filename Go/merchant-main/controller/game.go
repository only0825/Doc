package controller

import (
	"merchant/contrib/helper"
	"merchant/contrib/validator"
	"merchant/model"
	"strconv"
	"strings"

	g "github.com/doug-martin/goqu/v9"
	"github.com/valyala/fasthttp"
)

type SlotsController struct{}

// 游戏编辑参数
type GameListUpdateParams struct {
	OnLine     uint8  `rule:"digit" min:"0" max:"1" msg:"online error" name:"online"` // 上线 下线
	ClientType string `rule:"sDigit" msg:"client_type error" name:"client_type"`      // 1:pc 2:h5 3:app
	ShowType   string `name:"show_type" rule:"sDigit" msg:"show_type error"`          // 1 正常 2 热门 3 最新
	ImgCover   string `name:"img_cover" rule:"none"`                                  // 封面图
	ID         string `rule:"digit" msg:"id error" name:"id"`
	Sorting    int    `name:"sorting" rule:"digit" msg:"sorting error"`
	Name       string `name:"name" rule:"none"`
	EnName     string `name:"en_name" rule:"none"`
	VnAlias    string `name:"vn_alias" rule:"none"`
}

// 维护游戏on_line参数
type GameListOnlineParams struct {
	ID     string `json:"id" rule:"digit" msg:"id error" name:"id"`
	OnLine int64  `json:"online" rule:"digit" min:"0" max:"1" msg:"online error" name:"online"`
}

// 维护游戏的上线下线
func (that *SlotsController) UpdateState(ctx *fasthttp.RequestCtx) {

	params := GameListOnlineParams{}
	err := validator.Bind(ctx, &params)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	//params.ID = string(ctx.PostArgs().Peek("id"))
	//fmt.Printf("ID = %s\n", params.ID)
	//if !validator.CheckStringDigit(params.ID) {
	//	helper.Print(ctx, false, helper.ParamErr)
	//	return
	//}
	//
	//s := string(ctx.PostArgs().Peek("online"))
	//fmt.Println(s)
	//if !validator.CheckStringDigit(s) {
	//	helper.Print(ctx, false, helper.ParamErr)
	//	return
	//}
	//params.OnLine = int64(ctx.PostArgs().GetUintOrZero("online"))
	//
	//fmt.Println("DONE")

	game, err := model.GameFind(params.ID)
	if err != nil {
		helper.Print(ctx, false, helper.RecordNotExistErr)
		return
	}

	//if game.OnLine == params.OnLine {
	//	helper.Print(ctx, false, helper.StateParamErr)
	//	return
	//}

	err = model.GameListUpdate(game.ID, game.PlatformId, g.Record{"online": params.OnLine})
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// 游戏编辑
func (that *SlotsController) Update(ctx *fasthttp.RequestCtx) {

	params := GameListUpdateParams{}
	err := validator.Bind(ctx, &params)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	game, err := model.GameFind(params.ID)
	if err != nil {
		helper.Print(ctx, false, helper.RecordNotExistErr)
		return
	}

	if params.ImgCover == "" || !strings.Contains(params.ImgCover, ".") {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	// ShowType 1 正常 2 热门 3 最新
	// is_new 0=否,1=是 《------》 is_hot 0正常 1热门 2 全部
	isHot := 0
	isNew := 0
	showTypes := strings.Split(params.ShowType, ",")
	for _, v := range showTypes {
		if v == "2" {
			isHot = 1
		}

		if v == "3" {
			isNew = 1
		}
	}

	record := g.Record{
		"online":   params.OnLine,
		"is_new":   isNew,
		"is_hot":   isHot,
		"en_name":  params.EnName,
		"name":     params.Name,
		"vn_alias": params.VnAlias,
	}

	// 处理 支持类型 多选 1:pc 2:h5 4:app 0 全部
	clientType := 0
	clientTypes := strings.Split(params.ClientType, ",")
	for _, v := range clientTypes {
		val, _ := strconv.Atoi(v)
		clientType += val
	}

	if clientType == 7 {
		clientType = 0
	}

	record["client_type"] = clientType

	if params.ImgCover != "" {
		record["img_cover"] = params.ImgCover
		record["img_pc"] = params.ImgCover
		record["img_phone"] = params.ImgCover
	}

	record["sorting"] = params.Sorting
	err = model.GameListUpdate(game.ID, game.PlatformId, record)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// 游戏列表
func (that *SlotsController) List(ctx *fasthttp.RequestCtx) {

	name := string(ctx.PostArgs().Peek("name"))
	enName := string(ctx.PostArgs().Peek("en_name"))
	platformId := string(ctx.PostArgs().Peek("platform_id"))
	gameType := ctx.PostArgs().GetUintOrZero("game_type")
	gameCode := string(ctx.PostArgs().Peek("game_code"))
	clientType := ctx.PostArgs().GetUintOrZero("client_type")
	online := string(ctx.PostArgs().Peek("online"))
	page := ctx.PostArgs().GetUintOrZero("page")
	pageSize := ctx.PostArgs().GetUintOrZero("page_size")

	ex := g.Ex{}
	if name != "" {
		ex["name"] = name
	}

	if online != "" {
		ex["online"] = online
	}

	if enName != "" {
		ex["en_name"] = enName
	}

	if platformId != "" {
		ex["platform_id"] = platformId
	}

	if gameType != 0 {
		ex["game_type"] = gameType
	}

	if gameCode != "" {
		ex["game_code"] = gameCode
	}

	if clientType != 0 {
		clientTypes := map[int][]int{
			1: {0, 1, 3, 5},
			2: {0, 2, 3, 6},
			4: {0, 4, 6, 5},
		}
		if _, ok := clientTypes[clientType]; !ok {
			helper.Print(ctx, false, helper.DeviceErr)
			return
		}

		ex["client_type"] = clientTypes[clientType]
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 15
	}
	data, err := model.GameList(ex, uint(page), uint(pageSize))
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}
