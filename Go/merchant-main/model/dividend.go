package model

import (
	"errors"
	"fmt"
	g "github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/shopspring/decimal"
	"merchant/contrib/helper"
	"time"
)

func DividendInsert(data g.Record) error {

	data["prefix"] = meta.Prefix
	query, _, _ := dialect.Insert("tbl_member_dividend").Rows(data).ToSQL()
	fmt.Println(query)
	_, err := meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(fmt.Errorf("query [%s], error [%s]", query, err.Error()), helper.DBErr)
	}

	_ = PushMerchantNotify(dividendReviewFmt, data["apply_name"].(string), data["username"].(string), data["amount"].(string))

	return nil
}

func DividendList(page, pageSize int, startTime, endTime, reviewStartTime, reviewEndTime string, ex g.Ex) (DividendData, error) {

	data := DividendData{}
	// 没有查询条件  startTime endTime 必填
	if len(ex) == 0 && (startTime == "" || endTime == "") {
		return data, errors.New(helper.QueryTermsErr)
	}

	if startTime != "" && endTime != "" {

		startAt, err := helper.TimeToLocMs(startTime, loc)
		if err != nil {
			return data, errors.New(helper.DateTimeErr)
		}

		endAt, err := helper.TimeToLocMs(endTime, loc)
		if err != nil {
			return data, errors.New(helper.TimeTypeErr)
		}

		if startAt >= endAt {
			return data, errors.New(helper.QueryTimeRangeErr)
		}

		ex["apply_at"] = g.Op{"between": exp.NewRangeVal(startAt, endAt)}
	}

	if reviewStartTime != "" && reviewEndTime != "" {

		rStart, err := helper.TimeToLoc(reviewStartTime, loc)
		if err != nil {
			return data, errors.New(helper.TimeTypeErr)
		}

		rEnd, err := helper.TimeToLoc(reviewEndTime, loc)
		if err != nil {
			return data, errors.New(helper.TimeTypeErr)
		}

		if rStart >= rEnd {
			return data, errors.New(helper.QueryTimeRangeErr)
		}

		ex["review_at"] = g.Op{"between": exp.NewRangeVal(rStart, rEnd)}
	}
	ex["prefix"] = meta.Prefix

	t := dialect.From("tbl_member_dividend")
	if page == 1 {
		query, _, _ := t.Select(g.COUNT("id")).Where(ex).ToSQL()
		err := meta.MerchantDB.Get(&data.T, query)
		if err != nil {
			return data, pushLog(err, helper.DBErr)
		}

		if data.T == 0 {
			return data, nil
		}
	}

	offset := pageSize * (page - 1)
	query, _, _ := t.Select(colsDividend...).Where(ex).
		Offset(uint(offset)).Limit(uint(pageSize)).Order(g.C("apply_at").Desc()).ToSQL()
	err := meta.MerchantDB.Select(&data.D, query)
	if err != nil {
		return data, pushLog(err, helper.DBErr)
	}

	return data, nil
}

