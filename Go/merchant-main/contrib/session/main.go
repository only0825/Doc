package session

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/valyala/fasthttp"
	"strings"
	"time"
)

var (
	prefix string
	ctx    = context.Background()
	client *redis.ClusterClient
)

func New(reddb *redis.ClusterClient, p string) {

	client = reddb
	prefix = p
}

func AdminSet(value []byte, uid, deviceNo string) (string, error) {

	uuid := fmt.Sprintf("TD%s", uid)
	key := fmt.Sprintf("%d", Cputicks())

	val, err := client.Get(ctx, uuid).Result()

	pipe := client.TxPipeline()
	defer pipe.Close()

	if err != redis.Nil && len(val) > 0 {
		//同一个用户，一个时间段，只能登录一个
		results := strings.SplitN(val, ",", 3)
		pipe.Unlink(ctx, results[1])
	}

	v := fmt.Sprintf("%s,%s", deviceNo, key)
	pipe.Set(ctx, uuid, v, time.Duration(100)*time.Hour)
	pipe.SetNX(ctx, key, value, defaultExpires)

	_, err = pipe.Exec(ctx)

	/*
			results  := client.LRange(ctx, uuid, 0, -1).Val()
			n        := len(results)
			pipe := client.TxPipeline()
			defer pipe.Close()

			if n > 0 {
				val := strings.SplitN(results[0], ",", 3)
				if val[0] != deviceNo {

					for _, val := range results {
						v := strings.SplitN(val, ",", 3)
						pipe.Unlink(ctx, v[1])
					}
				}
			} else if n > 70 {
				for _, val := range results {
					v := strings.SplitN(val, ",", 3)
					pipe.Unlink(ctx, v[1])
				}
			}
			v := fmt.Sprintf("%s,%s", deviceNo, key)
			pipe.LPush(ctx, uuid, v)
			pipe.SetNX(ctx, key, value, defaultExpires)

		  _, err := pipe.Exec(ctx)
	*/
	return key, err
}

func Set(value []byte, uid string) (string, error) {

	uuid := fmt.Sprintf("TI%s", uid)
	key := fmt.Sprintf("%s:%d", prefix, Cputicks())

	val, err := client.Get(ctx, uuid).Result()

	pipe := client.TxPipeline()
	defer pipe.Close()

	if err != redis.Nil && len(val) > 0 {
		//同一个用户，一个时间段，只能登录一个
		pipe.Unlink(ctx, val)
	}

	pipe.Set(ctx, uuid, key, time.Duration(100)*time.Hour)
	pipe.SetNX(ctx, key, value, defaultExpires)

	_, err = pipe.Exec(ctx)

	return key, err
}

func Update(value []byte, uid uint64) bool {

	uuid := fmt.Sprintf("TI%d", uid)

	val := client.Get(ctx, uuid).Val()
	pipe := client.TxPipeline()
	defer pipe.Close()

	if len(val) == 0 {
		return false
	}

	pipe.Unlink(ctx, val)
	pipe.SetNX(ctx, val, value, renewExpires)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false
	}

	return true
}

func Offline(uids []string) {

	//if len(uids) == 0 {
	//	return
	//}
	//
	//var (
	//	uuids []string
	//	sKeys []string
	//)
	//for _, uid := range uids {
	//	uuids = append(uuids, fmt.Sprintf("TI%s", uid))
	//}
	//keys, err := client.MGet(ctx, uuids...).Result()
	//if err != nil {
	//	return
	//}
	//
	//for _, key := range keys {
	//	if k, ok := key.(string); ok {
	//		sKeys = append(sKeys, k)
	//	}
	//}
	//
	//fmt.Println(client.Unlink(ctx, sKeys...).Err())
}

func Destroy(ctx *fasthttp.RequestCtx) {

	key := string(ctx.Request.Header.Peek("t"))
	if len(key) == 0 {
		return
	}

	client.Unlink(ctx, key)
}

func Get(ctx *fasthttp.RequestCtx) ([]byte, error) {

	key := string(ctx.Request.Header.Peek("t"))
	if len(key) == 0 {

		key = string(ctx.QueryArgs().Peek("t"))
		if len(key) == 0 {
			return nil, errors.New("does not exist")
		}
	}

	pipe := client.TxPipeline()
	defer pipe.Close()

	data := pipe.Get(ctx, key)
	_ = pipe.ExpireAt(ctx, key, ctx.Time().Add(100*time.Hour))
	pipe.Exec(ctx)

	val, err := data.Bytes()
	if err != nil {
		return nil, err
	} else {
		return val, nil
	}
}

func GetByToken(token string) ([]byte, error) {

	val, err := client.Get(ctx, token).Bytes()
	if err == redis.Nil {
		return nil, errors.New("does not exist")
	} else if err != nil {
		return nil, err
	} else {
		return val, nil
	}
}
