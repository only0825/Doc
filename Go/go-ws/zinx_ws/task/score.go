package task

import (
	"encoding/json"
	"time"
	"zinx_ws/model"
	"zinx_ws/zlog"
)

// 足球比分
func Score() ([]byte, error) {
	cache := model.Rdb
	// 用阻塞方式弹出数据
	result, err := cache.BLPop(ctx, time.Duration(20)*time.Second, "scoreChange:Football").Result()
	if err != nil {
		return nil, err
	}

	zlog.Info.Println(result[1])

	map1 := make(map[string]interface{})
	map2 := make(map[int]interface{})
	json.Unmarshal([]byte(result[1]), &map1) // result[0]是列表名 result[1]是值
	map2[0] = "score"
	map2[1] = map1["changeList"]
	res, _ := json.Marshal(map2)
	if err != nil {
		return nil, err
	}

	return res, nil
}
