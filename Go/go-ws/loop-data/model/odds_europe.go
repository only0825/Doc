package model

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

var table1 = "hn_odds_europe"

// 根据比赛ID查询数据是否存在
func OeFind(matchId int) (bool, error) {
	oe := EuropeOdds{}
	db := DB.Table(table1).Where("match_id = ?", matchId).Find(&oe)
	err := db.Error
	if err != nil {
		msg := fmt.Sprintf("%s 表查询失败，match_id = %d, %s", table1, oe.MatchId, err)
		logrus.Error(msg)
		return false, err
	}

	affected := db.RowsAffected
	if affected == 0 {
		return false, nil
	}

	return true, nil
}

func OeUpdate(oe EuropeOdds, matchId int) error {
	err := DB.Table(table1).Where("match_id = ?", matchId).Save(oe).Error
	if err != nil {
		msg := fmt.Sprintf("%s 表更新失败，match_id = %d, %s", table1, matchId, err)
		logrus.Error(msg)
		return err
	}
	return nil
}

func OeAdd(oe EuropeOdds, matchId int) error {
	err := DB.Table(table1).Create(&oe).Error
	if err != nil {
		msg := fmt.Sprintf("%s 表添加失败，match_id = %d, %s", table1, matchId, err)
		logrus.Error(msg)
		return err
	}
	return nil
}
