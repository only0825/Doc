package controller

import (
	"fmt"
	"merchant/contrib/helper"
	"merchant/contrib/validator"
	"merchant/model"
	"strings"

	g "github.com/doug-martin/goqu/v9"
	"github.com/shopspring/decimal"
	"github.com/valyala/fasthttp"
)

type dividendInsertParam struct {
	Username      string `rule:"alnum" name:"username" min:"5" max:"14" msg:"username error"`
	Ty            int    `rule:"digit" name:"ty" min:"211" max:"222"  msg:"ty error"`
	WaterLimit    uint8  `rule:"digit" name:"water_limit" min:"1" max:"2"  msg:"water_limit error"`
	Pid           string `rule:"none" name:"pid" default:"0"`    //活动id，仅适用于静态展示页活动，只有发放活动红利是才需要选择
	PTitle        string `rule:"none" name:"ptitle" default:"0"` //活动名，仅适用于静态展示页活动，只有发放活动红利是才需要选择
	Amount        string `rule:"digit" name:"amount"  msg:"amount error"`
	WaterMultiple int64  `rule:"digit" name:"water_multiple" min:"0" max:"1000" default:"0" required:"0"  msg:"water_multiple error"`
	Remark        string `rule:"filter" name:"remark" min:"1" max:"300"  msg:"remark error"`
}

// 红利审核列表参数
type dividendUpdateParam struct {
	IDS          string `rule:"sDigit" name:"ids" msg:"ids error"`                             // 订单号
	ReviewRemark string `rule:"none" name:"review_remark" max:"300" msg:"review_remark error"` // 审核备注
	State        int    `name:"state" rule:"digit" min:"232" max:"233" msg:"state error"`      // 231 审核中 232 审核不通过  232 审核通过
}

// 修改红利审核备注
type dividendReviewRemarkUpdateParam struct {
	ID           string `rule:"digit" name:"id" msg:"id error"`
	ReviewRemark string `rule:"filter" name:"review_remark" min:"1" max:"100"`
}

type DividendController struct{}

func (that *DividendController) Insert(ctx *fasthttp.RequestCtx) {

	param := dividendInsertParam{}
	err := validator.Bind(ctx, &param)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	if param.Ty == helper.TransactionDividend {
		if param.Pid == "" || !validator.CtypeDigit(param.Pid) {
			helper.Print(ctx, false, helper.PlatIDErr)
			return
		}
	}

	admin, err := model.AdminToken(ctx)
	if err != nil {
		helper.Print(ctx, false, helper.AccessTokenExpires)
		return
	}

	// 仅中心钱包红利支持负数
	amount, ok := validator.CheckFloatScope(param.Amount, "-20000000.000", "20000000.000")
	if !ok {
		helper.Print(ctx, false, helper.AmountErr)
		return
	}

	waterFlow := decimal.Decimal{}
	// 需要流水限制
	if param.WaterLimit == 2 {

		// 红利金额为负数，不允许有流水
		if amount.Cmp(decimal.Zero) == -1 {
			helper.Print(ctx, false, helper.TurnoverMutiErr)
			return
		}

		if param.WaterMultiple == 0 {
			helper.Print(ctx, false, helper.TurnoverMutiErr)
			return
		}

		waterFlow = amount.Mul(decimal.NewFromInt(param.WaterMultiple))
	}

	m, err := model.MemberFindOne(param.Username)
	if err != nil || m.UID == "" {
		helper.Print(ctx, false, helper.UsernameErr)
		return
	}

	mb, err := decimal.NewFromString(m.Balance)
	if err != nil {
		helper.Print(ctx, false, helper.AmountErr)
		return
	}

	// 红利金额为负数，判断会员余额
	if amount.Cmp(decimal.Zero) == -1 {
		// 余额不够扣除
		if mb.Add(amount).Cmp(decimal.Zero) == -1 {
			helper.Print(ctx, false, fmt.Sprintf("%s,%s%s", helper.AmountErr, mb.String(), amount.String()))
			return
		}
	}

	data := g.Record{
		"id":             helper.GenId(),
		"uid":            m.UID,
		"pid":            0,
		"username":       param.Username,
		"level":          m.Level,
		"top_uid":        m.TopUid,
		"top_name":       m.TopName,
		"parent_uid":     m.ParentUid,
		"parent_name":    m.ParentName,
		"ty":             param.Ty,
		"water_limit":    param.WaterLimit,
		"water_multiple": param.WaterMultiple,
		"water_flow":     waterFlow.String(),
		"amount":         param.Amount,
		"remark":         param.Remark,
		"apply_at":       uint64(ctx.Time().UnixMilli()),
		"apply_uid":      admin["id"],
		"apply_name":     admin["name"],
		"automatic":      1, //手动发放
		"state":          model.DividendReviewing,
		"tester":         m.Tester,
	}

	if param.Ty == model.DividendPromo {
		data["pid"] = param.Pid
		data["ptitle"] = param.PTitle
	}

	err = model.DividendInsert(data)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// 红利列表
func (that *DividendController) List(ctx *fasthttp.RequestCtx) {

	id := ctx.PostArgs().GetUintOrZero("id")                            //订单号
	username := string(ctx.PostArgs().Peek("username"))                 //用户名
	ty := ctx.PostArgs().GetUintOrZero("ty")                            //红利类型
	waterLimit := ctx.PostArgs().GetUintOrZero("water_limit")           //流水限制 1无需流水限制 2需要流水限制
	applyName := string(ctx.PostArgs().Peek("apply_name"))              //申请人名字
	reviewName := string(ctx.PostArgs().Peek("review_name"))            //审核人名字
	startTime := string(ctx.PostArgs().Peek("start_time"))              //申请开始时间
	endTime := string(ctx.PostArgs().Peek("end_time"))                  //申请结束时间
	reviewStartTime := string(ctx.PostArgs().Peek("review_start_time")) //审核开始时间
	reviewEndTime := string(ctx.PostArgs().Peek("review_end_time"))     //审核结束时间
	remarkFlag := ctx.PostArgs().GetUintOrZero("remark_flag")           //0全部 1申请备注 2审核备注
	remark := string(ctx.PostArgs().Peek("remark"))                     //备注内容
	page := ctx.PostArgs().GetUintOrZero("page")                        //页数
	pageSize := ctx.PostArgs().GetUintOrZero("page_size")               //页大小
	flag := ctx.PostArgs().GetUintOrZero("flag")                        //0 所有 1 审核列表 2历史列表
	state := ctx.PostArgs().GetUintOrZero("state")                      //审核状态

	ex := g.Ex{}
	if username != "" {
		if !validator.CheckStringAlnum(username) {
			helper.Print(ctx, false, helper.UsernameErr)
			return
		}

		ex["username"] = username
	}

	if ty != 0 {
		if ty < model.DividendSite || ty > model.DividendAgency {
			helper.Print(ctx, false, helper.DividendTypeErr)
			return
		}

		ex["ty"] = ty
	}

	if waterLimit != 0 {
		ex["water_limit"] = waterLimit
	}

	if applyName != "" && validator.CheckStringAlnum(applyName) {
		ex["apply_name"] = applyName
	}

	if reviewName != "" && validator.CheckStringAlnum(reviewName) {
		ex["review_name"] = reviewName
	}

	if remarkFlag > 0 && remark != "" {

		if remarkFlag == 1 {
			ex["remark"] = remark
		}

		if remarkFlag == 2 {
			ex["review_remark"] = remark
		}
	}

	// id查询，其他条件失效
	if id != 0 {
		ex = g.Ex{
			"id": id,
		}
	}

	// 红利审核列表
	if flag == 1 {
		ex["state"] = model.DividendReviewing //红利审核中
	} else { //
		ex["tester"] = 1
		// 默认为红利历史列表
		s := map[int]bool{
			model.DividendReviewPass:   true, //红利审核通过
			model.DividendReviewReject: true, //红利审核拒绝
		}

		// 查询所有
		if flag == 0 {
			s[model.DividendReviewing] = true //红利审核中
		}
		if state > 0 {
			if _, ok := s[state]; !ok {
				helper.Print(ctx, false, helper.StateParamErr)
				return
			}

			ex["state"] = state
		}
	}

	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 15
	}
	data, err := model.DividendList(page, pageSize, startTime, endTime, reviewStartTime, reviewEndTime, ex)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}

