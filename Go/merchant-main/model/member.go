package model

import (
	"database/sql"
	"errors"
	"fmt"
	"merchant/contrib/helper"
	"merchant/contrib/session"
	"strconv"
	"strings"
	"time"

	g "github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/shopspring/decimal"
	"github.com/valyala/fasthttp"
	"github.com/wI2L/jettison"
)

var (
	PWD            = uint8(1) // 密码
	SMS            = uint8(2) // 短信
	WALLET         = uint8(3) // 钱包
	memberUnfreeze = map[uint8]string{
		PWD:    "MPE:%s",    // 密码尝试次数
		SMS:    "smsfgt:%s", // 短信发送次数
		WALLET: "P:%s:%s",   // 钱包限制
	}
	MemberHistoryField = map[string]bool{
		"realname": true, // 会员真实姓名
		"phone":    true, // 会员手机号
		"email":    true, // 会员邮箱
		"bankcard": true, // 会员银行卡查询
		"zalo":     true, //zalo
	}
)

type PlatBalance struct {
	ID      string `db:"id" json:"id"`
	Balance string `db:"balance" json:"balance"`
}

type memberTags struct {
	Uid     string `db:"uid" json:"uid"`
	TagId   int64  `db:"tag_id" json:"tags_id"`
	TagName string `db:"tag_name" json:"tag_name"`
}

type tag struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type memberTag struct {
	Uid  string `db:"uid" json:"uid"`
	Tags []tag  `json:"tags"`
}

/*
会员银行卡 校验日志数据概览
*/
type BankcardLog struct {
	Ts        string `db:"ts" json:"ts"`
	Username  string `db:"username" json:"username"`
	Uid       string `db:"uid" json:"uid"`
	BankName  string `db:"bank_name" json:"bank_name"`
	BankNo    string `db:"bankcard_no" json:"bankcard_no"`
	RealName  string `db:"realname" json:"realname"`
	Ip        string `db:"ip" json:"ip"`
	Status    int    `db:"status" json:"status"`
	Device    int    `db:"device" json:"device"`
	CreatedAt uint32 `db:"created_at" json:"created_at"`
}

type BankcardLogData struct {
	D []BankcardLog `json:"d"`
	T int64         `json:"t"`
	S uint          `json:"s"`
}

// MemberDataOverviewData 会员管理-会员列表-数据概览 response structure
type MemberDataOverviewData struct {
	NetAmount      float64 `json:"net_amount"`       // 总输赢
	ValidBetAmount float64 `json:"valid_bet_amount"` // 总有效投注
	Deposit        float64 `json:"deposit"`          // 总存款
	Withdraw       float64 `json:"withdraw"`         // 总提款
	Dividend       float64 `json:"dividend"`         // 总红利
	Rebate         float64 `json:"rebate"`           // 总返水
}

// MemberListData 会员列表
type MemberListData struct {
	T         int                      `json:"t"`
	S         int                      `json:"s"`
	EnableMod bool                     `json:"enable_mod"`
	D         []MemberListCol          `json:"d"`
	Agg       map[string]MemberAggData `json:"agg"`
	Info      map[string]memberInfo    `json:"info"`
}

type memberInfo struct {
	UID          string `db:"uid" json:"uid"`
	AgencyType   string `json:"agency_type" db:"agency_type"`
	Tester       string `json:"tester" db:"tester"`
	Username     string `db:"username" json:"username"`           //会员名
	LastLoginIp  string `db:"last_login_ip" json:"last_login_ip"` //最后登陆ip
	TopUid       string `db:"top_uid" json:"top_uid"`             //总代uid
	TopName      string `db:"top_name" json:"top_name"`           //总代代理
	ParentUid    string `db:"parent_uid" json:"parent_uid"`       //上级uid
	ParentName   string `db:"parent_name" json:"parent_name"`     //上级代理
	State        uint8  `db:"state" json:"state"`                 //状态 1正常 2禁用
	Remarks      string `db:"remarks" json:"remarks"`             //备注
	MaintainName string `db:"maintain_name" json:"maintain_name"` //
	GroupName    string `db:"group_name" json:"group_name"`       //团队名称
}

type MemberListCol struct {
	UID              string  `json:"uid" db:"uid"`
	DepositAmount    float64 `json:"deposit_amount" db:"deposit_amount"`
	WithdrawAmount   float64 `json:"withdrawal_amount" db:"withdrawal_amount"`
	ValidBetAmount   float64 `json:"valid_bet_amount" db:"valid_bet_amount"`
	RebateAmount     float64 `json:"rebate_amount" db:"rebate_amount"`
	RebatePoint      float64 `json:"rebate_point" db:"rebate_point"`
	CompanyNetAmount float64 `json:"company_net_amount" db:"company_net_amount"`
	TY               string  `json:"ty" db:"ty"`
	ZR               string  `json:"zr" db:"zr"`
	QP               string  `json:"qp" db:"qp"`
	DJ               string  `json:"dj" db:"dj"`
	DZ               string  `json:"dz" db:"dz"`
	CP               string  `json:"cp" db:"cp"`
	FC               string  `json:"fc" db:"fc"`
	BY               string  `json:"by" db:"by"`
	CgOfficialRebate string  `json:"cg_official_rebate" db:"cg_official_rebate"` //CG官方彩返点
	CgHighRebate     string  `json:"cg_high_rebate" db:"cg_high_rebate"`         //CG高频彩返点
	Lvl              int     `json:"lvl" db:"-"`
	PlanID           string  `json:"plan_id" db:"-"`
	PlanName         string  `json:"plan_name" db:"-"`
	MemCount         int64   `json:"user_count" db:"mem_count"`
}

type MemberAggData struct {
	MemCount       int    `db:"mem_count" json:"mem_count"`
	RegistCountNew int    `db:"regist_count" json:"regist_count"`
	ActiveCount    int    `db:"active_count" json:"active_count"`
	UID            string `db:"uid" json:"uid"`
}

type MemberRebateResult_t struct {
	ZR               decimal.Decimal
	QP               decimal.Decimal
	TY               decimal.Decimal
	DZ               decimal.Decimal
	DJ               decimal.Decimal
	CP               decimal.Decimal
	FC               decimal.Decimal
	BY               decimal.Decimal
	CGOfficialRebate decimal.Decimal
	CGHighRebate     decimal.Decimal
}

func MemberInsert(username, password, remark, maintainName, groupName, agencyType, tester string, createdAt uint32, mr MemberRebate) error {

	userName := strings.ToLower(username)
	if MemberExist(userName) {
		return errors.New(helper.UsernameExist)
	}

	uid := helper.GenId()
	mr.UID = uid
	mr.CreatedAt = createdAt
	mr.ParentUID = "0"
	mr.Prefix = meta.Prefix
	agtype, _ := strconv.ParseInt(agencyType, 10, 64)
	m := Member{
		UID:                 uid,
		Username:            userName,
		Password:            fmt.Sprintf("%d", MurmurHash(password, createdAt)),
		Prefix:              meta.Prefix,
		Birth:               "0",
		BirthHash:           "0",
		State:               1,
		CreatedAt:           createdAt,
		LastLoginIp:         "",
		LastLoginAt:         createdAt,
		LastLoginDevice:     "",
		LastLoginSource:     0,
		ParentUid:           "0",
		TopUid:              uid,
		TopName:             userName,
		FirstDepositAmount:  "0.000",
		FirstBetAmount:      "0.000",
		Balance:             "0.000",
		LockAmount:          "0.000",
		Commission:          "0.000",
		SecondDepositAmount: "0.000",
		Remarks:             remark,
		MaintainName:        maintainName,
		GroupName:           groupName,
		AgencyType:          agtype,
		EmailHash:           "0",
		RealnameHash:        "0",
		PhoneHash:           "0",
		ZaloHash:            "0",
		Level:               1,
		Tester:              tester,
	}

	tx, err := meta.MerchantDB.Begin() // 开启事务
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	query, _, _ := dialect.Insert("tbl_members").Rows(&m).ToSQL()
	_, err = tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		//fmt.Println("tbl_members query = ", query)
		//fmt.Println("tbl_members err = ", err.Error())
		return pushLog(err, helper.DBErr)
	}

	query, _, _ = dialect.Insert("tbl_member_rebate_info").Rows(&mr).ToSQL()
	_, err = tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		//fmt.Println("tbl_member_rebate_info query = ", query)
		//fmt.Println("tbl_member_rebate_info err = ", err.Error())
		return pushLog(err, helper.DBErr)
	}

	treeNode := MemberClosureInsert(uid, "0")
	_, err = tx.Exec(treeNode)
	if err != nil {
		_ = tx.Rollback()
		//fmt.Println("MemberClosureInsert err = ", err.Error())
		return pushLog(fmt.Errorf("sql : %s, error : %s", treeNode, err.Error()), helper.DBErr)
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("tx.Commit err = ", err.Error())
		return pushLog(err, helper.DBErr)
	}

	_ = MemberRebateUpdateCache2(uid, mr)
	_ = MemberUpdateCache(uid, "")
	_, err = session.Set([]byte(m.Username), m.UID)
	if err != nil {
		return errors.New(helper.SessionErr)
	}

	return nil
}

