package model

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/doug-martin/goqu/v9/exp"
	"merchant/contrib/helper"

	g "github.com/doug-martin/goqu/v9"
)

type Sms_t struct {
	Username  string `json:"username" db:"username"`
	Phone     string `json:"phone" db:"phone"`
	State     string `json:"state" db:"state"`
	Code      string `json:"code" db:"code"`
	IP        string `json:"ip" db:"ip"`
	Ty        string `json:"ty" db:"ty"`
	CreateAt  string `json:"create_at" db:"create_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
	Flags     string `json:"flags" db:"flags"`
	Source    string `json:"source" db:"source"`
	ID        string `json:"id" db:"id"`
}

type SmsData_t struct {
	D []Sms_t `json:"d"`
	T int64   `json:"t"`
	S uint    `json:"s"`
}

func SmsList(page, pageSize uint, start, end, username, phone string, state, ty int) (SmsData_t, error) {

	ex := g.Ex{}
	data := SmsData_t{}

	if username != "" {
		ex["username"] = username

	}
	if phone != "" {
		ex["phone"] = phone
	}

	if state != 0 {
		ex["state"] = state
	}

	if ty > 0 {
		ex["ty"] = ty
	}

	ex["prefix"] = meta.Prefix
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

	t := dialect.From("sms_log")
	if page == 1 {
		query, _, _ := t.Select(g.COUNT("*")).Where(ex).ToSQL()
		err := meta.MerchantTD.Get(&data.T, query)
		if err == sql.ErrNoRows {
			return data, nil
		}

		if err != nil {
			body := fmt.Errorf("%s,[%s]", err.Error(), query)
			return data, pushLog(body, helper.DBErr)
		}

		if data.T == 0 {
			return data, nil
		}
	}

	offset := (page - 1) * pageSize
	query, _, _ := t.Select("id", "ty", "state", "username", "ip", "code", "flags", "source", "phone", "create_at", "updated_at").Where(ex).Offset(offset).Limit(pageSize).Order(g.C("ts").Desc()).ToSQL()
	err := meta.MerchantTD.Select(&data.D, query)
	if err != nil {
		body := fmt.Errorf("%s,[%s]", err.Error(), query)
		return data, pushLog(body, helper.DBErr)
	}

	data.S = pageSize

	return data, nil
}
