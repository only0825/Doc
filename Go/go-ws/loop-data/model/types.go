package model

type ChangeList struct {
	ChangeList []Change
}

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

// 欧赔（胜平负）数据
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
}

// 大小球数据
type OverUnder struct {
	MatchId            int     // 比赛ID
	CompanyId          int     // 公司ID
	HandicapEarlyOdds  float64 // 初盘盘口赔率
	BigBallEarlyOdds   float64 // 初盘大球赔率
	SmallBallEarlyOdds float64 // 初盘小球赔率
	HandicapOdds       float64 // 即时盘盘口赔率
	BigBallOdds        float64 // 即时盘大球赔率
	SmallBallOdds      float64 // 即时盘小球赔率
	ChangeTime         string  // 变盘时间
	IsClose            int     // 是否封盘 临时性封盘或停止走地。
	OddsType           int     // 0无类型数据 1早餐盘 2赛前即时盘 3走地盘
}
