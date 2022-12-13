package model

import (
	"encoding/json"
	"go-data/common"
	"go-data/logger"
	"io"
	"net/http"
	"time"
)

type Change struct {
	MatchId       int    `json:"matchId"`       // 比赛ID
	State         int    `json:"state"`         // 比赛
	HomeScore     int    `json:"homeScore"`     // 主队得分
	AwayScore     int    `json:"awayScore"`     // 客队得分
	HomeHalfScore int    `json:"homeHalfScore"` // 主队上半场得分
	AwayHalfScore int    `json:"awayHalfScore"` // 客队上半场得分
	HomeRed       int    `json:"homeRed"`       // 主队红牌数
	AwayRed       int    `json:"awayRed"`       // 客队红牌数
	HomeYellow    int    `json:"homeYellow"`    // 主队黄牌数
	AwayYellow    int    `json:"awayYellow"`    // 客队红牌数
	HomeCorner    int    `json:"homeCorner"`    // 主队角球数
	AwayCorner    int    `json:"awayCorner"`    // 客队角球数
	HasLineup     string `json:"hasLineup"`     // 是否有阵容
	MatchTime     string `json:"matchTime"`     // 比赛时间
	StartTime     string `json:"startTime"`     // 开场时间
	Explain       string `json:"explain"`       // 比赛说明
	ExtraExplain  string `json:"extraExplain"`  // 比赛说明2
	InjuryTime    string `json:"injuryTime"`    // 上下半场补时时长
}

type ChangeList struct {
	ChangeList []Change
}

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
		//UpdateScore(string(msg))
		//fmt.Println("数据库更新成功")
	}
}
