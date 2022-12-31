package model

import (
	"database/sql"
	"errors"
	"fmt"
	"merchant/contrib/helper"
	"merchant/contrib/validator"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"

	g "github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

var (
	betTimeFlags = map[string]string{
		"1": "bet_time",
		"2": "settle_time",
		"3": "start_time",
	}
)

type GameGroupData struct {
	Agg map[string]string `json:"agg"`
	D   []GameGroup_t     `json:"d"`
	T   int               `json:"t"`
	S   int               `json:"s"`
}

type GameGroup_t struct {
	ApiName        string `json:"api_name" db:"api_type"`
	Total          string `json:"total" db:"total"`
	NetAmount      string `json:"net_amount" db:"net_amount"`
	ValidBetAmount string `json:"valid_bet_amount" db:"valid_bet_amount"`
	BetAmount      string `json:"bet_amount" db:"bet_amount"`
}

func RecordTransaction(page, pageSize int, startTime, endTime string, ex g.Ex) (TransactionData, error) {

	data := TransactionData{}
	ex["prefix"] = meta.Prefix
	if startTime != "" && endTime != "" {
		startAt, err := helper.TimeToLocMs(startTime, loc) // 毫秒级时间戳
		if err != nil {
			return data, errors.New(helper.TimeTypeErr)
		}

		endAt, err := helper.TimeToLocMs(endTime, loc) // 毫秒级时间戳
		if err != nil {
			return data, errors.New(helper.TimeTypeErr)
		}

		if startAt >= endAt {
			return data, errors.New(helper.QueryTimeRangeErr)
		}

		ex["created_at"] = g.Op{"between": exp.NewRangeVal(startAt, endAt)}
	}

	t := dialect.From("tbl_balance_transaction")
	if page == 1 {
		query, _, _ := t.Select(g.COUNT("id")).Where(ex).ToSQL()
		fmt.Println(query)
		err := meta.TiDB.Get(&data.T, query)
		if err != nil {
			return data, pushLog(err, helper.DBErr)
		}

		if data.T == 0 {
			return data, nil
		}

		query, _, _ = t.Select(g.L("sum(after_amount - before_amount)").As("agg")).Where(ex).ToSQL()
		//fmt.Println(query)
		err = meta.TiDB.Get(&data.Agg, query)
		if err != nil {
			return data, pushLog(err, helper.DBErr)
		}
	}

	offset := pageSize * (page - 1)
	query, _, _ := t.Select(colsTransaction...).Where(ex).
		Offset(uint(offset)).Limit(uint(pageSize)).Order(g.C("created_at").Desc()).ToSQL()
	fmt.Println(query)
	err := meta.TiDB.Select(&data.D, query)
	if err != nil && err != sql.ErrNoRows {
		return data, pushLog(err, helper.DBErr)
	}

	return data, nil
}

func RecordTransfer(page, pageSize int, startTime, endTime string, ex g.Ex) (TransferData, error) {

	data := TransferData{}
	if startTime != "" && endTime != "" {
		//判断日期
		startAt, err := helper.TimeToLocMs(startTime, loc) // 毫秒级时间戳
		if err != nil {
			return data, errors.New(helper.DateTimeErr)
		}
		endAt, err := helper.TimeToLocMs(endTime, loc) // 毫秒级时间戳
		if err != nil {
			return data, errors.New(helper.DateTimeErr)
		}

		if startAt >= endAt {
			return data, errors.New(helper.QueryTimeRangeErr)
		}

		ex["created_at"] = g.Op{"between": exp.NewRangeVal(startAt, endAt)}
	}
	ex["prefix"] = meta.Prefix
	t := dialect.From("tbl_member_transfer")
	if page == 1 {
		query, _, _ := t.Select(g.COUNT("id")).Where(ex).ToSQL()
		err := meta.MerchantDB.Get(&data.T, query)
		fmt.Println(query)
		if err != nil {
			return data, pushLog(err, helper.DBErr)
		}

		if data.T == 0 {
			return data, nil
		}

		query, _, _ = t.Select(g.SUM("amount").As("agg")).Where(ex).ToSQL()
		err = meta.MerchantDB.Get(&data.Agg, query)
		fmt.Println(query)
		if err != nil {
			return data, pushLog(err, helper.DBErr)
		}
	}

	offset := pageSize * (page - 1)
	query, _, _ := t.Select(colsTransfer...).Where(ex).
		Offset(uint(offset)).Limit(uint(pageSize)).Order(g.C("created_at").Desc()).ToSQL()
	fmt.Println(query)
	err := meta.MerchantDB.Select(&data.D, query)
	if err != nil && err != sql.ErrNoRows {
		return data, pushLog(err, helper.DBErr)
	}

	return data, nil
}

func Game(ty int, pageSize, page uint, params map[string]string) (GameRecordData, error) {

	data := GameRecordData{}
	if len(params["username"]) > 0 {
		username := strings.ToLower(params["username"])
		mb, err := memberInfoCache(username)
		if err != nil {
			return data, errors.New(helper.UsernameErr)
		}
		params["uid"] = mb.UID
		params["name"] = ""
	}
	if ty == GameTyValid {

		ex := g.Ex{
			"uid":    params["uid"],
			"flag":   "1",
			"prefix": meta.Prefix,
		}
		data, err := RecordAdminGame(params["time_flag"], params["start_time"], params["end_time"], page, pageSize, ex)
		if err != nil {
			return data, errors.New(helper.DBErr)
		}
		return data, nil
	}
	//查询条件
	ex := g.Ex{
		"prefix": meta.Prefix,
	}

	if params["pid"] != "" {
		if strings.Contains(params["pid"], ",") {
			pids := strings.Split(params["pid"], ",")

			var ids []interface{}
			for _, v := range pids {
				if validator.CtypeDigit(v) {
					ids = append(ids, v)
				}
			}

			ex["api_type"] = ids
		}

		if !strings.Contains(params["pid"], ",") {
			if validator.CtypeDigit(params["pid"]) {
				ex["api_type"] = params["pid"]
			}
		}
	}

	if params["flag"] != "" {
		ex["flag"] = params["flag"]
	}

	rangeField := ""
	if params["time_flag"] != "" {
		rangeField = betTimeFlags[params["time_flag"]]
	}

	if rangeField == "" {
		return data, errors.New(helper.QueryTermsErr)
	}

	if ty == GameMemberWinOrLose {

		ex["name"] = params["username"]
		data, err := RecordAdminGame(params["time_flag"], params["start_time"], params["end_time"], page, pageSize, ex)
		if err != nil {
			return data, errors.New(helper.DBErr)
		}
		return data, nil
	}

	if !validator.CtypeDigit(params["time_flag"]) {
		return data, errors.New(helper.QueryTermsErr)
	}

	if params["bet_min"] == "" && params["bet_max"] != "" {
		max, _ := strconv.ParseFloat(params["bet_max"], 64)
		ex["bet_amount"] = g.Op{"lt": max}
	}

	if params["bet_min"] != "" && params["bet_max"] == "" {
		min, _ := strconv.ParseFloat(params["bet_min"], 64)
		ex["bet_amount"] = g.Op{"gt": min}
	}

	if params["bet_min"] != "" && params["bet_max"] != "" {
		min, _ := strconv.ParseFloat(params["bet_min"], 64)
		max, _ := strconv.ParseFloat(params["bet_max"], 64)
		if max < min {
			return data, errors.New(helper.BetAmountRangeErr)
		}

		ex["bet_amount"] = g.Op{"between": exp.NewRangeVal(min, max)}
	}

	if params["uid"] != "" {
		ex["uid"] = params["uid"]
	}

	if params["plat_type"] != "" {
		ex["game_type"] = params["plat_type"]
	}

	if params["game_name"] != "" {
		ex["game_name"] = params["game_name"]
	}

	if params["username"] != "" {
		ex["name"] = params["username"]
	}

	if params["bill_no"] != "" {
		ex["bill_no"] = params["bill_no"]
	}

	if params["api_bill_no"] != "" {
		ex["api_bill_no"] = params["api_bill_no"]
	}

	if params["pre_settle"] != "" {
		early, _ := strconv.Atoi(params["pre_settle"])
		ex["presettle"] = early
	}

	if params["resettle"] != "" {
		second, _ := strconv.Atoi(params["resettle"])
		ex["resettle"] = second
	}

	if params["parent_name"] != "" {
		ex["parent_name"] = params["parent_name"]
	}

	if params["top_name"] != "" {
		ex["top_name"] = params["top_name"]
	}

	data, err := RecordAdminGame(params["time_flag"], params["start_time"], params["end_time"], page, pageSize, ex)
	if err != nil {
		return data, errors.New(helper.DBErr)
	}

	return data, nil
}

func GameGroup(ty, pageSize, page int, params map[string]string) (GameGroupData, error) {

	data := GameGroupData{}
	//判断日期
	startAt, err := helper.TimeToLocMs(params["start_time"], loc) // 毫秒级时间戳
	if err != nil {
		return data, errors.New(helper.DateTimeErr)
	}
	endAt, err := helper.TimeToLocMs(params["end_time"], loc) // 毫秒级时间戳
	if err != nil {
		return data, errors.New(helper.DateTimeErr)
	}

	if startAt >= endAt {
		return data, errors.New(helper.QueryTimeRangeErr)
	}

	//查询条件
	ex := g.Ex{
		"prefix": meta.Prefix,
	}
	if params["pid"] != "" {
		if strings.Contains(params["pid"], ",") {
			pids := strings.Split(params["pid"], ",")

			var ids []interface{}
			for _, v := range pids {
				if validator.CtypeDigit(v) {
					ids = append(ids, v)
				}
			}

			ex["api_type"] = ids
		}

		if !strings.Contains(params["pid"], ",") {
			if validator.CtypeDigit(params["pid"]) {
				ex["api_type"] = params["pid"]
			}
		}
	}

	if params["flag"] != "" {
		ex["flag"] = params["flag"]
	}

	rangeField := ""
	if params["time_flag"] != "" {
		rangeField = betTimeFlags[params["time_flag"]]
	}

	if rangeField == "" {
		return data, errors.New(helper.QueryTermsErr)
	}

	ex[rangeField] = g.Op{"between": exp.NewRangeVal(startAt, endAt)}
	if ty == GameMemberDayGroup {
		ex = g.Ex{
			"report_time": g.Op{"between": exp.NewRangeVal(startAt/1000, endAt/1000)},
		}
		if params["username"] != "" {
			username := strings.ToLower(params["username"])
			ex["username"] = username
			ex["report_type"] = 2
		}

		query, _, _ := dialect.From("tbl_report_game_user").Select(g.COUNT("id")).Where(ex).ToSQL()
		fmt.Println(query)
		err = meta.ReportDB.Get(&data.T, query)
		if err != nil {
			return data, pushLog(err, helper.DBErr)
		}
		if data.T == 0 {
			return data, nil
		}
		query, _, _ = dialect.From("tbl_report_game_user").Select(g.C("report_time").As("api_type"), g.SUM("bet_count").As("total"), g.L("sum(0-company_net_amount)").As("net_amount"), g.SUM("valid_bet_amount").As("valid_bet_amount"),
			g.SUM("bet_amount").As("bet_amount")).Where(ex).GroupBy("report_time").ToSQL()
		fmt.Println(query)
		err = meta.ReportDB.Select(&data.D, query)
		if err != nil {
			return data, pushLog(err, helper.DBErr)
		}

		return data, nil
	}

	if ty == GameMemberTransferGroup {

		if params["username"] != "" {
			username := strings.ToLower(params["username"])
			mb, err := memberInfoCache(username)
			if err != nil {
				return data, errors.New(helper.UsernameErr)
			}
			ex["uid"] = mb.UID
		}
		query, _, _ := dialect.From("tbl_game_record").Select(g.COUNT(g.DISTINCT("api_type"))).Where(ex).ToSQL()
		fmt.Println(query)
		err = meta.TiDB.Get(&data.T, query)
		if err != nil {
			return data, pushLog(err, helper.DBErr)
		}
		if data.T == 0 {
			return data, nil
		}
		query, _, _ = dialect.From("tbl_game_record").Select(g.C("api_type"), g.COUNT("id").As("total"), g.SUM("net_amount").As("net_amount"), g.SUM("valid_bet_amount").As("valid_bet_amount"),
			g.SUM("bet_amount").As("bet_amount")).Where(ex).GroupBy("api_type").ToSQL()
		fmt.Println(query)
		err = meta.TiDB.Select(&data.D, query)
		if err != nil {
			return data, pushLog(err, helper.DBErr)
		}

		return data, nil
	}

	return data, errors.New("ty error")
}

func RecordAdminGame(flag, startTime, endTime string, page, pageSize uint, ex g.Ex) (GameRecordData, error) {

	data := GameRecordData{}

	startAt, err := helper.TimeToLocMs(startTime, loc)
	if err != nil {
		return data, errors.New(helper.DateTimeErr)
	}

	endAt, err := helper.TimeToLocMs(endTime, loc)
	if err != nil {
		return data, errors.New(helper.DateTimeErr)
	}

	if startAt >= endAt {
		return data, errors.New(helper.QueryTimeRangeErr)
	}
	ex["prefix"] = meta.Prefix
	ex["tester"] = 1
	ex[betTimeFlags[flag]] = g.Op{"between": exp.NewRangeVal(startAt, endAt+999)}
	query, _, _ := dialect.From("tbl_game_record").Select(g.COUNT("id")).Where(ex).Limit(1).ToSQL()
	fmt.Println(query)
	err = meta.TiDB.Get(&data.T, query)
	if err != nil {
		return data, pushLog(err, helper.DBErr)
	}
	if data.T == 0 {
		return data, nil
	}

	offset := (page - 1) * pageSize
	query, _, _ = dialect.From("tbl_game_record").Select(colsGameRecord...).Where(ex).Order(g.C(betTimeFlags[flag]).Desc()).Offset(offset).Limit(pageSize).ToSQL()
	fmt.Println(query)
	err = meta.TiDB.Select(&data.D, query)
	if err != nil {
		return data, pushLog(err, helper.DBErr)
	}

	query, _, _ = dialect.From("tbl_game_record").Select(g.SUM("net_amount").As("net_amount"), g.SUM("valid_bet_amount").As("valid_bet_amount"),
		g.SUM("bet_amount").As("bet_amount"), g.SUM("rebate_amount").As("rebate_amount"),
	).Where(ex).ToSQL()
	fmt.Println(query)
	err = meta.TiDB.Get(&data.Agg, query)
	if err != nil {
		return data, pushLog(err, helper.DBErr)
	}

	return data, nil
}

func RecordDeposit(page, pageSize uint, startTime, endTime string, ex g.Ex) (FDepositData, error) {

	data := FDepositData{}
	ex["prefix"] = meta.Prefix
	ex["tester"] = 1
	if startTime != "" && endTime != "" {

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

		ex["created_at"] = g.Op{"between": exp.NewRangeVal(startAt, endAt)}
	}
	query, _, _ := dialect.From("tbl_deposit").Select(g.COUNT("id")).Where(ex).Limit(1).ToSQL()
	fmt.Println(query)
	err := meta.TiDB.Get(&data.T, query)
	if err != nil {
		return data, pushLog(err, helper.DBErr)
	}
	if data.T == 0 {
		return data, nil
	}

	offset := (page - 1) * pageSize
	query, _, _ = dialect.From("tbl_deposit").Select(colsDeposit...).Where(ex).Order(g.C("created_at").Desc()).Offset(offset).Limit(pageSize).ToSQL()
	fmt.Println(query)
	err = meta.TiDB.Select(&data.D, query)
	if err != nil {
		return data, pushLog(err, helper.DBErr)
	}

	return data, nil
}

func RecordDividend(page, pageSize uint, startTime, endTime string, ex g.Ex) (DividendData, error) {

	data := DividendData{}
	ex["state"] = DividendReviewPass
	ex["prefix"] = meta.Prefix
	ex["tester"] = 1
	if startTime != "" && endTime != "" {

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

		ex["review_at"] = g.Op{"between": exp.NewRangeVal(startAt, endAt)}
	}
	query, _, _ := dialect.From("tbl_member_dividend").Select(g.COUNT("id")).Where(ex).Limit(1).ToSQL()
	fmt.Println(query)
	err := meta.TiDB.Get(&data.T, query)
	if err != nil {
		return data, pushLog(err, helper.DBErr)
	}
	if data.T == 0 {
		return data, nil
	}

	offset := (page - 1) * pageSize
	query, _, _ = dialect.From("tbl_member_dividend").Select(colsDividend...).Where(ex).Order(g.C("review_at").Desc()).Offset(offset).Limit(pageSize).ToSQL()
	fmt.Println(query)
	err = meta.TiDB.Select(&data.D, query)
	if err != nil {
		return data, pushLog(err, helper.DBErr)
	}

	return data, nil
}

func RecordRebate(page, pageSize int, startTime, endTime string, ex g.Ex) (RebateData, error) {

	data := RebateData{}
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
	ex["prefix"] = meta.Prefix
	ex["tester"] = 1
	ex["cash_type"] = []int{161, 170, 642, 643, 644, 645, 646, 647, 648, 649}
	t := dialect.From("tbl_balance_transaction")
	if page == 1 {
		query, _, _ := t.Select(g.COUNT("id")).Where(ex).ToSQL()
		err := meta.TiDB.Get(&data.T, query)
		if err != nil {
			return data, pushLog(err, helper.DBErr)
		}

		if data.T == 0 {
			return data, nil
		}
	}

	offset := pageSize * (page - 1)
	query, _, _ := t.Select(colsTransaction...).Where(ex).
		Offset(uint(offset)).Limit(uint(pageSize)).Order(g.C("created_at").Desc()).ToSQL()
	err := meta.TiDB.Select(&data.D, query)
	if err != nil {
		return data, pushLog(err, helper.DBErr)
	}

	return data, nil
}

func RecordAdjust(page, pageSize int, startTime, endTime string, ex g.Ex) (AdjustData, error) {

	data := AdjustData{}
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

		ex["review_at"] = g.Op{"between": exp.NewRangeVal(startAt, endAt)}
	}
	ex["prefix"] = meta.Prefix
	ex["tester"] = 1
	t := dialect.From("tbl_member_adjust")
	if page == 1 {
		query, _, _ := t.Select(g.COUNT("id")).Where(ex).ToSQL()
		err := meta.TiDB.Get(&data.T, query)
		if err != nil {
			return data, pushLog(err, helper.DBErr)
		}

		if data.T == 0 {
			return data, nil
		}
	}

	offset := pageSize * (page - 1)
	query, _, _ := t.Select(colsMemberAdjust...).Where(ex).
		Offset(uint(offset)).Limit(uint(pageSize)).Order(g.C("review_at").Desc()).ToSQL()
	err := meta.TiDB.Select(&data.D, query)
	if err != nil {
		return data, pushLog(err, helper.DBErr)
	}

	return data, nil
}

