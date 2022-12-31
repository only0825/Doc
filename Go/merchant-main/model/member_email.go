package model

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/doug-martin/goqu/v9/exp"
	"merchant/contrib/helper"

	g "github.com/doug-martin/goqu/v9"
)

type Email_t struct {
	Username  string `json:"username" db:"username"`
	Mail      string `json:"mail" db:"mail"`
	State     string `json:"state" db:"state"`
	Code      string `json:"code" db:"code"`
	IP        string `json:"ip" db:"ip"`
	Ty        string `json:"ty" db:"ty"`
	CreateAt  string `json:"create_at" db:"create_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
	Source    string `json:"source" db:"source"`
	ID        string `json:"id" db:"id"`
}

type EmailData_t struct {
	D []Email_t `json:"d"`
	T int64     `json:"t"`
	S uint      `json:"s"`
}

func EmailList(page, pageSize uint, start, end, username, email string, state, ty int) (EmailData_t, error) {

	ex := g.Ex{}
	data := EmailData_t{}

	if username != "" {
		ex["username"] = username

	}
	if email != "" {
		ex["mail"] = email
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

	t := dialect.From("mail_log")
	if page == 1 {
		query, _, _ := t.Select(g.COUNT("*")).Where(ex).ToSQL()
		fmt.Println(query)
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
	query, _, _ := t.Select("id", "ty", "state", "username", "ip", "code", "source", "mail", "create_at", "updated_at").Where(ex).Offset(offset).Limit(pageSize).Order(g.C("ts").Desc()).ToSQL()
	fmt.Println(query)
	err := meta.MerchantTD.Select(&data.D, query)
	if err != nil {
		body := fmt.Errorf("%s,[%s]", err.Error(), query)
		return data, pushLog(body, helper.DBErr)
	}

	data.S = pageSize

	return data, nil
}
