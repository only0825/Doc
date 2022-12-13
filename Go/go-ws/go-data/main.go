package main

import (
	"go-data/common"
	"go-data/logger"
	"go-data/model"
	"net/http"
	"time"
)

func init() {
	rdbc, err := model.InitRedisCluster()
	if err != nil {
		logger.Error.Println("Redis初始化失败: ", err)
		return
	}

	db, err := model.InitDb("u_hainiu:hainiu@2022@tcp(43.135.76.182:3306)/db_hainiu?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		logger.Error.Println("数据库初始化失败: ", err)
		return
	}

	model.Rdbc = rdbc
	model.DB = db
}

func main() {
	//go oddsChange()
	go model.ScoreChange("http://api.wuhaicj.com/api/liveScore/change2", common.BasketBall)
	model.ScoreChange("http://api.wuhaicj.com/api/liveScore/change", common.Football)
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
