package model

import (
	"fmt"
	"merchant/contrib/helper"
	"time"
)

type tree_t struct {
	ID     int    `db:"id" json:"id"`         //
	Level  string `db:"level" json:"level"`   //分类等级
	Name   string `db:"name" json:"name"`     //分类名字
	Sort   int    `db:"sort" json:"sort"`     //排序
	Prefix string `db:"prefix" json:"prefix"` //排序
}

func LoadTrees() error {

	var parent []tree_t
	query := fmt.Sprintf("SELECT * FROM `tbl_tree` WHERE LENGTH(`level`) = 3 and prefix= '%s';", meta.Prefix)
	query = "/* master */ " + query
	fmt.Println(query)
	err := meta.MerchantDB.Select(&parent, query)
	if err != nil {
		fmt.Println("LoadTrees Select = ", err)
		return err
	}

	pipe := meta.MerchantRedis.TxPipeline()
	defer pipe.Close()

	for _, val := range parent {

		var data []tree_t
		key := fmt.Sprintf("%s:T:%s", meta.Prefix, val.Level)
		query := fmt.Sprintf("SELECT * FROM `tbl_tree`  WHERE prefix='%s' and `level` LIKE '%s%%' ORDER BY `level` ASC;", meta.Prefix, val.Level)
		query = "/* master */ " + query
		fmt.Println(query)
		err := meta.MerchantDB.Select(&data, query)
		if err != nil {
			return pushLog(err, helper.DBErr)
		}

		data = data[1:]
		b, _ := helper.JsonMarshal(data)
		pipe.Unlink(ctx, key)
		pipe.Set(ctx, key, string(b), time.Duration(100)*time.Hour)
		pipe.Persist(ctx, key)

	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return pushLog(err, helper.RedisErr)
	}

	return nil
}

func TreeList(level string) (string, error) {

	key := fmt.Sprintf("%s:T:%s", meta.Prefix, level)
	data, err := meta.MerchantRedis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return data, nil
}
