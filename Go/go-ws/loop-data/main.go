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

	rdb, err := utils.InitRedis()
	model.Rdb = rdb
	//rdb, err := utils.InitRedisCluster()
	//model.Rdb = rdb
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
	spec1 := "*/8 * * * * ?"  // 每隔8秒执行一次
	spec2 := "0 */5 * * * ?"  // 每隔5分钟执行一次
	spec3 := "0 */6 * * * ?"  // 每隔6分执行一次
	spec4 := "0 */4 * * * ?"  // 每隔4分钟执行一次
	spec5 := "*/10 * * * * ?" // 每隔10秒执行一次

	// 足球比分变量   Redis 1.推送数据  2.最新数据
	c.AddJob(spec1, task.ScoreChangeFootball{})

	// 足球指数变量   Redis 1.推送数据  2.最新数据
	c.AddJob(spec5, task.OddsChangeFootball{})
	// 足球指数全量   Mysql
	c.AddJob(spec2, task.OddsFootball{})

	// 篮球比分变量	Redis 1.推送数据  2.最新数据  TODO 暂时不上
	//c.AddJob(spec1, task.ScoreChangeBasketball{})
	// 篮球比分全量   Mysql
	c.AddJob(spec4, task.ScoreBasketball{})
	// 篮球 技术统计 （某场比赛的技术统计和球员统计） Redis
	c.AddJob(spec3, task.StatsBasketball{})

	// TODO 动画直播赛程赛果  这是给前端判断比赛有没有动画直播 下个版本才做
	//c.AddJob(spec4, task.AnimeFootball{})

	//启动计划任务
	c.Start()

	//关闭着计划任务, 但是不能关闭已经在执行中的任务.
	defer c.Stop()

	select {}

	// 实时比分
	//go model.ScoreChange("http://api.wuhaicj.com/api/liveScore/change2", utils.BasketBall)
	//model.ScoreChange("http://api.wuhaicj.com/api/liveScore/change", utils.Football)
}