func MemberFindByUid(uid string) (Member, error) {

	m := Member{}

	ex := g.Ex{
		"uid":    uid,
		"prefix": meta.Prefix,
	}
	query, _, _ := dialect.From("tbl_members").Select(colsMember...).Where(ex).Limit(1).ToSQL()
	err := meta.MerchantDB.Get(&m, query)
	if err != nil && err != sql.ErrNoRows {
		return m, pushLog(err, helper.DBErr)
	}

	if err == sql.ErrNoRows {
		return m, errors.New(helper.UsernameErr)
	}

	return m, nil
}

func MemberUpdateCache(uid, username string) error {

	var (
		err error
		dst Member
	)

	if helper.CtypeDigit(uid) {
		dst, err = MemberFindByUid(uid)
		if err != nil {
			return err
		}
	} else {
		dst, err = MemberFindOne(username)
		if err != nil {
			return err
		}
	}

	key := meta.Prefix + ":member:" + dst.Username
	fields := []interface{}{"uid", dst.UID, "username", dst.Username, "password", dst.Password, "birth", dst.Birth, "birth_hash", dst.BirthHash, "realname_hash", dst.RealnameHash, "email_hash", dst.EmailHash, "phone_hash", dst.PhoneHash, "zalo_hash", dst.ZaloHash, "prefix", dst.Prefix, "tester", dst.Tester, "withdraw_pwd", dst.WithdrawPwd, "regip", dst.Regip, "reg_device", dst.RegDevice, "reg_url", dst.RegUrl, "created_at", dst.CreatedAt, "last_login_ip", dst.LastLoginIp, "last_login_at", dst.LastLoginAt, "source_id", dst.SourceId, "first_deposit_at", dst.FirstDepositAt, "first_deposit_amount", dst.FirstDepositAmount, "first_bet_at", dst.FirstBetAt, "first_bet_amount", dst.FirstBetAmount, "", dst.SecondDepositAt, "", dst.SecondDepositAmount, "top_uid", dst.TopUid, "top_name", dst.TopName, "parent_uid", dst.ParentUid, "parent_name", dst.ParentName, "bankcard_total", dst.BankcardTotal, "last_login_device", dst.LastLoginDevice, "last_login_source", dst.LastLoginSource, "remarks", dst.Remarks, "state", dst.State, "level", dst.Level, "balance", dst.Balance, "lock_amount", dst.LockAmount, "commission", dst.Commission, "group_name", dst.GroupName, "agency_type", dst.AgencyType, "address", dst.Address, "avatar", dst.Avatar}

	pipe := meta.MerchantRedis.TxPipeline()
	pipe.Del(ctx, key)
	pipe.HMSet(ctx, key, fields...)
	pipe.Persist(ctx, key)
	pipe.Exec(ctx)
	pipe.Close()

	// 禁用
	if dst.State == 2 {
		session.Offline([]string{dst.UID})
	}

	return nil
}

/**
 * @Description: Transfer 会员列表-帐户信息
 * @Author: parker
 * @Date: 2021/4/7 10:43
 * @LastEditTime: 2021/4/7 10:43
 * @LastEditors: parker
 */
func MemberAccountInfo(username string) ([]PlatBalance, error) {

	var data []PlatBalance

	mb, err := MemberFindOne(username)
	if err != nil || len(mb.Username) == 0 {
		return data, errors.New(helper.UsernameErr)
	}

	//添加中心钱包余额
	data = append(data, PlatBalance{ID: "1", Balance: mb.Balance})
	data = append(data, memberPlatformBalance(username)...)

	return data, nil
}

func MemberExist(username string) bool {

	var uid uint64
	query, _, _ := dialect.From("tbl_members").Select("uid").Where(g.Ex{"username": username, "prefix": meta.Prefix}).ToSQL()
	err := meta.MerchantDB.Get(&uid, query)
	if err == nil && uid != 0 {
		return true
	}

	return false
}

