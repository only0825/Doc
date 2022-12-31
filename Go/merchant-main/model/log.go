package model

import (
	"database/sql"
	"errors"
	"fmt"
	g "github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/olivere/elastic/v7"
	"merchant/contrib/helper"
)

func AdminLoginLog(start, end string, page, pageSize int, ex g.Ex) (AdminLoginLogData, error) {

	data := AdminLoginLogData{}

	if start != "" && end != "" {

		startAt, err := helper.TimeToLoc(start, loc)
		if err != nil {
			return data, errors.New(helper.DateTimeErr)
		}

		endAt, err := helper.TimeToLoc(end, loc)
		if err != nil {
			return data, errors.New(helper.DateTimeErr)
		}

		ex["create_at"] = g.Op{"between": exp.NewRangeVal(startAt, endAt)}
	}

	ex["prefix"] = meta.Prefix

	t := dialect.From("member_login_log")
	if page == 1 {
		query, _, _ := t.Select(g.COUNT("*")).Where(ex).ToSQL()
		err := meta.MerchantTD.Get(&data.T, query)
		if err == sql.ErrNoRows {
			return data, nil
		}

		if err != nil {
			fmt.Println("Member Login Log err = ", err.Error())
			fmt.Println("Member Login Log query = ", query)
			body := fmt.Errorf("%s,[%s]", err.Error(), query)
			return data, pushLog(body, helper.DBErr)
		}
		if data.T == 0 {
			return data, nil
		}
	}

	offset := (page - 1) * pageSize
	query, _, _ := t.Select("username", "ip", "device_no", "create_at").Where(ex).Offset(uint(offset)).Limit(uint(pageSize)).Order(g.C("ts").Desc()).ToSQL()
	fmt.Println("Member Login Log query = ", query)
	err := meta.MerchantTD.Select(&data.D, query)
	if err != nil {
		return data, err
	}

	data.S = pageSize

	return data, nil
}

// 系统日志
func SystemLog(start, end string, page, pageSize int, query *elastic.BoolQuery) (SystemLogData, error) {

	data := SystemLogData{}
	//
	//if start != "" && end != "" {
	//
	//	startAt, err := helper.TimeToLoc(start, loc)
	//	if err != nil {
	//		return data, errors.New(helper.DateTimeErr)
	//	}
	//
	//	endAt, err := helper.TimeToLoc(end, loc)
	//	if err != nil {
	//		return data, errors.New(helper.DateTimeErr)
	//	}
	//
	//	query.Filter(
	//		elastic.NewRangeQuery("created_at").Gte(startAt).Lte(endAt),
	//	)
	//}
	//query.Filter(elastic.NewTermQuery("prefix", meta.Prefix))
	//
	//fields := []string{"uid", "name", "title", "ip", "content", "created_at", "prefix"}
	//total, result, _, err := EsQuerySearch(esPrefixIndex("system_log"), "@timestamp", page, pageSize, fields, query, nil)
	//if err != nil {
	//	return data, err
	//}
	//
	//data.T = total
	//data.S = pageSize
	//
	//for _, v := range result {
	//
	//	log := systemLog{}
	//	if err = helper.JsonUnmarshal(v.Source, &log); err != nil {
	//		return data, errors.New(helper.FormatErr)
	//	}
	//
	//	log.Id = v.Id
	//	data.D = append(data.D, log)
	//}

	return data, nil
}
