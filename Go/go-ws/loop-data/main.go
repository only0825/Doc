package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"loop-data/configs"
	"loop-data/model"
	"loop-data/task"
	"loop-data/utils"
	"os"
)

// 返回一个支持至 秒 级别的 cron
func newWithSeconds() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}

func main() {
	fmt.Println("staring loop ...")
	utils.LogInit("/loop-data.log")

	var err error
	// 获取命令行参数
	argc := len(os.Args)
	if argc != 2 {
		logrus.Error("运行格式错误，格式为 ./应用 <配置文件名称>")
		return
	}

	if err = configs.LoadConfig(os.Args[1]); err != nil {
		logrus.Error("Load config json error:", err)
		return
	}

	//rdb, err := utils.InitRedis()
	//model.Rdb = rdb
	rdb, err := utils.InitRedisCluster()
	model.Rdb = rdb
	if err != nil {
		logrus.Error("Redis初始化错误:", err)
		return
	}

	err = utils.InitMysql()
	if err != nil {
		logrus.Error("数据库初始化失败: ", err)
		return
	}

	c := newWithSeconds()
	spec1 := "*/3 * * * * ?" // 每隔5秒执行一次
	spec2 := "0 */1 * * * ?" // 每隔1分钟执行一次

	// 足球比分变量
	c.AddJob(spec1, task.ScoreChangeFootball{})

	// 足球指数变量
	c.AddJob(spec1, task.OddsChangeFootball{})

	// 足球 主盘口即时赔率（全量）
	c.AddJob(spec2, task.OddsFootball{})

	//启动计划任务
	c.Start()

	//关闭着计划任务, 但是不能关闭已经在执行中的任务.
	defer c.Stop()

	select {}

	// 实时比分
	//go model.ScoreChange("http://api.wuhaicj.com/api/liveScore/change2", utils.BasketBall)
	//model.ScoreChange("http://api.wuhaicj.com/api/liveScore/change", utils.Football)
}
