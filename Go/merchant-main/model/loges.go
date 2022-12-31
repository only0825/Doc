package model

import (
	"database/sql"
	"errors"
	"fmt"
	"merchant/contrib/helper"

	g "github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

func MemberRemarkLogList(uid, adminName, startTime, endTime string, page, pageSize int) (MemberRemarkLogData, error) {

	ex := g.Ex{
		"is_delete": g.Op{"neq": 1},
	}

	if uid != "" {
		ex["uid"] = uid
	}

	if adminName != "" {
		ex["created_name"] = adminName
	}

	data := MemberRemarkLogData{}
	if len(ex) == 0 && (startTime == "" || endTime == "") {
		return data, errors.New(helper.QueryTermsErr)
	}
	if startTime != "" && endTime != "" {

		startAt, err := helper.TimeToLoc(startTime, loc)
		if err != nil {
			return data, errors.New(helper.DateTimeErr)
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

	t := dialect.From("member_remarks_log")

	if page == 1 {
		query, _, _ := t.Select(g.COUNT("*")).Where(ex).ToSQL()

		fmt.Println(query)

		err := meta.MerchantTD.Get(&data.T, query)
		if err == sql.ErrNoRows {
			return data, nil
		}

		if err != nil {
			fmt.Println("Member Remarks Log err = ", err.Error())
			fmt.Println("Member Remarks Log query = ", query)
			body := fmt.Errorf("%s,[%s]", err.Error(), query)
			return data, pushLog(body, helper.DBErr)
		}
		if data.T == 0 {
			return data, nil
		}
	}

	offset := (page - 1) * pageSize
	query, _, _ := t.Select(colsMemberRemarksLog...).Where(ex).Offset(uint(offset)).Limit(uint(pageSize)).Order(g.C("ts").Desc()).ToSQL()
	fmt.Println("Member Remarks Log query = ", query)
	err := meta.MerchantTD.Select(&data.D, query)
	if err != nil {
		fmt.Println("Member Remarks Log err = ", err.Error())
		fmt.Println("Member Remarks Log query = ", query)
		body := fmt.Errorf("%s,[%s]", err.Error(), query)
		return data, pushLog(body, helper.DBErr)
	}

	data.S = pageSize

	return data, nil
}

func MemberLoginLogList(startTime, endTime string, page, pageSize int, ex g.Ex) (MemberLoginLogData, error) {

	data := MemberLoginLogData{}
	if len(ex) == 0 && (startTime == "" || endTime == "") {
		return data, errors.New(helper.QueryTermsErr)
	}
	if startTime != "" && endTime != "" {

		startAt, err := helper.TimeToLoc(startTime, loc)
		if err != nil {
			return data, errors.New(helper.DateTimeErr)
		}

		endAt, err := helper.TimeToLoc(endTime, loc)
		if err != nil {
			return data, errors.New(helper.TimeTypeErr)
		}

		if startAt >= endAt {
			return data, errors.New(helper.QueryTimeRangeErr)
		}

		ex["create_at"] = g.Op{"between": exp.NewRangeVal(startAt, endAt)}
	}
	ex["prefix"] = meta.Prefix

	t := dialect.From("member_login_log")
	fmt.Println(ex)
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
	query, _, _ := t.Select("username", "ip", "device", "device_no", "top_name", "parent_name", "create_at").Where(ex).Offset(uint(offset)).Limit(uint(pageSize)).Order(g.C("ts").Desc()).ToSQL()
	fmt.Println("Member Login Log query = ", query)

	err := meta.MerchantTD.Select(&data.D, query)
	if err != nil {
		fmt.Println("Member Login Log err = ", err.Error())
		fmt.Println("Member Login Log query = ", query)
		body := fmt.Errorf("%s,[%s]", err.Error(), query)
		return data, pushLog(body, helper.DBErr)
	}

	data.S = pageSize
	return data, nil
	//添加计数列
	//var countResult []IpUser
	//query, _, _ = t.Select(g.COUNT("username").As("username"), "ip").GroupBy("username", "ip").ToSQL()
	//err = meta.MerchantTD.Select(&countResult, query)
	//if err != nil {
	//	fmt.Printf("Member Login count name Log query = :%+v err:%+v\n", query, err.Error())
	//	body := fmt.Errorf("%s,[%s]", err.Error(), query)
	//	return data, pushLog(body, helper.DBErr)
	//}

	//var users = make([]string, len(data.D))
	//for _, d := range data.D {
	//	users = append(users, d.Username)
	//}
	//
	//mbs, err := memberFindBatch(users)
	//if err != nil {
	//	fmt.Println("Member detail err = ", err)
	//	return data, pushLog(err, helper.DBErr)
	//}
	//
	////更新用户团队信息
	//for i, d := range data.D {
	//	data.D[i].GroupName = mbs[d.Username].GroupName
	//}

	//fmt.Printf("DEBUG:memeber login  data:%+v\n", len(data.D))
	//
	////执行更新总数
	//if len(countResult) == 0 {
	//	fmt.Printf("DEBUG: countResult 0 return data:%+v\n", data.D)
	//	return data, nil
	//} else {
	//	for i, d := range data.D {
	//		for _, h := range countResult {
	//			if h.IP == d.IP {
	//				data.D[i].CountName += 1
	//			}
	//		}
	//	}
	//	return data, nil
	//}
}

//MemberAccessList 某ip对应的会员登陆信息
func MemberAccessList(page, pageSize int, ex g.Ex) (MemberAssocLogData, error) {

	data := MemberAssocLogData{S: pageSize}
	ex["prefix"] = meta.Prefix

	//去重总数
	t := dialect.From("member_login_log")
	var countResult []IpUser
	query, _, _ := t.Select(g.COUNT("username").As("username")).Where(ex).GroupBy("username").ToSQL()
	err := meta.MerchantTD.Select(&countResult, query)
	data.T = int64(len(countResult))
	if err == sql.ErrNoRows {
		return data, nil
	}

	if err != nil {
		fmt.Printf("Member Login Log err:%+v \n query:%+v\n", err.Error(), query)
		body := fmt.Errorf("%s,[%s]", err.Error(), query)
		return data, pushLog(body, helper.DBErr)
	}

	if data.T == 0 {
		return data, nil
	}

	//TD查去重
	//offset := (page - 1) * pageSize
	//query, _, _ = t.Select(g.DISTINCT("username").As("username")).Where(ex).Offset(uint(offset)).Limit(uint(pageSize)).ToSQL()
	//err = meta.MerchantTD.Select(&data.D, query)
	//if err != nil {
	//	fmt.Printf("Member Login Log err:%+v \n query:%+v\n", err.Error(), query)
	//	body := fmt.Errorf("%s,[%s]", err.Error(), query)
	//	return data, pushLog(body, helper.DBErr)
	//}
	//
	////获取用户名
	//var users = make([]string, len(data.D))
	//for _, d := range data.D {
	//	users = append(users, d.Username)
	//}
	//
	////更新用户状态 注册时间 等
	//mbs, err := memberFindBatch(users)
	//if err != nil {
	//	fmt.Println("Member detail err = ", err)
	//	return data, pushLog(err, helper.DBErr)
	//}
	//for i, d := range data.D {
	//	data.D[i].TopUID = mbs[d.Username].TopUid
	//	data.D[i].TopName = mbs[d.Username].TopName
	//	data.D[i].ParentName = mbs[d.Username].ParentName
	//	data.D[i].State = mbs[d.Username].State
	//	data.D[i].CreatedAt = mbs[d.Username].CreatedAt
	//	data.D[i].Remarks = mbs[d.Username].Remarks
	//	data.D[i].LastLoginAt = mbs[d.Username].LastLoginAt
	//	data.D[i].GroupName = mbs[d.Username].GroupName
	//}

	return data, nil
}
