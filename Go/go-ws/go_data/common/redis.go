package common

import (
	"context"
	"github.com/go-redis/redis/v9"
)

// 连接redis集群
func InitRedisCluster() (*redis.ClusterClient, error) {
	rdbc := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},
	})
	var ctx = context.Background()
	err := rdbc.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}
	return rdbc, nil
}

func initRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "", // no password set
		DB:           0,  // use default DB
		PoolSize:     15,
		MinIdleConns: 10, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
	})
	var ctx = context.Background()
	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
