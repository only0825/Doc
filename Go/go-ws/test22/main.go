package main

import (
	"fmt"
	"strconv"
	"time"
)

type bar struct {
	timestamp time.Time
	open      float64
	high      float64
	low       float64
	close     float64
	volume    float64
	cVolume   float64
}

func newBar(s []string) *bar {
	b := &bar{}

	if len(s) > 0 {
		b.timestamp, _ = time.Parse(time.RFC3339, s[0])
		fmt.Println(b.timestamp)
	}

	if len(s) > 1 {
		b.open, _ = strconv.ParseFloat(s[1], 64)
	}

	if len(s) > 2 {
		b.high, _ = strconv.ParseFloat(s[2], 64)
	}

	if len(s) > 3 {
		b.low, _ = strconv.ParseFloat(s[3], 64)
	}

	if len(s) > 4 {
		b.close, _ = strconv.ParseFloat(s[4], 64)
	}

	if len(s) > 5 {
		b.volume, _ = strconv.ParseFloat(s[5], 64)
	}

	if len(s) > 6 {
		b.cVolume, _ = strconv.ParseFloat(s[6], 64)
	}

	return b
}

func main() {
	s := [][]string{
		{
			"2019-03-20T16:00:00.000Z",
			"3.721",
			"3.743",
			"3.677",
			"3.708",
			"8422410",
			"22698348.04828491",
		},
		{
			"2019-03-19T16:00:00.000Z",
			"3.731",
			"3.799",
			"3.494",
			"3.72",
			"24912403",
			"67632347.24399722",
		},
	}

	var bs []*bar
	for _, data := range s {
		fmt.Println(data)
		bs = append(bs, newBar(data))
	}

	fmt.Println(bs[0].high)
	fmt.Printf("%v\n", bs)
}
