package model

import (
	"encoding/json"
	"fmt"
	"time"
)

var table = "hn_schedule_copy1"

type Schedule struct {
	MatchId       int    `json:"matchId" gorm:"column:matchId"`             // 比赛ID
	State         int    `json:"state" gorm:"column:state"`                 // 比赛
	HomeScore     int    `json:"homeScore" gorm:"column:homeScore"`         // 主队得分
	AwayScore     int    `json:"awayScore" gorm:"column:awayScore"`         // 客队得分
	HomeHalfScore int    `json:"homeHalfScore" gorm:"column:homeHalfScore"` // 主队上半场得分
	AwayHalfScore int    `json:"awayHalfScore" gorm:"column:awayHalfScore"` // 客队上半场得分
	HomeRed       int    `json:"homeRed" gorm:"column:homeRed"`             // 主队红牌数
	AwayRed       int    `json:"awayRed" gorm:"column:awayRed"`             // 客队红牌数
	HomeYellow    int    `json:"homeYellow" gorm:"column:homeYellow"`       // 主队黄牌数
	AwayYellow    int    `json:"awayYellow" gorm:"column:awayYellow"`       // 客队红牌数
	HomeCorner    int    `json:"homeCorner" gorm:"column:homeCorner"`       // 主队角球数
	AwayCorner    int    `json:"awayCorner" gorm:"column:awayCorner"`       // 客队角球数
	UpdateTime    string `json:"updateTime" gorm:"column:updateTime"`       // 客队角球数
}

func scheduleUpdate(sc Schedule) (err error) {
	err = DB.Table(table).Where("matchId = ?", sc.MatchId).Updates(sc).Error
	return err
}

func UpdateScore(msg string) {
	var changeList ChangeList
	// 将json消息转为数组
	err := json.Unmarshal([]byte(msg), &changeList)
	if err != nil {
		fmt.Println("JSON转换失败：", err)
	}

	var changes []Change
	changes = changeList.ChangeList
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
		err = scheduleUpdate(sc)
		if err != nil {
			fmt.Println("数据库更新分数错误：", err)
		}
	}

}
