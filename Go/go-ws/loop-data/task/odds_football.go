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

// 足球指数 全量
type OddsFootball struct {
}

func (this OddsFootball) Run() {
	//logrus.Info("足球指数全量 TaskOddsFootball start")
	odds(configs.Conf.ApiF.Odds)
}

// 足球指数全量，只存数据库
func odds(url string) {
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

	var obj Odds
	if err = json.Unmarshal(body, &obj); err != nil {
		logrus.Error("json反序列化失败: ", err)
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
			UpdateTime:         time.Now().Format("2006/01/02 15:04:05"),
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

	logrus.Info("足球-指数-全量 Mysql 存储成功！")
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
