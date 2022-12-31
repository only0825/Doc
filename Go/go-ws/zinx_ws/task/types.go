package task

// 足球比分 变量
type ScoreChange struct {
	MatchId       int    `json:"match_id"`        // 比赛ID
	State         int    `json:"state"`           // 比赛
	HomeScore     int    `json:"home_score"`      // 主队得分
	AwayScore     int    `json:"away_score"`      // 客队得分
	HomeHalfScore int    `json:"home_half_score"` // 主队上半场得分
	AwayHalfScore int    `json:"away_half_score"` // 客队上半场得分
	HomeRed       int    `json:"home_red"`        // 主队红牌数
	AwayRed       int    `json:"away_red"`        // 客队红牌数
	HomeYellow    int    `json:"home_yellow"`     // 主队黄牌数
	AwayYellow    int    `json:"away_yellow"`     // 客队红牌数
	HomeCorner    int    `json:"home_corner"`     // 主队角球数
	AwayCorner    int    `json:"away_corner"`     // 客队角球数
	HasLineup     string `json:"has_lineup"`      // 是否有阵容
	MatchTime     string `json:"match_time"`      // 比赛时间
	StartTime     string `json:"start_time"`      // 开场时间
	Explain       string `json:"explain"`         // 比赛说明
	ExtraExplain  string `json:"extra_explain"`   // 比赛说明2
	InjuryTime    string `json:"injury_time"`     // 上下半场补时时长
}

// 欧赔（胜平负） 变量
type EuropeOddsChange struct {
	MatchId         int     `json:"match_id"`           // 比赛ID
	CompanyId       int     `json:"company_id"`         // 公司ID
	HomeWinMainOdds float64 `json:"home_win_main_odds"` // 即时盘主胜赔率
	TieMainOdds     float64 `json:"tie_main_odds"`      // 即时盘和局赔率
	AwayWinMainOdds float64 `json:"away_win_main_odds"` // 即时盘客胜赔率
	ChangeTime      string  `json:"change_time"`        // 变盘时间
	IsClose         bool    `json:"is_close"`           // 是否封盘 临时性封盘或停止走地。
	OddsType        int     `json:"odds_type"`          // 0无类型数据 1早餐盘 2赛前即时盘 3走地盘
}

// 大小球数据 变量
type OverUnderChange struct {
	MatchId       int     `json:"match_id"`        // 比赛ID
	CompanyId     int     `json:"company_id"`      // 公司ID
	HandicapOdds  float64 `json:"handicap_odds"`   // 即时盘盘口赔率
	BigBallOdds   float64 `json:"big_ball_odds"`   // 即时盘大球赔率
	SmallBallOdds float64 `json:"small_ball_odds"` // 即时盘小球赔率
	ChangeTime    string  `json:"change_time"`     // 变盘时间
	IsClose       bool    `json:"is_close"`        // 是否封盘 临时性封盘或停止走地。
	OddsType      int     `json:"odds_type"`       // 0无类型数据 1早餐盘 2赛前即时盘 3走地盘
}

type respOddsChange struct {
	EuropeOdds []EuropeOddsChange `json:"europe_odds"` // 欧赔（胜平负） 变量
	OverUnder  []OverUnderChange  `json:"over_under"`  // 大小球数据 变量
}
