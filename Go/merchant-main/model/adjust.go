package model

import (
	"errors"
	"fmt"
	"merchant/contrib/helper"
	"time"

	g "github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/shopspring/decimal"
)

// AdjustData 账户调整列表response structure
type AdjustData struct {
	D      []MemberAdjust `json:"d"`
	T      int64          `json:"t"`
	S      uint           `json:"s"`
	Amount float64        `json:"amount"`
}

// AdjustInsert 新增账户调整申请
func AdjustInsert(data MemberAdjust) error {

	data.Prefix = meta.Prefix

	amount := decimal.NewFromFloat(data.Amount)
	// 上分不需要冻结中心钱包的钱
	if data.AdjustMode == AdjustUpMode {

		query, _, _ := dialect.Insert("tbl_member_adjust").Rows(data).ToSQL()
		_, err := meta.MerchantDB.Exec(query)
		if err != nil {
			return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
		}

		_ = PushMerchantNotify(adjustReviewFmt, data.ApplyName, data.Username, amount.String())

		return nil
	}

	record := g.Record{
		"id":             helper.GenId(),
		"uid":            data.UID,
		"username":       data.Username,
		"prefix":         data.Prefix,
		"ty":             0, // 后台调整
		"amount":         amount.String(),
		"adjust_type":    data.AdjustType,
		"adjust_mode":    data.AdjustMode,
		"is_turnover":    data.IsTurnover,
		"turnover_multi": data.TurnoverMulti,
		"apply_remark":   data.ApplyRemark,
		"images":         data.Images,
		"state":          AdjustReviewing, // 状态:256=审核中,257=同意, 258=拒绝
		"apply_at":       data.ApplyAt,
		"apply_uid":      data.ApplyUid,  // 申请人
		"apply_name":     data.ApplyName, // 申请人
		"top_uid":        data.TopUid,
		"top_name":       data.TopName,
		"parent_uid":     data.ParentUid,
		"parent_name":    data.ParentName,
		"tester":         data.Tester,
	}
	// 下分
	err := AdjustUpDownPoint(meta.Prefix, data.UID, data.Username, data.AdjustType, DownPointApply, amount, record)
	if err != nil {
		return err
	}

	_ = PushMerchantNotify(adjustReviewFmt, data.ApplyName, data.Username, amount.String())

	return nil
}

