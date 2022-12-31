package model

import (
	"database/sql"
	"fmt"
	"merchant/contrib/helper"
	"time"

	g "github.com/doug-martin/goqu/v9"
)

func SMSChannelList(ex g.Ex) ([]SMSChannel, error) {

	var data []SMSChannel
	ex["prefix"] = meta.Prefix
	t := dialect.From("tbl_sms")
	query, _, _ := t.Select("id", "name", "alias", "created_at", "txt", "voice", "remark", "created_name").
		Where(ex).Order(g.C("created_at").Desc()).ToSQL()
	fmt.Println(query)
	err := meta.MerchantDB.Select(&data, query)
	if err != nil && err != sql.ErrNoRows {
		return data, pushLog(err, helper.DBErr)
	}

	return data, nil
}

func SMSChannelUpdateState(cid string, txtState, voiceState int) error {

	ex := g.Ex{
		"id":     cid,
		"prefix": meta.Prefix,
	}

	rc := g.Record{}

	tm := map[int]int{
		0: 1,
		1: 2,
		2: 3,
	}

	if txtState != 0 {
		rc["txt"] = tm[txtState]
	}

	if voiceState != 0 {
		rc["voice"] = tm[voiceState]
	}

	query, _, _ := dialect.Update("tbl_sms").Set(rc).Where(ex).ToSQL()
	fmt.Println(query)
	_, err := meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	_ = LoadSMSChannels()

	return nil
}

func SMSChannelUpdate(cid string, rc g.Record) error {

	ex := g.Ex{
		"id":     cid,
		"prefix": meta.Prefix,
	}

	query, _, _ := dialect.Update("tbl_sms").Set(rc).Where(ex).ToSQL()
	fmt.Println(query)
	_, err := meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	_ = LoadSMSChannels()

	return nil
}

func SMSChannelInsert(data *SMSChannel) error {

	id := helper.GenId()
	data.Id = id
	data.CreatedAt = time.Now().Unix()
	data.Prefix = meta.Prefix

	query, _, _ := dialect.Insert("tbl_sms").Rows(data).ToSQL()
	fmt.Println(query)
	_, err := meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	_ = LoadSMSChannels()

	return nil
}

func LoadSMSChannels() (err error) {

	data := make([]SMSChannel, 0)
	ex := g.Or(
		g.Ex{"txt": "1"},
		g.Ex{"voice": "1"},
	)
	query, _, _ := dialect.From("tbl_sms").Select("name", "alias", "txt", "voice").Where(ex).ToSQL()
	query = "/* master */ " + query
	fmt.Println(query)
	err = meta.MerchantDB.Select(&data, query)
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	keyHead := fmt.Sprintf("%s:sms", meta.Prefix)
	_, err = meta.MerchantRedis.Del(ctx, keyHead+":text").Result()
	if err != nil {
		return err
	}

	_, err = meta.MerchantRedis.Del(ctx, keyHead+":voice").Result()
	if err != nil {
		return err
	}

	for _, v := range data {
		if v.Txt == "1" {
			fmt.Println(v.Name + "|sms" + v.Alias)
			_, err = meta.MerchantRedis.LPush(ctx, keyHead+":text", v.Name+"|sms"+v.Alias).Result()
			if err != nil {
				return err
			}
		}

		if v.Voice == "1" {
			fmt.Println(v.Name + "|vms" + v.Alias)
			_, err = meta.MerchantRedis.LPush(ctx, keyHead+":voice", v.Name+"|vms"+v.Alias).Result()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
