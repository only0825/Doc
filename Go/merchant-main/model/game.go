package model

import (
	"errors"
	"fmt"
	"merchant/contrib/helper"
	"strconv"
	"strings"

	g "github.com/doug-martin/goqu/v9"
)

func GameFind(id string) (GameLists, error) {

	var game GameLists
	query, _, _ := dialect.From("tbl_game_lists").Select(colsGameList...).Where(g.Ex{"id": id, "prefix": meta.Prefix}).ToSQL()
	err := meta.MerchantDB.Get(&game, query)
	if err != nil {
		return game, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	return game, nil
}

// 更新单条游戏数据
func GameListUpdate(id, pid string, record g.Record) error {

	query, _, _ := dialect.Update("tbl_game_lists").Set(record).Where(g.Ex{"id": id, "prefix": meta.Prefix}).ToSQL()
	_, err := meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	return LoadGameLists(pid)
}

// 游戏列表 分页查询
func GameList(ex g.Ex, page, pageSize uint) (GamePageList, error) {

	data := GamePageList{}
	ex["prefix"] = meta.Prefix
	t := dialect.From("tbl_game_lists")
	if page == 1 {
		query, _, _ := t.Select(g.COUNT(1)).Where(ex).ToSQL()
		err := meta.MerchantDB.Get(&data.T, query)
		if err != nil {
			return data, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
		}
	}

	offset := pageSize * (page - 1)
	query, _, _ := t.Select(colsGameList...).Where(ex).Order(g.I("sorting").Desc()).Order(g.I("id").Desc()).Offset(offset).Limit(pageSize).ToSQL()
	err := meta.MerchantDB.Select(&data.D, query)
	if err != nil {
		return data, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	// 处理client_type 给前端展示
	for key := range data.D {
		data.D[key].ClientType = strings.Join(gameIntClientToString(data.D[key].ClientType), ",")
	}

	data.S = pageSize

	return data, nil
}

func gameIntClientToString(client string) []string {

	clientType, err := strconv.ParseUint(client, 10, 64)
	if err != nil {
		return nil
	}

	if clientType == 0 {
		return []string{"1", "2", "4"}
	}

	var data []string
	clientString := []string{"", "1", "2", "4"}

	binary := strconv.FormatUint(clientType, 2)
	maxClient := len(binary)
	for key, v := range binary { // 0 = 48, 1 = 49

		if v == 49 {
			data = append(data, clientString[maxClient-key])
		}
	}

	return data
}

func LoadGameLists(pid ...string) error {

	var data []platJson
	ex := g.Ex{
		"state":     1,
		"game_type": []int64{3, 4, 7},
		"prefix":    meta.Prefix,
	}

	if len(pid) > 0 {
		if len(pid) == 1 {
			ex["id"] = pid[0]
		} else {
			ex["id"] = pid
		}
	}
	query, _, _ := dialect.From("tbl_platforms").Select(colsPlatJson...).Where(ex).Order(g.C("created_at").Asc()).ToSQL()
	query = "/* master */ " + query
	fmt.Println(query)
	err := meta.MerchantDB.Select(&data, query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	if len(data) == 0 {
		return errors.New(helper.RecordNotExistErr)
	}

	pipe := meta.MerchantRedis.TxPipeline()
	defer pipe.Close()

	for _, v := range data {

		var sg []showGameJson
		exG := g.Ex{
			"online":      GameOnline,
			"platform_id": v.ID,
			"prefix":      meta.Prefix,
		}
		query1, _, _ := dialect.From("tbl_game_lists").Select(colsShowGame...).Where(exG).Order(g.C("sorting").Asc()).ToSQL()
		query = "/* master */ " + query
		fmt.Println(query)
		err = meta.MerchantDB.Select(&sg, query1)
		if err != nil {
			return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
		}

		if len(sg) == 0 {
			continue
		}

		key := fmt.Sprintf("%s:game:%s", meta.Prefix, v.ID)
		hotKey := fmt.Sprintf("%s:game:h:%s", meta.Prefix, v.ID)
		newKey := fmt.Sprintf("%s:game:n:%s", meta.Prefix, v.ID)
		pipe.Unlink(ctx, key)
		pipe.Unlink(ctx, hotKey)
		pipe.Unlink(ctx, newKey)
		for _, val := range sg {
			b, err := helper.JsonMarshal(val)
			if err == nil {
				pipe.RPush(ctx, key, string(b))
				if val.IsNew == 1 {
					pipe.RPush(ctx, newKey, string(b))
				}
				if val.IsHot == 1 {
					pipe.RPush(ctx, hotKey, string(b))
				}
			}
		}
		pipe.Persist(ctx, key)
		pipe.Persist(ctx, hotKey)
		pipe.Persist(ctx, newKey)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return pushLog(err, helper.RedisErr)
	}

	return nil
}
