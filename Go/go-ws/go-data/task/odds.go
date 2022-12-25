package task

import (
	"encoding/json"
	"fmt"
	"go-data/common"
	"go-data/model"
	"go-data/zlog"
	"io"
	"net/http"
	"os"
	"time"
)

type OddsFootball struct {
}

func (this OddsFootball) Run() {
	zlog.Info.Println("TaskOddsFootball start")
	odds("http://api.wuhaicj.com/api/liveScore/odds", common.Football)
}

func newEuropeOdds(s []interface{}) *EuropeOdds {
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

func odds(url string, category int) {
	// 一、请求第三方接口
	res, err := http.Get(url)
	if err != nil {
		zlog.Error.Println("URL Request failed:", err)
		return
	}

	msg, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		zlog.Error.Println("Read body failed:", err)
		return
	}

	// 判断获取到的数据是否为空
	if len(msg) < 50 {
		zlog.Info.Println("list value empty")
		return
	}

	var obj Odds
	if err = json.Unmarshal(msg, &obj); err != nil {
		zlog.Error.Println("json反序列化失败: ", err)
		return
	}

	europeOdds := obj.List[0].EuropeOdds
	for i := range europeOdds {
		arr := europeOdds[i]
		e := newEuropeOdds(arr)
		fmt.Println(e.MatchId)
		fmt.Println(e.AwayWinMainOdds)

		os.Exit(1)
	}

	europeOddsByte, _ := json.Marshal(obj.List[0].EuropeOdds)
	saveToRedis(europeOddsByte, "europeOdds")

	overUnderByte, _ := json.Marshal(obj.List[0].OverUnder)
	saveToRedis(overUnderByte, "overUnder")
}

func saveToRedis(data []byte, changeType string) {
	model.Rdbc.Set(ctx, fmt.Sprintf("odds:football:%s", changeType), data, time.Duration(604800)*time.Second)
	zlog.Info.Println("odds redis 存储成功")
}

func typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}
