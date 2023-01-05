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

type ScoreChangeFootball struct {
}

func (this ScoreChangeFootball) Run() {
	scoreChange(configs.Conf.ApiF.ScoreChange, "Football")
}

// 接收第三方数据的结构体
type ScoreChangeList1 struct {
	ChangeList []struct {
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
	} `json:"changeList"`
}

// 足球比分 当天比赛的比分数据（20秒变量）
func scoreChange(url string, scType string) {
	var cache = model.Rdb
	// 一、请求第三方数据
	resp, err := http.Get(url)
	if err != nil {
		logrus.Error("足球-请求失败:", err.Error())
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("足球-io.ReadAll失败:", err.Error())
		return
	}

	// 判断获取到的数据是否为空
	if len(body) < 50 {
		//logrus.Info("changeList value empty")
		return
	}

	var scl = ScoreChangeList1{}
	err = json.Unmarshal(body, &scl)
	if err != nil {
		logrus.Error("足球-json 反序列化错误", err)
		return
	}

	// 转为下划线的JSON
	clByte, err := json.Marshal(utils.JsonSnakeCase{scl.ChangeList})
	if err != nil {
		logrus.Error("足球-json 序列化错误", err)
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
		logrus.Error("足球-Redis 事务报错:", err)
		return
	}

	logrus.Info("足球-比分-变量 Redis 存储成功！")

	for i := range scl.ChangeList {
		cl := scl.ChangeList[i]
		match, _ := json.Marshal(cl)
		cache.Set(ctx, "newest:score:f:"+strconv.Itoa(cl.MatchId), string(match), time.Duration(120)*time.Second)
	}

	// 三、更新数据库表
	//for i := range scl.ChangeList {
	//	var c = scl.ChangeList[i]
	//	var sc = model.Schedule1{
	//		MatchId:       c.MatchId,
	//		State:         c.State,
	//		HomeScore:     c.HomeScore,
	//		AwayScore:     c.AwayScore,
	//		HomeHalfScore: c.HomeHalfScore,
	//		AwayHalfScore: c.AwayHalfScore,
	//		HomeRed:       c.HomeRed,
	//		AwayRed:       c.AwayRed,
	//		HomeYellow:    c.HomeYellow,
	//		AwayYellow:    c.AwayYellow,
	//		HomeCorner:    c.HomeCorner,
	//		AwayCorner:    c.AwayCorner,
	//		UpdateTime:    time.Now().Format("2006/01/02 15:04:05"),
	//	}
	//	err = model.UpdateScore1(sc)
	//	if err != nil {
	//		logrus.Error("足球-数据库更新分数错误：", err)
	//		return
	//	}
	//}
	//
	//logrus.Info("足球-比分-变量 Mysql 更新成功！")
}
