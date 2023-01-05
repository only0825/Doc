package task

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"loop-data/configs"
	"loop-data/model"
	"loop-data/utils"
	"net/http"
	"strconv"
	"time"
)

type ScoreChangeBasketball struct {
}

func (this ScoreChangeBasketball) Run() {
	scoreChange2(configs.Conf.ApiB.ScoreChange, "Basketball")
}

type ScoreChangeList2 struct {
	ChangeList []struct {
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
		Explain       string `json:"explain"`
	} `json:"changeList"`
}

// 足球比分 当天比赛的比分数据（20秒变量）
func scoreChange2(url string, scType string) {
	var cache = model.Rdb
	// 一、请求第三方数据
	resp, err := http.Get(url)
	if err != nil {
		logrus.Error("篮球-数据请求失败:", err.Error())
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("篮球-io.ReadAll失败:", err.Error())
		return
	}

	// 判断获取到的数据是否为空
	if len(body) < 50 {
		//logrus.Info("changeList value empty")
		return
	}

	var scl = ScoreChangeList2{}
	err = json.Unmarshal(body, &scl)
	if err != nil {
		logrus.Error("篮球-json反序列化错误", err)
		return
	}

	// 转为下划线的JSON
	clByte, err := json.Marshal(utils.JsonSnakeCase{scl.ChangeList})
	if err != nil {
		logrus.Error("篮球-json序列化错误", err)
		return
	}

	// 二、存Redis （给推送服务用）  一分钟内相同的数据不写入scoreChange中
	isHave, _ := cache.Get(ctx, "scoreChangeTemp:"+scType).Result()
	if (isHave != "") && (isHave == utils.Md5String(string(clByte))) {
		return
	}

	// 开启Redis事务
	pipe := cache.TxPipeline()
	// 临时存放去重
	pipe.Set(ctx, "scoreChangeTemp:"+scType, utils.Md5String(string(clByte)), time.Duration(60)*time.Second)
	// 将获取到的数据存入到Redis队列
	pipe.LPush(ctx, "scoreChange:"+scType, string(clByte))
	_, err = pipe.Exec(ctx)
	if err != nil {
		logrus.Error("篮球-Redis 事务报错:", err)
		return
	}

	logrus.Info("篮球-比分-变量 Redis 存储成功！")

	// 根据比赛ID存缓存，如果HTTP接口会先查缓存中的数据，再查数据库
	for i := range scl.ChangeList {
		cl := scl.ChangeList[i]
		match, _ := json.Marshal(cl)
		cache.Set(ctx, "newest:score:b:"+strconv.Itoa(cl.MatchId), string(match), time.Duration(120)*time.Second)
	}

	// 三、更新数据库表
	//for i := range scl.ChangeList {
	//	var c = scl.ChangeList[i]
	//	var sc = model.Schedule2{
	//		MatchId:    c.MatchID,
	//		RemainTime: c.RemainTime,
	//		State:      c.MatchState,
	//		HomeScore:  convToInt(c.HomeScore),
	//		Home1:      convToInt(c.Home1),
	//		Home2:      convToInt(c.Home2),
	//		Home3:      convToInt(c.Home3),
	//		Home4:      convToInt(c.Home4),
	//		HomeOT1:    convToInt(c.HomeOT1),
	//		HomeOT2:    convToInt(c.HomeOT2),
	//		HomeOT3:    convToInt(c.HomeOT3),
	//		AwayScore:  convToInt(c.AwayScore),
	//		Away1:      convToInt(c.Away1),
	//		Away2:      convToInt(c.Away2),
	//		Away3:      convToInt(c.Away3),
	//		Away4:      convToInt(c.Away4),
	//		AwayOT1:    convToInt(c.AwayOT1),
	//		AwayOT2:    convToInt(c.AwayOT2),
	//		AwayOT3:    convToInt(c.AwayOT3),
	//		UpdateTime: time.Now().Format("2006/01/02 15:04:05"),
	//	}
	//	err = model.UpdateScore2(sc)
	//	if err != nil {
	//		logrus.Error("篮球-数据库更新分数错误：", err)
	//		return
	//	}
	//}
	//
	//logrus.Info("篮球-比分-变量 Mysql 更新成功！")
}

func convToInt(s string) int {
	digit := utils.CtypeDigit(s)
	if !digit {
		s = "0"
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		logrus.Error("篮球-string转int报错：", err)
	}
	return i
}
