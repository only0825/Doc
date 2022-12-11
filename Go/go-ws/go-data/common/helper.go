package common

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/go-redis/redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"log"
)

//type Books struct {
//	title   string
//	author  string
//	subject string
//	book_id int
//}

var DB *gorm.DB
var Ctx = context.Background()
var Rdbc *redis.ClusterClient
var rdb *redis.Client

func InitRedisCluster() *redis.ClusterClient {
	// 连接redis集群
	Rdbc = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},
	})
	err := Rdbc.Ping(Ctx).Err()
	if err == nil {
		log.Println("Redis cluster OK")
	} else {
		log.Println("Redis cluster wrong: ", err)
	}

	return Rdbc
}

func initRedis() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "", // no password set
		DB:           0,  // use default DB
		PoolSize:     15,
		MinIdleConns: 10, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
	})
	_, err = rdb.Ping(Ctx).Result()
	return err
}

func InitDb() {
	// 	db, err := sqlx.Connect("mysql", "u_hainiu:hainiu@2022@(43.135.76.182:3306)/db_hainiu")
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "u_hainiu:hainiu@2022@tcp(43.135.76.182:3306)/db_hainiu?charset=utf8&parseTime=True&loc=Local", // DSN data source name
		DefaultStringSize:         256,                                                                                            // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                                                                                           // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                                                                                           // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                                                                                           // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                                                                                          // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
	if err != nil {
		fmt.Println("Mysql Connect failed:", err)
		return
	}
	DB = db
}

func Md5String(str string) string {
	m := md5.New()
	_, err := io.WriteString(m, str)
	if err != nil {
		log.Fatal(err)
	}
	arr := m.Sum(nil)
	return fmt.Sprintf("%x", arr)
}

/*
Redis 操作：
   err := rdb.Set(ctx, "key", "value", 0).Err()
   if err != nil {
       panic(err)
   }

	// 过期时间
	err := rdbc.Set(ctx, "key111", "value111", time.Duration(20)*time.Second).Err()
	if err != nil {
		panic(err)
	}

   val, err := rdb.Get(ctx, "key").Result()
   if err != nil {
       panic(err)
   }
   fmt.Println("key", val)

   val2, err := rdb.Get(ctx, "key2").Result()
   if err == redis.Nil {
       fmt.Println("key2 does not exist")
   } else if err != nil {
       panic(err)
   } else {
       fmt.Println("key2", val2)
   }
   // Output: key value
   // key2 does not exist


		// 将json消息转为数组
		//err = json.Unmarshal(msg, &m)
		//if err != nil {
		//	fmt.Println("err :", err)
		//	break
		//}
*/
