package task

import (
	"github.com/valyala/fasthttp"
	"go-data/common"
	"go-data/configs"
	"go-data/model"
	"go-data/zlog"
	"time"
)

type ScoreFootball struct {
}

func (this ScoreFootball) Run() {
	zlog.Info.Println("足球比分变量 TaskScoreFootball start")
	scoreChange(configs.Conf.ApiF.ScoreChange, "Football")
}

type ScoreBasketBall struct {
}

func (this ScoreBasketBall) Run() {
	zlog.Info.Println("篮球比分变量 TaskScoreBasketBall start")
	scoreChange(configs.Conf.ApiB.ScoreChange, "BasketBall")
}

// 足球比分 当天比赛的比分数据（20秒变量）
func scoreChange(url string, scType string) {
	var cache = model.Rdb

	status, resp, err := fasthttp.Get(nil, url)
	if err != nil {
		zlog.Error.Println("请求失败:", err.Error())
		return
	}

	if status != fasthttp.StatusOK {
		zlog.Error.Println("请求没有成功:", status)
		return
	}

	// 判断获取到的数据是否为空
	if len(resp) < 50 {
		zlog.Info.Println("changeList value empty")
		return
	}

	// 二、存Redis （给推送服务用）  一分钟内相同的数据不写入scoreChange中
	isHave, _ := cache.Get(ctx, "scoreChangeTemp:"+scType).Result()
	if (isHave != "") && (isHave == common.Md5String(string(resp))) {
		//zlog.Info.Println("重复数据")
		return
	}
	// 开启Redis事务
	pipe := cache.TxPipeline()
	// 临时存放去重
	pipe.Set(ctx, "scoreChangeTemp:"+scType, common.Md5String(string(resp)), time.Duration(60)*time.Second)
	// 将获取到的数据存入到Redis队列
	pipe.LPush(ctx, "scoreChange:"+scType, string(resp))
	_, err = pipe.Exec(ctx)
	if err != nil {
		zlog.Error.Println("Redis 事务报错:", err)
		return
	}
	zlog.Info.Println("足球比分变量 Redis 存储成功！\r\n")

	// 三、更新数据库表
	//UpdateScore(changes)
	//fmt.Println("数据库更新成功")
}
