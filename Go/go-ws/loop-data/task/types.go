package task

import (
	"context"
)

var ctx = context.Background()

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
