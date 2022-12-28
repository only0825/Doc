package task

import (
	"go-data/configs"
	"go-data/zlog"
)

type ScoreBasketBall struct {
}

func (this ScoreBasketBall) Run() {
	zlog.Info.Println("篮球比分变量 TaskScoreBasketBall start")
	scoreChange(configs.Conf.ApiB.ScoreChange, "BasketBall")
}
