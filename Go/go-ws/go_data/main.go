package main

import (
	"github.com/robfig/cron/v3"
	"go_data/common"
	"go_data/logger"
	"go_data/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
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

// 返回一个支持至 秒 级别的 cron
func newWithSeconds() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}

func main() {
	//go oddsChange()

	i := 0
	c := newWithSeconds()
	//AddFunc
	spec := "*/5 * * * * ?"
	c.AddFunc(spec, func() {
		i++
		logger.Info.Println("cron running:", i)
	})

	//AddJob方法
	c.AddJob(spec, model.TaskScoreFootball{})
	//c.AddJob(spec, model.TaskScoreBasketBall{})

	//启动计划任务
	c.Start()

	//关闭着计划任务, 但是不能关闭已经在执行中的任务.
	defer c.Stop()

	select {}

	// 实时比分
	//go model.ScoreChange("http://api.wuhaicj.com/api/liveScore/change2", common.BasketBall)
	//model.ScoreChange("http://api.wuhaicj.com/api/liveScore/change", common.Football)
}

func oddsChange() {
	res, err := http.Get("http://api.wuhaicj.com/api/liveScore/oddsChange")
	if err != nil {
		logger.Error.Println("URL Request failed:", err)
		return
	}
	logger.Info.Println(res)
}
