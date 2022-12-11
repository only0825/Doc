package main

import (
	"fmt"
	"go-data/common"
	"go-data/model"
	"net/http"
	"time"
)

type Content struct {
	id      int
	title   string
	content string
}

var m = make(map[string]interface{})

func init() {
	common.InitRedisCluster()
	common.InitDb()
}

func main() {
	//go oddsChange()
	//model.ScoreChange()

	str := "{\"changeList\":[{\"matchId\":2318090,\"state\":3,\"homeScore\":0,\"awayScore\":1,\"homeHalfScore\":0,\"awayHalfScore\":1,\"homeRed\":0,\"awayRed\":0,\"homeYellow\":0,\"awayYellow\":0,\"homeCorner\":0,\"awayCorner\":0,\"hasLineup\":\"\",\"matchTime\":\"2022-12-11 18:45:00\",\"startTime\":\"2022-12-11 19:56:57\",\"explain\":\"\",\"extraExplain\":\"\",\"injuryTime\":\"\"},{\"matchId\":2267198,\"state\":3,\"homeScore\":3,\"awayScore\":0,\"homeHalfScore\":2,\"awayHalfScore\":0,\"homeRed\":0,\"awayRed\":0,\"homeYellow\":2,\"awayYellow\":2,\"homeCorner\":3,\"awayCorner\":3,\"hasLineup\":\"\",\"matchTime\":\"2022-12-11 19:00:00\",\"startTime\":\"2022-12-11 20:02:13\",\"explain\":\"\",\"extraExplain\":\"\",\"injuryTime\":\"\"},{\"matchId\":2286901,\"state\":3,\"homeScore\":2,\"awayScore\":3,\"homeHalfScore\":0,\"awayHalfScore\":2,\"homeRed\":0,\"awayRed\":0,\"homeYellow\":2,\"awayYellow\":3,\"homeCorner\":5,\"awayCorner\":0,\"hasLineup\":\"\",\"matchTime\":\"2022-12-11 19:00:00\",\"startTime\":\"2022-12-11 20:06:13\",\"explain\":\"\",\"extraExplain\":\"\",\"injuryTime\":\"\"},{\"matchId\":2318072,\"state\":3,\"homeScore\":0,\"awayScore\":1,\"homeHalfScore\":0,\"awayHalfScore\":0,\"homeRed\":0,\"awayRed\":0,\"homeYellow\":0,\"awayYellow\":1,\"homeCorner\":1,\"awayCorner\":4,\"hasLineup\":\"\",\"matchTime\":\"2022-12-11 19:00:00\",\"startTime\":\"2022-12-11 20:06:42\",\"explain\":\"\",\"extraExplain\":\"\",\"injuryTime\":\"\"},{\"matchId\":2318119,\"state\":3,\"homeScore\":3,\"awayScore\":1,\"homeHalfScore\":2,\"awayHalfScore\":0,\"homeRed\":0,\"awayRed\":0,\"homeYellow\":4,\"awayYellow\":1,\"homeCorner\":4,\"awayCorner\":8,\"hasLineup\":\"\",\"matchTime\":\"2022-12-11 19:00:00\",\"startTime\":\"2022-12-11 20:00:15\",\"explain\":\"\",\"extraExplain\":\"\",\"injuryTime\":\"\"}]}"
	model.UpdateScore(str)
}

func oddsChange() {
	for {
		time.Sleep(time.Duration(3) * time.Second)
		res, err := http.Get("http://api.wuhaicj.com/api/liveScore/oddsChange")
		if err != nil {
			fmt.Println("URL Request failed:", err)
			break
		}
		fmt.Println(res)
	}
}