// 批量获取会员标签
func MemberBatchTag(uids []string) (string, error) {

	var tags []memberTags

	t := dialect.From("tbl_member_tags")
	query, _, _ := t.Select("uid", "tag_id", "tag_name").Where(g.Ex{"uid": uids, "prefix": meta.Prefix}).ToSQL()
	err := meta.MerchantDB.Select(&tags, query)
	if err != nil {
		return "", pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	var result []memberTag
	for _, id := range uids {
		item := memberTag{Uid: id, Tags: []tag{}}
		for _, v := range tags {
			if strings.EqualFold(v.Uid, id) {
				item.Tags = append(item.Tags, tag{ID: v.TagId, Name: v.TagName})
			}
		}

		result = append(result, item)
	}

	data, err := jettison.Marshal(result)
	if err != nil {
		return "", errors.New(helper.FormatErr)
	}

	return string(data), nil
}

// 更新用户状态
func MemberUpdateState(usernames []string, state int8) error {

	query, _, _ := dialect.Update("tbl_members").
		Set(g.Record{"state": state}).Where(g.Ex{"username": usernames, "prefix": meta.Prefix}).ToSQL()
	_, err := meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	/*
		// 更新用户redis state
		pipe := meta.MerchantRedis.TxPipeline()
		defer pipe.Close()

		for _, v := range usernames {
			pipe.HSet(ctx, usernames[v], "state", state)
			pipe.Persist(ctx, usernames[v])
		}

		_, err = pipe.Exec(ctx)
		if err != nil {
			return pushLog(err, helper.RedisErr)
		}
	*/
	for _, v := range usernames {
		MemberUpdateCache("", v)
	}

	return nil
}

/**
 * @Description: MemberList 会员列表
 * @Author: parker
 * @Date: 2021/4/14 16:38
 * @LastEditTime: 2021/4/14 19:00
 * @LastEditors: parker
 */
func MemberList(page, pageSize int, tag, startTime, endTime string, ex g.Ex) (MemberPageData, error) {

	data := MemberPageData{}
	var err error
	if startTime != "" && endTime != "" {

		startAt, err := helper.TimeToLoc(startTime, loc)
		if err != nil {
			return data, errors.New(helper.TimeTypeErr)
		}

		endAt, err := helper.TimeToLoc(endTime, loc)
		if err != nil {
			return data, errors.New(helper.TimeTypeErr)
		}
		if startAt >= endAt {
			return data, errors.New(helper.QueryTimeRangeErr)
		}

		ex["created_at"] = g.Op{"between": exp.NewRangeVal(startAt, endAt)}
	}

	if tag != "" {

		ex1 := g.Ex{
			"prefix":   meta.Prefix,
			"tag_name": tag,
		}
		var ids []uint64
		query, _, _ := dialect.From("tbl_member_tags").Select(g.DISTINCT(g.C("uid"))).Where(ex1).Order(g.C("uid").Desc()).Limit(100).ToSQL()
		fmt.Println(query)
		err := meta.TiDB.Select(&ids, query)
		if err != nil {
			return data, pushLog(err, helper.DBErr)
		}

		if len(ids) > 0 {
			ex["uid"] = ids
		}
	}

	data, err = memberList(page, pageSize, ex)
	if err != nil {
		return data, err
	}

	if len(data.D) < 1 {
		return data, nil
	}

	var (
		uids []string
	)
	for _, v := range data.D {
		uids = append(uids, v.UID)
	}

	d, err := grpc_t.DecryptAll(uids, true, []string{"realname", "email", "phone", "zalo"})

	fmt.Println("grpc_t.Decrypt uids = ", uids)
	fmt.Println("grpc_t.Decrypt d = ", d)

	if err != nil {
		fmt.Println("grpc_t.Decrypt err = ", err)
		return data, errors.New(helper.GetRPCErr)
	}

	for k, v := range data.D {
		data.D[k].Password = ""
		data.D[k].RealName = d[v.UID]["realname"]
		data.D[k].Email = d[v.UID]["email"]
		data.D[k].Phone = d[v.UID]["phone"]
		data.D[k].Zalo = d[v.UID]["zalo"]
	}

	return data, nil
}

func memberList(page, pageSize int, ex g.Ex) (MemberPageData, error) {

	data := MemberPageData{}
	ex["prefix"] = meta.Prefix
	t := dialect.From("tbl_members")
	if page == 1 {
		totalQuery, _, _ := t.Select(g.COUNT("uid")).Where(ex).ToSQL()
		err := meta.MerchantDB.Get(&data.T, totalQuery)
		if err != nil && err != sql.ErrNoRows {
			return data, pushLog(err, helper.DBErr)
		}

		if data.T == 0 {
			return data, nil
		}
	}
	offset := (page - 1) * pageSize
	var d []Member
	query, _, _ := t.Select(colsMember...).Where(ex).Order(g.I("created_at").Desc()).Offset(uint(offset)).Limit(uint(pageSize)).ToSQL()
	err := meta.MerchantDB.Select(&d, query)
	if err != nil {
		return data, pushLog(err, helper.DBErr)
	}

	if len(d) == 0 {
		return data, nil
	}

	for _, v := range d {
		val := MemberData{Member: v}
		data.D = append(data.D[0:], val)
	}
	data.S = uint(pageSize)
	return data, nil
}

func AgencyList(ex exp.ExpressionList, parentID, username, startTime, endTime, sortField string, isAsc, page, pageSize int, agencyType string) (MemberListData, error) {

	data := MemberListData{}
	startAt, err := helper.TimeToLoc(startTime, loc)
	if err != nil {
		return data, errors.New(helper.DateTimeErr)
	}

	endAt, err := helper.TimeToLoc(endTime, loc)
	if err != nil {
		return data, errors.New(helper.DateTimeErr)
	}

	if startAt >= endAt {
		return data, errors.New(helper.QueryTimeRangeErr)
	}

	data.S = pageSize

	if sortField != "" && username == "" && agencyType != "391" { // 排序
		data.D, data.T, err = memberListSort(ex, parentID, sortField, startAt, endAt, isAsc, page, pageSize)
		if err != nil {
			return data, err
		}
	} else {
		data.D, data.T, err = agencyList(ex, startAt, endAt, page, pageSize, parentID, agencyType)
		if err != nil {
			return data, err
		}
	}

	if len(data.D) == 0 {
		return data, nil
	}

	var (
		ids []string
	)
	for _, v := range data.D {
		ids = append(ids, v.UID)

	}

	// 获取用户状态 最后登录ip
	members, err := memberInfoFindBatch(ids)
	if err != nil {
		return data, err
	}

	data.Info = members

	// 获取用户的反水比例
	rebate, err := MemberRebateSelect(ids)
	if err != nil {
		return data, err
	}

	lvParams := make(map[string]string)
	for _, member := range members {
		lvParams[member.UID] = member.TopUid
	}

	// 获取代理层级  佣金方案
	lvls := memberLvl(lvParams)

	for i, v := range data.D {
		if rb, ok := rebate[v.UID]; ok {
			data.D[i].DJ = rb.DJ
			data.D[i].TY = rb.TY
			data.D[i].ZR = rb.ZR
			data.D[i].QP = rb.QP
			data.D[i].DZ = rb.DZ
			data.D[i].CP = rb.CP
			data.D[i].FC = rb.FC
			data.D[i].BY = rb.BY
			data.D[i].CgOfficialRebate = rb.CgOfficialRebate
			data.D[i].CgHighRebate = rb.CgHighRebate
		}

		if lv, ok := lvls[v.UID]; ok {
			data.D[i].Lvl = lv
		}

	}

	// 直属下级人数 新增注册人数
	data.Agg, err = MemberAgg(ids, startAt, endAt)
	if err != nil {
		return data, err
	}

	key := fmt.Sprintf("%s:rebate:enablemod", meta.Prefix)
	if meta.MerchantRedis.Exists(ctx, key).Val() > 0 {
		data.EnableMod = true
	}

	return data, nil
}

// 获取佣金方案
func memberPlan(ids []string) (map[string]map[string]string, error) {

	ex := g.Ex{
		"uid":    ids,
		"prefix": meta.Prefix,
	}

	var conf []CommssionConf
	query, _, _ := dialect.From("tbl_commission_conf").Select("id", "uid", "plan_id").Where(ex).ToSQL()
	err := meta.MerchantDB.Select(&conf, query)
	if err != nil && err != sql.ErrNoRows {
		return nil, pushLog(err, helper.DBErr)
	}

	if len(conf) == 0 {
		return nil, nil
	}

	var planID []string
	for _, v := range conf {
		planID = append(planID, v.PlanID)
	}

	var plans []CommissionPlan
	ex = g.Ex{
		"id":     planID,
		"prefix": meta.Prefix,
	}
	query, _, _ = dialect.From("tbl_commission_plan").Select(colsCommPlan...).Where(ex).ToSQL()
	err = meta.MerchantDB.Select(&plans, query)
	if err != nil && err != sql.ErrNoRows {
		return nil, pushLog(err, helper.DBErr)
	}

	if len(plans) == 0 {
		return nil, nil
	}

	planMap := make(map[string]string)
	for _, v := range plans {
		planMap[v.ID] = v.Name
	}

	data := make(map[string]map[string]string)
	for _, v := range conf {

		data[v.UID] = map[string]string{
			"plan_id": v.PlanID,
		}

		if name, ok := planMap[v.PlanID]; ok {
			data[v.UID]["name"] = name
		}

	}

	return data, nil
}

// 获取代理层级
func memberLvl(params map[string]string) map[string]int {

	var or []exp.Expression

	for k, v := range params {
		or = append(or, g.And(
			g.C("ancestor").Eq(v),   // 总代id
			g.C("descendant").Eq(k), // 代理id
		))
	}

	var trees []MembersTree
	query, _, _ := dialect.From("tbl_members_tree").Where(g.Or(or...)).ToSQL()
	err := meta.MerchantDB.Select(&trees, query)
	if err != nil {
		return nil
	}

	data := make(map[string]int, len(trees))
	for _, v := range trees {
		data[v.Descendant] = v.Lvl
	}

	return data
}

func memberListSort(ex exp.ExpressionList, parentID, sortField string, startAt, endAt int64, isAsc, page, pageSize int) ([]MemberListCol, int, error) {

	var data []MemberListCol

	exC := g.Ex{
		"report_time": g.Op{"between": exp.NewRangeVal(startAt, endAt)},
		"report_type": 2, // 投注时间2结算时间3投注时间月报4结算时间月报
		"prefix":      meta.Prefix,
	}

	ex = ex.Append(exC)

	number := 0
	if page == 1 {

		query, _, _ := dialect.From("tbl_report_agency").Select(g.COUNT(g.DISTINCT("uid"))).Where(ex).ToSQL()
		err := meta.ReportDB.Get(&number, query)
		if err != nil && err != sql.ErrNoRows {
			return data, 0, pushLog(err, helper.DBErr)
		}

		if number == 0 {
			return data, 0, nil
		}
	}

	orderField := g.L("report_time")
	if sortField != "" {
		orderField = g.L(fmt.Sprintf(`sum(%s)`, sortField))
	}

	orderBy := orderField.Desc()
	if isAsc == 1 {
		orderBy = orderField.Asc()
	}

	and := g.And(ex, g.C("uid").Neq(g.C("parent_uid")))
	if parentID != "" {
		and = g.And(
			exC,
			g.Or(
				g.And(
					g.C("uid").In(parentID),
					g.C("data_type").Eq(2),
				),
				g.And(
					g.C("data_type").Eq(1),
					g.C("parent_uid").Eq(parentID),
				),
			),
		)
	}

	offset := (page - 1) * pageSize
	query, _, _ := dialect.From("tbl_report_agency").Select(
		"uid",
		g.SUM("deposit_amount").As("deposit_amount"),
		g.SUM("withdrawal_amount").As("withdrawal_amount"),
		g.SUM("valid_bet_amount").As("valid_bet_amount"),
		g.SUM("rebate_amount").As("rebate_amount"),
		g.SUM("rebate_point").As("rebate_point"),
		g.SUM("company_net_amount").As("company_net_amount"),
	).GroupBy("uid").
		Where(and).
		Offset(uint(offset)).
		Limit(uint(pageSize)).
		Order(orderBy).
		ToSQL()
	fmt.Println(query)
	err := meta.ReportDB.Select(&data, query)
	if err != nil {
		return data, number, pushLog(err, helper.DBErr)
	}

	return data, number, nil
}

func agencyList(ex exp.ExpressionList, startAt, endAt int64, page, pageSize int, parentID, agencyType string) ([]MemberListCol, int, error) {

	var data []MemberListCol
	number := 0
	ex = ex.Append(g.C("prefix").Eq(meta.Prefix))
	if agencyType == "391" {
		ex = ex.Append(g.C("agency_type").Eq(391))
		ex = ex.Append(g.C("tester").Eq(1))
	}
	if page == 1 {
		query, _, _ := dialect.From("tbl_members").Select(g.COUNT(1)).Where(ex).ToSQL()
		fmt.Println(query)
		err := meta.MerchantDB.Get(&number, query)
		if err != nil && err != sql.ErrNoRows {
			return data, number, pushLog(err, helper.DBErr)
		}

		if number == 0 {
			return data, number, nil
		}
	}

	var members []Member
	offset := (page - 1) * pageSize
	query, _, _ := dialect.From("tbl_members").Select("uid").Where(ex).Offset(uint(offset)).
		Limit(uint(pageSize)).Order(g.I("created_at").Desc()).ToSQL()
	fmt.Println(query)
	err := meta.MerchantDB.Select(&members, query)
	if err != nil {
		return data, number, pushLog(err, helper.DBErr)
	}

	// 补全数据
	var ids []string
	idMap := make(map[string]bool, len(members))
	for _, member := range members {
		if member.UID != parentID {
			ids = append(ids, member.UID)
		}
		idMap[member.UID] = true
	}

	// 获取统计数据
	and := g.And(
		g.C("report_type").Eq(2),
		g.C("prefix").Eq(meta.Prefix),
		g.C("report_time").Between(exp.NewRangeVal(startAt, endAt)),
	)

	if parentID == "" {
		and = and.Append(
			g.And(
				g.C("uid").Neq(g.C("parent_uid")),
				g.C("uid").In(ids),
				g.C("data_type").Eq("1"),
			),
		)
	} else {

		or := g.Or(
			g.And(
				g.C("uid").In(parentID),
				g.C("data_type").Eq("2"),
			),
		)

		if len(ids) > 0 {
			or = or.Append(
				g.And(
					g.C("data_type").Eq("1"),
					g.C("uid").In(ids),
				),
			)
		}

		and = and.Append(or)
	}
	query, _, _ = dialect.From("tbl_report_agency").Where(and).
		Select(
			"uid",
			g.MAX("mem_count").As("mem_count"),
			g.SUM("deposit_amount").As("deposit_amount"),
			g.MAX("mem_count").As("mem_count"),
			g.SUM("withdrawal_amount").As("withdrawal_amount"),
			g.SUM("valid_bet_amount").As("valid_bet_amount"),
			g.SUM("rebate_amount").As("rebate_amount"),
			g.SUM("rebate_point").As("rebate_point"),
			g.SUM("company_net_amount").As("company_net_amount"),
		).GroupBy("uid").
		ToSQL()
	err = meta.ReportDB.Select(&data, query)
	fmt.Println(query)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err.Error())
		return data, number, pushLog(err, helper.DBErr)
	}

	if len(ids) == len(data) {
		return data, number, nil
	}

	// 可能有会员未生成报表数据 这时需要给未生成报表的会员 赋值默认返回值
	//否则会出现total和data length 不一致的问题
	for _, v := range data {
		delete(idMap, v.UID)
	}

	for id := range idMap {
		ele := MemberListCol{UID: id}
		data = append([]MemberListCol{ele}, data...)
	}

	return data, number, nil
}

