package task

import (
	"github.com/sirupsen/logrus"
	"loop-data/configs"
)

type ScoreBasketBall struct {
}

func (this ScoreBasketBall) Run() {
	logrus.Info("篮球比分变量 TaskScoreBasketBall start")
	scoreChange(configs.Conf.ApiB.ScoreChange, "BasketBall")
}