// AdjustList 账户调整列表
func AdjustList(startTime, endTime string, ex g.Ex, page, pageSize int) (AdjustData, error) {

	data := AdjustData{}

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

	ex["apply_at"] = g.Op{
		"between": exp.NewRangeVal(startAt, endAt),
	}
	ex["prefix"] = meta.Prefix

	t := dialect.From("tbl_member_adjust")
	if page == 1 {
		query, _, _ := t.Select(g.COUNT(1).As("t"), g.COALESCE(g.SUM("amount"), 0).As("amount")).Where(ex).ToSQL()
		err = meta.MerchantDB.Get(&data, query)
		if err != nil {
			return data, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
		}

		if data.T == 0 {
			return data, nil
		}
	}

	offset := (page - 1) * pageSize
	query, _, _ := t.Select(colsMemberAdjust...).Where(ex).Order(g.C("apply_at").Desc()).Offset(uint(offset)).Limit(uint(pageSize)).ToSQL()
	err = meta.MerchantDB.Select(&data.D, query)
	if err != nil {
		return data, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	data.S = uint(pageSize)

	return data, nil
}

// 中心钱包 下分申请 下分申请通过 下分申请拒绝
// 上分申请通过 adjust 参数为更新调整记录状态
// 下分申请 adjust 参数为新增调整记录
// 下分申请通过 adjust 参数为更新调整记录状态
// 下分申请拒绝 adjust 参数为更新调整记录状态
func AdjustUpDownPoint(prefix, uid, username string, adjustType, flag int, money decimal.Decimal, adjust g.Record) error {

	var (
		balanceAfter decimal.Decimal
		balance      decimal.Decimal
		lockAmount   decimal.Decimal
	)

	//1、判断金额是否合法
	if money.Cmp(zero) == -1 {
		return errors.New(helper.AmountErr)
	}

	flags := map[int]bool{
		DownPointApply:       true,
		DownPointApplyPass:   true,
		DownPointApplyReject: true,
		UpPointApplyPass:     true,
	}

	// 非中心钱包转锁定钱包且非锁定钱包转中心钱包
	if _, ok := flags[flag]; !ok {
		return errors.New(helper.ParamErr)
	}

	tx, err := meta.MerchantDB.Begin()
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	record := g.Record{}
	switch flag {
	case UpPointApplyPass:
		// 获取钱包余额
		mb, err := MemberBalance(username)
		if err != nil {
			return err
		}

		// 中心钱包余额
		balance, _ = decimal.NewFromString(mb.Balance)
		// 中心钱包转出
		balanceAfter = balance.Add(money)

		r := g.Record{
			"state":          AdjustReviewPass, //审核状态:1=审核中,2=审核通过,3=审核未通过
			"hand_out_state": AdjustSuccess,    //上下分状态 1 失败 2成功 3场馆上分处理中
			"review_remark":  adjust["review_remark"],
			"review_at":      adjust["review_at"],
			"review_uid":     adjust["review_uid"],
			"review_name":    adjust["review_name"],
		}
		// 更新调整记录状态
		query, _, _ := dialect.Update("tbl_member_adjust").Set(r).Where(g.Ex{"id": adjust["id"]}).ToSQL()
		_, err = tx.Exec(query)
		if err != nil {
			_ = tx.Rollback()
			return pushLog(err, helper.DBErr)
		}

		//4、新增账变记录
		trans := MemberTransaction{
			AfterAmount:  balanceAfter.String(),
			Amount:       money.String(),
			BeforeAmount: balance.String(),
			BillNo:       adjust["id"].(string),
			CreatedAt:    time.Now().UnixMilli(),
			ID:           helper.GenId(),
			CashType:     helper.TransactionUpPoint,
			UID:          uid,
			Username:     username,
			Prefix:       meta.Prefix,
		}
		// 离线转卡存款
		if adjustType == 3 {
			trans.CashType = helper.TransactionOfflineDeposit
		}
		query, _, _ = dialect.Insert("tbl_balance_transaction").Rows(trans).ToSQL()
		_, err = tx.Exec(query)
		if err != nil {
			_ = tx.Rollback()
			return pushLog(err, helper.DBErr)
		}

		member, err := MemberFindOne(username)
		if err != nil {
			_ = tx.Rollback()
			return pushLog(err, helper.DBErr)
		}

		// 离线转卡存款
		if adjustType == 3 {
			now := fmt.Sprintf("%d", time.Now().Unix())
			dr := g.Record{
				"id":            helper.GenId(),
				"prefix":        prefix,
				"oid":           0,
				"uid":           uid,
				"username":      username,
				"channel_id":    11, // 通道id
				"cid":           12, // 线下转卡类型
				"pid":           0,
				"amount":        money.String(),
				"state":         DepositSuccess,
				"automatic":     "1",
				"created_at":    now,
				"created_uid":   "0",
				"created_name":  "",
				"confirm_at":    now,
				"confirm_uid":   "0",
				"confirm_name":  "",
				"review_remark": "",
				"tester":        member.Tester,
				"flag":          3,
			}
			query, _, _ = dialect.Insert("tbl_deposit").Rows(dr).ToSQL()
			_, err = tx.Exec(query)
			if err != nil {
				_ = tx.Rollback()
				return pushLog(err, helper.DBErr)
			}

			// 存款成功发送到队列

			//param := map[string]interface{}{
			//	"bean_ty":            "4",
			//	"username":           username,
			//	"amount":             money.String(),
			//	"deposit_created_at": now,
			//	"deposit_success_at": now,
			//}
			//
			//_, err = BeanPut("promo", param, 0)
			//if err != nil {
			//	fmt.Println("user invite BeanPut err:", err.Error())
			//}
			member, err := MemberFindOne(username)
			if err != nil {
				if member.FirstDepositAt == 0 {
					rec := g.Record{
						"first_deposit_at":     now,
						"first_deposit_amount": money,
					}
					ex := g.Ex{
						"uid": member.UID,
					}
					query, _, _ := dialect.Update("tbl_members").Set(rec).Where(ex).ToSQL()
					fmt.Printf("memberFirstDeposit Update: %v\n", query)

					_, err := meta.MerchantDB.Exec(query)
					if err != nil {
						fmt.Println("update member first_amount err:", err.Error())
					}
				} else if member.SecondDepositAt == 0 {
					rec := g.Record{
						"second_deposit_at":     now,
						"second_deposit_amount": money,
					}
					ex := g.Ex{
						"uid": member.UID,
					}
					query, _, _ := dialect.Update("tbl_members").Set(rec).Where(ex).ToSQL()
					fmt.Printf("memberSecondDeposit Update: %v\n", query)

					_, err := meta.MerchantDB.Exec(query)
					if err != nil {
						fmt.Println("update member second_amount err:", err.Error())
					}
				}

			}
		}

		// 中心钱包上分
		record["balance"] = g.L(fmt.Sprintf("balance+%s", money.String()))
	case DownPointApply:
		// 获取钱包余额
		mb, err := MemberBalance(username)
		if err != nil {
			return err
		}

		// 检查中心钱包余额
		balance, _ = decimal.NewFromString(mb.Balance)
		if balance.Sub(money).IsNegative() {
			return errors.New(helper.LackOfBalance)
		}

		// 中心钱包转入
		balanceAfter = balance.Sub(money)
		adjust["amount"] = "-" + adjust["amount"].(string)
		// 新增调整记录
		query, _, _ := dialect.Insert("tbl_member_adjust").Rows(adjust).ToSQL()
		_, err = tx.Exec(query)
		if err != nil {
			_ = tx.Rollback()
			return pushLog(err, helper.DBErr)
		}

		//4、新增账变记录
		trans := MemberTransaction{
			AfterAmount:  balanceAfter.String(),
			Amount:       money.String(),
			BeforeAmount: balance.String(),
			BillNo:       adjust["id"].(string),
			CreatedAt:    time.Now().UnixNano() / 1e6,
			ID:           helper.GenId(),
			CashType:     helper.TransactionDownPoint,
			UID:          uid,
			Username:     username,
			Prefix:       meta.Prefix,
		}

		query, _, _ = dialect.Insert("tbl_balance_transaction").Rows(trans).ToSQL()
		_, err = tx.Exec(query)
		if err != nil {
			_ = tx.Rollback()
			return pushLog(err, helper.DBErr)
		}

		// 中心钱包下分
		record["balance"] = g.L(fmt.Sprintf("balance-%s", money.String()))
		// 锁定钱包上分
		record["lock_amount"] = g.L(fmt.Sprintf("lock_amount+%s", money.String()))

	case DownPointApplyPass:
		// 获取钱包余额
		mb, err := MemberBalance(username)
		if err != nil {
			return err
		}

		// 检查锁定钱包余额
		balance, _ = decimal.NewFromString(mb.LockAmount)
		if balance.Sub(money).IsNegative() {
			return errors.New(helper.LackOfBalance)
		}

		r := g.Record{
			"state":          AdjustReviewPass, //审核状态:1=审核中,2=审核通过,3=审核未通过
			"hand_out_state": AdjustSuccess,    //上下分状态 1 失败 2成功 3场馆上分处理中
			"review_remark":  adjust["review_remark"],
			"review_at":      adjust["review_at"],
			"review_uid":     adjust["review_uid"],
			"review_name":    adjust["review_name"],
		}
		// 更新调整记录状态
		query, _, _ := dialect.Update("tbl_member_adjust").Set(r).Where(g.Ex{"id": adjust["id"]}).ToSQL()
		_, err = tx.Exec(query)
		if err != nil {
			_ = tx.Rollback()
			return pushLog(err, helper.DBErr)
		}

		// 锁定钱包下分
		record["lock_amount"] = g.L(fmt.Sprintf("lock_amount-%s", money.String()))
	case DownPointApplyReject:
		// 获取钱包余额
		mb, err := MemberBalance(username)
		if err != nil {
			return err
		}

		// 检查锁定钱包余额
		balance, _ = decimal.NewFromString(mb.Balance)
		lockAmount, _ = decimal.NewFromString(mb.LockAmount)
		if lockAmount.Sub(money).IsNegative() {
			return errors.New(helper.LackOfBalance)
		}

		// 中心钱包转出
		balanceAfter = balance.Add(money)
		r := g.Record{
			"state":         AdjustReviewReject, //审核状态:1=审核中,2=审核通过,3=审核未通过
			"review_remark": adjust["review_remark"],
			"review_at":     adjust["review_at"],
			"review_uid":    adjust["review_uid"],
			"review_name":   adjust["review_name"],
		}
		// 更新调整记录状态
		query, _, _ := dialect.Update("tbl_member_adjust").Set(r).Where(g.Ex{"id": adjust["id"]}).ToSQL()
		_, err = tx.Exec(query)
		if err != nil {
			_ = tx.Rollback()
			return pushLog(err, helper.DBErr)
		}

		//4、新增账变记录
		trans := MemberTransaction{
			AfterAmount:  balanceAfter.String(),
			Amount:       money.String(),
			BeforeAmount: balance.String(),
			BillNo:       adjust["id"].(string),
			CreatedAt:    time.Now().UnixNano() / 1e6,
			ID:           helper.GenId(),
			CashType:     helper.TransactionDownPointBack,
			UID:          uid,
			Username:     username,
			Prefix:       meta.Prefix,
		}

		query, _, _ = dialect.Insert("tbl_balance_transaction").Rows(trans).ToSQL()
		_, err = tx.Exec(query)
		if err != nil {
			_ = tx.Rollback()
			return pushLog(err, helper.DBErr)
		}

		// 中心钱包上分
		record["balance"] = g.L(fmt.Sprintf("balance+%s", money.String()))
		// 锁定钱包下分
		record["lock_amount"] = g.L(fmt.Sprintf("lock_amount-%s", money.String()))
	}

	ex := g.Ex{
		"uid": uid,
	}
	query, _, _ := dialect.Update("tbl_members").Set(record).Where(ex).ToSQL()
	res, err := tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		return pushLog(err, helper.DBErr)
	}

	if r, _ := res.RowsAffected(); r == 0 {
		_ = tx.Rollback()
		return pushLog(err, helper.DBErr)
	}

	err = tx.Commit()
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	_ = MemberUpdateCache(uid, "")

	return nil
}

