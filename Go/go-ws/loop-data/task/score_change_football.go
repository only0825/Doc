package task

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"loop-data/common"
	"loop-data/configs"
	"loop-data/model"
	"net/http"
	"time"
)

type ScoreChangeFootball struct {
}

func (this ScoreChangeFootball) Run() {
	//logrus.Info("足球比分变量 TaskScoreFootball start")
	scoreChange(configs.Conf.ApiF.ScoreChange, "Football")
}

// 足球比分 当天比赛的比分数据（20秒变量）
func scoreChange(url string, scType string) {
	var cache = model.Rdbc

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
		//logrus.Info("changeList value empty")
		return
	}

	var scl = ScoreChangeList{}
	err = json.Unmarshal(body, &scl)
	if err != nil {
		logrus.Error("json 反序列化错误", err)
		return
	}

	clArr := []ChangeListJson{}
	for i := range scl.ChangeList {
		sclData := scl.ChangeList[i]
		cl := ChangeListJson{
			MatchId:       sclData.MatchId,
			State:         sclData.State,
			HomeScore:     sclData.HomeScore,
			AwayScore:     sclData.AwayScore,
			HomeHalfScore: sclData.HomeHalfScore,
			AwayHalfScore: sclData.AwayHalfScore,
			HomeRed:       sclData.HomeCorner,
			AwayRed:       sclData.AwayRed,
			HomeYellow:    sclData.HomeYellow,
			AwayYellow:    sclData.AwayYellow,
			HomeCorner:    sclData.HomeCorner,
			AwayCorner:    sclData.AwayCorner,
			HasLineup:     sclData.HasLineup,
			MatchTime:     sclData.MatchTime,
			StartTime:     sclData.StartTime,
			Explain:       sclData.Explain,
			ExtraExplain:  sclData.ExtraExplain,
			InjuryTime:    sclData.InjuryTime,
		}
		clArr = append(clArr, cl)
	}

	clByte, err := json.Marshal(clArr)
	if err != nil {
		logrus.Error("json 序列化错误", err)
		return
	}

	// 二、存Redis （给推送服务用）  一分钟内相同的数据不写入scoreChange中
	isHave, _ := cache.Get(ctx, "scoreChangeTemp:"+scType).Result()
	if (isHave != "") && (isHave == common.Md5String(string(clByte))) {
		return
	}

	// 开启Redis事务
	pipe := cache.TxPipeline()
	// 临时存放去重
	pipe.Set(ctx, "scoreChangeTemp:"+scType, common.Md5String(string(clByte)), time.Duration(60)*time.Second)
	// 将获取到的数据存入到Redis队列
	pipe.LPush(ctx, "scoreChange:"+scType, string(clByte))
	_, err = pipe.Exec(ctx)
	if err != nil {
		logrus.Error("Redis 事务报错:", err)
		return
	}

	logrus.Info("足球比分变量 Redis 存储成功！")

	// 三、更新数据库表
	//UpdateScore(changes)
	//fmt("数据库更新成功")
}
