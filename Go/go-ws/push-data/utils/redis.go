package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"push-data/configs"
)

var ctx = context.Background()

// 连接redis集群
func InitRedisCluster() (*redis.ClusterClient, error) {
	rdcInfo := configs.Conf.RedisCluster
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{rdcInfo[0], rdcInfo[1], rdcInfo[2], rdcInfo[3], rdcInfo[4], rdcInfo[5]},
	})
	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}

func InitRedis() (*redis.Client, error) {
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
		return nil, err
	}
	return rdb, nil
}