// MemberAgg 获取直属下级人数
func MemberAgg(ids []string, startTIme, endTime int64) (map[string]MemberAggData, error) {

	aggs := make(map[string]MemberAggData)
	var data []MemberAggData

	for _, id := range ids {
		aggs[id] = MemberAggData{UID: id}
	}

	ex := g.Ex{
		"uid":         ids,
		"report_time": g.Op{"between": exp.NewRangeVal(startTIme, endTime)},
		"report_type": 2,
		"data_type":   1,
		"prefix":      meta.Prefix,
	}

	query, _, _ := dialect.From("tbl_report_agency").Select(g.MAX("subordinate_count").As("mem_count"),
		g.SUM("regist_count").As("regist_count"), g.MAX("active_count").As("active_count"), "uid").Where(ex).GroupBy("uid").ToSQL()
	err := meta.ReportDB.Select(&data, query)
	if err != nil && err != sql.ErrNoRows {
		return nil, pushLog(err, helper.DBErr)
	}

	if len(data) == 0 {
		return aggs, nil
	}

	for _, v := range data {
		aggs[v.UID] = v
	}

	return aggs, nil
}

func MemberRebateSelect(ids []string) (map[string]MemberRebate, error) {

	var own []MemberRebate
	query, _, _ := dialect.From("tbl_member_rebate_info").Select(colsMemberRebate...).Where(g.Ex{"uid": ids, "prefix": meta.Prefix}).ToSQL()
	err := meta.MerchantDB.Select(&own, query)
	if err != nil {
		return nil, pushLog(err, helper.DBErr)
	}

	data := make(map[string]MemberRebate)
	for _, v := range own {
		data[v.UID] = v
	}
	return data, nil
}

