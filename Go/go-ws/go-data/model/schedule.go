package model

import (
	"fmt"
	"time"
)

func UpdateScore(changes []Change) {
	var table = "hn_schedule_copy1"
	for i := range changes {
		var c = changes[i]
		var sc = Schedule{
			MatchId:       c.MatchId,
			State:         c.State,
			HomeScore:     c.HomeScore,
			AwayScore:     c.AwayScore,
			HomeHalfScore: c.HomeHalfScore,
			AwayHalfScore: c.AwayHalfScore,
			HomeRed:       c.HomeRed,
			AwayRed:       c.AwayRed,
			HomeYellow:    c.HomeYellow,
			AwayYellow:    c.AwayYellow,
			HomeCorner:    c.HomeCorner,
			AwayCorner:    c.AwayCorner,
			UpdateTime:    time.Now().Format("2006/01/02 15:04:05"),
		}

		err := DB.Table(table).Where("matchId = ?", sc.MatchId).Updates(sc).Error
		if err != nil {
			fmt.Println("数据库更新分数错误：", err)
		}
	}
}
