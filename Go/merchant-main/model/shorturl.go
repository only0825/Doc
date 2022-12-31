package model

import (
	"github.com/go-redis/redis/v8"
	"merchant/contrib/helper"
	"time"
)

func ShortURLSet(uri string) error {

	key := meta.Prefix + ":shorturl:domain"
	pipe := meta.MerchantRedis.TxPipeline()
	defer pipe.Close()

	pipe.Set(ctx, key, uri, 100*time.Hour)
	pipe.Persist(ctx, key)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return pushLog(err, helper.RedisErr)
	}

	return nil
}

func ShortURLGet() (string, error) {

	key := meta.Prefix + ":shorturl:domain"
	cmd := meta.MerchantRedis.Get(ctx, key)
	//fmt.Println(cmd.String())
	resc, err := cmd.Result()
	if err != nil && err != redis.Nil {
		return "", pushLog(err, helper.RedisErr)
	}

	if err == redis.Nil {
		return "", redis.Nil
	}

	return resc, nil
}