// 更新用户信息
func MemberUpdate(username, adminID string, param map[string]string, tagsId []string) error {

	var phoneHash string

	mb, err := MemberFindOne(username)
	if err != nil {
		return err
	}

	if len(mb.Username) == 0 {
		return errors.New(helper.UsernameErr)
	}

	param["uid"] = mb.UID
	record := g.Record{}
	if gender, ok := param["gender"]; ok {
		if gender != "0" {
			record["gender"] = param["gender"]
		}
	}

	encFields := [][]string{}
	if _, ok := param["realname"]; ok {

		realNameHash := fmt.Sprintf("%d", MurmurHash(param["realname"], 0))
		if realNameHash != mb.RealnameHash {
			record["realname_hash"] = realNameHash
			encFields = append(encFields, []string{"realname", param["realname"]})
		}
	}

	if _, ok := param["phone"]; ok {

		phoneHash = fmt.Sprintf("%d", MurmurHash(param["phone"], 0))
		ok, err = memberBindCheck(g.Ex{"phone_hash": phoneHash})
		if err != nil {
			return err
		}

		if ok {
			return errors.New(helper.PhoneExist)
		}

		if phoneHash != mb.PhoneHash {
			record["phone_hash"] = phoneHash
			encFields = append(encFields, []string{"phone", param["phone"]})
		}
	}

	if _, ok := param["email"]; ok {

		emailHash := fmt.Sprintf("%d", MurmurHash(param["email"], 0))
		ok, err = memberBindCheck(g.Ex{"email_hash": emailHash})
		if err != nil {
			return err
		}

		if ok {
			return errors.New(helper.EmailExist)
		}

		if emailHash != mb.EmailHash {
			record["email_hash"] = emailHash
			encFields = append(encFields, []string{"email", param["email"]})
		}
	}

	if _, ok := param["zalo"]; ok {

		zaloHash := fmt.Sprintf("%d", MurmurHash(param["zalo"], 0))
		ok, err = memberBindCheck(g.Ex{"zalo_hash": zaloHash})
		if err != nil {
			return err
		}

		if ok {
			return errors.New(helper.ZaloExist)
		}

		record["zalo_hash"] = zaloHash
		encFields = append(encFields, []string{"zalo", param["zalo"]})

	}

	if _, ok := param["address"]; ok {
		record["address"] = param["address"]
	}

	if _, ok := param["birth"]; ok {
		record["birth"] = param["birth"]
	}

	if _, ok := param["birth_hash"]; ok {
		record["birth_hash"] = param["birth_hash"]
	}

	tags := map[string]string{}
	if len(tagsId) > 0 {

		var tagls []Tags
		// 查询标签
		query, _, _ := dialect.From("tbl_tags").Select(colsTags...).Where(g.Ex{"id": tagsId, "prefix": meta.Prefix}).ToSQL()
		err = meta.MerchantDB.Select(&tagls, query)
		if err != nil {
			return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
		}

		for _, v := range tagls {
			tags[fmt.Sprintf("%d", v.ID)] = v.Name
		}
	}

	tx, err := meta.MerchantDB.Begin()
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	uid := param["uid"]
	ex := g.Ex{"uid": uid, "prefix": meta.Prefix}
	if len(record) > 0 {
		// 更新会员信息
		query, _, _ := dialect.Update("tbl_members").Set(record).Where(ex).ToSQL()
		fmt.Println(query)
		_, err = tx.Exec(query)
		if err != nil {
			return pushLog(err, helper.DBErr)
		}
	}

	// 删除该用户的所有标签
	query, _, _ := dialect.Delete("tbl_member_tags").Where(ex).ToSQL()
	_, err = tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		return pushLog(err, helper.DBErr)
	}

	if len(tags) > 0 {

		var data []MemberTags
		for k, v := range tags {
			tag := MemberTags{
				ID:        helper.GenId(),
				UID:       uid,
				AdminID:   adminID,
				TagID:     k,
				TagName:   v,
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
				Prefix:    meta.Prefix,
			}
			data = append(data, tag)
		}

		query, _, _ = dialect.Insert("tbl_member_tags").Rows(data).ToSQL()
		_, err = tx.Exec(query)
		if err != nil {
			_ = tx.Rollback()
			return pushLog(err, helper.DBErr)
		}
	}

	err = grpc_t.Encrypt(mb.UID, encFields)
	if err != nil {
		_ = tx.Rollback()
		fmt.Println("grpc_t.Encrypt = ", err)
		return errors.New(helper.UpdateRPCErr)
	}

	err = tx.Commit()
	if err != nil {
		return pushLog(err, helper.DBErr)
	}
	MemberUpdateCache(uid, "")

	if _, ok := param["phone"]; ok {
		key := fmt.Sprintf("%s:phoneExist", meta.Prefix)
		meta.MerchantRedis.SAdd(ctx, key, phoneHash).Val()
	}

	return nil
}

// 会员管理-会员列表-解除密码限制/解除短信限制
func MemberRetryReset(username string, ty uint8, pid string) error {

	if _, ok := memberUnfreeze[ty]; !ok {
		return errors.New(helper.UnfreezeTyErr)
	}

	switch ty {
	case WALLET: // 解锁钱包限制
		return memberPlatformRetryReset(username, pid)

	case PWD, SMS: // 解除密码限制/解除短信限制
		err := meta.MerchantRedis.Del(ctx, fmt.Sprintf(memberUnfreeze[ty], username)).Err()
		if err != nil {
			return pushLog(err, "redis")
		}
	}

	return nil
}