// 代理管理-记录管理-提款
func RecordWithdraw(page, pageSize int, startTime, endTime, applyStartTime, applyEndTime string, ex g.Ex) (FWithdrawData, error) {

	data := FWithdrawData{}
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

		ex["withdraw_at"] = g.Op{"between": exp.NewRangeVal(startAt, endAt)}
	}
	if applyStartTime != "" && applyEndTime != "" {

		startAt, err := helper.TimeToLoc(applyStartTime, loc)
		if err != nil {
			return data, errors.New(helper.TimeTypeErr)
		}

		endAt, err := helper.TimeToLoc(applyEndTime, loc)
		if err != nil {
			return data, errors.New(helper.TimeTypeErr)
		}

		if startAt >= endAt {
			return data, errors.New(helper.QueryTimeRangeErr)
		}

		ex["created_at"] = g.Op{"between": exp.NewRangeVal(startAt, endAt)}
	}
	ex["prefix"] = meta.Prefix
	ex["tester"] = 1
	t := dialect.From("tbl_withdraw")
	if page == 1 {
		query, _, _ := t.Select(g.COUNT("id")).Where(ex).ToSQL()
		err := meta.TiDB.Get(&data.T, query)
		if err != nil {
			return data, pushLog(err, helper.DBErr)
		}

		if data.T == 0 {
			return data, nil
		}
	}

	offset := pageSize * (page - 1)
	query, _, _ := t.Select(colsWithdraw...).Where(ex).
		Offset(uint(offset)).Limit(uint(pageSize)).Order(g.C("created_at").Desc()).ToSQL()
	err := meta.TiDB.Select(&data.D, query)
	if err != nil {
		return data, pushLog(err, helper.DBErr)
	}

	return data, nil
}

