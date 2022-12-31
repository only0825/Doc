package task

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"loop-data/configs"
	"loop-data/model"
	"loop-data/utils"
	"net/http"
	"time"
)

// 足球指数 变量
type OddsChangeFootball struct {
}

func (this OddsChangeFootball) Run() {
	//logrus.Info("足球指数变量 TaskOddsFootball start")
	oddsChange(configs.Conf.ApiF.OddsChange, "Football")
}

func oddsChange(url string, odType string) {
	var cache = model.Rdb

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

	var obj OddsChange
	if err = json.Unmarshal(body, &obj); err != nil {
		logrus.Error("json 反序列化失败: ", err)
		return
	}

	var eocArr []EuropeOddsChange
	europeOdds := obj.ChangeList[0].EuropeOdds

	for i := range europeOdds {
		arr := europeOdds[i]
		neo := newOddsEuropeChange(arr)
		var oeNew = EuropeOddsChange{
			MatchId:         neo.MatchId,
			CompanyId:       neo.CompanyId,
			HomeWinMainOdds: neo.HomeWinMainOdds,
			TieMainOdds:     neo.TieMainOdds,
			AwayWinMainOdds: neo.AwayWinMainOdds,
			ChangeTime:      neo.ChangeTime,
			IsClose:         neo.IsClose,
			OddsType:        neo.OddsType,
		}
		eocArr = append(eocArr, oeNew)
	}

	var oucArr []OverUnderChange
	overUnder := obj.ChangeList[0].OverUnder
	for i := range overUnder {
		arr := overUnder[i]
		nouc := newOverUnderChange(arr)
		var ouNew = OverUnderChange{
			MatchId:       nouc.MatchId,
			CompanyId:     nouc.CompanyId,
			HandicapOdds:  nouc.HandicapOdds,
			BigBallOdds:   nouc.BigBallOdds,
			SmallBallOdds: nouc.SmallBallOdds,
			ChangeTime:    nouc.ChangeTime,
			IsClose:       nouc.IsClose,
			OddsType:      nouc.OddsType,
		}
		oucArr = append(oucArr, ouNew)
	}

	rc := RespOddsChange{
		EuropeOdds: eocArr,
		OverUnder:  oucArr,
	}

	saveByte, err := json.Marshal(rc)
	if err != nil {
		logrus.Error("json 序列化错误:", err)
		return
	}

	// 二、存Redis （给推送服务用）  一分钟内相同的数据不写入oddsChange中
	isHave, _ := cache.Get(ctx, "oddsChangeTemp:"+odType).Result()
	if (isHave != "") && (isHave == utils.Md5String(string(saveByte))) {
		return
	}

	// 开启Redis事务
	pipe := cache.TxPipeline()
	// 临时存放去重
	pipe.Set(ctx, "oddsChangeTemp:"+odType, utils.Md5String(string(saveByte)), time.Duration(60)*time.Second)
	// 将获取到的数据存入到Redis队列
	pipe.LPush(ctx, "oddsChange:"+odType, string(saveByte))
	_, err = pipe.Exec(ctx)
	if err != nil {
		logrus.Error("Redis 事务报错:", err)
		return
	}

	logrus.Info("足球指数变量 Redis 存储成功！")
}

func newOddsEuropeChange(s []interface{}) *EuropeOddsChange {
	eo := &EuropeOddsChange{}

	if len(s) > 0 {
		eo.MatchId = int(s[0].(float64))
	}

	if len(s) > 1 {
		eo.CompanyId = int(s[1].(float64))
	}

	if len(s) > 2 {
		eo.HomeWinMainOdds = s[2].(float64)
	}

	if len(s) > 3 {
		eo.TieMainOdds = s[3].(float64)
	}

	if len(s) > 4 {
		eo.AwayWinMainOdds = s[4].(float64)
	}

	if len(s) > 5 {
		eo.ChangeTime = s[5].(string)
	}

	if len(s) > 6 {
		eo.IsClose = s[6].(bool)
	}

	if len(s) > 7 {
		eo.OddsType = int(s[7].(float64))
	}

	return eo
}

func newOverUnderChange(s []interface{}) *OverUnderChange {
	ouc := &OverUnderChange{}

	if len(s) > 0 {
		ouc.MatchId = int(s[0].(float64))
	}

	if len(s) > 1 {
		ouc.CompanyId = int(s[1].(float64))
	}

	if len(s) > 2 {
		ouc.HandicapOdds = s[2].(float64)
	}

	if len(s) > 3 {
		ouc.BigBallOdds = s[3].(float64)
	}

	if len(s) > 4 {
		ouc.SmallBallOdds = s[4].(float64)
	}

	if len(s) > 5 {
		ouc.ChangeTime = s[5].(string)
	}

	if len(s) > 6 {
		ouc.IsClose = s[6].(bool)
	}

	if len(s) > 7 {
		ouc.OddsType = int(s[7].(float64))
	}

	return ouc
}
