package model

import (
	"database/sql"
	"errors"
	"fmt"
	"merchant/contrib/helper"
	"merchant/contrib/validator"
	"strings"
	"time"

	g "github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

//MessageInsert  站内信新增
func MessageInsert(record g.Record, sendAt string) error {

	stAt, err := helper.TimeToLoc(sendAt, loc)
	if err != nil {
		return errors.New(helper.DateTimeErr)
	}

	record["send_at"] = stAt
	record["prefix"] = meta.Prefix
	query, _, _ := dialect.Insert("tbl_messages").Rows(record).ToSQL()
	fmt.Println(query)
	_, err = meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	return nil
}

//MessageList  站内信列表
func MessageList(page, pageSize int, sendStartTime, sendEndTime,
	startTime, endTime, reviewStartTime, reviewEndTime string, ex g.Ex) (MessageData, error) {

	data := MessageData{}
	// 没有查询条件  startTime endTime 必填
	if len(ex) == 0 && (startTime == "" || endTime == "") {
		return data, errors.New(helper.QueryTermsErr)
	}

	if sendStartTime != "" && sendEndTime != "" {

		startAt, err := helper.TimeToLoc(sendStartTime, loc)
		if err != nil {
			return data, errors.New(helper.DateTimeErr)
		}

		endAt, err := helper.TimeToLoc(sendEndTime, loc)
		if err != nil {
			return data, errors.New(helper.TimeTypeErr)
		}

		if startAt >= endAt {
			return data, errors.New(helper.QueryTimeRangeErr)
		}

		ex["send_at"] = g.Op{"between": exp.NewRangeVal(startAt, endAt)}
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

	t := dialect.From("tbl_messages")
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
	query, _, _ := t.Select(colsMessage...).Where(ex).
		Offset(uint(offset)).Limit(uint(pageSize)).Order(g.C("apply_at").Desc()).ToSQL()
	fmt.Println(query)
	err := meta.MerchantDB.Select(&data.D, query)
	if err != nil {
		return data, pushLog(err, helper.DBErr)
	}

	return data, nil
}

//MessageUpdate  站内信更新
func MessageUpdate(id, sendAt string, record g.Record) error {

	ex := g.Ex{
		"id":     id,
		"prefix": meta.Prefix,
	}
	data := Message{}
	t := dialect.From("tbl_messages")
	query, _, _ := t.Select(colsMessage...).Where(ex).ToSQL()
	fmt.Println(query)
	err := meta.MerchantDB.Get(&data, query)
	if err != nil && err != sql.ErrNoRows {
		return pushLog(err, helper.DBErr)
	}

	if err == sql.ErrNoRows {
		return errors.New(helper.RecordNotExistErr)
	}

	if data.State != 1 {
		return errors.New(helper.NoDataUpdate)
	}

	stAt, err := helper.TimeToLoc(sendAt, loc)
	if err != nil {
		return errors.New(helper.DateTimeErr)
	}

	record["send_at"] = stAt
	query, _, _ = dialect.Update("tbl_messages").Set(record).Where(ex).ToSQL()
	fmt.Println(query)
	_, err = meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	return nil
}

//MessageReview  站内信审核
func MessageReview(id string, state, flag int, admin map[string]string) error {

	ex := g.Ex{
		"id":     id,
		"prefix": meta.Prefix,
	}
	data := Message{}
	query, _, _ := dialect.From("tbl_messages").Select(colsMessage...).Where(ex).ToSQL()
	fmt.Println(query)
	err := meta.MerchantDB.Get(&data, query)
	if err != nil && err != sql.ErrNoRows {
		return pushLog(err, helper.DBErr)
	}

	if err == sql.ErrNoRows {
		return errors.New(helper.RecordNotExistErr)
	}

	if data.State != 1 {
		return errors.New(helper.NoDataUpdate)
	}

	ns := time.Now().Unix()
	// 审核通过时已经超过了发送时间，记录作废
	if ns > data.SendAt && state == 2 {

		record := g.Record{
			"state":       4,
			"review_at":   time.Now().Unix(),
			"review_uid":  admin["id"],
			"review_name": admin["name"],
		}
		query, _, _ = dialect.Update("tbl_messages").Set(record).Where(ex).ToSQL()
		fmt.Println(query)
		_, err = meta.MerchantDB.Exec(query)
		if err != nil {
			return pushLog(err, helper.DBErr)
		}

		return errors.New(helper.RecordExpired)
	}

	record := g.Record{
		"state":       state,
		"review_at":   time.Now().Unix(),
		"review_uid":  admin["id"],
		"review_name": admin["name"],
	}
	query, _, _ = dialect.Update("tbl_messages").Set(record).Where(ex).ToSQL()
	fmt.Println(query)
	_, err = meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	// 审核通过
	if state == 2 {

		param := map[string]interface{}{
			"flag":       "1",                            //发送站内信
			"message_id": data.ID,                        //id
			"title":      data.Title,                     //标题
			"content":    data.Content,                   //内容
			"is_top":     fmt.Sprintf("%d", data.IsTop),  //0不置顶 1置顶
			"is_push":    fmt.Sprintf("%d", data.IsPush), //0不推送 1推送
			"is_vip":     fmt.Sprintf("%d", data.IsVip),  //0普通站内信 1vip站内信
			"ty":         fmt.Sprintf("%d", data.Ty),     //1站内消息 2活动消息
			"send_name":  data.SendName,                  //发送人名
			"prefix":     meta.Prefix,                    //商户前缀
		}
		if data.IsVip == 1 {
			param["level"] = data.Level
		} else {
			param["usernames"] = data.Usernames
		}

		sDelay := int64(0)
		// 1 定时发送 2立即发送
		if flag == 2 {
			sDelay = data.SendAt - ns
			err = BeanPutDelay("message", param, int(sDelay))
		} else {
			err = BeanPut("message", param)
		}
		if err != nil {
			_ = pushLog(err, helper.ServerErr)
		}
	}

	return nil
}

//MessageDetail  已发站内信详情
func MessageDetail(id string, page, pageSize int) (MessageTDData, error) {

	data := MessageTDData{
		S: pageSize,
	}

	ex := g.Ex{
		"id":     id,
		"prefix": meta.Prefix,
	}
	var sendState int
	query, _, _ := dialect.From("tbl_messages").Select("send_state").Where(ex).ToSQL()
	fmt.Println(query)
	err := meta.MerchantDB.Get(&sendState, query)
	if err != nil && err != sql.ErrNoRows {
		return data, pushLog(err, helper.DBErr)
	}

	if err == sql.ErrNoRows || sendState != 2 {
		return data, errors.New(helper.RecordNotExistErr)
	}

	ex = g.Ex{
		"prefix":     meta.Prefix,
		"message_id": id,
		"is_delete":  0,
	}
	t := dialect.From("messages")
	if page == 1 {
		query, _, _ = t.Select(g.COUNT("ts")).Where(ex).ToSQL()
		fmt.Println(query)
		err = meta.MerchantTD.Get(&data.T, query)
		if err != nil && err != sql.ErrNoRows {
			return data, pushLog(err, helper.DBErr)
		}

		if data.T == 0 {
			return data, nil
		}
	}

	offset := (page - 1) * pageSize
	query, _, _ = t.Select(colsMessageTD...).Where(ex).Offset(uint(offset)).Limit(uint(pageSize)).Order(g.C("ts").Desc()).ToSQL()
	fmt.Println(query)
	err = meta.MerchantTD.Select(&data.D, query)
	if err != nil {
		return data, pushLog(err, helper.DBErr)
	}

	return data, nil
}

//MessageSystemList  已发站系统内信列表
func MessageSystemList(startTime, endTime string, page, pageSize int, ex g.Ex) (MessageTDData, error) {

	data := MessageTDData{
		S: pageSize,
	}
	startAt, err := helper.TimeToLoc(startTime, loc) // 毫秒级时间戳
	if err != nil {
		return data, errors.New(helper.DateTimeErr)
	}
	endAt, err := helper.TimeToLoc(endTime, loc) // 毫秒级时间戳
	if err != nil {
		return data, errors.New(helper.DateTimeErr)
	}

	if startAt >= endAt {
		return data, errors.New(helper.QueryTimeRangeErr)
	}

	ex["prefix"] = meta.Prefix
	ex["send_name"] = "system"
	ex["send_at"] = g.Op{"between": exp.NewRangeVal(startAt, endAt)}
	t := dialect.From("messages")
	if page == 1 {
		query, _, _ := t.Select(g.COUNT("ts")).Where(ex).ToSQL()
		fmt.Println(query)
		err = meta.MerchantTD.Get(&data.T, query)
		if err != nil && err != sql.ErrNoRows {
			return data, pushLog(err, helper.DBErr)
		}

		if data.T == 0 {
			return data, nil
		}
	}

	offset := (page - 1) * pageSize
	query, _, _ := t.Select(colsMessageTD...).Where(ex).Offset(uint(offset)).Limit(uint(pageSize)).Order(g.C("ts").Desc()).ToSQL()
	fmt.Println(query)
	err = meta.MerchantTD.Select(&data.D, query)
	if err != nil {
		return data, pushLog(err, helper.DBErr)
	}

	return data, nil
}

//MessageDelete  站内信删除
func MessageDelete(id, tss string) error {

	param := map[string]interface{}{
		"flag": "2", //删除站内信
	}
	if tss == "" {
		if !validator.CtypeDigit(id) {
			return errors.New(helper.IDErr)
		}

		ex := g.Ex{
			"id":     id,
			"prefix": meta.Prefix,
		}
		record := g.Record{
			"state": 4,
		}
		query, _, _ := dialect.Update("tbl_messages").Set(record).Where(ex).ToSQL()
		fmt.Println(query)
		_, err := meta.MerchantDB.Exec(query)
		if err != nil {
			return pushLog(err, helper.DBErr)
		}

		param["message_id"] = id //站内信id
		err = BeanPut("message", param)
		if err != nil {
			_ = pushLog(err, helper.ServerErr)
		}

		return nil
	}

	var records []g.Record
	for _, v := range strings.Split(tss, ",") {
		// 2022-06-07T16:28:26.285+07:00
		t, err := time.ParseInLocation("2006-01-02T15:04:05.999999+07:00", v, loc)
		if err != nil {
			return errors.New(helper.ParamErr)
		}
		record := g.Record{
			"ts":        t.UnixMicro(),
			"is_delete": 1,
		}
		records = append(records, record)
	}
	query, _, _ := dialect.Insert("messages").Rows(records).ToSQL()
	fmt.Println(query)
	_, err := meta.MerchantTD.Exec(query)
	if err != nil {
		fmt.Println("insert messages = ", err.Error(), records)
	}

	return nil
}

// 发送站内信
func messageSend(msgID, title, subTitle, content, sendName, prefix string, isTop, isVip, ty int, names []string) error {

	record := g.Record{
		"message_id": msgID,
		"title":      title,
		"sub_title":  subTitle,
		"content":    content,
		"send_name":  sendName,
		"prefix":     prefix,
		"is_top":     isTop,
		"is_vip":     isVip,
		"is_read":    0,
		"is_delete":  0,
		"send_at":    time.Now().Unix(),
		"ty":         ty,
	}
	var records []g.Record
	for _, v := range names {
		ts := time.Now()
		record["ts"] = ts.UnixMilli()
		record["username"] = v
		records = append(records, record)
	}

	query, _, _ := dialect.Insert("messages").Rows(records).ToSQL()
	fmt.Println(query)
	_, err := meta.MerchantTD.Exec(query)
	if err != nil {
		fmt.Println("insert messages = ", err.Error(), records)
	}

	return nil
}
