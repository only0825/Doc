package model

import (
	"encoding/json"
	"fmt"
	"go_data/common"
	"go_data/logger"
	"io"
	"net/http"
	"time"
)

// 足球比分 当天比赛的比分数据（20秒变量）
func ScoreChange(url string, category int) {
	// 判断足球还是篮球
	var name = "Football"
	if category == common.BasketBall {
		name = "Basketball"
	}

	for {
		// 一、请求第三方接口
		time.Sleep(time.Duration(3) * time.Second)
		res, err := http.Get(url)
		if err != nil /**/ {
			logger.Error.Println("URL Request failed:", err)
			continue
		}
		msg, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			logger.Error.Println("Read body failed:", err)
			continue
		}
		// 判断获取到的数据是否为空
		var m = make(map[string]interface{})
		err = json.Unmarshal(msg, &m)
		if err != nil {
			logger.Error.Println("json.Unmarshal failed:", err)
			continue
		}
		if m["changeList"] == "" {
			logger.Error.Println("changeList empty")
			continue
		}

		// 二、存Redis （给推送服务用）  一分钟内相同的数据不写入scoreChange中
		logger.Info.Println(name + "，数据获取成功")
		isHave, _ := Rdbc.Get(Ctx, "scoreChangeTemp:"+name).Result()
		if (isHave != "") && (isHave == common.Md5String(string(msg))) {
			logger.Info.Println(name + "，重复数据")
			continue
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
			continue
		}
		logger.Info.Println(name + "，Redis存储成功")

		// 三、更新数据库表
		//model.UpdateScore(string(msg))
		//fmt.Println("数据库更新成功")
	}
}

func UpdateScore(msg string) {
	var table = "hn_schedule_copy1"
	var changeList ChangeList
	// 将json消息转为数组
	err := json.Unmarshal([]byte(msg), &changeList)
	if err != nil {
		fmt.Println("JSON转换失败：", err)
	}

	var changes []Change
	changes = changeList.ChangeList
	for i := range changes {
		var c = changes[i]
		var sc = Schedule{
			MatchId:       c.MatchId,
			State:         c.State,
			HomeScore:     c.HomeScore,
			AwayScore:     c.AwayScore,
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

		err = DB.Table(table).Where("matchId = ?", sc.MatchId).Updates(sc).Error
		if err != nil {
			fmt.Println("数据库更新分数错误：", err)
		}
	}
}