// 处理 提款订单返回数据
func WithdrawDealListData(data FWithdrawData) (WithdrawListData, error) {

	result := WithdrawListData{
		T:   data.T,
		Agg: data.Agg,
	}

	if len(data.D) == 0 {
		return result, nil
	}

	var (
		bids []string
		uids []string
	)

	encFields := []string{"realname"}

	for _, v := range data.D {
		bids = append(bids, v.BID)
		uids = append(uids, v.UID)

		encFields = append(encFields, "bankcard"+v.BID)
	}

	bankcards, err := bankcardListDBByIDs(bids)
	if err != nil {
		return result, pushLog(err, helper.DBErr)
	}

	//fmt.Println("bids = ", bids)
	//fmt.Println("uids = ", uids)

	recs, err := grpc_t.DecryptAll(uids, true, encFields)
	if err != nil {
		fmt.Println("grpc_t.Decrypt err = ", err)
		return result, errors.New(helper.GetRPCErr)
	}

	/*
		d1, err := grpc_t.DecryptAll(uids, true, []string{"realname"})
		if err != nil {
			fmt.Println("grpc_t.Decrypt err = ", err)
			return result, errors.New(helper.GetRPCErr)
		}

		d2, err := grpc_t.DecryptAll(bids, true, []string{"bankcard"})
		if err != nil {
			fmt.Println("grpc_t.Decrypt err = ", err)
			return result, errors.New(helper.GetRPCErr)
		}
	*/
	// 处理返回前端的数据
	for _, v := range data.D {
		w := withdrawCols{
			Withdraw:           v,
			MemberBankNo:       recs[v.UID]["bankcard"+v.BID],
			MemberBankRealName: recs[v.UID]["realname"],
			MemberRealName:     recs[v.UID]["realname"],
		}

		card, ok := bankcards[v.BID]
		if ok {
			w.MemberBankID = card.BankID
			w.MemberBankAddress = card.BankAddress
		}

		result.D = append(result.D, w)
	}

	return result, nil
}

