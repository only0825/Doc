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

	"github.com/doug-martin/goqu/v9/exp"

	"github.com/shopspring/decimal"
	"github.com/wI2L/jettison"

	g "github.com/doug-martin/goqu/v9"
	"github.com/valyala/fasthttp"
)

type MemberController struct{}

type memberStateParam struct {
	Username string `rule:"none" name:"username"`                                 // 用户username  批量用逗号隔开
	State    int8   `rule:"digit" name:"state" min:"1" max:"2" msg:"state error"` // 状态： 1 正常 2 禁用
	Remark   string `rule:"filter" name:"remark" max:"300" msg:"remark error"`    // 备注
}

// setTagParam 设置/批量设置用户标签，取消用户标签
type setTagParam struct {
	Batch int    `rule:"digit" min:"0" max:"1" default:"0" msg:"batch error" name:"batch"` // 1批量添加 0编辑单个用户标签
	uid   string `rule:"sDigit" msg:"uid error" name:"uid"`
	tags  string `rule:"sDigit" min:"1" msg:"tags error" name:"tags"`
}

// setSVipParam 解除密码限制/解除短信限制 parameters structure
type retryResetParam struct {
	Username string `rule:"uname" min:"5" max:"14" msg:"username error" name:"username"`
	Ty       uint8  `rule:"digit" min:"1" max:"3" msg:"ty error" name:"ty"` // 1解除密码限制 2解除短信限制 3解除场馆钱包限制
	Pid      string `rule:"none" msg:"pid error" name:"pid"`                // 场馆id(解除场馆钱包限制时需要)
}

// 用户备注参数
type remarkLogParams struct {
	Username string `rule:"none" name:"username" msg:"username error"`
	File     string `rule:"none" name:"file" msg:"file error" default:""`
	Msg      string `rule:"none" name:"msg" max:"300"`
}

func (that *MemberController) Detail(ctx *fasthttp.RequestCtx) {

	username := string(ctx.QueryArgs().Peek("username"))
	username = strings.ToLower(username)
	if !validator.CheckUName(username, 5, 14) {
		helper.Print(ctx, false, helper.UsernameErr)
		return
	}

	data, err := model.MemberInfo(username)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}

// GetAccountInfo 会员列表-帐户信息
func (that *MemberController) AccountInfo(ctx *fasthttp.RequestCtx) {

	username := string(ctx.QueryArgs().Peek("username"))
	username = strings.ToLower(username)
	if !validator.CheckUName(username, 5, 14) {
		helper.Print(ctx, false, helper.UsernameErr)
		return
	}

	data, err := model.MemberAccountInfo(username)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}

// Balance 查询钱包余额
func (that *MemberController) BalanceBatch(ctx *fasthttp.RequestCtx) {

	uids := string(ctx.PostArgs().Peek("uids"))
	if !validator.CheckStringCommaDigit(uids) {
		helper.Print(ctx, false, helper.UIDErr)
		return
	}

	s := strings.Split(uids, ",")
	balance, err := model.MemberBalanceBatch(s)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.PrintJson(ctx, true, balance)
}

func (that *MemberController) TagBatch(ctx *fasthttp.RequestCtx) {

	uids := string(ctx.PostArgs().Peek("uids"))
	if !validator.CheckStringCommaDigit(uids) {
		helper.Print(ctx, false, helper.UIDErr)
		return
	}

	s := strings.Split(uids, ",")
	balance, err := model.MemberBatchTag(s)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.PrintJson(ctx, true, balance)
}

