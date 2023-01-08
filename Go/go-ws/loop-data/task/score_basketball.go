package task

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"loop-data/configs"
	"loop-data/model"
	"net/http"
	"time"
)

// 篮球比分 全量
type ScoreBasketball struct {
}

func (this ScoreBasketball) Run() {
	//logrus.Info("足球指数全量 TaskOddsFootball start")
	score2(configs.Conf.ApiB.Score)
}

type ScoreList2 struct {
	MatchList []struct {
		MatchId       int    `json:"matchId"`
		MatchState    int    `json:"matchState"`
		RemainTime    string `json:"remainTime"`
		OvertimeCount int    `json:"overtimeCount"`
		HomeScore     string `json:"homeScore"`
		Home1         string `json:"home1"`
		Home2         string `json:"home2"`
		Home3         string `json:"home3"`
		Home4         string `json:"home4"`
		HomeOT1       string `json:"homeOT1"`
		HomeOT2       string `json:"homeOT2"`
		HomeOT3       string `json:"homeOT3"`
		AwayScore     string `json:"awayScore"`
		Away1         string `json:"away1"`
		Away2         string `json:"away2"`
		Away3         string `json:"away3"`
		Away4         string `json:"away4"`
		AwayOT1       string `json:"awayOT1"`
		AwayOT2       string `json:"awayOT2"`
		AwayOT3       string `json:"awayOT3"`
		HasStats      bool   `json:"hasStats"`
		ExplainCn     string `json:"explainCn"`
	} `json:"matchList"`
}

// 篮球比分全量，只存数据库
func score2(url string) {
	resp, err := http.Get(url)
	if err != nil {
		logrus.Error("请求失败:", err.Error())
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("io.ReadAll失败:", err.Error())
		return
	}

	// 判断获取到的数据是否为空
	if len(body) < 50 {
		//logrus.Info("list value empty")
		return
	}

	var sl2 ScoreList2
	if err = json.Unmarshal(body, &sl2); err != nil {
		logrus.Error("json反序列化失败: ", err)
		return
	}

	// 更新数据库表
	for i := range sl2.MatchList {
		var c = sl2.MatchList[i]
		var sc = model.Schedule2{
			RemainTime: c.RemainTime,
			State:      c.MatchState,
			HomeScore:  convToInt(c.HomeScore),
			Home1:      convToInt(c.Home1),
			Home2:      convToInt(c.Home2),
			Home3:      convToInt(c.Home3),
			Home4:      convToInt(c.Home4),
			HomeOT1:    convToInt(c.HomeOT1),
			HomeOT2:    convToInt(c.HomeOT2),
			HomeOT3:    convToInt(c.HomeOT3),
			AwayScore:  convToInt(c.AwayScore),
			Away1:      convToInt(c.Away1),
			Away2:      convToInt(c.Away2),
			Away3:      convToInt(c.Away3),
			Away4:      convToInt(c.Away4),
			AwayOT1:    convToInt(c.AwayOT1),
			AwayOT2:    convToInt(c.AwayOT2),
			AwayOT3:    convToInt(c.AwayOT3),
			UpdateTime: time.Now().Format("2006/01/02 15:04:05"),
		}
		err = model.UpdateScore2(sc, c.MatchId)
		if err != nil {
			logrus.Error("篮球-数据库更新分数错误：", err)
			return
		}
	}

	logrus.Info("篮球-比分-全量 Mysql 更新成功！")
}
