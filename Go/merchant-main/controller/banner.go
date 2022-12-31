package controller

import (
	"fmt"
	"merchant/contrib/helper"
	"merchant/contrib/validator"
	"merchant/model"
	"net/url"
	"strconv"
	"strings"
	"time"

	g "github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/valyala/fasthttp"
)

type BannerController struct{}

type bannerListParam struct {
	Flags     uint8  `rule:"digit" min:"1" max:"5" msg:"flags error,val [1-5]" name:"flags"`                   //1:APP闪屏页广告、2:轮播图广告、3:WEB体育场馆广告、4:站点广告位。
	Device    string `rule:"none" name:"device"`                                                               //设备号，多个设备用户逗号分隔
	StartTime string `rule:"none" msg:"start_time error" name:"start_time"`                                    //开始时间
	EndTime   string `rule:"none" msg:"end_time error" name:"end_time"`                                        //结束时间
	State     uint8  `rule:"digit" default:"0" min:"0" max:"3" msg:"state val[0-3]" required:"0" name:"state"` //状态 1待发布 2开启 3停用
	PageSize  uint   `rule:"digit" default:"10" min:"10" max:"200" msg:"page_size error" name:"page_size"`     //每页数量
	Page      uint   `rule:"digit" default:"1" min:"1" msg:"page error" name:"page"`                           //页码
}

type bannerUpdateParam struct {
	ID          string `json:"id" db:"id" rule:"digit" msg:"id error" name:"id"`                                                          //
	Title       string `json:"title" db:"title"  rule:"filter" msg:"title error" required:"0" name:"title"`                               //标题
	Device      string `json:"device" db:"device" rule:"sDigit" msg:"device error" required:"0" name:"device"`                            //设备类型(1,2)
	RedirectURL string `json:"redirect_url" db:"redirect_url" rule:"none" msg:"redirect_url error" required:"0" name:"redirect_url"`      //跳转地址
	Images      string `json:"images" db:"images" rule:"none" msg:"images error" required:"0" name:"images"`                              //图片路径
	Seq         string `json:"seq" db:"seq" rule:"digit" min:"1" max:"100" msg:"seq error" required:"0" name:"seq"`                       //排序
	Flags       string `json:"flags" db:"flags" rule:"digit" min:"1" max:"10" msg:"flags error" required:"0" name:"flags"`                //广告类型
	ShowType    string `json:"show_type" db:"show_type" rule:"digit" min:"1" max:"2" msg:"show_type error" required:"0" name:"show_type"` //1 永久有效 2 指定时间
	ShowAt      string `json:"show_at" db:"show_at" rule:"none" msg:"show_at error" required:"0" name:"show_at"`                          //开始展示时间
	HideAt      string `json:"hide_at" db:"hide_at" rule:"none" msg:"hide_at error" required:"0" name:"hide_at"`                          //结束展示时间
	URLType     string `json:"url_type" db:"url_type" rule:"digit" min:"0" max:"3" msg:"url_type error" required:"0" name:"url_type"`     //链接类型 1站内 2站外
}

func (that *BannerController) List(ctx *fasthttp.RequestCtx) {

	params := bannerListParam{}
	err := validator.Bind(ctx, &params)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	exs := exp.NewExpressionList(exp.AndType)
	ex := g.Ex{"flags": params.Flags}
	if params.Device != "" {
		//ex["device"] = params.Device
		exs = exs.Append(g.Or(g.Ex{"device": g.Op{"like": params.Device}}, g.Ex{"device": 0}))
	}

	if params.State > 0 {
		ex["state"] = params.State
	}

	exs = exs.Append(ex)

	data, err := model.BannerList(params.StartTime, params.EndTime, params.Page, params.PageSize, exs)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}

