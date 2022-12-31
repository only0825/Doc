package model

import (
	"database/sql"
	"errors"
	"fmt"
	g "github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/valyala/fasthttp"
	"merchant/contrib/helper"
)

// 黑名单列表
func BlacklistList(page, pageSize uint, startTime, endTime string, ty int, ex g.Ex) (BlacklistData, error) {

	data := BlacklistData{}

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
	ex["prefix"] = meta.Prefix
	t := dialect.From("tbl_blacklist")
	if page == 1 {

		query, _, _ := t.Select(g.COUNT(1)).Where(ex).ToSQL()
		err := meta.MerchantDB.Get(&data.T, query)
		if err != nil {
			return data, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
		}

		if data.T == 0 {
			return data, nil
		}
	}

	data.S = pageSize
	offset := (page - 1) * pageSize
	query, _, _ := t.Select(colsBlacklist...).Where(ex).
		Order(g.C("created_at").Desc()).Offset(offset).Limit(pageSize).ToSQL()
	err := meta.MerchantDB.Select(&data.D, query)
	if err != nil {
		return data, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	return data, nil
}

// 黑名单添加
func BlacklistInsert(fCtx *fasthttp.RequestCtx, ty int, value string, record g.Record) error {

	var (
		data []BankCard_t
		key  string
	)
	user, err := AdminToken(fCtx)
	if err != nil {
		return errors.New(helper.AccessTokenExpires)
	}

	ex := g.Ex{
		"ty":     ty,
		"value":  value,
		"prefix": meta.Prefix,
	}
	ok, err := BlacklistExist(ex)
	if err != nil {
		return err
	}

	if ok {
		return errors.New(helper.RecordExistErr)
	}

	record["created_at"] = fCtx.Time().In(loc).Unix()
	record["created_uid"] = user["id"]
	record["created_name"] = user["name"]
	record["prefix"] = meta.Prefix

	query, _, _ := dialect.Insert("tbl_blacklist").Rows(record).ToSQL()
	_, err = meta.MerchantDB.Exec(query)
	if err != nil {
		return errors.New(helper.DBErr)
	}

	switch ty {
	case TyDevice:
		key = fmt.Sprintf("%s:merchant:device_blacklist", meta.Prefix)
	case TyIP:
		key = fmt.Sprintf("%s:merchant:ip_blacklist", meta.Prefix)
	case TyPhone:
		key = fmt.Sprintf("%s:merchant:phone_blacklist", meta.Prefix)
	case TyBankcard:
		key = fmt.Sprintf("%s:merchant:bankcard_blacklist", meta.Prefix)
	case TyRebate:
		key = fmt.Sprintf("%s:merchant:rebate_blacklist", meta.Prefix)
	case TyCGRebate:
		key = fmt.Sprintf("%s:merchant:cgrebate_blacklist", meta.Prefix)
	case TyPromoteLink:
		key = fmt.Sprintf("%s:merchant:link_blacklist", meta.Prefix)
	case TyWhiteIP:
		key = fmt.Sprintf("%s:merchant:ip_whitelist", meta.Prefix)
	}

	meta.MerchantRedis.SAdd(ctx, key, value).Val()

	if ty == TyBankcard {
		ex = g.Ex{
			"prefix":         meta.Prefix,
			"bank_card_hash": fmt.Sprintf("%d", MurmurHash(value, 0)),
		}
		recs := g.Record{
			"state": "3",
		}
		query, _, _ = dialect.Update("tbl_member_bankcard").Set(recs).Where(ex).ToSQL()
		_, err = meta.MerchantDB.Exec(query)
		if err != nil {
			return errors.New(helper.DBErr)
		}

		query, _, _ = dialect.From("tbl_member_bankcard").Select(colsBankcard...).Where(ex).ToSQL()
		err = meta.MerchantDB.Select(&data, query)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		for _, v := range data {
			BankcardUpdateCache(v.Username)
		}
	}

	return nil
}

// 黑名单更新备注
func BlacklistUpdate(id, remark string) error {

	ex := g.Ex{
		"id":     id,
		"prefix": meta.Prefix,
	}
	record := g.Record{
		"remark": remark,
	}
	query, _, _ := dialect.Update("tbl_blacklist").Set(record).Where(ex).ToSQL()
	_, err := meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	return nil
}

// 删除记录
func BlacklistDelete(id string) error {

	ex := g.Ex{
		"id":     id,
		"prefix": meta.Prefix,
	}

	data := Blacklist{}
	query, _, _ := dialect.From("tbl_blacklist").Select(colsBlacklist...).Where(ex).ToSQL()
	err := meta.MerchantDB.Get(&data, query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	query, _, _ = dialect.Delete("tbl_blacklist").Where(ex).ToSQL()
	_, err = meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	// id 银行卡删除黑名单
	if data.Ty == TyBankcard {
		// 银行卡从黑名单移出后，修改卡状态为停用
		valueHash := fmt.Sprintf("%d", MurmurHash(data.Value, 0))
		ex = g.Ex{
			"prefix":         meta.Prefix,
			"bank_card_hash": valueHash,
		}
		recs := g.Record{
			"state": "2",
		}
		query, _, _ = dialect.Update("tbl_member_bankcard").Set(recs).Where(ex).ToSQL()
		_, err := meta.MerchantDB.Exec(query)
		if err != nil {
			return errors.New(helper.DBErr)
		}

		// 从黑名单删除银行卡后，更新redis 黑名单的银行卡信息=
		key := fmt.Sprintf("%s:merchant:bankcard_blacklist", meta.Prefix)
		cmd := meta.MerchantRedis.SRem(ctx, key, data.Value)
		err = cmd.Err()
		if err != nil {
			return errors.New(err.Error())
		}
	}

	// 更新结束
	_ = LoadBlacklists(data.Ty)

	return nil
}

// 满足条件的黑名单数量
func BlacklistExist(ex g.Ex) (bool, error) {

	var count int
	query, _, _ := dialect.From("tbl_blacklist").Select(g.COUNT("id")).Where(ex).Limit(1).ToSQL()
	err := meta.MerchantDB.Get(&count, query)
	if err != nil {
		return false, pushLog(err, helper.DBErr)
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func LoadBlacklists(ty int) error {

	var data []Blacklist

	if ty != 0 {
		ex := g.Ex{"ty": ty}
		query, _, _ := dialect.From("tbl_blacklist").Select(colsBlacklist...).Where(ex).ToSQL()
		query = "/* master */ " + query
		fmt.Println(query)
		err := meta.MerchantDB.Select(&data, query)
		if err != nil {
			return pushLog(err, helper.DBErr)
		}
	} else {
		query, _, _ := dialect.From("tbl_blacklist").Select(colsBlacklist...).ToSQL()
		query = "/* master */ " + query
		fmt.Println(query)
		err := meta.MerchantDB.Select(&data, query)
		if err != nil {
			return pushLog(err, helper.DBErr)
		}
	}

	pipe := meta.MerchantRedis.Pipeline()
	defer pipe.Close()

	deviceKey := fmt.Sprintf("%s:merchant:device_blacklist", meta.Prefix)
	ipKey := fmt.Sprintf("%s:merchant:ip_blacklist", meta.Prefix)
	phoneKey := fmt.Sprintf("%s:merchant:phone_blacklist", meta.Prefix)
	bankcardKey := fmt.Sprintf("%s:merchant:bankcard_blacklist", meta.Prefix)
	rebateKey := fmt.Sprintf("%s:merchant:rebate_blacklist", meta.Prefix)
	cgRebateKey := fmt.Sprintf("%s:merchant:cgrebate_blacklist", meta.Prefix)
	linkKey := fmt.Sprintf("%s:merchant:link_blacklist", meta.Prefix)
	ipWhiteKey := fmt.Sprintf("%s:merchant:ip_whitelist", meta.Prefix)

	if ty != 0 {
		key := ""
		switch ty {
		case TyDevice:
			key = deviceKey
		case TyIP:
			key = ipKey
		case TyPhone:
			key = phoneKey
		case TyBankcard:
			key = bankcardKey
		case TyRebate:
			key = rebateKey
		case TyCGRebate:
			key = cgRebateKey
		case TyPromoteLink:
			key = linkKey
		case TyWhiteIP:
			key = ipWhiteKey
		}
		pipe.Unlink(ctx, key)
	} else {
		pipe.Unlink(ctx, deviceKey)
		pipe.Unlink(ctx, ipKey)
		pipe.Unlink(ctx, phoneKey)
		pipe.Unlink(ctx, bankcardKey)
		pipe.Unlink(ctx, rebateKey)
		pipe.Unlink(ctx, cgRebateKey)
		pipe.Unlink(ctx, linkKey)
		pipe.Unlink(ctx, ipWhiteKey)
	}

	for _, v := range data {
		key := ""
		switch v.Ty {
		case TyDevice:
			key = deviceKey
		case TyIP:
			key = ipKey
		case TyPhone:
			key = phoneKey
		case TyBankcard:
			key = bankcardKey
		case TyRebate:
			key = rebateKey
		case TyCGRebate:
			key = cgRebateKey
		case TyPromoteLink:
			key = linkKey
		case TyWhiteIP:
			key = ipWhiteKey
		}

		pipe.SAdd(ctx, key, v.Value)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return pushLog(err, helper.RedisErr)
	}

	return nil
}

// 解锁手机号码（确认没有会员账号绑定的情况）
func BlacklistClearPhone(phone string) error {

	var count uint64
	ex := g.Ex{
		"phone_hash": fmt.Sprintf("%d", MurmurHash(phone, 0)),
		"prefix":     meta.Prefix,
	}
	query, _, _ := dialect.From("tbl_members").Select(g.COUNT("uid")).Where(ex).ToSQL()
	fmt.Println(phone, query)
	err := meta.MerchantDB.Get(&count, query)
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	if count == 0 {
		key := fmt.Sprintf("%s:phoneExist", meta.Prefix)
		meta.MerchantRedis.SRem(ctx, key, phone).Val()
	} else {
		return errors.New(helper.PhoneBindAlreadyErr)
	}

	return nil
}
