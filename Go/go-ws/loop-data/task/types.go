package task

import (
	"context"
)

var ctx = context.Background()

type ScoreChangeList struct {
	ChangeList []ChangeList `json:"changeList"`
}

type ChangeList struct {
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

type ChangeListJson struct {
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

type Odds struct {
	List []struct {
		Handicap      [][]interface{} `json:"handicap"`
		EuropeOdds    [][]interface{} `json:"europeOdds"`
		OverUnder     [][]interface{} `json:"overUnder"`
		HandicapHalf  [][]interface{} `json:"handicapHalf"`
		OverUnderHalf [][]interface{} `json:"overUnderHalf"`
	} `json:"list"`
}

type OddsChange struct {
	ChangeList []struct {
		Handicap      [][]interface{} `json:"handicap"`
		EuropeOdds    [][]interface{} `json:"europeOdds"`
		OverUnder     [][]interface{} `json:"overUnder"`
		HandicapHalf  [][]interface{} `json:"handicapHalf"`
		OverUnderHalf [][]interface{} `json:"overUnderHalf"`
	} `json:"changeList"`
}

type EuropeOddsArr struct {
	EuropeOdds []EuropeOdds
}

// 欧赔（胜平负） 全量
type EuropeOdds struct {
	MatchId          int     `json:"match_id"`            // 比赛ID
	CompanyId        int     `json:"company_id"`          // 公司ID
	HomeWinEarlyOdds float64 `json:"home_win_early_odds"` // 初盘主胜赔率
	TieEarlyOdds     float64 `json:"tie_early_odds"`      // 初盘和局赔率
	AwayWinEarlyOdds float64 `json:"away_win_early_odds"` // 初盘客胜赔率
	HomeWinMainOdds  float64 `json:"home_win_main_odds"`  // 即时盘主胜赔率
	TieMainOdds      float64 `json:"tie_main_odds"`       // 即时盘和局赔率
	AwayWinMainOdds  float64 `json:"away_win_main_odds"`  // 即时盘客胜赔率
	ChangeTime       string  `json:"change_time"`         // 变盘时间
	IsClose          bool    `json:"is_close"`            // 是否封盘 临时性封盘或停止走地。
	OddsType         int     `json:"odds_type"`           // 0无类型数据 1早餐盘 2赛前即时盘 3走地盘
}
type EuropeOddsChangeArr struct {
	EuropeOddsChange []EuropeOddsChange
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

type OverUnderArr struct {
	OverUnder []OverUnder
}

// 大小球数据 全量
type OverUnder struct {
	MatchId            int     `json:"match_id"`              // 比赛ID
	CompanyId          int     `json:"company_id"`            // 公司ID
	HandicapEarlyOdds  float64 `json:"handicap_early_odds"`   // 初盘盘口赔率
	BigBallEarlyOdds   float64 `json:"big_ball_early_odds"`   // 初盘大球赔率
	SmallBallEarlyOdds float64 `json:"small_ball_early_odds"` // 初盘小球赔率
	HandicapOdds       float64 `json:"handicap_odds"`         // 即时盘盘口赔率
	BigBallOdds        float64 `json:"big_ball_odds"`         // 即时盘大球赔率
	SmallBallOdds      float64 `json:"small_ball_odds"`       // 即时盘小球赔率
	ChangeTime         string  `json:"change_time"`           // 变盘时间
	IsClose            bool    `json:"is_close"`              // 是否封盘 临时性封盘或停止走地。
	OddsType           int     `json:"odds_type"`             // 0无类型数据 1早餐盘 2赛前即时盘 3走地盘
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

type RespOddsChange struct {
	EuropeOdds []EuropeOddsChange `json:"europe_odds"` // 欧赔（胜平负） 变量
	OverUnder  []OverUnderChange  `json:"over_under"`  // 大小球数据 变量
}