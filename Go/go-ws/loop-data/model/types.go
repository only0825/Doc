package model

// schedule 足球表
type Schedule1 struct {
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
	UpdateTime    string `json:"updateTime" gorm:"column:updateTime"`       // 数据更新时间
}

// schedule 篮球表
type Schedule2 struct {
	//MatchId    int    `json:"matchId" gorm:"column:matchId"`       // 比赛ID
	RemainTime string `json:"remainTime" gorm:"column:remainTime"` // 小节剩余时间
	State      int    `json:"state" gorm:"column:state"`           // 比赛状态 0:未开赛 1:一节 2:二节 3:三节 4:四节  5:1OT  6:2OT  7:3OT  50:中场  -1:完场  -2:待定  -3:中断  -4:取消  -5:推迟（1OT是第一次加时赛）
	HomeScore  int    `json:"homeScore" gorm:"column:homeScore"`   // 主队总分
	Home1      int    `json:"home1" gorm:"column:home1"`           // 主队第一节得分
	Home2      int    `json:"home2" gorm:"column:home2"`           // 主队第二节得分
	Home3      int    `json:"home3" gorm:"column:home3"`           // 主队第三节得分
	Home4      int    `json:"home4" gorm:"column:home4"`           // 主队第四节得分
	HomeOT1    int    `json:"homeOT1" gorm:"column:homeOT1"`       // 主队第一加时得分
	HomeOT2    int    `json:"homeOT2" gorm:"column:homeOT2"`       // 主队第二加时得分
	HomeOT3    int    `json:"homeOT3" gorm:"column:homeOT3"`       // 主队第三加时得分
	AwayScore  int    `json:"awayScore" gorm:"column:awayScore"`   // 客队总分
	Away1      int    `json:"away1" gorm:"column:away1"`           // 主队第一节得分
	Away2      int    `json:"away2" gorm:"column:away2"`           // 主队第二节得分
	Away3      int    `json:"away3" gorm:"column:away3"`           // 主队第三节得分
	Away4      int    `json:"away4" gorm:"column:away4"`           // 主队第四节得分
	AwayOT1    int    `json:"awayOT1" gorm:"column:awayOT1"`       // 主队第一加时得分
	AwayOT2    int    `json:"awayOT2" gorm:"column:awayOT2"`       // 主队第二加时得分
	AwayOT3    int    `json:"awayOT3" gorm:"column:awayOT3"`       // 主队第三加时得分
	HasStats   int    `json:"hasStats" gorm:"column:hasStats"`     // 是否有技术统计1是 0否
	explainCn  string `json:"explainCn" gorm:"column:explainCn"`   // 比赛说明、备注-中文
	UpdateTime string `json:"updateTime" gorm:"column:updateTime"` // 数据更新时间
}

// hn_odds_europe表 欧赔（胜平负）
type EuropeOdds struct {
	MatchId          int     `json:"match_id" gorm:"column:match_id"`                       // 比赛ID
	CompanyId        int     `json:"company_id" gorm:"column:company_id"`                   // 公司ID
	HomeWinEarlyOdds float64 `json:"home_win_early_odds" gorm:"column:home_win_early_odds"` // 初盘主胜赔率
	TieEarlyOdds     float64 `json:"tie_early_odds" gorm:"column:tie_early_odds"`           // 初盘和局赔率
	AwayWinEarlyOdds float64 `json:"away_win_early_odds" gorm:"column:away_win_early_odds"` // 初盘客胜赔率
	HomeWinMainOdds  float64 `json:"home_win_main_odds" gorm:"column:home_win_main_odds"`   // 即时盘主胜赔率
	TieMainOdds      float64 `json:"tie_main_odds" gorm:"column:tie_main_odds"`             // 即时盘和局赔率
	AwayWinMainOdds  float64 `json:"away_win_main_odds" gorm:"column:away_win_main_odds"`   // 即时盘客胜赔率
	ChangeTime       string  `json:"change_time" gorm:"column:change_time"`                 // 变盘时间
	IsClose          int     `json:"is_close" gorm:"column:is_close"`                       // 是否封盘 临时性封盘或停止走地。
	OddsType         int     `json:"odds_type" gorm:"column:odds_type"`                     // 0无类型数据 1早餐盘 2赛前即时盘 3走地盘
	UpdateTime       string  `json:"update_time" gorm:"column:update_time"`                 // 数据更新时间
}

// hn_odds_over_under表 大小球
type OverUnder struct {
	//MatchId            int     `json:"match_id" gorm:"column:match_id"`                           // 比赛ID
	CompanyId          int     `json:"company_id" gorm:"column:company_id"`                       // 公司ID
	HandicapEarlyOdds  float64 `json:"handicap_early_odds" gorm:"column:handicap_early_odds"`     // 初盘盘口赔率
	BigBallEarlyOdds   float64 `json:"big_ball_early_odds" gorm:"column:big_ball_early_odds"`     // 初盘大球赔率
	SmallBallEarlyOdds float64 `json:"small_ball_early_odds" gorm:"column:small_ball_early_odds"` // 初盘小球赔率
	HandicapOdds       float64 `json:"handicap_odds" gorm:"column:handicap_odds"`                 // 即时盘盘口赔率
	BigBallOdds        float64 `json:"big_ball_odds" gorm:"column:big_ball_odds"`                 // 即时盘大球赔率
	SmallBallOdds      float64 `json:"small_ball_odds" gorm:"column:small_ball_odds"`             // 即时盘小球赔率
	ChangeTime         string  `json:"change_time" gorm:"column:change_time"`                     // 变盘时间
	IsClose            int     `json:"is_close" gorm:"column:is_close"`                           // 是否封盘 临时性封盘或停止走地。
	OddsType           int     `json:"odds_type" gorm:"column:odds_type"`                         // 0无类型数据 1早餐盘 2赛前即时盘 3走地盘
	UpdateTime         string  `json:"update_time" gorm:"column:update_time"`                     // 数据更新时间
}
