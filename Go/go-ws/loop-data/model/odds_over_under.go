package model

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

var table2 = "hn_odds_over_under"

// 根据比赛ID查询数据是否存在
func OuFind(matchId int) (bool, error) {
	oe := OverUnder{}
	db := DB.Table(table2).Where("match_id = ?", matchId).Find(&oe)
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

func OuUpdate(oe OverUnder) error {
	err := DB.Table(table2).Where("match_id = ?", oe.MatchId).Save(oe).Error
	if err != nil {
		msg := fmt.Sprintf("%s 表更新失败，match_id = %d, %s", table1, oe.MatchId, err)
		logrus.Error(msg)
		return err
	}
	return nil
}

func OuAdd(oe OverUnder) error {
	err := DB.Table(table2).Create(&oe).Error
	if err != nil {
		msg := fmt.Sprintf("%s 表添加失败，match_id = %d, %s", table1, oe.MatchId, err)
		logrus.Error(msg)
		return err
	}
	return nil
}
