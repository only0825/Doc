package main

import (
	"github.com/robfig/cron/v3"
	"go-data/common"
	"go-data/configs"
	"go-data/model"
	"go-data/task"
	"go-data/zlog"
)

func init() {
	if err := configs.LoadConfig(); err != nil {
		zlog.Error.Println("Load config json error:", err)
		return
	}

	rdbc, err := common.InitRedisCluster()
	if err != nil {
		zlog.Error.Println("Redis初始化失败: ", err)
		return
	}

	db, err := common.InitMysql()
	if err != nil {
		zlog.Error.Println("数据库初始化失败: ", err)
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
		zlog.Info.Println("cron running:", i)
	})

	// 足球变量
	c.AddJob(spec, task.ScoreFootball{}) // 足球比分
	//c.AddJob(spec, task.OddsFootball{})  // 足球指数

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
