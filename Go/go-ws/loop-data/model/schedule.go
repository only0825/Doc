package model

// 更新足球赛程表
func UpdateScore1(sc Schedule1) error {
	var table = "hn_schedule"
	err := DB.Table(table).Where("matchId = ?", sc.MatchId).Updates(sc).Error
	if err != nil {
		return err
	}
	return nil
}

// 更新篮球赛程表
func UpdateScore2(sc Schedule2, matchId int) error {
	var table = "hn_basketball_schedule"
	err := DB.Table(table).Where("matchId = ?", matchId).Updates(sc).Error
	if err != nil {
		return err
	}
	return nil
}
