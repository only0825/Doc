package task

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"loop-data/configs"
	"loop-data/model"
	"net/http"
	"strconv"
	"time"
)

type StatsBasketball struct {
}

func (this StatsBasketball) Run() {
	stats(configs.Conf.ApiB.Stats, "")
}

func stats(url string, scType string) {
	var cache = model.Rdb

	resp, err := http.Get(url)
	if err != nil {
		logrus.Error("请求失败:", err.Error())
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("io.ReadAll失败:", err.Error())
		return
	}

	// 判断获取到的数据是否为空
	if len(body) < 50 {
		//logrus.Info("changeList value empty")
		return
	}

	var mla = MatchListArr{}
	err = json.Unmarshal(body, &mla)
	if err != nil {
		logrus.Error("json 反序列化错误", err)
		return
	}

	for i := range mla.MatchList {
		matchId := mla.MatchList[i].MatchID
		marshal, _ := json.Marshal(mla.MatchList[i])
		length := strconv.Itoa(len(marshal))
		value := "s:" + length + ":\"" + string(marshal) + "\";" // PHP存取的时候必须要是这个结构
		cache.Set(ctx, "basketball:stats:"+strconv.Itoa(matchId), value, time.Duration(604800)*time.Second)
	}

	logrus.Info("篮球技术统计 Redis 存储成功！")
}

type MatchListArr struct {
	MatchList []MatchList `json:"matchList"`
}
type HomePlayerList struct {
	PlayerID         int    `json:"playerId"`
	PlayerEn         string `json:"playerEn"`
	PlayerChs        string `json:"playerChs"`
	PlayerCht        string `json:"playerCht"`
	PositionEn       string `json:"positionEn"`
	PositionCn       string `json:"positionCn"`
	Playtime         string `json:"playtime"`
	ShootHit         string `json:"shootHit"`
	Shoot            string `json:"shoot"`
	ThreePointHit    string `json:"threePointHit"`
	ThreePointShoot  string `json:"threePointShoot"`
	FreeThrowHit     string `json:"freeThrowHit"`
	FreeThrowShoot   string `json:"freeThrowShoot"`
	OffensiveRebound string `json:"offensiveRebound"`
	DefensiveRebound string `json:"defensiveRebound"`
	Assist           string `json:"assist"`
	Foul             string `json:"foul"`
	Steal            string `json:"steal"`
	Turnover         string `json:"turnover"`
	Block            string `json:"block"`
	Score            string `json:"score"`
	IsOnFloor        bool   `json:"isOnFloor"`
}
type AwayPlayerList struct {
	PlayerID         int    `json:"playerId"`
	PlayerEn         string `json:"playerEn"`
	PlayerChs        string `json:"playerChs"`
	PlayerCht        string `json:"playerCht"`
	PositionEn       string `json:"positionEn"`
	PositionCn       string `json:"positionCn"`
	Playtime         string `json:"playtime"`
	ShootHit         string `json:"shootHit"`
	Shoot            string `json:"shoot"`
	ThreePointHit    string `json:"threePointHit"`
	ThreePointShoot  string `json:"threePointShoot"`
	FreeThrowHit     string `json:"freeThrowHit"`
	FreeThrowShoot   string `json:"freeThrowShoot"`
	OffensiveRebound string `json:"offensiveRebound"`
	DefensiveRebound string `json:"defensiveRebound"`
	Assist           string `json:"assist"`
	Foul             string `json:"foul"`
	Steal            string `json:"steal"`
	Turnover         string `json:"turnover"`
	Block            string `json:"block"`
	Score            string `json:"score"`
	IsOnFloor        bool   `json:"isOnFloor"`
}
type MatchList struct {
	MatchID        int              `json:"matchId"`
	HomeTeamEn     string           `json:"homeTeamEn"`
	HomeTeamCn     string           `json:"homeTeamCn"`
	AwayTeamEn     string           `json:"awayTeamEn"`
	AwayTeamCn     string           `json:"awayTeamCn"`
	CostTime       string           `json:"costTime"`
	HomeScore      string           `json:"homeScore"`
	HomeFast       string           `json:"homeFast"`
	HomeInside     string           `json:"homeInside"`
	HomeExceed     string           `json:"homeExceed"`
	HomeTotalmis   string           `json:"homeTotalmis"`
	AwayScore      string           `json:"awayScore"`
	AwayFast       string           `json:"awayFast"`
	AwayInside     string           `json:"awayInside"`
	AwayExceed     string           `json:"awayExceed"`
	AwayTotalmis   string           `json:"awayTotalmis"`
	HomePlayerList []HomePlayerList `json:"homePlayerList"`
	AwayPlayerList []AwayPlayerList `json:"awayPlayerList"`
}
