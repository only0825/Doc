package model

import (
	"database/sql"
	"errors"
	"fmt"
	g "github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/wI2L/jettison"
	"merchant/contrib/helper"
	"time"
)

// 公告添加
func NoticeInsert(data Notice) error {

	data.Prefix = meta.Prefix
	query, _, _ := dialect.Insert("tbl_notices").Rows(&data).ToSQL()
	_, err := meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	return nil
}

// 公告列表
func NoticeList(page, pageSize uint, startTime, endTime string, ex g.Ex) (NoticeData, error) {

	data := NoticeData{}

	if startTime != "" && endTime != "" {
		st, err := time.Parse("2006-01-02 15:04:05", startTime)
		if err != nil {
			return data, errors.New(helper.TimeTypeErr)
		}

		et, err := time.Parse("2006-01-02 15:04:05", endTime)
		if err != nil {
			return data, errors.New(helper.TimeTypeErr)
		}

		ex["created_at"] = g.Op{"between": exp.NewRangeVal(st.Unix(), et.Unix())}
	}
	ex["prefix"] = meta.Prefix
	t := dialect.From("tbl_notices")
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
	query, _, _ := t.Select(colsNotice...).Where(ex).Order(g.C("created_at").Desc()).Offset(offset).Limit(pageSize).ToSQL()
	err := meta.MerchantDB.Select(&data.D, query)
	if err != nil && err != sql.ErrNoRows {
		return data, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	return data, nil
}

// 公告 启用 停用
func NoticeUpdateState(id string, state int) error {

	data := Notice{}
	ex := g.Ex{
		"id":     id,
		"prefix": meta.Prefix,
	}
	query, _, _ := dialect.From("tbl_notices").Select(colsNotice...).Where(ex).ToSQL()
	err := meta.MerchantDB.Get(&data, query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	record := g.Record{
		"state":  state,
		"prefix": meta.Prefix,
	}
	query, _, _ = dialect.Update("tbl_notices").Set(record).Where(ex).ToSQL()
	_, err = meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	_ = noticeRefreshToCache()

	return nil
}

// 公告删除
func NoticeDelete(id string) error {

	ex := g.Ex{
		"id":     id,
		"state":  1, // 开启
		"prefix": meta.Prefix,
	}
	query, _, _ := dialect.Delete("tbl_notices").Where(ex).ToSQL()
	_, err := meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	return nil
}

// 公告更新
func NoticeUpdate(ex g.Ex, record g.Record) error {

	ex["prefix"] = meta.Prefix
	query, _, _ := dialect.Update("tbl_notices").Set(record).Where(ex).ToSQL()
	_, err := meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	_ = noticeRefreshToCache()

	return nil
}

func noticeRefreshToCache() error {

	var notices []Notice
	ex := g.Ex{
		"state":  2, //开启的
		"prefix": meta.Prefix,
	}
	query, _, _ := dialect.From("tbl_notices").Select(colsNotice...).Where(ex).Order(g.C("created_at").Desc()).ToSQL()
	err := meta.MerchantDB.Select(&notices, query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	data, err := jettison.Marshal(notices)
	if err != nil {
		return errors.New(helper.FormatErr)
	}

	pipe := meta.MerchantRedis.TxPipeline()
	defer pipe.Close()

	pipe.Unlink(ctx, "notices")
	pipe.Set(ctx, "notices", string(data), 100*time.Hour)
	pipe.Persist(ctx, "notices")

	_, err = pipe.Exec(ctx)

	return err
}