func (that *BannerController) Insert(ctx *fasthttp.RequestCtx) {

	device := string(ctx.PostArgs().Peek("device"))
	flags := string(ctx.PostArgs().Peek("flags"))
	hideAt := string(ctx.PostArgs().Peek("hide_at"))
	images := string(ctx.PostArgs().Peek("images"))
	redirectUrl := string(ctx.PostArgs().Peek("redirect_url"))
	seq := ctx.PostArgs().GetUintOrZero("seq")
	showAt := string(ctx.PostArgs().Peek("show_at"))
	showType := string(ctx.PostArgs().Peek("show_type"))
	title := string(ctx.PostArgs().Peek("title"))
	urlType := string(ctx.PostArgs().Peek("url_type"))

	title = strings.ReplaceAll(title, "&nbsp;", " ")
	if len(title) < 3 || len(title) > 200 {
		helper.Print(ctx, false, helper.ContentLengthErr)
		return
	}

	ds := strings.Split(device, ",")
	if len(ds) > 0 {
		for _, v := range ds {
			i, _ := strconv.Atoi(v)
			if _, ok := model.DeviceMap[i]; !ok {
				helper.Print(ctx, false, helper.DeviceTypeErr)
				return
			}
		}
	}

	if seq < 1 || seq > 100 {
		helper.Print(ctx, false, helper.PlatSeqErr)
		return
	}

	if showType == model.BannerShowTypeSpecify {
		_, err := time.Parse("2006-01-02 15:04:05", showAt)
		if err != nil {
			helper.Print(ctx, false, helper.DateTimeErr)
			return
		}

		_, err = time.Parse("2006-01-02 15:04:05", hideAt)
		if err != nil {
			helper.Print(ctx, false, helper.DateTimeErr)
			return
		}
	}

	data, err := model.AdminToken(ctx)
	if err != nil {
		helper.Print(ctx, false, helper.AccessTokenExpires)
		return
	}

	switch urlType {
	case "1": //站内链接
		if !strings.HasPrefix(redirectUrl, "/") {
			helper.Print(ctx, false, helper.URLErr)
			return
		}
	case "2": //站外链接
		_, err := url.Parse(redirectUrl)
		if err != nil {
			helper.Print(ctx, false, helper.URLErr)
			return
		}
	case "3": //场馆链接
		if !helper.CtypeDigit(redirectUrl) {
			helper.Print(ctx, false, helper.URLErr)
			return
		}
	default:
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	banner := model.Banner{
		ID:          helper.GenId(),
		Title:       title,
		Images:      images,
		Flags:       flags,
		Device:      device,
		URLType:     urlType,
		RedirectURL: redirectUrl,
		Seq:         fmt.Sprintf("%d", seq),
		HideAt:      hideAt,
		ShowType:    showType,
		ShowAt:      showAt,
		UpdatedAt:   fmt.Sprintf("%d", ctx.Time().Unix()),
		UpdatedName: data["name"],
		UpdatedUID:  data["id"],
		State:       model.BannerStateWait,
	}
	err = model.BannerInsert(banner)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

func (that *BannerController) Update(ctx *fasthttp.RequestCtx) {

	id := string(ctx.PostArgs().Peek("id"))
	device := string(ctx.PostArgs().Peek("device"))
	flags := string(ctx.PostArgs().Peek("flags"))
	hideAt := string(ctx.PostArgs().Peek("hide_at"))
	images := string(ctx.PostArgs().Peek("images"))
	redirectUrl := string(ctx.PostArgs().Peek("redirect_url"))
	seq := ctx.PostArgs().GetUintOrZero("seq")
	showAt := string(ctx.PostArgs().Peek("show_at"))
	showType := string(ctx.PostArgs().Peek("show_type"))
	title := string(ctx.PostArgs().Peek("title"))
	urlType := string(ctx.PostArgs().Peek("url_type"))

	title = strings.ReplaceAll(title, "&nbsp;", " ")
	record := g.Record{}
	if title != "" {
		if len(title) < 3 || len(title) > 200 {
			helper.Print(ctx, false, helper.ContentLengthErr)
			return
		}

		record["title"] = title
	}

	if device != "" {
		ds := strings.Split(device, ",")
		if len(ds) > 0 {
			for _, v := range ds {
				i, _ := strconv.Atoi(v)
				if _, ok := model.DeviceMap[i]; !ok {
					helper.Print(ctx, false, helper.DeviceTypeErr)
					return
				}
			}
		}

		record["device"] = device
	}

	if redirectUrl != "" {
		switch urlType {
		case "1": //站内链接
			if !strings.HasPrefix(redirectUrl, "/") {
				helper.Print(ctx, false, helper.ParamErr)
				return
			}
		case "2": //站外链接
			_, err := url.Parse(redirectUrl)
			if err != nil {
				helper.Print(ctx, false, helper.ParamErr)
				return
			}
		case "3": //场馆链接
			if !helper.CtypeDigit(redirectUrl) {
				helper.Print(ctx, false, helper.URLErr)
				return
			}
		default:
			helper.Print(ctx, false, helper.ParamErr)
			return
		}

		record["redirect_url"] = redirectUrl
	}

	if images != "" {
		record["images"] = images
	}

	if seq != 0 {
		if seq < 1 || seq > 100 {
			helper.Print(ctx, false, helper.PlatSeqErr)
			return
		}

		record["seq"] = seq
	}

	if flags != "" {
		record["flags"] = flags
	}

	if showType != "" {
		if showType != "1" && showType != "2" {
			helper.Print(ctx, false, helper.ParamErr)
			return
		}

		record["show_type"] = showType
	}

	if urlType != "" {
		record["url_type"] = urlType
	}

	if len(record) == 0 {
		helper.Print(ctx, false, helper.NoDataUpdate)
		return
	}

	data, err := model.AdminToken(ctx)
	if err != nil {
		helper.Print(ctx, false, helper.AccessTokenExpires)
		return
	}

	record["updated_name"] = data["name"]
	record["updated_uid"] = data["id"]
	record["updated_at"] = fmt.Sprintf("%d", ctx.Time().Unix())
	err = model.BannerUpdate(showAt, hideAt, id, record)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

func (that *BannerController) Delete(ctx *fasthttp.RequestCtx) {

	id := string(ctx.QueryArgs().Peek("id"))
	if !validator.CheckStringDigit(id) {
		helper.Print(ctx, false, helper.IDErr)
		return
	}

	err := model.BannerDelete(id)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

func (that *BannerController) UpdateState(ctx *fasthttp.RequestCtx) {

	id := string(ctx.PostArgs().Peek("id"))
	if !validator.CheckStringDigit(id) {
		helper.Print(ctx, false, helper.IDErr)
		return
	}

	stateVal := string(ctx.PostArgs().Peek("state"))
	if !validator.CheckStringDigit(stateVal) &&
		!validator.CheckIntScope(stateVal, model.BannerStateWait, model.BannerStateEnd) {
		helper.Print(ctx, false, helper.StateParamErr)
		return
	}

	state, _ := strconv.Atoi(stateVal)
	err := model.BannerUpdateState(id, uint8(state))
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}
