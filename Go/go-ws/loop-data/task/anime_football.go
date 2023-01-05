package task

import (
	"loop-data/configs"
)

type AnimeFootball struct {
}

func (this AnimeFootball) Run() {
	anime1(configs.Conf.ApiB.Stats, "")
}

func anime1(url string, scType string) {

}