func bankcardListDBByIDs(ids []string) (map[string]BankCard_t, error) {

	data := make(map[string]BankCard_t)
	if len(ids) == 0 {
		return nil, errors.New(helper.UsernameErr)
	}

	ex := g.Ex{"id": ids}
	bankcards, _, err := BankcardsList(ex)
	if err != nil {
		return data, pushLog(err, helper.DBErr)
	}

	for _, v := range bankcards {
		data[v.ID] = v
	}

	return data, nil
}

func BankcardsList(ex g.Ex) ([]BankCard_t, string, error) {

	var data []BankCard_t
	t := dialect.From("tbl_member_bankcard")
	query, _, _ := t.Select(colsBankcard...).Where(ex).Order(g.C("created_at").Desc()).ToSQL()
	err := meta.MerchantDB.Select(&data, query)
	if err != nil && err != sql.ErrNoRows {
		return data, "db", err
	}

	return data, "", nil
}

func RecordGroup(page, pageSize int, startTime, endTime string, ex g.Ex, parentName string) (AgencyTransferRecordData, error) {

	data := AgencyTransferRecordData{}
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

		ex["updated_at"] = g.Op{"between": g.Range(startAt, endAt)}
	}
	orEx := g.Or()
	if parentName != "" {

		orEx = g.Or(
			g.Ex{"after_name": parentName},
			g.Ex{"before_name": parentName},
		)
	}
	ex["prefix"] = meta.Prefix
	t := dialect.From("tbl_agency_transfer_record")
	if page == 1 {
		query, _, _ := t.Select(g.COUNT(1)).Where(g.And(ex, orEx)).ToSQL()
		err := meta.MerchantDB.Get(&data.T, query)
		if err != nil {
			return data, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
		}

		if data.T == 0 {
			return data, nil
		}
	}

	offset := (page - 1) * pageSize
	query, _, _ := t.Select(colsAgencyTransferRecord...).Where(g.And(ex, orEx)).
		Order(g.C("updated_at").Desc()).Offset(uint(offset)).Limit(uint(pageSize)).ToSQL()
	err := meta.MerchantDB.Select(&data.D, query)
	if err != nil && err != sql.ErrNoRows {
		return data, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	return data, nil
}

