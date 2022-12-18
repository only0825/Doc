package model

type ChangeList struct {
	ChangeList []Change
}

type Change struct {
	MatchId       int    `json:"matchId"`       // 比赛ID
	State         int    `json:"state"`         // 比赛
	HomeScore     string `json:"homeScore"`     // 主队得分
	AwayScore     string `json:"awayScore"`     // 客队得分
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
