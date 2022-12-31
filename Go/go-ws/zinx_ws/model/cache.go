package model

//
//// 使用单例模式进行封装
//var once sync.Once
//
//// 定义缓存接口，提供自定义的 get、set 方法
//type ICache interface {
//	Set(key string, value interface{}, ttl time.Duration)
//	Get(key string) interface{}
//}
//
//type CacheRedis struct {
//	Client *redis.Client
//}
//
//var credis *CacheRedis
//
//func (credis *CacheRedis) Set(key string, value interface{}, ttl time.Duration) {
//	credis.Client.Set(ctx, key, value, ttl)
//}
//func (credis CacheRedis) Get(key string) interface{} {
//	return credis.Client.Get(ctx, key)
//}
//
//func initRedis() *redis.Client {
//	rcInfo := configs.Conf.Redis
//	rdb := redis.NewClient(&redis.Options{
//		Addr:     fmt.Sprintf("%s:%d", rcInfo.Host, rcInfo.Port),
//		Password: rcInfo.Password, // no password set
//		DB:       rcInfo.Select,   // use default DB
//		//PoolSize:     15,
//		//MinIdleConns: 10, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
//	})
//	once.Do(func() {
//		credis.Client = rdb
//	})
//	return rdb
//}
//
//type CacheRedisCluster struct {
//	Client *redis.ClusterClient
//}
//
//var credisCluster *CacheRedisCluster
//
//func (credis *CacheRedisCluster) Set(key string, value interface{}, ttl time.Duration) {
//	credis.Client.Set(ctx, key, value, ttl)
//}
//func (credis CacheRedisCluster) Get(key string) interface{} {
//	return credis.Client.Get(ctx, key)
//}
//
//// 连接redis集群
//func initRedisCluster() *redis.ClusterClient {
//
//	rdcInfo := configs.Conf.RedisCluster
//	rdbc := redis.NewClusterClient(&redis.ClusterOptions{
//		Addrs: []string{rdcInfo[0], rdcInfo[1], rdcInfo[2], rdcInfo[3], rdcInfo[4], rdcInfo[5]},
//	})
//	err := rdbc.Ping(ctx).Err()
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	once.Do(func() {
//		cluster := &CacheRedisCluster{}
//		cluster.Client = rdbc
//		credisCluster = cluster
//	})
//	return rdbc
//}
