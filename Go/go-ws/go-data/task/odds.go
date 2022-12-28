package task

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"go-data/common"
	"go-data/configs"
	"go-data/model"
	"go-data/zlog"
	"time"
)

type OddsFootball struct {
}

func (this OddsFootball) Run() {
	zlog.Info.Println("足球指数全量 TaskOddsFootball start")
	odds(configs.Conf.ApiF.Odds)
}

type OddsChangeFootball struct {
}

func (this OddsChangeFootball) Run() {
	zlog.Info.Println("足球指数变量 TaskOddsFootball start")
	oddsChange(configs.Conf.ApiF.OddsChange, "Football")
}

// 足球指数全量，只存数据库
func odds(url string) {

	status, resp, err := fasthttp.Get(nil, url)
	if err != nil {
		zlog.Error.Println("请求失败:", err.Error())
		return
	}

	if status != fasthttp.StatusOK {
		zlog.Error.Println("请求没有成功:", status)
		return
	}

	// 判断获取到的数据是否为空
	if len(resp) < 50 {
		zlog.Info.Println("list value empty")
		return
	}

	var obj Odds
	if err = json.Unmarshal(resp, &obj); err != nil {
		zlog.Error.Println("json反序列化失败: ", err)
		return
	}

	europeOdds := obj.List[0].EuropeOdds
	for i := range europeOdds {
		arr := europeOdds[i]
		eo := newOddsEurope(arr)

		isClose := 2
		if eo.IsClose {
			isClose = 1
		}
		var oeNew = model.EuropeOdds{
			MatchId:          eo.MatchId,
			CompanyId:        eo.CompanyId,
			HomeWinEarlyOdds: eo.HomeWinEarlyOdds,
			TieEarlyOdds:     eo.TieEarlyOdds,
			AwayWinEarlyOdds: eo.AwayWinEarlyOdds,
			HomeWinMainOdds:  eo.HomeWinEarlyOdds,
			TieMainOdds:      eo.TieMainOdds,
			AwayWinMainOdds:  eo.AwayWinMainOdds,
			ChangeTime:       eo.ChangeTime,
			IsClose:          isClose,
			OddsType:         eo.OddsType,
		}
		find, err := model.OeFind(oeNew.MatchId)
		if err != nil {
			break
		}
		if find {
			err := model.OeUpdate(oeNew)
			if err != nil {
				break
			}
		} else {
			err := model.OeAdd(oeNew)
			if err != nil {
				break
			}
		}
	}

	overUnder := obj.List[0].OverUnder
	for i := range overUnder {
		arr := overUnder[i]
		ou := newOverUnder(arr)

		isClose := 2
		if ou.IsClose {
			isClose = 1
		}
		var ouNew = model.OverUnder{
			MatchId:            ou.MatchId,
			CompanyId:          ou.CompanyId,
			HandicapEarlyOdds:  ou.HandicapEarlyOdds,
			BigBallEarlyOdds:   ou.BigBallEarlyOdds,
			SmallBallEarlyOdds: ou.SmallBallEarlyOdds,
			HandicapOdds:       ou.HandicapOdds,
			BigBallOdds:        ou.BigBallOdds,
			SmallBallOdds:      ou.SmallBallOdds,
			ChangeTime:         ou.ChangeTime,
			IsClose:            isClose,
			OddsType:           ou.OddsType,
		}
		find, err := model.OuFind(ouNew.MatchId)
		if err != nil {
			break
		}
		if find {
			err := model.OuUpdate(ouNew)
			if err != nil {
				break
			}
		} else {
			err := model.OuAdd(ouNew)
			if err != nil {
				break
			}
		}
	}

	zlog.Info.Println("足球指数全量 Mysql 存储成功！ \r\n")
}

func oddsChange(url string, odType string) {
	var cache = model.Rdb

	status, resp, err := fasthttp.Get(nil, url)
	if err != nil {
		zlog.Error.Println("请求失败:", err.Error())
		return
	}

	if status != fasthttp.StatusOK {
		zlog.Error.Println("请求没有成功:", status)
		return
	}

	// 判断获取到的数据是否为空
	if len(resp) < 50 {
		zlog.Info.Println("changeList value empty")
		return
	}

	// 二、存Redis （给推送服务用）  一分钟内相同的数据不写入oddsChange中
	isHave, _ := cache.Get(ctx, "oddsChangeTemp:"+odType).Result()
	if (isHave != "") && (isHave == common.Md5String(string(resp))) {
		return
	}
	// 开启Redis事务
	pipe := cache.TxPipeline()
	// 临时存放去重
	pipe.Set(ctx, "oddsChangeTemp:"+odType, common.Md5String(string(resp)), time.Duration(60)*time.Second)
	// 将获取到的数据存入到Redis队列
	pipe.LPush(ctx, "oddsChange:"+odType, string(resp))
	_, err = pipe.Exec(ctx)
	if err != nil {
		zlog.Error.Println("Redis 事务报错:", err)
		return
	}
	zlog.Info.Println("足球比分变量 Redis 存储成功！\r\n")

}

func newOddsEurope(s []interface{}) *EuropeOdds {
	eo := &EuropeOdds{}

	if len(s) > 0 {
		eo.MatchId = int(s[0].(float64))
	}

	if len(s) > 1 {
		eo.CompanyId = int(s[1].(float64))
	}

	if len(s) > 2 {
		eo.HomeWinEarlyOdds = s[2].(float64)
	}

	if len(s) > 3 {
		eo.TieEarlyOdds = s[3].(float64)
	}

	if len(s) > 4 {
		eo.AwayWinEarlyOdds = s[4].(float64)
	}

	if len(s) > 5 {
		eo.HomeWinMainOdds = s[5].(float64)
	}

	if len(s) > 6 {
		eo.TieMainOdds = s[6].(float64)
	}

	if len(s) > 7 {
		eo.AwayWinMainOdds = s[7].(float64)
	}

	if len(s) > 8 {
		eo.ChangeTime = s[8].(string)
	}

	if len(s) > 9 {
		eo.IsClose = s[9].(bool)
	}

	if len(s) > 10 {
		eo.OddsType = int(s[10].(float64))
	}

	return eo
}

func newOverUnder(s []interface{}) *OverUnder {
	ou := &OverUnder{}

	if len(s) > 0 {
		ou.MatchId = int(s[0].(float64))
	}

	if len(s) > 1 {
		ou.CompanyId = int(s[1].(float64))
	}

	if len(s) > 2 {
		ou.HandicapEarlyOdds = s[2].(float64)
	}

	if len(s) > 3 {
		ou.BigBallEarlyOdds = s[3].(float64)
	}

	if len(s) > 4 {
		ou.SmallBallEarlyOdds = s[4].(float64)
	}

	if len(s) > 5 {
		ou.HandicapOdds = s[5].(float64)
	}

	if len(s) > 6 {
		ou.BigBallOdds = s[6].(float64)
	}

	if len(s) > 7 {
		ou.SmallBallOdds = s[7].(float64)
	}

	if len(s) > 8 {
		ou.ChangeTime = s[8].(string)
	}

	if len(s) > 9 {
		ou.IsClose = s[9].(bool)
	}

	if len(s) > 10 {
		ou.OddsType = int(s[10].(float64))
	}

	return ou
}