// 会员红利列表
func (that *DividendController) MemberList(ctx *fasthttp.RequestCtx) {

	username := string(ctx.FormValue("username"))
	ty := ctx.PostArgs().GetUintOrZero("ty")
	state := ctx.PostArgs().GetUintOrZero("state")
	startTime := string(ctx.FormValue("start_time"))
	endTime := string(ctx.FormValue("end_time"))
	page := ctx.PostArgs().GetUintOrZero("page")
	pageSize := ctx.PostArgs().GetUintOrZero("page_size")

	username = strings.ToLower(username)
	if username == "" || !validator.CheckUName(username, 5, 14) {
		helper.Print(ctx, false, helper.UsernameErr)
		return
	}

	ex := g.Ex{
		"username": username,
	}

	if page == 0 {
		page = 1
	}

	if pageSize == 0 {
		pageSize = 15
	}

	if ty != 0 {
		if ty < model.DividendSite || ty > model.DividendAgency {
			helper.Print(ctx, false, helper.DividendTypeErr)
			return
		}

		ex["ty"] = ty
	}

	if state != 0 {
		if state < model.DividendReviewing || state > model.DividendReviewReject {
			helper.Print(ctx, false, helper.StateParamErr)
			return
		}
		ex["state"] = state
	}

	data, err := model.DividendList(page, pageSize, startTime, endTime, "", "", ex)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}

// 红利审核
func (that *DividendController) Update(ctx *fasthttp.RequestCtx) {

	param := dividendUpdateParam{}
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

	ids := strings.Split(param.IDS, ",")
	err = model.DividendReview(param.State, ctx.Time().Unix(), admin["id"], admin["name"], param.ReviewRemark, ids)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// 更新  红利审核备注
func (that *DividendController) ReviewRemarkUpdate(ctx *fasthttp.RequestCtx) {

	param := dividendReviewRemarkUpdateParam{}
	err := validator.Bind(ctx, &param)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	state, err := model.DividendGetState(param.ID)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	if state == model.DividendReviewing {
		helper.Print(ctx, false, helper.ReviewStateErr)
		return
	}

	ex := g.Ex{
		"id": param.ID,
	}
	record := g.Record{
		"review_remark": param.ReviewRemark,
	}
	err = model.DividendUpdate(ex, record)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}
