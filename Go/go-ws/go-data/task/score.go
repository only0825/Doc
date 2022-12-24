package task

import (
	"go-data/common"
	"go-data/configs"
	"go-data/model"
	"go-data/zlog"
	"io"
	"net/http"
	"time"
)

type ScoreFootball struct {
}

func (this ScoreFootball) Run() {
	zlog.Info.Println("TaskScoreFootball start")
	scoreChange(configs.Conf.ApiF.ScoreChange, "Football")
}

type ScoreBasketBall struct {
}

func (this ScoreBasketBall) Run() {
	zlog.Info.Println("TaskScoreBasketBall start")
	scoreChange(configs.Conf.ApiB.ScoreChange, "BasketBall")
}

// 足球比分 当天比赛的比分数据（20秒变量）
func scoreChange(url string, scType string) {
	var cache = model.Rdbc

	// 一、请求第三方接口
	res, err := http.Get(url)
	if err != nil /**/ {
		zlog.Error.Println("URL Request failed:", err)
		return
	}

	msg, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		zlog.Error.Println("Read body failed:", err)
		return
	}

	// 判断获取到的数据是否为空
	if len(msg) < 50 {
		zlog.Info.Println("changeList value empty")
		return
	}

	// 二、存Redis （给推送服务用）  一分钟内相同的数据不写入scoreChange中
	zlog.Info.Println("数据获取成功")
	isHave, _ := cache.Get(ctx, "scoreChangeTemp:"+scType).Result()
	if (isHave != "") && (isHave == common.Md5String(string(msg))) {
		zlog.Info.Println("重复数据")
		return
	}
	// 开启Redis事务
	pipe := cache.TxPipeline()
	// 临时存放去重
	pipe.Set(ctx, "scoreChangeTemp:"+scType, common.Md5String(string(msg)), time.Duration(60)*time.Second)
	// 将获取到的数据存入到Redis队列
	pipe.LPush(ctx, "scoreChange:"+scType, string(msg))
	_, err = pipe.Exec(ctx)
	if err != nil {
		zlog.Error.Println("Redis 事务报错:", err)
		return
	}
	zlog.Info.Println("Redis存储成功")

	// 三、更新数据库表
	//UpdateScore(changes)
	//fmt.Println("数据库更新成功")
}