func RecordIssuse(id string) ([]string, error) {

	tableName := "tbl_vncp_plan_issues"
	var result []string

	ex := g.Ex{
		"plan_id": id,
	}
	build := dialect.From(tableName).Where(ex)

	build = build.Select(
		"id",
	).Order(g.C("created_at").Desc())
	query, _, _ := build.ToSQL()
	fmt.Println(query)
	err := meta.MerchantDB.Select(&result, query)
	if err != nil {
		return result, err
	}
	return result, nil
}

func RecordOrder(page, pageSize int, ex g.Ex) (OrderData, error) {

	data := OrderData{}
	ex["prefix"] = meta.Prefix
	t := dialect.From("tbl_vncp_orders")
	if page == 1 {
		query, _, _ := t.Select(g.COUNT("id")).Where(ex).ToSQL()
		fmt.Println(query)
		err := meta.MerchantDB.Get(&data.T, query)
		if err != nil {
			return data, pushLog(err, helper.DBErr)
		}

		if data.T == 0 {
			return data, nil
		}
	}

	offset := pageSize * (page - 1)
	query, _, _ := t.Select(g.C("username"), g.C("pay_amount"), g.C("bonus").As("bonus")).Where(ex).
		Offset(uint(offset)).Limit(uint(pageSize)).Order(g.C("created_at").Desc()).ToSQL()
	fmt.Println(query)
	err := meta.MerchantDB.Select(&data.D, query)
	if err != nil {
		return data, pushLog(err, helper.DBErr)
	}

	for i := 0; i < len(data.D); i++ {
		pay, _ := decimal.NewFromString(data.D[i].PayAmount)
		bonus, _ := decimal.NewFromString(data.D[i].Bonus)
		data.D[i].NetAmount = bonus.Sub(pay).StringFixed(4)
	}
	return data, nil
}
