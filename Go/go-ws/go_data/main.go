package main

import (
	"go_data/common"
	"go_data/logger"
	"go_data/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func init() {

	rdbc, err := common.InitRedisCluster()
	if err != nil {
		logger.Error.Println("Redis初始化失败: ", err)
		return
	}

	//db, err := model.InitDb("u_hainiu:hainiu@2022@tcp(43.135.76.182:3306)/db_hainiu?charset=utf8&parseTime=True&loc=Local")
	dsn := "root:password@tcp(127.0.0.1:3306)/hainiu?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error.Println("数据库初始化失败: ", err)
		return
	}

	model.DB = db
	model.Rdbc = rdbc
}

func main() {
	//go oddsChange()
	//
	//j := "{\"changeList\":[]}"
	//fmt.Println(j)

	// 实时比分
	//go model.ScoreChange("http://api.wuhaicj.com/api/liveScore/change2", common.BasketBall)
	//model.ScoreChange("http://api.wuhaicj.com/api/liveScore/change", common.Football)
}

func oddsChange() {
	for {
		time.Sleep(time.Duration(3) * time.Second)
		res, err := http.Get("http://api.wuhaicj.com/api/liveScore/oddsChange")
		if err != nil {
			logger.Error.Println("URL Request failed:", err)
			break
		}
		logger.Info.Println(res)
	}
}
