package main

import (
	"github.com/robfig/cron/v3"
	"go-data/common"
	"go-data/configs"
	"go-data/task"
	"go-data/zlog"
	"os"
)

// 返回一个支持至 秒 级别的 cron
func newWithSeconds() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}

func main() {

	// 获取命令行参数
	argc := len(os.Args)
	if argc != 2 {
		zlog.Error.Println("运行格式错误，格式为 ./应用 <配置文件名称>")
		return
	}

	if err := configs.LoadConfig(os.Args[1]); err != nil {
		zlog.Error.Println("Load config json error:", err)
		return
	}

	cache := configs.Conf.Cache
	err := common.InitCache(cache)
	if err != nil {
		zlog.Error.Println("Redis初始化错误:", err)
		return
	}

	err = common.InitMysql()
	if err != nil {
		zlog.Error.Println("数据库初始化失败: ", err)
		return
	}

	c := newWithSeconds()

	//spec1 := "*/5 * * * * ?" // 每隔5秒执行一次
	//spec2 := "0 */1 * * * ?" // 每隔1分钟执行一次
	spec2 := "*/1 * * * * ?" //

	// 足球比分变量
	//c.AddJob(spec1, task.ScoreFootball{})

	// 足球 主盘口即时赔率（全量）
	c.AddJob(spec2, task.OddsFootball{})

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
