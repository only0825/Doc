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

// 足球指数 变量
type OddsChangeFootball struct {
}

func (this OddsChangeFootball) Run() {
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
		// 根据比赛ID存缓存，如果HTTP接口会先查缓存中的数据，再查数据库
		match, _ := json.Marshal(oeNew)
		cache.Set(ctx, "newest:odds:f:europe:"+strconv.Itoa(oeNew.MatchId), string(match), time.Duration(120)*time.Second)
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
		// 根据比赛ID存缓存，如果HTTP接口会先查缓存中的数据，再查数据库
		match, _ := json.Marshal(ouNew)
		cache.Set(ctx, "newest:odds:f:overunder:"+strconv.Itoa(ouNew.MatchId), string(match), time.Duration(120)*time.Second)
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

	logrus.Info("足球-指数-变量 Redis 存储成功！")

	// 三、更新数据库数据
	//updateChangeOdds(rc)
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

//func updateChangeOdds(rc RespOddsChange) {
//	for i := range rc.EuropeOdds {
//		var e = rc.EuropeOdds[i]
//		var eo = model.EuropeOdds{
//			MatchId:         e.MatchId,
//			HomeWinMainOdds: e.HomeWinMainOdds,
//			TieMainOdds:     e.TieMainOdds,
//			AwayWinMainOdds: e.AwayWinMainOdds,
//			ChangeTime:      e.ChangeTime,
//			IsClose:         boolToInt(e.IsClose),
//			OddsType:        e.OddsType,
//			UpdateTime:      time.Now().Format("2006/01/02 15:04:05"),
//		}
//		err := model.OeUpdate(eo)
//		if err != nil {
//			return
//		}
//	}
//
//	for i := range rc.OverUnder {
//		var o = rc.OverUnder[i]
//		var ou = model.OverUnder{
//			MatchId:       o.MatchId,
//			HandicapOdds:  o.HandicapOdds,
//			BigBallOdds:   o.BigBallOdds,
//			SmallBallOdds: o.SmallBallOdds,
//			ChangeTime:    o.ChangeTime,
//			IsClose:       boolToInt(o.IsClose),
//			OddsType:      o.OddsType,
//			UpdateTime:    time.Now().Format("2006/01/02 15:04:05"),
//		}
//		err := model.OuUpdate(ou)
//		if err != nil {
//			return
//		}
//	}
//}

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

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 2
}
