package model

import (
	"database/sql"
	"errors"
	"fmt"
	"merchant/contrib/helper"
	"strings"
	"time"

	g "github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/valyala/fastjson"
)

func BannerUpdateState(id string, state uint8) error {

	banner := Banner{}
	ex := g.Ex{
		"id":     id,
		"prefix": meta.Prefix,
	}
	query, _, _ := dialect.Select(colsBanner...).From("tbl_banner").Where(ex).ToSQL()
	err := meta.MerchantDB.Get(&banner, query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	if banner.State == state || (banner.State == BannerStateProcessing && state == BannerStateWait) {
		return errors.New(helper.StateParamErr)
	}

	data := map[uint8]uint8{
		1: BannerStateProcessing,
		3: BannerStateProcessing,
		2: BannerStateEnd,
	}

	status, ok := data[banner.State]
	if !ok {
		return errors.New(helper.StateParamErr)
	}

	tx, err := meta.MerchantDB.Beginx()
	if err != nil {
		return errors.New(helper.TransErr)
	}

	if banner.Flags == "1" && status == BannerStateProcessing {

		exFlags := g.Ex{
			"flags":  1,
			"prefix": meta.Prefix,
		}
		record := g.Record{
			"state": BannerStateEnd,
		}
		query, _, _ = dialect.Update("tbl_banner").Set(record).Where(exFlags).ToSQL()
		_, err = tx.Exec(query)
		if err != nil {
			_ = tx.Rollback()
			return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
		}
	}

	record := g.Record{
		"state": status,
	}
	query, _, _ = dialect.Update("tbl_banner").Set(record).Where(ex).ToSQL()
	_, err = tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	err = tx.Commit()
	if err == nil {
		return bannerRefreshToCache(id)
	}

	return nil
}

func BannerDelete(id string) error {

	ex := g.Ex{
		"id":    id,
		"state": []int{BannerStateWait, BannerStateEnd}, // 启用的时候不能进行编辑和删除
	}
	query, _, _ := dialect.Delete("tbl_banner").Where(ex).ToSQL()
	_, err := meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	bannerRefreshToCache(id)

	//fmt.Println("BannerDelete", bannerRefreshToCache(id))

	return nil
}

func BannerList(startTime, endTime string, page, pageSize uint, exs exp.ExpressionList) (BannerData, error) {

	data := BannerData{}
	//ex["prefix"] = meta.Prefix
	exs = exs.Append(g.Ex{"prefix": meta.Prefix})

	orEx := g.Or()
	if startTime != "" && endTime != "" {
		startAt, err := helper.TimeToLoc(startTime, loc)
		if err != nil {
			return data, errors.New(helper.DateTimeErr)
		}

		endAt, err := helper.TimeToLoc(endTime, loc)
		if err != nil {
			return data, errors.New(helper.DateTimeErr)
		}

		if startAt >= endAt {
			return data, errors.New(helper.QueryTimeRangeErr)
		}

		orEx = g.Or(
			g.Ex{"show_at": g.Op{"between": exp.NewRangeVal(startAt, endAt)}},
			g.Ex{"hide_at": g.Op{"between": exp.NewRangeVal(startAt, endAt)}},
			g.And(g.Ex{"show_at": g.Op{"lt": startAt}}, g.Ex{"hide_at": g.Op{"gt": endAt}}),
			g.Ex{"show_type": 1},
		)

		exs = exs.Append(orEx)
	}

	t := dialect.From("tbl_banner")
	if page == 1 {
		//query, _, _ := t.Select(g.COUNT(1)).Where(g.And(ex, orEx)).ToSQL()
		query, _, _ := t.Select(g.COUNT(1)).Where(exs).ToSQL()
		err := meta.MerchantDB.Get(&data.T, query)
		if err != nil {
			return data, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
		}

		if data.T == 0 {
			return data, nil
		}
	}

	data.S = pageSize
	offset := (page - 1) * pageSize
	//query, _, _ := t.Select(colsBanner...).Where(g.And(ex, orEx)).
	//	Order(g.C("updated_at").Desc()).Offset(offset).Limit(pageSize).ToSQL()
	query, _, _ := t.Select(colsBanner...).Where(exs).
		Order(g.C("updated_at").Desc()).Offset(offset).Limit(pageSize).ToSQL()
	err := meta.MerchantDB.Select(&data.D, query)
	if err != nil && err != sql.ErrNoRows {
		return data, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	return data, nil
}

func BannerUpdate(showAt, hideAt, id string, record g.Record) error {

	var (
		err error
		st  time.Time
		et  time.Time
		now = time.Now()
	)

	showType, _ := record["show_type"].(string)
	if showType == BannerShowTypeSpecify {
		st, err = time.Parse("2006-01-02 15:04:05", showAt)
		if err != nil {
			return errors.New(helper.TimeTypeErr)
		}

		et, err = time.Parse("2006-01-02 15:04:05", hideAt)
		if err != nil {
			return errors.New(helper.TimeTypeErr)
		}
	}

	record["show_at"] = "0"
	record["hide_at"] = "0"
	if showType == BannerShowTypeSpecify {
		record["show_at"] = fmt.Sprintf("%d", st.Unix())
		record["hide_at"] = fmt.Sprintf("%d", et.Unix())
	}
	ex := g.Ex{
		"id": id,
		"state": []int{
			BannerStateWait,
			BannerStateEnd,
		},
		"prefix": meta.Prefix,
	}
	query, _, _ := dialect.Update("tbl_banner").Set(record).Where(ex).ToSQL()
	_, err = meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	if showType == BannerShowTypeSpecify {
		// banner自动展示
		sDelay := st.Sub(now).Seconds()
		BeanPutDelay("Banner", map[string]interface{}{"id": id, "ty": "2"}, int(sDelay)-5)

		// banner自动隐藏
		eDelay := et.Sub(now).Seconds()
		BeanPutDelay("Banner", map[string]interface{}{"id": id, "ty": "3"}, int(eDelay)-5)
	}

	fmt.Println("BannerUpdate", bannerRefreshToCache(id))

	return nil
}

func BannerInsert(record Banner) error {

	var (
		err error
		st  time.Time
		et  time.Time
		now = time.Now()
	)

	if record.ShowType == BannerShowTypeSpecify {
		st, err = time.Parse("2006-01-02 15:04:05", record.ShowAt)
		if err != nil {
			return errors.New(helper.TimeTypeErr)
		}

		et, err = time.Parse("2006-01-02 15:04:05", record.HideAt)
		if err != nil {
			return errors.New(helper.TimeTypeErr)
		}
	}

	record.ShowAt = "0"
	record.HideAt = "0"
	record.Prefix = meta.Prefix
	if record.ShowType == BannerShowTypeSpecify {
		record.ShowAt = fmt.Sprintf("%d", st.Unix())
		record.HideAt = fmt.Sprintf("%d", et.Unix())
	}

	query, _, _ := dialect.Insert("tbl_banner").Rows(&record).ToSQL()
	_, err = meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	if record.ShowType == BannerShowTypeSpecify {
		// banner自动展示
		sDelay := st.Sub(now).Seconds()
		BeanPutDelay("Banner", map[string]interface{}{"id": record.ID, "ty": "2"}, int(sDelay)-5)

		// banner自动隐藏
		eDelay := et.Sub(now).Seconds()
		BeanPutDelay("Banner", map[string]interface{}{"id": record.ID, "ty": "3"}, int(eDelay)-5)
	}

	fmt.Println("BannerInsert", bannerRefreshToCache(record.ID))

	return nil
}

// 加载banner到缓存
func LoadBanners() error {

	single := []int{1}
	array := []int{2, 5}
	pipe := meta.MerchantRedis.TxPipeline()
	defer pipe.Close()

	for _, v := range single {
		for k := range DeviceMap {
			key := fmt.Sprintf("%s:banner:G%d%d", meta.Prefix, v, k)
			pipe.Unlink(ctx, key)
		}

		singleBanner := Banner{}
		ex := g.Ex{
			"state":  2,
			"flags":  v,
			"prefix": meta.Prefix,
		}
		query, _, _ := dialect.From("tbl_banner").Select(colsBanner...).Where(ex).ToSQL()
		query = "/* master */ " + query
		fmt.Println(query)
		err := meta.MerchantDB.Get(&singleBanner, query)
		if err != nil {
			if err != sql.ErrNoRows {
				_ = pushLog(err, helper.DBErr)
			}
			continue
		}

		base := fastjson.MustParse(singleBanner.Images)
		base.Set("id", fastjson.MustParse(fmt.Sprintf(`"%s"`, singleBanner.ID)))
		base.Set("url", fastjson.MustParse(fmt.Sprintf(`"%s"`, singleBanner.RedirectURL)))
		base.Set("flags", fastjson.MustParse(fmt.Sprintf(`"%s"`, singleBanner.URLType)))

		if singleBanner.Device == "0" {
			for k := range DeviceMap {
				key := fmt.Sprintf("%s:banner:G%d%d", meta.Prefix, v, k)
				pipe.Set(ctx, key, base.String(), 100*time.Hour)
				pipe.Persist(ctx, key)
			}
		} else {
			di := strings.SplitN(singleBanner.Device, ",", 8)
			for _, val := range di {
				key := fmt.Sprintf("%s:banner:G%d%s", meta.Prefix, v, val)
				pipe.Set(ctx, key, base.String(), 100*time.Hour)
				pipe.Persist(ctx, key)
			}
		}
	}

	for _, v := range array {
		for k := range DeviceMap {
			key := fmt.Sprintf("%s:banner:G%d%d", meta.Prefix, v, k)
			pipe.Unlink(ctx, key)
		}

		var recs []Banner
		ex := g.Ex{
			"state":  2,
			"flags":  v,
			"prefix": meta.Prefix,
		}
		query, _, _ := dialect.From("tbl_banner").Select(colsBanner...).Where(ex).Order(g.C("seq").Asc()).ToSQL()
		query = "/* master */ " + query
		fmt.Println(query)
		err := meta.MerchantDB.Select(&recs, query)
		if err != nil {
			if err != sql.ErrNoRows {
				_ = pushLog(err, helper.DBErr)
			}
			continue
		}

		results := map[string][]string{}
		for _, val := range recs {
			val.Title = strings.ReplaceAll(val.Title, "&nbsp;", " ")
			img := fastjson.GetBytes([]byte(val.Images), "ad")
			str := fmt.Sprintf(`{"id":"%s", "title":"%s", "url":"%s", "sort":"%s", "img":"%s", "flags":"%s"}`,
				val.ID, val.Title, val.RedirectURL, val.Seq, string(img), val.URLType)

			if val.Device == "0" {
				for k := range DeviceMap {
					key := fmt.Sprintf("%s:banner:G%d%d", meta.Prefix, v, k)
					results[key] = append(results[key], str)
				}
			} else {
				di := strings.SplitN(val.Device, ",", 8)
				for _, d := range di {
					key := fmt.Sprintf("%s:banner:G%d%s", meta.Prefix, v, d)
					results[key] = append(results[key], str)
				}
			}

			arr := new(fastjson.Arena)
			for key, value := range results {

				aa := arr.NewArray()
				for k, v := range value {
					aa.SetArrayItem(k, fastjson.MustParse(v))
				}

				pipe.Set(ctx, key, aa.String(), 100*time.Hour)
				pipe.Persist(ctx, key)
				arr.Reset()
			}
			arr = nil
		}
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return pushLog(err, helper.RedisErr)
	}

	return nil
}

func bannerRefreshToCache(id string) error {

	single := map[string]bool{
		"1": true,
	}
	array := map[string]bool{
		"2": true,
		"5": true,
	}
	record := Banner{}
	ex := g.Ex{
		"id":     id,
		"prefix": meta.Prefix,
	}
	query, _, _ := dialect.From("tbl_banner").Select(colsBanner...).Where(ex).ToSQL()
	err := meta.MerchantDB.Get(&record, query)
	if err != nil {
		if err != sql.ErrNoRows {
			err = pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
		}
		return err
	}

	pipe := meta.MerchantRedis.TxPipeline()
	defer pipe.Close()

	flags := record.Flags
	// 单一banner
	ok := single[record.Flags]
	if ok {
		for k := range DeviceMap {
			key := fmt.Sprintf("%s:banner:G%s%d", meta.Prefix, flags, k)
			pipe.Unlink(ctx, key)
		}

		recs := Banner{}
		ex = g.Ex{
			"state":  2,
			"flags":  flags,
			"prefix": meta.Prefix,
		}
		query, _, _ = dialect.From("tbl_banner").Select(colsBanner...).Where(ex).ToSQL()
		err = meta.MerchantDB.Get(&recs, query)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if err == sql.ErrNoRows {

			//首页弹窗只有一条生效，则在手动关闭后，查无数据，即需要提交redis更新
			_, err = pipe.Exec(ctx)
			if err != nil {
				return pushLog(err, helper.RedisErr)
			}

			return nil
		}

		base := fastjson.MustParse(recs.Images)
		base.Set("url", fastjson.MustParse(fmt.Sprintf(`"%s"`, recs.RedirectURL)))
		base.Set("flags", fastjson.MustParse(fmt.Sprintf(`"%s"`, recs.URLType)))

		// 全端支持
		if recs.Device == "0" {
			for k := range DeviceMap {
				key := fmt.Sprintf("%s:banner:G%s%d", meta.Prefix, flags, k)
				pipe.Set(ctx, key, base.String(), 100*time.Hour)
				pipe.Persist(ctx, key)
			}
		} else {
			di := strings.SplitN(recs.Device, ",", 8)
			for _, val := range di {
				key := fmt.Sprintf("%s:banner:G%s%s", meta.Prefix, flags, val)
				pipe.Set(ctx, key, base.String(), 100*time.Hour)
				pipe.Persist(ctx, key)
			}
		}
	}

	// 数组banner
	ok = array[record.Flags]
	if ok {
		for k := range DeviceMap {
			key := fmt.Sprintf("%s:banner:G%s%d", meta.Prefix, flags, k)
			pipe.Unlink(ctx, key)
		}

		var recs []Banner
		ex = g.Ex{
			"state":  2,
			"flags":  flags,
			"prefix": meta.Prefix,
		}
		query, _, _ = dialect.From("tbl_banner").Select(colsBanner...).Where(ex).Order(g.C("seq").Asc()).ToSQL()
		err = meta.MerchantDB.Select(&recs, query)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if err == sql.ErrNoRows {
			return nil
		}

		results := map[string][]string{}
		for _, val := range recs {

			img := fastjson.GetBytes([]byte(val.Images), "ad")
			str := fmt.Sprintf(`{"id":"%s", "title":"%s", "url":"%s", "sort":"%s", "img":"%s", "flags":"%s"}`, val.ID, val.Title, val.RedirectURL, val.Seq, string(img), val.URLType)

			if val.Device == "0" {
				for k := range DeviceMap {
					key := fmt.Sprintf("%s:banner:G%s%d", meta.Prefix, flags, k)
					results[key] = append(results[key], str)
				}
			} else {
				di := strings.SplitN(val.Device, ",", 8)
				for _, d := range di {
					key := fmt.Sprintf("%s:banner:G%s%s", meta.Prefix, flags, d)
					results[key] = append(results[key], str)
				}
			}
		}

		arr := new(fastjson.Arena)
		for key, value := range results {

			aa := arr.NewArray()
			for k, v := range value {
				aa.SetArrayItem(k, fastjson.MustParse(v))
			}

			pipe.Set(ctx, key, aa.String(), 100*time.Hour)
			pipe.Persist(ctx, key)
			arr.Reset()
		}
		arr = nil
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return pushLog(err, helper.RedisErr)
	}

	return nil
}
