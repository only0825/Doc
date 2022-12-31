package controller

import (
	g "github.com/doug-martin/goqu/v9"
	"github.com/shopspring/decimal"
	"github.com/valyala/fasthttp"
	"merchant/contrib/helper"
	"merchant/contrib/validator"
	"merchant/model"
	"strconv"
	"strings"
)

type AdjustController struct{}

type adjustParam struct {
	UID           string `rule:"digit" name:"uid" msg:"uid error"`
	Username      string `rule:"alnum" name:"username" min:"2" max:"30" msg:"username error"`               //
	AdjustType    int    `rule:"digit" min:"1" max:"3" name:"adjust_type" msg:"adjust type error"`          // 调整类型:1=系统调整,2=输赢调整,3=线下转卡充值
	AdjustMode    int    `rule:"digit" min:"251" max:"252" name:"adjust_mode" msg:"adjust mode error"`      // 调整方式:251=上分,252=下分
	IsTurnover    int    `rule:"digit" min:"0" max:"1" name:"is_turnover" msg:"turnover error"`             // 是否需要流水限制:1=需要,0=不需要
	TurnoverMulti int    `rule:"none" min:"1" name:"turnover_multi" default:"0" msg:"turnover multi error"` // 流水倍数
	Amount        string `rule:"float" name:"amount" msg:"amount error"`                                    // 调整金额
	ApplyRemark   string `rule:"filter" min:"1" max:"100" name:"apply_remark" msg:"apply remark error"`     // 申请备注
	Images        string `rule:"none" min:"1" max:"100" name:"images" msg:"images error"`                   // 图片地址
}

type adjustReviewParam struct {
	ID     string `rule:"digit" name:"id" msg:"id error"`
	State  int    `rule:"digit" min:"257" max:"258" name:"state" msg:"state error"`             // 状态:256=审核中,257=同意, 258=拒绝
	Remark string `rule:"filter" min:"1" max:"100" name:"remark" msg:"remark error" default:""` // 审核备注
}

