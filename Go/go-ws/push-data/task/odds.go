package task

import (
	"encoding/json"
	"push-data/model"
	"time"
)

func Odds() ([]byte, error) {
	cache := model.Rdb
	// 用阻塞方式弹出数据
	result, err := cache.BLPop(ctx, time.Duration(30)*time.Second, "oddsChange:Football").Result()
	if err != nil {
		return nil, err
	}

	respMap := make(map[int]interface{})
	arr := respOddsChange{}
	json.Unmarshal([]byte(result[1]), &arr) // result[0]是列表名 result[1]是值
	respMap[0] = "odds"
	respMap[1] = arr
	res, _ := json.Marshal(respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}
