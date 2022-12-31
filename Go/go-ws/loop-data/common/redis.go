package common

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"loop-data/configs"
	"loop-data/model"
)

var ctx = context.Background()

func InitCache(cacheType string) error {
	if cacheType == "redisCluster" {
		err := initRedisCluster()
		return err
	} else if cacheType == "redis" {
		err := initRedis()
		return err
	}
	return errors.New("cacheType 错误，只能是redis和 redisCluster")
}

// 连接redis集群
func initRedisCluster() error {
	rdcInfo := configs.Conf.RedisCluster
	rdbc := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{rdcInfo[0], rdcInfo[1], rdcInfo[2], rdcInfo[3], rdcInfo[4], rdcInfo[5]},
	})
	err := rdbc.Ping(ctx).Err()
	if err != nil {
		return err
	}
	model.Rdbc = rdbc
	return nil
}

func initRedis() error {
	rcInfo := configs.Conf.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", rcInfo.Host, rcInfo.Port),
		Password: rcInfo.Password, // no password set
		DB:       rcInfo.Select,   // use default DB
		//PoolSize:     15,
		//MinIdleConns: 10, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
	})
	err := rdb.Ping(ctx).Err()
	if err != nil {
		return err
	}
	model.Rdb = rdb
	return nil
}

type Cache struct {
	redis        redis.Client
	redisCluster redis.ClusterClient
}

func GetRedisInstance() interface{} {
	cacheType := configs.Conf.Cache
	if cacheType == "redisCluster" {
		return model.Rdbc
	} else if cacheType == "redis" {
		return model.Rdb
	} else {
		return nil
	}
}