// Insert 会员列表-账户调整
func (that *AdjustController) Insert(ctx *fasthttp.RequestCtx) {

	params := adjustParam{}
	err := validator.Bind(ctx, &params)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	// amount范围在1-2E越南盾(1k-200000k)
	_, ok := validator.CheckFloatScope(params.Amount, "1", "200000")
	if !ok {
		helper.Print(ctx, false, helper.AmountErr)
		return
	}

	// 系统调整-上分-有流水限制 和输赢调整-上分-有流水限制
	if params.AdjustMode == model.AdjustUpMode && params.IsTurnover == 1 {
		// 校验流水倍数  min:"1" max:"100" name:"turnover_multi"
		if params.TurnoverMulti < 1 || params.TurnoverMulti > 100 {
			helper.Print(ctx, false, helper.TurnoverMutiErr)
			return
		}
	} else {
		params.TurnoverMulti = 0
	}

	// 校验图片
	if params.Images != "" {
		if len(params.Images) < 5 {
			helper.Print(ctx, false, helper.ImagesURLErr)
			return
		}
	}
	//if params.Images != "" {
	//	if !validator.CheckUrl(params.Images) {
	//		helper.Print(ctx, false, helper.ImagesURLErr)
	//		return
	//	}
	//}

	// get member info
	m, err := model.MemberFindOne(params.Username)
	if err != nil || m.UID != params.UID {
		helper.Print(ctx, false, helper.UIDErr)
		return
	}

	admin, err := model.AdminToken(ctx)
	if err != nil {
		helper.Print(ctx, false, helper.AccessTokenExpires)
		return
	}

	amount, _ := decimal.NewFromString(params.Amount)
	data := model.MemberAdjust{
		ID:            helper.GenId(),
		UID:           params.UID,
		Ty:            0, // 后台调整
		Username:      params.Username,
		TopUid:        m.TopUid,
		TopName:       m.TopName,
		ParentUid:     m.ParentUid,
		ParentName:    m.ParentName,
		AdjustType:    params.AdjustType,
		AdjustMode:    params.AdjustMode,
		IsTurnover:    params.IsTurnover,
		TurnoverMulti: params.TurnoverMulti,
		ApplyRemark:   params.ApplyRemark,
		Images:        params.Images,
		State:         model.AdjustReviewing, // 状态:256=审核中,257=同意, 258=拒绝
		ApplyAt:       ctx.Time().Unix(),
		ApplyUid:      admin["id"],   // 申请人
		ApplyName:     admin["name"], // 申请人
		ReviewUid:     "0",
		Tester:        m.Tester,
	}
	data.Amount, _ = amount.Float64()
	err = model.AdjustInsert(data)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// List 会员管理-账户调整审核列表
func (that *AdjustController) List(ctx *fasthttp.RequestCtx) {

	id := string(ctx.QueryArgs().Peek("id"))
	username := string(ctx.QueryArgs().Peek("username"))            // 会员账号
	applyName := string(ctx.QueryArgs().Peek("apply_name"))         // 申请人
	reviewName := string(ctx.QueryArgs().Peek("review_name"))       // 审核人
	adjustType := ctx.QueryArgs().GetUintOrZero("adjust_type")      // 调整类型:1=系统调整,2=输赢调整,3=线下转卡充值
	adjustMode := ctx.QueryArgs().GetUintOrZero("adjust_mode")      // 调整方式:251=上方,252=下分
	state := ctx.QueryArgs().GetUintOrZero("state")                 // 审核状态
	handOutState := ctx.QueryArgs().GetUintOrZero("hand_out_state") // 发放状态
	ty := string(ctx.QueryArgs().Peek("ty"))                        // 记录来源 0 后台调整 1财务下分
	startTime := string(ctx.QueryArgs().Peek("start_time"))
	endTime := string(ctx.QueryArgs().Peek("end_time"))
	sPage := string(ctx.QueryArgs().Peek("page"))
	sPageSize := string(ctx.QueryArgs().Peek("page_size"))

	if !validator.CheckStringDigit(sPage) || !validator.CheckIntScope(sPageSize, 10, 200) {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	page, _ := strconv.Atoi(sPage)
	pageSize, _ := strconv.Atoi(sPageSize)
	//_, err := time.Parse("2006-01-02 15:04:05", startTime)
	//if err != nil {
	//	helper.Print(ctx, false, helper.TimeTypeErr)
	//	return
	//}
	//
	//_, err = time.Parse("2006-01-02 15:04:05", endTime)
	//if err != nil {
	//	helper.Print(ctx, false, helper.TimeTypeErr)
	//	return
	//}

	ex := g.Ex{}
	if id != "" {
		if !validator.CheckStringDigit(id) {
			helper.Print(ctx, false, helper.IDErr)
			return
		}

		ex["id"] = id
	}

	if username != "" {
		username = strings.ToLower(username)
		if !validator.CheckUName(username, 5, 14) {
			helper.Print(ctx, false, helper.UsernameErr)
			return
		}

		ex["username"] = username
	}

	if applyName != "" {
		if !validator.CheckAName(applyName, 5, 20) {
			helper.Print(ctx, false, helper.ApplyNameErr)
			return
		}

		ex["apply_name"] = applyName
	}

	if reviewName != "" {

		if !validator.CheckAName(reviewName, 5, 20) {
			helper.Print(ctx, false, helper.ReviewNameErr)
			return
		}

		ex["review_name"] = reviewName
	}

	if state > 0 {
		s := map[int]bool{
			256: true, //AdjustReviewing
			257: true, //AdjustReviewPass
			258: true, //AdjustReviewReject
		}
		if _, ok := s[state]; !ok {
			helper.Print(ctx, false, helper.ReviewStateErr)
			return
		}

		ex["state"] = state
	}

	if ty != "" {
		if !validator.CheckIntScope(ty, 0, 1) {
			helper.Print(ctx, false, helper.AdjustTyErr)
			return
		}
		ex["ty"] = ty
	}

	// 只有审核通过的记录才有发放状态
	if state == 257 {
		if handOutState > 0 {
			s := map[int]bool{
				261: true, //AdjustFailed
				262: true, //AdjustSuccess
				263: true, //AdjustPlatDealing
			}
			if _, ok := s[state]; !ok {
				helper.Print(ctx, false, helper.HandOutStateErr)
				return
			}

			ex["hand_out_state"] = handOutState
		}
	}

	if adjustType > 0 {
		s := map[int]bool{
			1: true, //系统调整
			2: true, //输赢调整
			3: true, //线下转卡充值
		}
		if _, ok := s[adjustType]; !ok {
			helper.Print(ctx, false, helper.HandOutStateErr)
			return
		}

		ex["adjust_type"] = adjustType
	}

	if adjustMode > 0 {
		s := map[int]bool{
			251: true, //AdjustUpMode
			252: true, //AdjustDownMode
		}
		if _, ok := s[adjustMode]; !ok {
			helper.Print(ctx, false, helper.HandOutStateErr)
			return
		}

		ex["adjust_mode"] = adjustMode
	}

	data, err := model.AdjustList(startTime, endTime, ex, page, pageSize)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}

// Review 会员管理-账户调整审核
func (that *AdjustController) Review(ctx *fasthttp.RequestCtx) {

	param := adjustReviewParam{}
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

	record := g.Record{
		"id":            param.ID,
		"review_remark": param.Remark,
		"state":         param.State,
		"review_at":     ctx.Time().Unix(),
		"review_uid":    admin["id"],
		"review_name":   admin["name"],
	}
	err = model.AdjustReview(param.State, record)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}
