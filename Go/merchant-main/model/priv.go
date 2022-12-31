package model

import (
	"errors"
	"fmt"
	g "github.com/doug-martin/goqu/v9"
	"github.com/go-redis/redis/v8"
	"merchant/contrib/helper"
)

func PrivList(gid, adminGid string) (string, error) {

	if gid != "" {
		ok, err := groupSubCheck(gid, adminGid)
		if err != nil {
			return "[]", err
		}

		if !ok {
			return "[]", errors.New(helper.MethodNoPermission)
		}
	} else {
		gid = adminGid
	}

	gKey := fmt.Sprintf("%s:priv:list:GM%s", meta.Prefix, gid)
	// 运营总监分组
	//if gid == "2" || gid == "10" {
	//	gKey = fmt.Sprintf("%s:priv:PrivAll", meta.Prefix)
	//}
	cmd := meta.MerchantRedis.Get(ctx, gKey)
	//fmt.Println(cmd.String())
	val, err := cmd.Result()
	if err != nil && err != redis.Nil {
		return val, pushLog(err, helper.RedisErr)
	}

	return val, nil
}

func PrivCheck(uri, gid string) error {

	key := fmt.Sprintf("%s:priv:PrivMap", meta.Prefix)
	cmd := meta.MerchantRedis.HGet(ctx, key, uri)
	//fmt.Println(cmd.String())
	privId, err := cmd.Result()
	if err != nil {
		return pushLog(err, helper.RedisErr)
	}

	id := fmt.Sprintf("%s:priv:GM%s", meta.Prefix, gid)
	hcmd := meta.MerchantRedis.HExists(ctx, id, privId)
	//fmt.Println(hcmd.String())
	exists := hcmd.Val()
	if !exists {
		return errors.New("404")
	}

	return nil
}

/**
 * @Description: 刷新缓存
 * @Author: carl
 */
func LoadPrivs() error {

	var records []Priv

	query, _, _ := dialect.From("tbl_admin_priv").
		Select("pid", "state", "id", "name", "sortlevel", "module").Where(g.Ex{"prefix": meta.Prefix}).Order(g.C("sortlevel").Asc()).ToSQL()
	query = "/* master */ " + query
	fmt.Println(query)
	err := meta.MerchantDB.Select(&records, query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	recs, err := helper.JsonMarshal(records)
	if err != nil {
		return errors.New(helper.FormatErr)
	}

	pipe1 := meta.MerchantRedis.TxPipeline()
	defer pipe1.Close()

	privAllKey := fmt.Sprintf("%s:priv:PrivAll", meta.Prefix)
	pipe1.Unlink(ctx, privAllKey)
	pipe1.Set(ctx, privAllKey, string(recs), 0)
	pipe1.Persist(ctx, privAllKey)

	_, err = pipe1.Exec(ctx)
	if err != nil {
		return pushLog(err, helper.RedisErr)
	}

	pipe2 := meta.MerchantRedis.TxPipeline()
	defer pipe2.Close()

	privMapKey := fmt.Sprintf("%s:priv:PrivMap", meta.Prefix)
	pipe2.Unlink(ctx, privMapKey)
	for _, val := range records {
		pipe2.HSet(ctx, privMapKey, val.Module, val.ID)
		fmt.Println(privMapKey, val.Module, val.ID)
	}
	pipe2.Persist(ctx, privMapKey)

	_, err = pipe2.Exec(ctx)
	if err != nil {
		return pushLog(err, helper.RedisErr)
	}

	return nil
}
