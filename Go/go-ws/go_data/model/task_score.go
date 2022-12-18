package model

import (
	"encoding/json"
	"fmt"
	"go_data/common"
	"go_data/logger"
	"io"
	"net/http"
	"strconv"
	"time"
)

type TaskScoreFootball struct {
}

func (this TaskScoreFootball) Run() {
	logger.Info.Println("TaskScoreFootball start")
	ScoreChange("http://api.wuhaicj.com/api/liveScore/change", common.Football)
}

type TaskScoreBasketBall struct {
}

func (this TaskScoreBasketBall) Run() {
	logger.Info.Println("TaskScoreBasketBall start")
	ScoreChange("http://api.wuhaicj.com/api/liveScore/change2", common.BasketBall)
}

// 足球比分 当天比赛的比分数据（20秒变量）
func ScoreChange(url string, category int) {
	// 判断足球还是篮球
	var name = "Football"
	if category == common.BasketBall {
		name = "Basketball"
	}

	// 一、请求第三方接口
	res, err := http.Get(url)
	if err != nil /**/ {
		logger.Error.Println("URL Request failed:", err)
		return
	}
	msg, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		logger.Error.Println("Read body failed:", err)
		return
	}
	// 判断获取到的数据是否为空
	var changeList ChangeList
	err = json.Unmarshal(msg, &changeList)
	if err != nil {
		logger.Error.Println("json.Unmarshal failed:", err)
		return
	}
	changes := changeList.ChangeList
	if len(changes) == 0 {
		logger.Error.Println("changeList value empty")
		return
	}

	// 二、存Redis （给推送服务用）  一分钟内相同的数据不写入scoreChange中
	logger.Info.Println(name + "，数据获取成功")
	isHave, _ := Rdbc.Get(Ctx, "scoreChangeTemp:"+name).Result()
	if (isHave != "") && (isHave == common.Md5String(string(msg))) {
		logger.Info.Println(name + "，重复数据")
		return
	}
	// 开启Redis事务
	pipe := Rdbc.TxPipeline()
	// 临时存放去重
	pipe.Set(Ctx, "scoreChangeTemp:"+name, common.Md5String(string(msg)), time.Duration(60)*time.Second)
	// 将获取到的数据存入到Redis队列
	pipe.LPush(Ctx, "scoreChange:"+name, string(msg))
	_, err = pipe.Exec(Ctx)
	if err != nil {
		logger.Error.Println(name+"Redis 事务报错:", err)
		return
	}
	logger.Info.Println(name + "，Redis存储成功")

	// 三、更新数据库表
	//UpdateScore(changes)
	//fmt.Println("数据库更新成功")
}

func UpdateScore(changes []Change) {
	var table = "hn_schedule_copy1"
	for i := range changes {
		var c = changes[i]
		homeScore, _ := strconv.Atoi(c.HomeScore)
		awayScore, _ := strconv.Atoi(c.AwayScore)
		var sc = Schedule{
			MatchId:       c.MatchId,
			State:         c.State,
			HomeScore:     homeScore,
			AwayScore:     awayScore,
			HomeHalfScore: c.HomeHalfScore,
			AwayHalfScore: c.AwayHalfScore,
			HomeRed:       c.HomeRed,
			AwayRed:       c.AwayRed,
			HomeYellow:    c.HomeYellow,
			AwayYellow:    c.AwayYellow,
			HomeCorner:    c.HomeCorner,
			AwayCorner:    c.AwayCorner,
			UpdateTime:    time.Now().Format("2006/01/02 15:04:05"),
		}

		err := DB.Table(table).Where("matchId = ?", sc.MatchId).Updates(sc).Error
		if err != nil {
			fmt.Println("数据库更新分数错误：", err)
		}
	}
}
