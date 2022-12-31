package model

import (
	"database/sql"
	"errors"
	"fmt"
	"merchant/contrib/helper"

	g "github.com/doug-martin/goqu/v9"
)

type Platform struct {
	ID              string `db:"id" json:"id"`
	Name            string `db:"name" json:"name"`
	Prefix          string `db:"prefix" json:"prefix"`
	GameType        int    `db:"game_type" json:"game_type"`
	Flags           int    `db:"flags" json:"flags"`
	State           int    `db:"state" json:"state"`
	Maintained      int    `db:"maintained" json:"maintained"`
	MaintainedStart string `db:"maintained_start" json:"maintained_start"`
	MaintainedEnd   string `db:"maintained_end" json:"maintained_end"`
	Seq             int    `db:"seq" json:"seq"`
	CreatedAt       int32  `db:"created_at" json:"created_at"`
	UpdatedAt       int32  `db:"updated_at" json:"updated_at"`
}

type platJson struct {
	ID         string `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	GameType   int    `db:"game_type" json:"game_type"`
	Maintained int    `db:"maintained" json:"maintained"`
	Flags      int    `db:"flags" json:"flags"`
	State      int    `db:"state" json:"state"`
	Seq        int    `db:"seq" json:"seq"`
}

type PlatformData struct {
	T int64      `json:"t"`
	D []Platform `json:"d"`
	S int        `json:"s"`
}

type navJson struct {
	Cate
	L []platJson `json:"l"`
}

func PlatformUpdate(ex g.Ex, record g.Record) error {

	query, _, _ := dialect.Update("tbl_platforms").Set(record).Where(ex).Order(g.C("created_at").Desc()).ToSQL()
	_, err := meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	return LoadPlatforms()
}

func PlatformFindOne(ex g.Ex) (Platform, error) {

	data := Platform{}
	ex["prefix"] = meta.Prefix
	query, _, _ := dialect.From("tbl_platforms").Select(colsPlatform...).Where(ex).Limit(1).ToSQL()
	err := meta.MerchantDB.Get(&data, query)
	if err != nil && err != sql.ErrNoRows {
		return data, pushLog(err, helper.DBErr)
	}

	if err == sql.ErrNoRows {
		return data, pushLog(err, helper.RecordNotExistErr)
	}

	return data, nil
}

func PlatformList(ex g.Ex, page, pageSize int) (PlatformData, error) {

	data := PlatformData{
		S: pageSize,
	}
	ex["prefix"] = meta.Prefix
	offset := (page - 1) * pageSize
	t := dialect.From("tbl_platforms")
	if page == 1 {
		query, _, _ := t.Select(g.COUNT("id")).Where(ex).ToSQL()
		err := meta.MerchantDB.Get(&data.T, query)
		if err != nil {
			return data, pushLog(err, helper.DBErr)
		}

		if data.T == 0 {
			return data, nil
		}

	}

	query, _, _ := t.Select(colsPlatform...).Where(ex).Order(g.C("created_at").Desc()).Offset(uint(offset)).Limit(uint(pageSize)).ToSQL()
	err := meta.MerchantDB.Select(&data.D, query)
	if err != nil {
		return data, pushLog(err, helper.DBErr)
	}

	return data, nil
}

func LoadPlatforms() error {

	var data []platJson

	ex := g.Ex{
		"state":  1,
		"prefix": meta.Prefix,
	}
	query, _, _ := dialect.From("tbl_platforms").
		Select(colsPlatJson...).Where(ex).Order(g.C("seq").Asc()).ToSQL()
	query = "/* master */ " + query
	fmt.Println(query)
	err := meta.MerchantDB.Select(&data, query)
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	if len(data) == 0 {
		return errors.New(helper.RecordNotExistErr)
	}

	b, err := helper.JsonMarshal(data)
	if err != nil {
		return errors.New(helper.FormatErr)
	}

	navJ, err := NavMinio()
	if err != nil {
		return err
	}

	navB, err := helper.JsonMarshal(navJ)
	if err != nil {
		return errors.New(helper.FormatErr)
	}

	//fmt.Println(string(navB))
	//fmt.Println(string(b))

	pipe := meta.MerchantRedis.TxPipeline()
	defer pipe.Close()

	for _, val := range data {
		k := fmt.Sprintf("%s:plat:%s", meta.Prefix, val.ID)
		b1, err := helper.JsonMarshal(val)
		if err != nil {
			_ = pushLog(err, helper.FormatErr)
			continue
		}

		fmt.Println(k, string(b1))
		pipe.Unlink(ctx, k)
		pipe.Set(ctx, k, string(b1), 0)
		pipe.Persist(ctx, k)
	}

	navKey := fmt.Sprintf("%s:nav", meta.Prefix)
	platKey := fmt.Sprintf("%s:plat", meta.Prefix)
	pipe.Unlink(ctx, navKey)
	pipe.Unlink(ctx, platKey)
	pipe.Set(ctx, navKey, string(navB), 0)
	pipe.Persist(ctx, navKey)
	pipe.Set(ctx, platKey, string(b), 0)
	pipe.Persist(ctx, platKey)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return pushLog(err, helper.RedisErr)
	}

	return nil
}

func NavMinio() ([]navJson, error) {

	var top []Cate
	query, _, _ := dialect.From("tbl_tree").
		Where(g.C("level").ILike("0010%"), g.C("prefix").Eq(meta.Prefix)).Order(g.C("sort").Asc()).ToSQL()
	fmt.Println(query)
	err := meta.MerchantDB.Select(&top, query)
	if err != nil {
		return nil, pushLog(err, helper.DBErr)
	}

	topLen := len(top)
	if topLen == 0 {
		fmt.Println("NavMinio query = ", query)
		return nil, errors.New(helper.RecordNotExistErr)
	}

	data := make([]navJson, topLen)
	for k, v := range top {

		data[k].Cate = v
		ex := g.Ex{
			"state":     1,
			"game_type": v.ID,
			"prefix":    meta.Prefix,
		}
		query, _, _ = dialect.From("tbl_platforms").Select(colsPlatJson...).Where(ex).Order(g.C("seq").Asc()).ToSQL()
		fmt.Println(query)
		err = meta.MerchantDB.Select(&data[k].L, query)
		if err != nil {
			_ = pushLog(err, helper.DBErr)
			continue
		}
	}

	return data, nil
}
