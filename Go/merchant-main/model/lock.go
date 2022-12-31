package model

import (
	"errors"
	"fmt"
	"merchant/contrib/helper"
	"time"
)

func Lock(id string) error {

	val := fmt.Sprintf("%s%s", defaultRedisKeyPrefix, id)
	ok, err := meta.MerchantRedis.SetNX(ctx, val, "1", LockTimeout).Result()
	if err != nil {
		return pushLog(err, helper.RedisErr)
	}
	if !ok {
		return errors.New(helper.RequestBusy)
	}

	return nil
}

func LockWait(id string, ttl time.Duration) error {

	val := fmt.Sprintf("%s%s", defaultRedisKeyPrefix, id)

	for {
		ok, err := meta.MerchantRedis.SetNX(ctx, val, "1", ttl).Result()
		if err != nil {
			return pushLog(err, helper.RedisErr)
		}

		if !ok {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		return nil
	}
}

func LockTTL(id string, ttl time.Duration) error {

	val := fmt.Sprintf("%s%s", defaultRedisKeyPrefix, id)
	ok, err := meta.MerchantRedis.SetNX(ctx, val, "1", ttl).Result()
	if err != nil || !ok {
		return pushLog(err, helper.RedisErr)
	}

	return nil
}

func LockSetExpire(id string, expiration time.Duration) error {

	val := fmt.Sprintf("%s%s", defaultRedisKeyPrefix, id)
	ok, err := meta.MerchantRedis.Expire(ctx, val, expiration).Result()
	if err != nil || !ok {
		return pushLog(err, helper.RedisErr)
	}

	return nil
}

func Unlock(id string) {

	val := fmt.Sprintf("%s%s", defaultRedisKeyPrefix, id)
	res, err := meta.MerchantRedis.Unlink(ctx, val).Result()
	if err != nil || res != 1 {
		fmt.Println("Unlock res = ", res)
		fmt.Println("Unlock err = ", err)
	}
}
