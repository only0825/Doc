package task

import (
	"encoding/json"
	"time"
	"zinx_ws/model"
)

// 足球比分
func Score() ([]byte, error) {
	cache := model.Rdbc
	// 用阻塞方式弹出数据
	result, err := cache.BLPop(ctx, time.Duration(20)*time.Second, "scoreChange:Football").Result()
	if err != nil {
		return nil, err
	}

	respMap := make(map[int]interface{})
	arr := []ScoreChange{}
	json.Unmarshal([]byte(result[1]), &arr) // result[0]是列表名 result[1]是值
	respMap[0] = "score"
	respMap[1] = arr
	res, _ := json.Marshal(respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}