func (that *MemberController) Insert(ctx *fasthttp.RequestCtx) {

	name := string(ctx.PostArgs().Peek("username"))
	password := string(ctx.PostArgs().Peek("password"))
	maintainName := string(ctx.PostArgs().Peek("maintain_name"))
	groupName := string(ctx.PostArgs().Peek("group_name"))
	agencyType := string(ctx.PostArgs().Peek("agency_type")) //391团队393普通
	remark := string(ctx.PostArgs().Peek("remark"))
	tester := string(ctx.PostArgs().Peek("tester"))

	tyTemp := string(ctx.PostArgs().Peek("ty"))
	zrTemp := string(ctx.PostArgs().Peek("zr"))
	qpTemp := string(ctx.PostArgs().Peek("qp"))
	djTemp := string(ctx.PostArgs().Peek("dj"))
	dzTemp := string(ctx.PostArgs().Peek("dz"))
	cpTemp := string(ctx.PostArgs().Peek("cp"))
	fcTemp := string(ctx.PostArgs().Peek("fc"))
	byTemp := string(ctx.PostArgs().Peek("by"))
	cgHighRebateTemp := string(ctx.PostArgs().Peek("cg_high_rebate"))
	cgOfficialRebateTemp := string(ctx.PostArgs().Peek("cg_official_rebate"))

	if len(maintainName) == 0 {
		maintainName = ""
	}

	vs := model.MemberRebateScale()
	ty, err := decimal.NewFromString(tyTemp)
	if err != nil || ty.IsNegative() || ty.GreaterThan(vs.TY) {
		helper.Print(ctx, false, helper.RebateOutOfRange)
	}

	zr, err := decimal.NewFromString(zrTemp)
	if err != nil || zr.IsNegative() || zr.GreaterThan(vs.ZR) {
		helper.Print(ctx, false, helper.RebateOutOfRange)
	}

	qp, err := decimal.NewFromString(qpTemp)
	if err != nil || qp.IsNegative() || qp.GreaterThan(vs.QP) {
		helper.Print(ctx, false, helper.RebateOutOfRange)
	}

	dj, err := decimal.NewFromString(djTemp)
	if err != nil || dj.IsNegative() || dj.GreaterThan(vs.DJ) {
		helper.Print(ctx, false, helper.RebateOutOfRange)
	}

	dz, err := decimal.NewFromString(dzTemp)
	if err != nil || dz.IsNegative() || dz.GreaterThan(vs.DZ) {
		helper.Print(ctx, false, helper.RebateOutOfRange)
	}

	cp, err := decimal.NewFromString(cpTemp)
	if err != nil || cp.IsNegative() || cp.GreaterThan(vs.CP) {
		helper.Print(ctx, false, helper.RebateOutOfRange)
	}

	fc, err := decimal.NewFromString(fcTemp)
	if err != nil || fc.IsNegative() || fc.GreaterThan(vs.FC) {
		helper.Print(ctx, false, helper.RebateOutOfRange)
	}

	by, err := decimal.NewFromString(byTemp)
	if err != nil || by.IsNegative() || by.GreaterThan(vs.BY) {
		helper.Print(ctx, false, helper.RebateOutOfRange)
	}

	nine := decimal.NewFromInt(9.00)
	cgHighRebate, err := decimal.NewFromString(cgHighRebateTemp)
	if err != nil || fc.IsNegative() ||
		cgHighRebate.GreaterThan(vs.CGHighRebate) ||
		nine.GreaterThan(cgHighRebate) {
		helper.Print(ctx, false, helper.RebateOutOfRange)
	}
	cgOfficialRebate, err := decimal.NewFromString(cgOfficialRebateTemp)
	if err != nil || fc.IsNegative() ||
		cgOfficialRebate.GreaterThan(vs.CGOfficialRebate) ||
		nine.GreaterThan(cgOfficialRebate) {
		helper.Print(ctx, false, helper.RebateOutOfRange)
	}

	name = strings.ToLower(name)
	if !validator.CheckUName(name, 5, 14) {
		helper.Print(ctx, false, helper.UsernameErr)
		return
	}

	if maintainName != "" && !validator.CtypeAlnum(maintainName) {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	if !validator.CheckUPassword(password, 8, 20) {
		helper.Print(ctx, false, helper.PasswordFMTErr)
		return
	}

	if agencyType != "391" && agencyType != "393" {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	if agencyType == "391" && len(groupName) < 1 {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	if tester != "0" {
		tester = "1"
	}

	mr := model.MemberRebate{
		TY:               ty.StringFixed(1),
		ZR:               zr.StringFixed(1),
		QP:               qp.StringFixed(1),
		DJ:               dj.StringFixed(1),
		DZ:               dz.StringFixed(1),
		CP:               cp.StringFixed(1),
		FC:               fc.StringFixed(1),
		BY:               by.StringFixed(1),
		CgOfficialRebate: cgOfficialRebate.StringFixed(2),
		CgHighRebate:     cgHighRebate.StringFixed(2),
	}
	createdAt := uint32(ctx.Time().Unix())

	// 添加下级代理
	err = model.MemberInsert(name, password, remark, maintainName, groupName, agencyType, tester, createdAt, mr)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// Balance 查询钱包余额
func (that *MemberController) Balance(ctx *fasthttp.RequestCtx) {

	username := string(ctx.QueryArgs().Peek("username"))
	username = strings.ToLower(username)
	if !validator.CheckUName(username, 5, 14) {
		helper.Print(ctx, false, helper.UsernameErr)
		return
	}

	balance, err := model.MemberBalance(username)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	data, err := jettison.Marshal(balance)
	if err != nil {
		helper.Print(ctx, false, helper.FormatErr)
		return
	}

	helper.PrintJson(ctx, true, string(data))
}

// 修改用户状态
func (that *MemberController) UpdateState(ctx *fasthttp.RequestCtx) {

	params := memberStateParam{}
	err := validator.Bind(ctx, &params)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	// 验证用户名
	params.Username = strings.ToLower(params.Username)
	names := strings.Split(params.Username, ",")
	for _, v := range names {
		if !validator.CheckUName(v, 5, 14) {
			helper.Print(ctx, false, helper.UsernameErr)
			return
		}
	}

	admin, err := model.AdminToken(ctx)
	if err != nil {
		helper.Print(ctx, false, helper.AccessTokenExpires)
		return
	}

	err = model.MemberRemarkInsert("", params.Remark, admin["name"], names, ctx.Time().Unix())
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	err = model.MemberUpdateState(names, params.State)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

/**
 * @Description: List 会员列表
 * @Author: parker
 * @Date: 2021/4/14 16:38
 * @LastEditTime: 2021/4/14 19:00
 * @LastEditors: parker
 */
func (that *MemberController) List(ctx *fasthttp.RequestCtx) {

	ty := ctx.PostArgs().GetUintOrZero("ty")                  //1 批量匹配
	username := string(ctx.PostArgs().Peek("username"))       //会员帐号
	realname := string(ctx.PostArgs().Peek("realname"))       //会员姓名
	phone := string(ctx.PostArgs().Peek("phone"))             //手机号
	agent := string(ctx.PostArgs().Peek("agent"))             //代理帐号
	tag := string(ctx.PostArgs().Peek("tag"))                 //会员标签
	state := ctx.PostArgs().GetUintOrZero("state")            //状态 0:全部,1:启用,2:禁用
	regStartTime := string(ctx.PostArgs().Peek("start_time")) //注册开始时间
	regEndTime := string(ctx.PostArgs().Peek("end_time"))     //注册结束时间
	email := string(ctx.PostArgs().Peek("email"))             //邮箱
	ipFlag := ctx.PostArgs().GetUintOrZero("ip_flag")         //1:最近登录ip,2:注册IP
	ip := string(ctx.PostArgs().Peek("ip"))                   //精确ip
	deviceFlag := ctx.PostArgs().GetUintOrZero("device_flag") //设备类型1:登录设备号,2:注册设备号
	device := string(ctx.PostArgs().Peek("device"))           //设备号
	page := ctx.PostArgs().GetUintOrZero("page")              //页码
	pageSize := ctx.PostArgs().GetUintOrZero("page_size")     //每页数量
	level := string(ctx.PostArgs().Peek("level"))             //设备号

	if ty != 1 {

		r := map[int]bool{
			1: true,
			2: true,
		}
		if state > 0 {
			if _, ok := r[state]; !ok {
				helper.Print(ctx, false, helper.ParamErr)
				return
			}
		}

		if ipFlag > 0 {
			if _, ok := r[ipFlag]; !ok {
				helper.Print(ctx, false, helper.ParamErr)
				return
			}
		}

		if deviceFlag > 0 {
			if _, ok := r[deviceFlag]; !ok {
				helper.Print(ctx, false, helper.ParamErr)
				return
			}
		}
	}

	if page < 1 {
		page = 1
	}

	if pageSize < 10 || pageSize > 200 {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	ex := g.Ex{}
	if username != "" {

		// 多个会员名用,分隔
		username = strings.ToLower(username)
		sName := strings.Split(username, ",")
		var usernames []string
		for _, name := range sName {
			if !validator.CheckUName(name, 5, 14) {
				helper.Print(ctx, false, helper.UsernameErr)
				return
			}

			usernames = append(usernames, name)
		}

		if ty == 0 && len(usernames) > 10 {
			ex["username"] = usernames[:10]
		} else {
			ex["username"] = usernames
		}

		data, err := model.MemberList(page, pageSize, "", "", "", ex)
		if err != nil {
			helper.Print(ctx, false, err.Error())
			return
		}

		helper.Print(ctx, true, data)
		return
	}

	if agent != "" {
		agent = strings.ToLower(agent)
		if !validator.CheckUName(agent, 5, 14) && len(username) > 50 {
			helper.Print(ctx, false, helper.AgentNameErr)
			return
		}

		ex["parent_name"] = agent
	}

	if ip != "" {

		if ipFlag == 2 {
			ex["regip"] = ip
		} else {
			ex["last_login_ip"] = ip
		}

	}

	if phone != "" {
		if !validator.IsVietnamesePhone(phone) {
			helper.Print(ctx, false, helper.PhoneFMTErr)
			return
		}

		ex["phone_hash"] = fmt.Sprintf("%d", model.MurmurHash(phone, 0))
	}

	if email != "" {
		if !strings.Contains(email, "@") {
			helper.Print(ctx, false, helper.EmailFMTErr)
			return
		}

		ex["email_hash"] = fmt.Sprintf("%d", model.MurmurHash(email, 0))
	}

	if state > 0 {
		ex["state"] = state
	}

	if realname != "" {
		ex["realname_hash"] = fmt.Sprintf("%d", model.MurmurHash(realname, 0))
	}

	if device != "" {
		// 最后登录设备号
		if deviceFlag == 1 {
			ex["last_login_device"] = device
		} else if deviceFlag == 2 { // 注册设备号
			ex["reg_device"] = device
		}
	}

	if level != "" {
		if len(level) == 1 && validator.CtypeDigit(level) {
			ex["level"] = level
		} else {
			levels := strings.Split(level, ",")
			ex["level"] = levels
		}
	}

	if level != "" {
		if len(level) == 1 && validator.CtypeDigit(level) {
			ex["level"] = level
		} else {
			levels := strings.Split(level, ",")
			ex["level"] = levels
		}
	}

	data, err := model.MemberList(page, pageSize, tag, regStartTime, regEndTime, ex)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}

func (that *MemberController) Agency(ctx *fasthttp.RequestCtx) {

	username := string(ctx.PostArgs().Peek("agency_name"))       //会员帐号
	maintainName := string(ctx.PostArgs().Peek("maintain_name")) //维护人
	groupName := string(ctx.PostArgs().Peek("group_name"))       //团队名称
	state := ctx.PostArgs().GetUintOrZero("state")               //状态 0:全部,1:启用,2:禁用
	regStartTime := string(ctx.PostArgs().Peek("start_time"))    //注册开始时间
	regEndTime := string(ctx.PostArgs().Peek("end_time"))        //注册结束时间
	page := ctx.PostArgs().GetUintOrZero("page")                 //页码
	pageSize := ctx.PostArgs().GetUintOrZero("page_size")        //每页数量
	parentID := string(ctx.PostArgs().Peek("uid"))
	sortField := string(ctx.PostArgs().Peek("sort_field"))
	isAsc := ctx.PostArgs().GetUintOrZero("is_asc")
	agencyType := string(ctx.PostArgs().Peek("agency_type"))

	if page < 1 {
		page = 1
	}

	if pageSize < 10 || pageSize > 200 {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	var press = exp.NewExpressionList(exp.AndType, g.C("uid").Eq(g.C("top_uid")))
	if username != "" {
		username = strings.ToLower(username)
		if !validator.CheckUName(username, 5, 14) && username != "root" {
			helper.Print(ctx, false, helper.UsernameErr)
			return
		}
		press = exp.NewExpressionList(exp.AndType, g.C("username").Eq(username))
	}

	if parentID != "" {
		press = exp.NewExpressionList(exp.AndType, g.Or(g.C("parent_uid").Eq(parentID), g.C("uid").Eq(parentID)))
	}

	if state > 0 {
		press = press.Append(g.C("state").Eq(state))
	}

	if maintainName != "" {
		press = press.Append(g.C("maintain_name").Eq(maintainName))
	}

	if agencyType == "391" && groupName != "" {
		press = press.Append(g.C("group_name").Eq(groupName))
	}

	if sortField != "" {
		sortFields := map[string]bool{
			"deposit_amount":     true,
			"withdrawal_amount":  true,
			"dividend_amount":    true,
			"valid_bet_amount":   true,
			"rebate_amount":      true,
			"company_net_amount": true,
			"rebate_point":       true,
		}

		if _, ok := sortFields[sortField]; !ok {
			helper.Print(ctx, false, helper.ParamErr)
			return
		}

		if !validator.CheckIntScope(strconv.Itoa(isAsc), 0, 1) {
			helper.Print(ctx, false, helper.ParamErr)
			return
		}
	}

	data, err := model.AgencyList(press, parentID, username, regStartTime, regEndTime, sortField, isAsc, page, pageSize, agencyType)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}

// 修改会员信息
func (that *MemberController) Update(ctx *fasthttp.RequestCtx) {

	tagsID := string(ctx.PostArgs().Peek("tags_id"))
	realname := string(ctx.PostArgs().Peek("realname"))
	username := string(ctx.PostArgs().Peek("username"))
	birth := string(ctx.PostArgs().Peek("birth"))
	phone := string(ctx.PostArgs().Peek("phone"))
	email := string(ctx.PostArgs().Peek("email"))
	zalo := string(ctx.PostArgs().Peek("zalo"))
	address := string(ctx.PostArgs().Peek("address")) //收货地址
	address, _ = url.QueryUnescape(address)

	username = strings.ToLower(username)
	if !validator.CheckUName(username, 5, 14) {
		helper.Print(ctx, false, helper.UsernameErr)
		return
	}

	param := map[string]string{}
	if birth != "" {
		t, err := time.Parse("2006-01-02", birth)
		if err != nil {
			helper.Print(ctx, false, helper.TimeTypeErr)
			return
		}

		param["birth"] = fmt.Sprintf("%d", t.Unix())
		param["birth_hash"] = fmt.Sprintf("%d", model.MurmurHash(birth, 0))
	}

	if realname != "" {
		if helper.CtypePunct(param["realname"]) {
			helper.Print(ctx, false, helper.RealNameFMTErr)
			return
		}

		param["realname"] = realname
	}

	if phone != "" {
		if !validator.IsVietnamesePhone(phone) {
			helper.Print(ctx, false, helper.PhoneFMTErr)
			return
		}

		param["phone"] = phone
	}

	if email != "" {
		if !strings.Contains(email, "@") {
			helper.Print(ctx, false, helper.EmailFMTErr)
			return
		}

		param["email"] = email
	}

	if zalo != "" {
		if !helper.CtypeDigit(zalo) {
			helper.Print(ctx, false, helper.ZaloFMTErr)
			return
		}

		param["zalo"] = zalo
	}

	if address != "" {
		if len(strings.Split(address, "|")) != 4 {
			helper.Print(ctx, false, helper.AddressFMTErr)
			return
		}

		param["address"] = validator.FilterInjection(address)
	}

	var userTagsId []string
	if tagsID != "" {
		if !validator.CheckStringCommaDigit(tagsID) {
			helper.Print(ctx, false, helper.UserTagErr)
			return
		}

		userTagsId = strings.Split(tagsID, ",")
	}

	admin, err := model.AdminToken(ctx)
	if err != nil {
		helper.Print(ctx, false, helper.AccessTokenExpires)
		return
	}

	//fmt.Println("param = ", param)
	err = model.MemberUpdate(username, admin["id"], param, userTagsId)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// 会员管理-会员列表-解除密码限制/解除短信限制/场馆钱包限制
func (that *MemberController) RetryReset(ctx *fasthttp.RequestCtx) {

	param := retryResetParam{}
	err := validator.Bind(ctx, &param)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	if param.Ty == model.WALLET {
		if !validator.CtypeDigit(param.Pid) {
			helper.Print(ctx, false, helper.PlatIDErr)
			return
		}
	}

	err = model.MemberRetryReset(param.Username, param.Ty, param.Pid)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// 会员备注添加
func (that *MemberController) RemarkLogInsert(ctx *fasthttp.RequestCtx) {

	params := remarkLogParams{}
	err := validator.Bind(ctx, &params)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	if !validator.CheckStringLength(params.Msg, 1, 300) {
		helper.Print(ctx, false, helper.ContentLengthErr)
		return
	}

	if params.File != "" {
		if len(params.File) < 5 {
			helper.Print(ctx, false, helper.FileURLErr)
			return
		}
	}

	// 验证用户名
	params.Username = strings.ToLower(params.Username)
	names := strings.Split(params.Username, ",")
	for _, v := range names {
		if !validator.CheckUName(v, 5, 14) {
			helper.Print(ctx, false, helper.UsernameErr)
			return
		}
	}

	admin, err := model.AdminToken(ctx)
	if err != nil {
		helper.Print(ctx, false, helper.AccessTokenExpires)
		return
	}

	err = model.MemberRemarkInsert(params.File, params.Msg, admin["name"], names, ctx.Time().Unix())
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// 会员备注修改
func (that *MemberController) RemarkLogUpdate(ctx *fasthttp.RequestCtx) {

	ts := string(ctx.PostArgs().Peek("ts"))
	file := string(ctx.PostArgs().Peek("file"))
	msg := string(ctx.PostArgs().Peek("msg"))

	if !validator.CheckStringLength(msg, 1, 300) {
		helper.Print(ctx, false, helper.ContentLengthErr)
		return
	}

	if file != "" && len(file) < 5 {
		helper.Print(ctx, false, helper.FileURLErr)
		return
	}

	admin, err := model.AdminToken(ctx)
	if err != nil {
		helper.Print(ctx, false, helper.AccessTokenExpires)
		return
	}

	err = model.MemberRemarkUpdate(ts, file, msg, admin["name"], ctx.Time().Unix())
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// 会员备注删除
func (that *MemberController) RemarkLogDelete(ctx *fasthttp.RequestCtx) {

	ts := string(ctx.PostArgs().Peek("ts"))
	admin, err := model.AdminToken(ctx)
	if err != nil {
		helper.Print(ctx, false, helper.AccessTokenExpires)
		return
	}

	err = model.MemberRemarkDelete(ts, admin["name"], ctx.Time().Unix())
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// 会员管理-会员列表-数据概览
func (that MemberController) Overview(ctx *fasthttp.RequestCtx) {

	username := string(ctx.QueryArgs().Peek("username"))
	startTime := string(ctx.QueryArgs().Peek("start_time"))
	endTime := string(ctx.QueryArgs().Peek("end_time"))

	username = strings.ToLower(username)
	if !validator.CheckUName(username, 5, 14) {
		helper.Print(ctx, false, helper.UsernameErr)
		return
	}

	_, err := time.Parse("2006-01-02 15:04:05", startTime)
	if err != nil {
		helper.Print(ctx, false, helper.DateTimeErr)
		return
	}

	_, err = time.Parse("2006-01-02 15:04:05", endTime)
	if err != nil {
		helper.Print(ctx, false, helper.DateTimeErr)
		return
	}

	data, err := model.MemberDataOverview(username, startTime, endTime)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}

// 会员日志列表
func (that *MemberController) RemarkLogList(ctx *fasthttp.RequestCtx) {

	uid := string(ctx.QueryArgs().Peek("uid"))
	adminName := string(ctx.QueryArgs().Peek("admin_name"))
	startTime := string(ctx.QueryArgs().Peek("start_time"))
	endTime := string(ctx.QueryArgs().Peek("end_time"))
	sPage := string(ctx.QueryArgs().Peek("page"))
	sPageSize := string(ctx.QueryArgs().Peek("page_size"))

	if !validator.CheckStringDigit(uid) {
		helper.Print(ctx, false, helper.UIDErr)
		return
	}

	if !validator.CheckStringDigit(sPage) || !validator.CheckStringDigit(sPageSize) {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	if adminName != "" && !validator.CheckAName(adminName, 5, 20) {
		helper.Print(ctx, false, helper.AdminNameErr)
		return
	}

	if startTime != "" {
		_, err := time.Parse("2006-01-02 15:04:05", startTime)
		if err != nil {
			helper.Print(ctx, false, helper.DateTimeErr)
			return
		}
	}

	if endTime != "" {
		_, err := time.Parse("2006-01-02 15:04:05", endTime)
		if err != nil {
			helper.Print(ctx, false, helper.DateTimeErr)
			return
		}
	}

	page, _ := strconv.Atoi(sPage)
	pageSize, _ := strconv.Atoi(sPageSize)
	data, err := model.MemberRemarkLogList(uid, adminName, startTime, endTime, page, pageSize)

	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}

func (that *MemberController) UpdatePwd(ctx *fasthttp.RequestCtx) {

	username := string(ctx.PostArgs().Peek("username"))
	pwd := string(ctx.PostArgs().Peek("pwd"))
	ty := ctx.PostArgs().GetUintOrZero("ty")
	if username == "" || pwd == "" {
		helper.Print(ctx, false, helper.ParamNull)
		return
	}

	if ty > 1 {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	// 会员名校验
	username = strings.ToLower(username)
	if !validator.CheckUName(username, 5, 14) {
		helper.Print(ctx, false, helper.UsernameErr)
		return
	}

	// 会员密码校验
	if !validator.CheckUPassword(pwd, 8, 20) {
		helper.Print(ctx, false, helper.PasswordFMTErr)
		return
	}

	err := model.MemberUpdatePwd(username, pwd, ty, ctx)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// History 查询用户真实姓名/邮箱/手机号/银行卡号修改历史
func (that *MemberController) History(ctx *fasthttp.RequestCtx) {

	id := string(ctx.PostArgs().Peek("id"))
	field := string(ctx.PostArgs().Peek("field"))
	encrypt := ctx.PostArgs().GetBool("encrypt")

	if !validator.CheckStringDigit(id) {
		helper.Print(ctx, false, helper.IDErr)
		return
	}

	if _, ok := model.MemberHistoryField[field]; !ok {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	data, err := model.MemberHistory(id, field, encrypt)
	if err != nil {
		helper.Print(ctx, false, helper.ServerErr)
		return
	}

	helper.Print(ctx, true, data)
}

// History 查询用户真实姓名/邮箱/手机号/银行卡号修改历史
func (that *MemberController) HistoryField(ctx *fasthttp.RequestCtx) {

	id := string(ctx.PostArgs().Peek("id"))
	encrypt := ctx.PostArgs().GetBool("encrypt")
	field := ctx.UserValue("field").(string)

	if _, ok := model.MemberHistoryField[field]; !ok {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	if field == "bankcard" {
		field = string(ctx.PostArgs().Peek("field"))
	}

	if !validator.CheckStringDigit(id) {
		helper.Print(ctx, false, helper.IDErr)
		return
	}

	data, err := model.MemberHistory(id, field, encrypt)
	if err != nil {
		helper.Print(ctx, false, helper.ServerErr)
		return
	}

	helper.Print(ctx, true, data)
}

// Full 查询用户真实姓名/邮箱/手机号/银行卡号明文信息
func (that *MemberController) Full(ctx *fasthttp.RequestCtx) {

	id := string(ctx.PostArgs().Peek("id"))
	field := string(ctx.PostArgs().Peek("field"))
	if !validator.CheckStringDigit(id) {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	if _, ok := model.MemberHistoryField[field]; !ok {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	data, err := model.MemberFull(id, []string{field})
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	//fmt.Println("grpc_t.Decrypt data = ", data)
	//fmt.Println("grpc_t.Decrypt field = ", field)

	helper.Print(ctx, true, data[field])
}

// Full 查询用户真实姓名/邮箱/手机号/银行卡号明文信息
func (that *MemberController) FullField(ctx *fasthttp.RequestCtx) {

	id := string(ctx.PostArgs().Peek("id"))
	field := ctx.UserValue("field").(string)
	if _, ok := model.MemberHistoryField[field]; !ok {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	if field == "bankcard" {
		field = string(ctx.PostArgs().Peek("field"))
	}
	if !validator.CheckStringDigit(id) {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	data, err := model.MemberFull(id, []string{field})
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	//fmt.Println("grpc_t.Decrypt data = ", data)
	//fmt.Println("grpc_t.Decrypt field = ", field)

	helper.Print(ctx, true, data[field])
}

func (that *MemberController) SetBalanceZero(ctx *fasthttp.RequestCtx) {

	username := string(ctx.PostArgs().Peek("username"))
	remark := string(ctx.PostArgs().Peek("remark"))

	username = strings.ToLower(username)
	if !validator.CheckUName(username, 5, 14) {
		helper.Print(ctx, false, helper.UsernameErr)
		return
	}

	if remark == "" {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	admin, err := model.AdminToken(ctx)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	err = model.MemberBalanceZero(username, remark, admin["id"], admin["name"])
	if err != nil {
		helper.Print(ctx, false, err.Error())
	}

	helper.Print(ctx, true, "success")
}

// UpdateTopMember 修改密码以及返水比例
func (that *MemberController) UpdateTopMember(ctx *fasthttp.RequestCtx) {

	username := string(ctx.PostArgs().Peek("username"))
	password := string(ctx.PostArgs().Peek("password"))
	remarks := string(ctx.PostArgs().Peek("remarks"))
	groupName := string(ctx.PostArgs().Peek("group_name"))
	state := ctx.PostArgs().GetUintOrZero("state") // 状态 1正常 2禁用

	username = strings.ToLower(username)
	if !validator.CheckUName(username, 5, 14) && username != "root" {
		helper.Print(ctx, false, helper.UsernameErr)
		return
	}

	mb, err := model.MemberFindOne(username)
	if err != nil {
		helper.Print(ctx, false, helper.UsernameErr)
		return
	}

	recd := g.Record{}
	if password != "" {
		if !validator.CheckUPassword(password, 8, 20) {
			helper.Print(ctx, false, helper.PasswordFMTErr)
			return
		}
		recd["password"] = fmt.Sprintf("%d", model.MurmurHash(password, mb.CreatedAt))
	}

	if state != 0 {
		if state > 2 || state < 1 {
			helper.Print(ctx, false, helper.PasswordFMTErr)
			return
		}
		recd["state"] = state
	}

	if remarks != "" {
		recd["remarks"] = remarks
	}

	if groupName != "" {
		if len(groupName) > 50 {
			helper.Print(ctx, false, helper.ParamErr)
			return
		}
		recd["group_name"] = groupName
	}

	if len(recd) == 0 {
		helper.Print(ctx, false, helper.NoDataUpdate)
		return
	}

	// 更新代理
	err = model.MemberUpdateInfo(mb.UID, recd)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// UpdateMemberRebate 修改返水比例
func (that *MemberController) UpdateMemberRebate(ctx *fasthttp.RequestCtx) {

	username := string(ctx.PostArgs().Peek("username"))
	tyTemp := string(ctx.PostArgs().Peek("ty"))
	zrTemp := string(ctx.PostArgs().Peek("zr"))
	qpTemp := string(ctx.PostArgs().Peek("qp"))
	djTemp := string(ctx.PostArgs().Peek("dj"))
	dzTemp := string(ctx.PostArgs().Peek("dz"))
	cpTemp := string(ctx.PostArgs().Peek("cp"))
	fcTemp := string(ctx.PostArgs().Peek("fc"))
	byTemp := string(ctx.PostArgs().Peek("by"))
	cgHighRebateTemp := string(ctx.PostArgs().Peek("cg_high_rebate"))
	cgOfficialRebateTemp := string(ctx.PostArgs().Peek("cg_official_rebate"))

	username = strings.ToLower(username)
	if !validator.CheckUName(username, 5, 14) && username != "root" {
		helper.Print(ctx, false, helper.UsernameErr)
		return
	}

	mb, err := model.MemberFindOne(username)
	if err != nil {
		helper.Print(ctx, false, helper.UsernameErr)
		return
	}

	ty, err := decimal.NewFromString(tyTemp) //下级会员体育返水比例
	if err != nil || ty.IsNegative() {
		fmt.Println("sty = ", tyTemp)
		helper.Print(ctx, false, helper.RebateOutOfRange)
		return
	}

	zr, err := decimal.NewFromString(zrTemp) //下级会员真人返水比例
	if err != nil || zr.IsNegative() {
		fmt.Println("szr = ", zrTemp)
		helper.Print(ctx, false, helper.RebateOutOfRange)
		return
	}

	qp, err := decimal.NewFromString(qpTemp) //下级会员棋牌返水比例
	if err != nil || qp.IsNegative() {
		fmt.Println("sqp = ", qpTemp)
		helper.Print(ctx, false, helper.RebateOutOfRange)
		return
	}

	dj, err := decimal.NewFromString(djTemp) //下级会员电竞返水比例
	if err != nil || dj.IsNegative() {
		fmt.Println("sdj = ", djTemp)
		helper.Print(ctx, false, helper.RebateOutOfRange)
		return
	}

	dz, err := decimal.NewFromString(dzTemp) //下级会员电子返水比例
	if err != nil || dz.IsNegative() {
		fmt.Println("sdz = ", dzTemp)
		helper.Print(ctx, false, helper.RebateOutOfRange)
		return
	}

	cp, err := decimal.NewFromString(cpTemp) //下级会员彩票返水比例
	if err != nil || cp.IsNegative() {
		fmt.Println("scp = ", cpTemp)
		helper.Print(ctx, false, helper.RebateOutOfRange)
		return
	}

	fc, err := decimal.NewFromString(fcTemp) //下级会员斗鸡返水比例
	if err != nil || fc.IsNegative() {
		fmt.Println("sfc = ", fcTemp)
		helper.Print(ctx, false, helper.RebateOutOfRange)
		return
	}

	by, err := decimal.NewFromString(byTemp) //下级会员捕鱼返水比例
	if err != nil || by.IsNegative() {
		fmt.Println("sby = ", byTemp)
		helper.Print(ctx, false, helper.RebateOutOfRange)
		return
	}

	nine := decimal.NewFromFloat(9.0)
	cgHighRebate, err := decimal.NewFromString(cgHighRebateTemp) //下级最高cg高频彩返点
	if err != nil || cgHighRebate.LessThan(nine) {
		fmt.Println("cgHighRebateTemp = ", cgHighRebateTemp)
		helper.Print(ctx, false, helper.RebateOutOfRange)
		return
	}

	cgOfficialRebate, err := decimal.NewFromString(cgOfficialRebateTemp) //下级最高cg官方彩返点
	if err != nil || cgOfficialRebate.LessThan(nine) {
		fmt.Println("cgOfficialRebateTemp = ", cgOfficialRebateTemp)
		helper.Print(ctx, false, helper.RebateOutOfRange)
		return
	}

	mr := model.MemberRebateResult_t{
		TY:               ty,
		ZR:               zr,
		QP:               qp,
		DJ:               dj,
		DZ:               dz,
		CP:               cp,
		FC:               fc,
		BY:               by,
		CGHighRebate:     cgHighRebate,
		CGOfficialRebate: cgOfficialRebate,
	}
	// 获取下级最高返点
	maxSubRebate, err := model.MemberMaxRebateFindOne(mb.UID)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	// 低于下级最高返点
	if ok := model.MemberRebateCmp(maxSubRebate, mr); !ok {
		helper.Print(ctx, false, helper.RebateOutOfRange)
		return
	}

	// 非总代的返水比例判断，不能高于上级
	if mb.ParentUid != "0" && mb.ParentUid != "" {
		ParentRebate, err := model.MemberParentRebate(mb.ParentUid)
		if err != nil {
			helper.Print(ctx, false, err.Error())
			return
		}

		// 高于上级返水比例
		if ok := model.MemberRebateCmp(mr, ParentRebate); !ok {
			helper.Print(ctx, false, helper.RebateOutOfRange)
			return
		}
	} else { //总代
		maxScale := model.MemberRebateScale()
		// 高于最高返水比例
		if ok := model.MemberRebateCmp(mr, maxScale); !ok {
			helper.Print(ctx, false, helper.RebateOutOfRange)
			return
		}
	}

	// 更新代理返水比例
	err = model.MemberUpdateRebateInfo(mb.UID, mr)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// UpdateMaintainName 修改维护人
func (that *MemberController) UpdateMaintainName(ctx *fasthttp.RequestCtx) {

	username := string(ctx.PostArgs().Peek("username"))
	maintainName := string(ctx.PostArgs().Peek("maintain_name"))

	username = strings.ToLower(username)
	if !validator.CheckUName(username, 5, 14) {
		helper.Print(ctx, false, helper.UsernameErr)
		return
	}

	if !validator.CtypeAlnum(maintainName) {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	mb, err := model.MemberFindOne(username)
	if err != nil {
		helper.Print(ctx, false, helper.UsernameErr)
		return
	}

	// 更新代理
	err = model.MemberUpdateMaintainName(mb.UID, maintainName)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

func (that *MemberController) MemberList(ctx *fasthttp.RequestCtx) {

	param := model.MemberListParam{}
	err := validator.Bind(ctx, &param)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	if param.Username != "" {
		param.Username = strings.ToLower(param.Username)
		if !validator.CheckUName(param.Username, 5, 14) {
			helper.Print(ctx, false, helper.UsernameErr)
			return
		}
	}

	if param.ParentName != "" {
		param.ParentName = strings.ToLower(param.ParentName)
		if !validator.CheckUName(param.ParentName, 5, 14) {
			helper.Print(ctx, false, helper.AgentNameErr)
			return
		}
	}

	data, err := model.AgencyMemberList(param)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}