// 会员列表 用户日志写入
func MemberRemarkInsert(file, msg, adminName string, names []string, createdAt int64) error {

	// 获取所有用户的uid
	members, err := memberFindBatch(names)
	if err != nil {
		return err
	}

	if len(members) != len(names) {
		return errors.New(helper.UsernameErr)
	}

	for username, member := range members {

		rc := g.Record{
			"id":           helper.GenId(),
			"uid":          member.UID,
			"username":     username,
			"msg":          msg,
			"file":         file,
			"created_at":   createdAt,
			"created_name": adminName,
			"updated_at":   0,
			"updated_name": "",
			"is_delete":    0,
			"prefix":       meta.Prefix,
			"ts":           time.Now().In(loc).UnixMicro(),
		}
		query, _, _ := dialect.Insert("member_remarks_log").Rows(&rc).ToSQL()
		fmt.Println(query)
		_, err = meta.MerchantTD.Exec(query)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return nil
}

// 会员备注修改
func MemberRemarkUpdate(ts, file, msg, adminName string, updatedAt int64) error {

	t, err := time.ParseInLocation("2006-01-02T15:04:05.999999+07:00", ts, loc)
	if err != nil {
		return pushLog(err, helper.TimeTypeErr)
	}

	rc := g.Record{
		"msg":          msg,
		"file":         file,
		"updated_at":   updatedAt,
		"updated_name": adminName,
		"ts":           t.UnixMicro(),
	}
	query, _, _ := dialect.Insert("member_remarks_log").Rows(&rc).ToSQL()
	fmt.Println(query)
	_, err = meta.MerchantTD.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

// 会员备注删除
func MemberRemarkDelete(ts, adminName string, updatedAt int64) error {

	t, err := time.ParseInLocation("2006-01-02T15:04:05.999999+07:00", ts, loc)
	if err != nil {
		return pushLog(err, helper.TimeTypeErr)
	}

	rc := g.Record{
		"updated_at":   updatedAt,
		"updated_name": adminName,
		"is_delete":    1,
		"ts":           t.UnixMicro(),
	}
	query, _, _ := dialect.Insert("member_remarks_log").Rows(&rc).ToSQL()
	fmt.Println(query)
	_, err = meta.MerchantTD.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

// 会员管理-会员列表-数据概览
func MemberDataOverview(username, startTime, endTime string) (MemberDataOverviewData, error) {

	data := MemberDataOverviewData{}

	// 获取uid
	mb, err := MemberFindOne(username)
	if err != nil {
		return data, err
	}

	ss, err := helper.TimeToLoc(startTime, loc)
	if err != nil {
		return data, errors.New(helper.TimeTypeErr)
	}

	se, err := helper.TimeToLoc(endTime, loc)
	if err != nil {
		return data, errors.New(helper.TimeTypeErr)
	}

	// 毫秒级时间戳
	mss, err := helper.TimeToLocMs(startTime, loc)
	if err != nil {
		return data, errors.New(helper.TimeTypeErr)
	}

	mse, err := helper.TimeToLocMs(endTime, loc)
	if err != nil {
		return data, errors.New(helper.TimeTypeErr)
	}

	if mss > mse {
		return data, errors.New(helper.QueryTimeRangeErr)
	}

	// 总输赢 && 总有效投注
	ex := g.Ex{}
	ex["flag"] = 1
	ex["name"] = username
	r := GameResult_t{}
	ex["bet_time"] = g.Op{"between": exp.NewRangeVal(mss, mse)}

	query, _, _ := dialect.From("tbl_game_record").Select(g.SUM("valid_bet_amount").As("valid_bet_amount"),
		g.SUM("net_amount").As("net_amount")).Where(ex).Limit(1).ToSQL()
	err = meta.TiDB.Get(&r, query)
	if err != nil {
		return data, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}
	data.NetAmount = r.NetAmount.Float64
	data.ValidBetAmount = r.ValidBetAmount.Float64

	// 总存款
	ex = g.Ex{"uid": mb.UID,
		"prefix":     meta.Prefix,
		"state":      DepositSuccess,
		"confirm_at": g.Op{"between": exp.NewRangeVal(ss, se)},
	}
	query, _, _ = dialect.From("tbl_deposit").
		Select(g.COALESCE(g.SUM("amount"), 0).As("dividend")).Where(ex).ToSQL()
	err = meta.MerchantDB.Get(&data.Deposit, query)
	if err != nil {
		return data, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	// 总提款
	wex := g.Ex{"uid": mb.UID,
		"prefix":     meta.Prefix,
		"state":      WithdrawSuccess,
		"confirm_at": g.Op{"between": exp.NewRangeVal(ss, se)},
	}
	query, _, _ = dialect.From("tbl_withdraw").
		Select(g.COALESCE(g.SUM("amount"), 0).As("withdraw")).Where(wex).ToSQL()
	err = meta.MerchantDB.Get(&data.Withdraw, query)
	if err != nil {
		return data, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	// 总红利
	dex := g.Ex{"uid": mb.UID,
		"prefix":   meta.Prefix,
		"state":    DividendReviewPass,
		"apply_at": g.Op{"between": exp.NewRangeVal(mss, mse)},
	}
	query, _, _ = dialect.From("tbl_member_dividend").
		Select(g.COALESCE(g.SUM("amount"), 0).As("dividend")).Where(dex).ToSQL()
	err = meta.MerchantDB.Get(&data.Dividend, query)
	if err != nil {
		return data, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	// 总返水
	rex := g.Ex{
		"uid":    mb.UID,
		"prefix": meta.Prefix,
		"cash_type": []int{
			helper.TransactionRebateCasino,       //真人返水
			helper.TransactionRebateLottery,      //彩票返水
			helper.TransactionRebateSport,        //体育返水
			helper.TransactionRebateDesk,         //棋牌返水
			helper.TransactionRebateESport,       //电竞返水
			helper.TransactionRebateCockFighting, //斗鸡返水
			helper.TransactionRebateFishing,      //捕鱼返水
			helper.TransactionRebateLott,         //电游返水
			helper.TransactionRebateCGLottery,    //彩票返点
			helper.TransactionSubRebate,          //下级返水
		},
		"created_at": g.Op{"between": exp.NewRangeVal(ss*1000, se*1000)},
	}
	query, _, _ = dialect.From("tbl_balance_transaction").
		Select(g.COALESCE(g.SUM("amount"), 0).As("rebate")).Where(rex).ToSQL()
	fmt.Println(query)
	err = meta.TiDB.Get(&data.Rebate, query)
	if err != nil {
		return data, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	return data, nil
}

func MemberUpdatePwd(username, pwd string, ty int, ctx *fasthttp.RequestCtx) error {

	mb, err := MemberFindOne(username)
	if err != nil || mb.Username == "" {
		return errors.New(helper.UsernameErr)
	}

	admin, err := AdminToken(ctx)
	if err != nil || admin["name"] == "" {
		return errors.New(helper.AccessTokenExpires)
	}

	record := g.Record{}
	if ty == 1 {
		record["withdraw_pwd"] = fmt.Sprintf("%d", MurmurHash(pwd, mb.CreatedAt))
	} else {
		record["password"] = fmt.Sprintf("%d", MurmurHash(pwd, mb.CreatedAt))
	}
	query, _, _ := dialect.Update("tbl_members").Set(record).Where(g.Ex{"uid": mb.UID}).ToSQL()
	_, err = meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	MemberUpdateCache("", username)
	return nil
}

func MemberHistory(id, field string, encrypt bool) ([]string, error) {

	history, err := grpc_t.View(id, field, encrypt)
	if err != nil {
		fmt.Println("grpc_t.View err = ", err)
		return nil, err
	}

	//fmt.Println("grpc_t.View history = ", history)
	return history, nil
}

func MemberFull(id string, field []string) (map[string]string, error) {

	recs, err := grpc_t.Decrypt(id, false, field)
	if err != nil {
		fmt.Println("grpc_t.Decrypt err = ", err)
		return nil, err
	}

	return recs, nil
}

func MemberBalanceZero(username, remark, adminID, adminName string) error {

	mb, err := MemberBalance(username)
	if err != nil {
		return err
	}

	balance, err := decimal.NewFromString(mb.Balance)
	if err != nil {
		return err
	}
	// 余额大于0，不清零
	if balance.Cmp(decimal.Zero) > 0 {
		return nil
	}

	tx, err := meta.MerchantDB.Begin()
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	record := g.Record{
		"balance": "0.00",
	}
	query, _, _ := dialect.Update("tbl_members").Set(record).Where(g.Ex{"uid": mb.UID}).ToSQL()
	//fmt.Println(query)
	_, err = tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	id := helper.GenId()
	trans := MemberTransaction{
		AfterAmount:  "0.00",
		Amount:       balance.Abs().String(),
		BeforeAmount: balance.String(),
		BillNo:       id,
		CreatedAt:    time.Now().UnixMilli(),
		ID:           id,
		CashType:     helper.TransactionSetBalanceZero,
		UID:          mb.UID,
		Username:     username,
		Prefix:       meta.Prefix,
		Remark:       remark,
	}
	query, _, _ = dialect.Insert("tbl_balance_transaction").Rows(trans).ToSQL()
	_, err = tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	_ = tx.Commit()
	//fmt.Println(query)
	_ = MemberUpdateCache(mb.UID, "")
	return nil
}

//根据 uid数组，redis批量获取用户余额
func MemberBalance(username string) (MBBalance, error) {

	mb := MBBalance{}

	t := dialect.From("tbl_members")
	query, _, _ := t.Select(colsMemberBalance...).Where(g.Ex{"username": username, "prefix": meta.Prefix}).ToSQL()
	err := meta.MerchantDB.Get(&mb, query)
	if err != nil && err != sql.ErrNoRows {
		return mb, pushLog(err, helper.DBErr)
	}

	if err == sql.ErrNoRows {
		return mb, errors.New(helper.UsernameErr)
	}

	return mb, nil
}

//根据 uid数组，redis批量获取用户余额
func MemberBalanceBatch(uids []string) (string, error) {

	var mbs []MBBalance
	t := dialect.From("tbl_members")
	query, _, _ := t.Select(colsMemberBalance...).Where(g.Ex{"uid": uids}).ToSQL()
	err := meta.MerchantDB.Select(&mbs, query)
	if err != nil {
		return "", pushLog(err, helper.DBErr)
	}

	data, err := jettison.Marshal(mbs)
	if err != nil {
		return "", errors.New(helper.FormatErr)
	}

	return string(data), nil
}

//根据 uid数组，redis批量获取用户余额
func memMapBalanceBatch(uids []string) (map[string]MBBalance, error) {

	var mbs []MBBalance
	mbsm := map[string]MBBalance{}
	t := dialect.From("tbl_members")
	query, _, _ := t.Select(colsMemberBalance...).Where(g.Ex{"uid": uids}).ToSQL()
	err := meta.MerchantDB.Select(&mbs, query)
	if err != nil {
		return nil, pushLog(err, helper.DBErr)
	}

	for _, v := range mbs {
		mbsm[v.UID] = v
	}

	return mbsm, nil
}

// 解锁场馆钱包限制
func memberPlatformRetryReset(username, pid string) error {

	user, err := MemberFindOne(username)
	if err != nil {
		return err
	}

	param, err := memberPlatPromoInfo(user.UID, pid)
	if err != nil {
		return pushLog(err, helper.RedisErr)
	}

	if param == nil {
		return errors.New(helper.PlatNoPromoApply)
	}

	// 余额解锁活动
	key := fmt.Sprintf(memberUnfreeze[WALLET], user.UID, pid)
	meta.MerchantRedis.Unlink(ctx, key)

	// 活动状态变更类型为解锁活动
	param["alter_ty"] = fmt.Sprintf("%d", PromoUnlock)
	// 解锁类型为余额解锁
	param["unlock_ty"] = fmt.Sprintf("%d", PromoUnlockAdmin)
	// 投递消息队列，异步处理会员场馆活动解锁
	//err = BeanPut("promo", param)
	//if err != nil {
	//	_ = pushLog(err, helper.ServerErr)
	//}

	return nil
}

// 检测手机号，email，是否已经被会员绑定
// 仅用来检测会员信息绑定
func memberBindCheck(ex g.Ex) (bool, error) {

	var id string

	t := dialect.From("tbl_members")
	query, _, _ := t.Select("uid").Where(ex).Limit(1).ToSQL()
	fmt.Println("memberBindCheck", query)
	err := meta.MerchantDB.Get(&id, query)
	if err != nil && err != sql.ErrNoRows {
		return false, pushLog(err, helper.DBErr)
	}

	return err != sql.ErrNoRows, nil
}

// 通过用户名获取用户在redis中的数据
func MemberFindOne(name string) (Member, error) {

	m := Member{}
	if name == "" {
		return m, errors.New(helper.UsernameErr)
	}

	t := dialect.From("tbl_members")
	query, _, _ := t.Select(colsMember...).Where(g.Ex{"username": name, "prefix": meta.Prefix}).ToSQL()
	err := meta.MerchantDB.Get(&m, query)
	if err != nil && err != sql.ErrNoRows {
		return m, pushLog(err, helper.DBErr)
	}

	if err == sql.ErrNoRows {
		return m, errors.New(helper.UsernameErr)
	}

	return m, nil
}

func memberFindBatch(names []string) (map[string]Member, error) {

	data := map[string]Member{}

	if len(names) == 0 {
		return data, errors.New(helper.ParamNull)
	}

	var mbs []Member
	t := dialect.From("tbl_members")
	query, _, _ := t.Select(colsMember...).Where(g.Ex{"username": names, "prefix": meta.Prefix}).ToSQL()
	err := meta.MerchantDB.Select(&mbs, query)
	if err != nil {
		return data, pushLog(err, helper.DBErr)
	}

	if len(mbs) > 0 {
		for _, v := range mbs {
			if v.Username != "" {
				data[v.Username] = v
			}
		}
	}

	return data, nil
}

func memberInfoFindBatch(ids []string) (map[string]memberInfo, error) {

	if len(ids) == 0 {
		return nil, errors.New(helper.ParamNull)
	}

	var mbs []memberInfo
	query, _, _ := dialect.From("tbl_members").Select(colsMemberInfo...).Where(g.Ex{"uid": ids}).ToSQL()
	err := meta.MerchantDB.Select(&mbs, query)
	if err != nil {
		return nil, pushLog(err, helper.DBErr)
	}

	if len(mbs) == 0 {
		return nil, nil
	}

	data := make(map[string]memberInfo, len(mbs))
	for _, v := range mbs {
		data[v.UID] = v
	}

	return data, nil
}

// 获取会员指定场馆的活动锁定信息
func memberPlatPromoInfo(uid, pid string) (map[string]interface{}, error) {

	param := map[string]interface{}{}
	key := fmt.Sprintf("P:%s:%s", uid, pid)
	fields := []string{
		"pid",
		"pname",
		"cash_type",
		"apply_at",
		"water_flow",
	}
	rs, err := meta.MerchantRedis.HMGet(ctx, key, fields...).Result()
	if err != nil {
		return nil, err
	}

	if len(fields) != len(rs) {
		return nil, err
	}

	for k, v := range rs {
		if v == nil {
			return nil, err
		}
		param[fields[k]] = v
	}

	return param, nil
}

//获取会员可转账的场馆
func memberPlatformBalance(username string) []PlatBalance {

	var p []PlatBalance
	t := dialect.From("tbl_member_platform")
	ex := g.Ex{
		"username": username,
		"id": []string{
			"6798510151614082003",
			"7219886347116135962",
			"5864536520308659696",
			"1958997188942770517",
			"2658175169982643138",
			"2306868265751172637",
			"6238858173568905466",
			"1055235995899664907",
			"1371916058167324188",
			"2854120181948444476",
			"1846182857231915191",
			"2299282204811996672",
			"7591876028427516934",
			"8840968482572372234",
			"6982718883667836955",
			"1386624620395927266",
			"1794601907316741515",
			"6861705028422769039",
		},
		"prefix": meta.Prefix,
	}
	query, _, _ := t.Select(colsPlatBalance...).Where(ex).ToSQL()
	err := meta.MerchantDB.Select(&p, query)
	if err != nil {
		_ = pushLog(err, helper.DBErr)
	}

	return p
}

func MemberUpdateInfo(uid string, record g.Record) error {

	ex := g.Ex{
		"uid":    uid,
		"prefix": meta.Prefix,
	}
	query, _, _ := dialect.Update("tbl_members").Set(&record).Where(ex).ToSQL()
	fmt.Println(query)
	_, err := meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	_ = MemberUpdateCache(uid, "")
	return nil
}

func MemberUpdateRebateInfo(uid string, mr MemberRebateResult_t) error {

	//key := fmt.Sprintf("%s:rebate:enablemod", meta.Prefix)
	//if meta.MerchantRedis.Exists(ctx, key).Val() == 0 {
	//	return errors.New(helper.MemberRebateModDisable)
	//}

	ex := g.Ex{
		"uid":    uid,
		"prefix": meta.Prefix,
	}
	recd := g.Record{
		"ty":                 mr.TY.StringFixed(1),
		"zr":                 mr.ZR.StringFixed(1),
		"qp":                 mr.QP.StringFixed(1),
		"dj":                 mr.DJ.StringFixed(1),
		"dz":                 mr.DZ.StringFixed(1),
		"cp":                 mr.CP.StringFixed(1),
		"fc":                 mr.FC.StringFixed(1),
		"by":                 mr.BY.StringFixed(1),
		"cg_high_rebate":     mr.CGHighRebate.StringFixed(2),
		"cg_official_rebate": mr.CGOfficialRebate.StringFixed(2),
	}
	query, _, _ := dialect.Update("tbl_member_rebate_info").Set(&recd).Where(ex).ToSQL()
	fmt.Println(query)
	_, err := meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	_ = MemberRebateUpdateCache1(uid, mr)
	return nil
}

func MemberUpdateMaintainName(uid, maintainName string) error {

	tx, err := meta.MerchantDB.Begin() // 开启事务
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	subEx := g.Ex{
		"uid": uid,
	}
	recd := g.Record{
		"maintain_name": maintainName,
	}
	query, _, _ := dialect.Update("tbl_members").Set(&recd).Where(subEx).ToSQL()
	_, err = tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		return pushLog(err, helper.DBErr)
	}

	err = tx.Commit()
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	MemberUpdateCache(uid, "")
	return nil
}

func MemberMaxRebateFindOne(uid string) (MemberRebateResult_t, error) {

	data := MemberMaxRebate{}
	res := MemberRebateResult_t{}

	t := dialect.From("tbl_member_rebate_info")
	query, _, _ := t.Select(
		g.MAX("zr").As("zr"),
		g.MAX("qp").As("qp"),
		g.MAX("dz").As("dz"),
		g.MAX("dj").As("dj"),
		g.MAX("ty").As("ty"),
		g.MAX("cp").As("cp"),
		g.MAX("fc").As("fc"),
		g.MAX("by").As("by"),
		g.MAX("cg_high_rebate").As("cg_high_rebate"),
		g.MAX("cg_official_rebate").As("cg_official_rebate"),
	).Where(g.Ex{"parent_uid": uid, "prefix": meta.Prefix}).ToSQL()
	fmt.Println(query)
	err := meta.MerchantDB.Get(&data, query)
	if err == sql.ErrNoRows {

		res.ZR = decimal.NewFromInt(0).Truncate(1)
		res.QP = decimal.NewFromInt(0).Truncate(1)
		res.TY = decimal.NewFromInt(0).Truncate(1)
		res.DJ = decimal.NewFromInt(0).Truncate(1)
		res.DZ = decimal.NewFromInt(0).Truncate(1)
		res.CP = decimal.NewFromInt(0).Truncate(1)
		res.FC = decimal.NewFromInt(0).Truncate(1)
		res.BY = decimal.NewFromInt(0).Truncate(1)
		res.CGHighRebate = decimal.NewFromFloat(9.00).Truncate(2)
		res.CGOfficialRebate = decimal.NewFromFloat(9.00).Truncate(2)

		return res, nil
	}
	if err != nil {
		return res, pushLog(err, helper.DBErr)
	}

	res.ZR = decimal.NewFromFloat(data.ZR.Float64).Truncate(1)
	res.QP = decimal.NewFromFloat(data.QP.Float64).Truncate(1)
	res.TY = decimal.NewFromFloat(data.TY.Float64).Truncate(1)
	res.DJ = decimal.NewFromFloat(data.DJ.Float64).Truncate(1)
	res.DZ = decimal.NewFromFloat(data.DZ.Float64).Truncate(1)
	res.CP = decimal.NewFromFloat(data.CP.Float64).Truncate(1)
	res.FC = decimal.NewFromFloat(data.FC.Float64).Truncate(1)
	res.BY = decimal.NewFromFloat(data.BY.Float64).Truncate(1)
	res.CGHighRebate = decimal.NewFromFloat(data.CgHighRebate.Float64).Truncate(2)
	res.CGOfficialRebate = decimal.NewFromFloat(data.CgOfficialRebate.Float64).Truncate(2)

	return res, nil
}

func MemberParentRebate(uid string) (MemberRebateResult_t, error) {

	data := MemberMaxRebate{}
	res := MemberRebateResult_t{}

	t := dialect.From("tbl_member_rebate_info")
	query, _, _ := t.Select(
		g.C("zr").As("zr"),
		g.C("qp").As("qp"),
		g.C("dz").As("dz"),
		g.C("dj").As("dj"),
		g.C("ty").As("ty"),
		g.C("cp").As("cp"),
		g.C("fc").As("fc"),
		g.C("by").As("by"),
		g.MAX("cg_high_rebate").As("cg_high_rebate"),
		g.MAX("cg_official_rebate").As("cg_official_rebate"),
	).Where(g.Ex{"uid": uid, "prefix": meta.Prefix}).ToSQL()
	err := meta.MerchantDB.Get(&data, query)
	if err != nil {
		return res, pushLog(err, helper.DBErr)
	}

	res.ZR = decimal.NewFromFloat(data.ZR.Float64).Truncate(1)
	res.QP = decimal.NewFromFloat(data.QP.Float64).Truncate(1)
	res.TY = decimal.NewFromFloat(data.TY.Float64).Truncate(1)
	res.DJ = decimal.NewFromFloat(data.DJ.Float64).Truncate(1)
	res.DZ = decimal.NewFromFloat(data.DZ.Float64).Truncate(1)
	res.CP = decimal.NewFromFloat(data.CP.Float64).Truncate(1)
	res.FC = decimal.NewFromFloat(data.FC.Float64).Truncate(1)
	res.BY = decimal.NewFromFloat(data.BY.Float64).Truncate(1)
	res.CGHighRebate = decimal.NewFromFloat(data.CgHighRebate.Float64).Truncate(2)
	res.CGOfficialRebate = decimal.NewFromFloat(data.CgOfficialRebate.Float64).Truncate(2)

	return res, nil
}

//代理管理 下级会员
func AgencyMemberList(param MemberListParam) (AgencyMemberData, error) {

	res := AgencyMemberData{}
	//查询MySQL,必须是代理的下级会员
	ex := g.Ex{"tester": 1}

	if param.ParentName != "" {
		ex["parent_name"] = param.ParentName
	}

	if param.State != 0 {
		ex["state"] = param.State
	}

	if param.Username != "" {
		ex["username"] = param.Username
	}

	if param.RegStart != "" && param.RegEnd != "" {

		startAt, err := helper.DayOfStart(param.RegStart, loc)
		if err != nil {
			return res, errors.New(helper.TimeTypeErr)
		}

		endAt, err := helper.DayOfEnd(param.RegEnd, loc)
		if err != nil {
			return res, errors.New(helper.TimeTypeErr)
		}

		ex["created_at"] = g.Op{"between": exp.NewRangeVal(startAt, endAt)}
	}

	t := dialect.From("tbl_members")
	if param.Page == 1 {
		countQuery, _, _ := t.Select(g.COUNT(1)).Where(ex).ToSQL()
		err := meta.MerchantDB.Get(&res.T, countQuery)
		if err != nil {
			return res, pushLog(fmt.Errorf("%s,[%s]", err.Error(), countQuery), helper.DBErr)
		}

		if res.T == 0 {
			return res, nil
		}
	}

	var mbList []memberListShow
	offset := (param.Page - 1) * param.PageSize
	query, _, _ := t.Select(colsMemberListShow...).
		Where(ex).Offset(uint(offset)).Limit(uint(param.PageSize)).Order(g.C("created_at").Desc()).ToSQL()
	err := meta.MerchantDB.Select(&mbList, query)
	if err != nil {
		return res, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	var (
		uids        []string
		agencyNames []string
	)

	for _, val := range mbList {
		uids = append(uids, val.UID)

		if val.ParentName != "" && val.ParentName != "root" {
			agencyNames = append(agencyNames, val.ParentName)
		}
	}

	// 用户中心钱包余额
	balanceMap, err := memMapBalanceBatch(uids)
	if err != nil {
		return res, err
	}

	rangeParam := map[string][]interface{}{}
	if param.StartAt != "" && param.EndAt != "" {

		startAt, err := helper.DayOfStart(param.StartAt, loc)
		if err != nil {
			return res, errors.New(helper.TimeTypeErr)
		}

		endAt, err := helper.DayOfEnd(param.EndAt, loc)
		if err != nil {
			return res, errors.New(helper.TimeTypeErr)
		}

		rangeParam["report_time"] = []interface{}{startAt, endAt}
	}

	// 获取用户数据
	md, err := MemberSumByRange(param.StartAt, param.EndAt, uids)
	if err != nil {
		return res, err
	}

	for _, m := range mbList {

		val := memberListData{memberListShow: m}
		if md, ok := md[m.UID]; ok {
			val.CompanyNetAmount, _ = decimal.NewFromFloat(md.CompanyNetAmount).Truncate(4).Float64()
			val.DepositAmount, _ = decimal.NewFromFloat(md.DepositAmount).Truncate(4).Float64()
			val.WithdrawalAmount, _ = decimal.NewFromFloat(md.WithdrawAmount).Truncate(4).Float64()
			val.ValidBetAmount, _ = decimal.NewFromFloat(md.ValidBetAmount).Truncate(4).Float64()
			val.RebateAmount, _ = decimal.NewFromFloat(md.RebateAmount).Truncate(4).Float64()
			val.DividendAmount, _ = decimal.NewFromFloat(md.DividendAmount).Truncate(4).Float64()
			val.DividendAgency, _ = decimal.NewFromFloat(md.DividendAgency).Truncate(4).Float64()
		}

		if _, o := balanceMap[m.UID]; o {
			val.Balance = balanceMap[m.UID].Balance
		}

		res.D = append(res.D, val)
	}

	return res, nil
}

func MemberSumByRange(start, end string, uids []string) (map[string]AgencyBaseSumField, error) {

	if start != "" && end != "" {

		startAt, err := helper.DayOfStart(start, loc)
		if err != nil && err != sql.ErrNoRows {
			return nil, errors.New(helper.TimeTypeErr)
		}

		endAt, err := helper.DayOfEnd(end, loc)
		if err != nil && err != sql.ErrNoRows {
			return nil, errors.New(helper.TimeTypeErr)
		}

		var (
			result = map[string]AgencyBaseSumField{}
			data   []MemReport
			num    int
		)
		ex := g.Ex{
			"uid":         uids,
			"report_time": g.Op{"between": g.Range(startAt, endAt)},
			"report_type": 2,
			"data_type":   2,
		}
		query, _, _ := dialect.From("tbl_report_agency").
			Select(g.COUNT("uid").As("num")).Where(ex).Order(g.C("uid").Desc()).ToSQL()
		fmt.Println(query)
		err = meta.ReportDB.Get(&num, query)
		if num > 0 {

			query, _, _ = dialect.From("tbl_report_agency").
				Select(g.C("uid").As("uid"), g.SUM("deposit_amount").As("deposit_amount"), g.SUM("withdrawal_amount").As("withdrawal_amount"),
					g.SUM("adjust_amount").As("adjust_amount"), g.SUM("valid_bet_amount").As("valid_bet_amount"),
					g.SUM("company_net_amount").As("company_net_amount"), g.SUM("dividend_amount").As("dividend_amount"),
					g.SUM("rebate_amount").As("rebate_amount"),
				).Where(ex).GroupBy("uid").Order(g.C("uid").Desc()).ToSQL()
			fmt.Println(query)
			err = meta.ReportDB.Select(&data, query)
			if err != nil {
				return result, err
			}
			for _, v := range data {
				obj := AgencyBaseSumField{
					DepositAmount:    v.DepositAmount,
					WithdrawAmount:   v.WithdrawalAmount,
					ValidBetAmount:   v.ValidBetAmount,
					CompanyNetAmount: v.CompanyNetAmount,
					DividendAmount:   v.DividendAmount,
					RebateAmount:     v.RebateAmount,
					AdjustAmount:     v.AdjustAmount,
				}
				result[v.Uid] = obj
			}
			return result, nil
		}
	}

	return nil, nil
}
