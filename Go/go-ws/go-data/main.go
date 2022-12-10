package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"time"
)

var m = make(map[string]interface{})

var rdb *redis.Client

func main() {
	//go oddsChange()
	scoreChange()
}

func scoreChange() {
	for {
		time.Sleep(time.Duration(3) * time.Second)
		res, err := http.Get("http://api.wuhaicj.com/api/liveScore/change")
		if err != nil {
			fmt.Println("URL Request failed:", err)
			break
		}
		msg, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Read body failed:", err)
			break
		}
		// 将json消息转为数组
		err = json.Unmarshal(msg, &m)
		if err != nil {
			fmt.Println("err :", err)
			break
		}

		fmt.Println(m["changeList"])
		defer res.Body.Close()
	}

}

func oddsChange() {
	for {
		time.Sleep(time.Duration(3) * time.Second)
		res, err := http.Get("http://api.wuhaicj.com/api/liveScore/oddsChange")
		if err != nil {
			fmt.Println("URL Request failed:", err)
			break
		}
		fmt.Println(res)
	}
}

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "", // no password set
		DB:           0,  // use default DB
		PoolSize:     15,
		MinIdleConns: 10, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
	})

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local", // DSN data source name
		DefaultStringSize:         256,                                                                        // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                                                                       // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                                                                       // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                                                                       // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                                                                      // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})

	fmt.Println(db)
	fmt.Println(err)
}
