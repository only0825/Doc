package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-data/common"
	"go-data/model"
	"go-data/zlog"
	"io"
	"net/http"
	"time"
)

type OddsFootball struct {
}

func (this OddsFootball) Run() {
	zlog.Info.Println("TaskScoreFootball start")
	oddsChange("http://api.wuhaicj.com/api/liveScore/oddsChange", common.Football)
}

func oddsChange(url string, category int) {
	// 一、请求第三方接口
	res, err := http.Get(url)
	if err != nil {
		zlog.Error.Println("URL Request failed:", err)
		return
	}

	msg, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		zlog.Error.Println("Read body failed:", err)
		return
	}

	// 判断获取到的数据是否为空
	if len(msg) < 50 {
		zlog.Info.Println("changeList value empty")
		return
	}

	var obj OddsChange
	if err = json.Unmarshal(msg, &obj); err != nil {
		zlog.Error.Println("json反序列化失败: ", err)
		return
	}

	europeOddsByte, _ := json.Marshal(obj.ChangeList[0].EuropeOdds)
	saveToRedis(europeOddsByte, "europeOdds")

	overUnderByte, _ := json.Marshal(obj.ChangeList[0].OverUnder)
	saveToRedis(overUnderByte, "overUnder")
}

func saveToRedis(data []byte, changeType string) error {
	// 二、存Redis （给推送服务用）  一分钟内相同的数据不写入scoreChange中
	isHave, _ := model.Rdbc.Get(ctx, fmt.Sprintf("oddsChangeTemp:football:%s", changeType)).Result()
	if (isHave != "") && (isHave == common.Md5String(string(data))) {
		zlog.Info.Println("重复数据")
		return errors.New("重复数据")
	} else {
		// 开启Redis事务
		pipe := model.Rdbc.TxPipeline()
		// 临时存放去重
		pipe.Set(ctx, fmt.Sprintf("oddsChangeTemp:football:%s", changeType), common.Md5String(string(data)), time.Duration(60)*time.Second)
		// 将获取到的数据存入到Redis队列
		pipe.LPush(ctx, fmt.Sprintf("oddsChange:football:%s", changeType), data)
		_, err := pipe.Exec(ctx)
		if err != nil {
			zlog.Error.Println("Redis 事务报错:", err)
			return err
		}
		zlog.Info.Println("Redis存储成功")
		return nil
	}
}