func DividendReview(state int, ts int64, adminID, adminName, reviewRemark string, ids []string) error {

	ex := g.Ex{
		"id":     ids,
		"prefix": meta.Prefix,
		"state":  DividendReviewing,
	}
	// 批量/单条不通过
	if state == DividendReviewReject {
		record := g.Record{
			"state":         DividendReviewReject,
			"review_remark": reviewRemark,
			"review_at":     ts,
			"review_uid":    adminID,
			"review_name":   adminName,
		}
		query, _, _ := dialect.Update("tbl_member_dividend").Set(record).Where(ex).ToSQL()
		fmt.Println(query)
		_, err := meta.MerchantDB.Exec(query)
		if err != nil {
			return pushLog(err, helper.DBErr)
		}

		return nil
	}

	var (
		data []MemberDividend
	)
	query, _, _ := dialect.From("tbl_member_dividend").Where(ex).ToSQL()
	fmt.Println(query)
	err := meta.MerchantDB.Select(&data, query)
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	var (
		pids      []string
		usernames []string
	)
	for _, v := range data {

		fmt.Println(v)
		mb, err := MemberBalance(v.Username)
		if err != nil {
			_ = pushLog(err, helper.BalanceErr)
			continue
		}

		tx, err := meta.MerchantDB.Begin()
		if err != nil {
			return pushLog(err, helper.DBErr)
		}

		record := g.Record{
			"state":         DividendReviewPass,
			"review_remark": reviewRemark,
			"review_at":     ts,
			"review_uid":    adminID,
			"review_name":   adminName,
		}
		query, _, _ = dialect.Update("tbl_member_dividend").Set(record).Where(g.Ex{"id": v.ID}).ToSQL()
		fmt.Println(query)
		_, err = tx.Exec(query)
		if err != nil {
			_ = tx.Rollback()
			_ = pushLog(err, helper.DBErr)
			continue
		}

		// 活动红利
		if v.PID != "" && v.Ty == DividendPromo {
			r := g.Record{
				"id":            helper.GenId(),
				"flag":          "static",
				"prefix":        meta.Prefix,
				"uid":           v.UID,           // 会员账号
				"username":      v.Username,      // 会员名称
				"level":         v.Level,         // 会员等级
				"top_uid":       v.TopUid,        // 总代账号uid
				"top_name":      v.TopName,       // 总代账号名
				"parent_uid":    v.ParentUid,     // 上级代理账号uid
				"parent_name":   v.ParentName,    // 上级代理账号名
				"pid":           v.PID,           // 活动ID
				"title":         v.PTitle,        // 活动名称
				"amount":        v.Amount,        // 活动金额
				"bonus_type":    1,               // 彩金类型 1 固定金额 2 百分比
				"bonus_rate":    0,               // 彩金比例 bonus_type=2时使用
				"bonus":         v.Amount,        // 活动彩金
				"flow":          v.WaterFlow,     // 流水金额
				"multiple":      v.WaterMultiple, // 流水倍数
				"state":         2,               // 已发放
				"created_at":    ts,              // 记录的创建时间
				"review_at":     ts,              // 记录的最后更新时间
				"review_uid":    adminID,         // 活动审批人员ID 如果需要人工审批同意，维护这个字段
				"review_name":   adminName,       // 审核人名
				"inspect_at":    0,               // 稽查时间
				"inspect_uid":   "0",             // 流水稽查ID
				"inspect_name":  "",              // 稽查人名
				"inspect_state": 1,               // 稽查状态 1 待稽查 2 完成流水 3 余额清零
			}
			query, _, _ = dialect.Insert("tbl_promo_record").Rows(r).ToSQL()
			fmt.Println(query)
			_, err = tx.Exec(query)
			if err != nil {
				_ = tx.Rollback()
				_ = pushLog(err, helper.DBErr)
				continue
			}

			pids = append(pids, v.PID)
			usernames = append(usernames, v.Username)
		}

		balance, _ := decimal.NewFromString(mb.Balance)
		amount := decimal.NewFromFloat(v.Amount)

		// 中心钱包转出
		balanceAfter := balance.Add(amount)

		//1、判断金额是否合法
		if balanceAfter.IsNegative() {
			_ = tx.Rollback()
			_ = pushLog(fmt.Errorf("after amount : %s less than 0", balanceAfter.String()), helper.BalanceErr)
			continue
		}

		trans := MemberTransaction{
			AfterAmount:  balanceAfter.String(),
			Amount:       amount.String(),
			BeforeAmount: balance.String(),
			BillNo:       v.ID,
			CreatedAt:    time.Now().UnixMilli(),
			ID:           helper.GenId(),
			CashType:     helper.TransactionDividend,
			UID:          v.UID,
			Username:     v.Username,
			Prefix:       meta.Prefix,
		}
		query, _, _ = dialect.Insert("tbl_balance_transaction").Rows(trans).ToSQL()
		_, err = tx.Exec(query)
		fmt.Println(query)
		if err != nil {
			_ = tx.Rollback()
			_ = pushLog(err, helper.DBErr)
			continue
		}

		op := "+"
		// 红利金额为负数
		if amount.IsNegative() {
			op = "-"
		}
		// 中心钱包上下分
		record = g.Record{
			"balance": g.L(fmt.Sprintf("balance%s%s", op, amount.Abs().String())),
		}
		ex = g.Ex{
			"uid": v.UID,
		}
		query, _, _ = dialect.Update("tbl_members").Set(record).Where(ex).ToSQL()
		fmt.Println(query)
		_, err = tx.Exec(query)
		if err != nil {
			_ = tx.Rollback()
			_ = pushLog(err, helper.DBErr)
			continue
		}

		_ = tx.Commit()
		key := meta.Prefix + ":member:" + v.Username
		_ = meta.MerchantRedis.HSet(ctx, key, "balance", balanceAfter.String()).Err()
	}

	if len(usernames) > 0 {
		for k, v := range usernames {
			title := "Quý Khách Của P3 Thân Mến"
			content := "Khuyến Mãi Đã Được Tặng Vào Tài Khoản Của Bạn,Vui Lòng KIểm Tra Ngay,Nếu Bạn Có Bất Cứ Thắc Mắc Vấn Đề Gì Vui Lòng Liên Hệ CSKH Để Biết Thêm Chi Tiết .\n【P3】Chúc Bạn Cược Đâu Thắng Đó !!"
			err = messageSend(pids[k], title, "", content, "system", meta.Prefix, 0, 0, 2, []string{v})
			if err != nil {
				_ = pushLog(err, helper.ESErr)
			}
		}
	}

	return nil
}

// 更新红利
func DividendUpdate(ex g.Ex, record g.Record) error {

	return dividendUpdate(ex, record)
}

func DividendGetState(id string) (int, error) {

	var state int
	query, _, _ := dialect.From("tbl_member_dividend").Select("state").Where(g.Ex{"id": id}).ToSQL()
	err := meta.MerchantDB.Get(&state, query)
	if err != nil {
		return state, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	return state, nil
}

// 会员红利记录更新
func dividendUpdate(ex g.Ex, record g.Record) error {

	ex["prefix"] = meta.Prefix
	t := dialect.Update("tbl_member_dividend")
	query, _, _ := t.Set(record).Where(ex).ToSQL()
	_, err := meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	return nil
}