// AdjustReview 账户调整-审核
func AdjustReview(state int, record g.Record) error {

	data := MemberAdjust{}

	query, _, _ := dialect.From("tbl_member_adjust").
		Select(colsMemberAdjust...).Where(g.Ex{"id": record["id"]}).Limit(1).ToSQL()
	err := meta.MerchantDB.Get(&data, query)
	if err != nil {
		return errors.New(helper.DBErr)
	}

	// 只有审核中状态的才能操作
	if data.State != AdjustReviewing {
		return errors.New(helper.StateParamErr)
	}

	amount := decimal.NewFromFloat(data.Amount)
	// 下分只有中心钱包
	if data.AdjustMode == AdjustDownMode {

		flag := DownPointApplyReject
		if state == AdjustReviewPass {
			flag = DownPointApplyPass
		}

		return AdjustUpDownPoint(meta.Prefix, data.UID, data.Username, data.AdjustType, flag, amount.Abs(), record)
	}

	// 后面都是上分 先处理拒绝的业务逻辑 再分别处理场馆和中心钱包同意的逻辑
	if state == AdjustReviewReject {

		// 更新调整记录状态
		query, _, _ = dialect.Update("tbl_member_adjust").Set(record).Where(g.Ex{"id": record["id"]}).ToSQL()
		_, err = meta.MerchantDB.Exec(query)
		if err != nil {
			return errors.New(helper.DBErr)
		}

		return nil
	}

	// 中心钱包上分申请同意逻辑
	err = AdjustUpDownPoint(meta.Prefix, data.UID, data.Username, data.AdjustType, UpPointApplyPass, amount, record)
	if err != nil {
		return err
	}

	return nil
}
