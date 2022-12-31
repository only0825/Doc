package model

import (
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Rdbc *redis.ClusterClient
var Rdb *redis.Client

//var ctx = context.Background()
//
//type CLIENT_TYPE string
//
//const (
//	CACHE_REDIS         CLIENT_TYPE = "CACHE_REDIS"
//	CACHE_REDIS_CLUSTER             = "CACHE_REDIS_CLUSTER"
//)
//
//type CacheClient struct {
//	ClientType    CLIENT_TYPE
//	CRedisCluster *CacheRedisCluster
//	CRedis        *CacheRedis
//}
//
//var CacheCli *CacheClient
//var ccache *CacheRedis
//var lcache *CacheRedisCluster
//
//func InitInstance(cacheType CLIENT_TYPE) {
//	client := &CacheClient{
//		ClientType:    cacheType,
//		CRedisCluster: credisCluster,
//		CRedis:        credis,
//	}
//	CacheCli = client
//	switch cacheType {
//	case CACHE_REDIS:
//		CacheCli.CRedis.Client = initRedis()
//	case CACHE_REDIS_CLUSTER:
//		initRedisCluster()
//		//client.CRedisCluster.Client = initRedisCluster()
//		//CacheCli = client
//	default:
//		break
//	}
//}
//
//func (client *CacheClient) getDriver() ICache {
//	if client.CRedisCluster.Client != nil {
//		return client.CRedisCluster
//	}
//	return client.CRedis
//}
//
//func (credis *CacheClient) Set(key string, value interface{}, ttl time.Duration) {
//	credis.getDriver().Set(key, value, ttl)
//}
//func (credis *CacheClient) Get(key string) interface{} {
//	return credis.getDriver().Get(key)
//}
