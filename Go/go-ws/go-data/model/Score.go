package model

import (
	"encoding/json"
	"fmt"
	"go-data/common"
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
func ScoreChange() {
	for {
		time.Sleep(time.Duration(3) * time.Second)
		res, err := http.Get("http://api.wuhaicj.com/api/liveScore/change")
		if err != nil /**/ {
			fmt.Println("URL Request failed:", err)
			continue
		}
		msg, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Read body failed:", err)
			continue
		}

		// 一分钟内相同的数据不写入scoreChange中
		isHave, _ := common.Rdbc.Get(common.Ctx, "scoreChangeTemp").Result()
		if (isHave != "") && (isHave == common.Md5String(string(msg))) {
			continue
		}
		// 临时存放去重
		common.Rdbc.Set(common.Ctx, "scoreChangeTemp", common.Md5String(string(msg)), time.Duration(60)*time.Second)
		// 将获取到的数据存入到Redis队列
		common.Rdbc.LPush(common.Ctx, "scoreChange", string(msg))

		defer res.Body.Close()
	}

}

func UpdateScore(msg string) {
	var changeList ChangeList
	// 将json消息转为数组
	err := json.Unmarshal([]byte(msg), &changeList)
	fmt.Println(err)
	fmt.Println(changeList)
	//arr2 := arr["changeList"]
	//for s := range arr2 {
	//	fmt.Println(s)
	//}
	// works with Take
	//result := map[string]interface{}{}
	//DB.Table("hn_content").Take(&result)
	//fmt.Println(result)
	//
}
