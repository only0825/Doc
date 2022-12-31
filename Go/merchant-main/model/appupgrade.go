package model

import (
	"database/sql"
	"fmt"
	"merchant/contrib/helper"
	"time"

	g "github.com/doug-martin/goqu/v9"
)

// app 自动升级配置
type AppUpgrade struct {
	ID          string `db:"id" json:"id"`
	Platform    string `db:"platform" json:"platform"`
	Version     string `db:"version" json:"version"`
	IsForce     uint8  `db:"is_force" json:"is_force"`
	Content     string `db:"content" json:"content"`
	URL         string `db:"url" json:"url"`
	UpdatedAt   uint32 `db:"updated_at" json:"updated_at"`
	UpdatedUid  string `db:"updated_uid" json:"updated_uid"`
	UpdatedName string `db:"updated_name" json:"updated_name"`
	Prefix      string `db:"prefix" json:"prefix"`
}

// 升级信息
type UpgradeInfo struct {
	Platform string `db:"platform" json:"platform"`
	Version  string `db:"version" json:"version"`
	IsForce  uint8  `db:"is_force" json:"is_force"`
	Content  string `db:"content" json:"content"`
	URL      string `db:"url" json:"url"`
}

// 更新升级配置
func AppUpgradeUpdate(device, version, content, url, isForce string, admin map[string]string) error {

	//recs := AppUpgrade{}
	//query, _, _ := dialect.From("tbl_app_upgrade").Rows(recs).Where(g.Ex{"prefix": meta.Prefix}).ToSQL()

	query := fmt.Sprintf(`INSERT INTO tbl_app_upgrade (id,platform,version,content,url,is_force,updated_at,updated_uid,updated_name,prefix) VALUES('%s','%s','%s','%s','%s','%s','%d','%s','%s','%s') on duplicate key update version = '%s',content = '%s',url = '%s',is_force = '%s',updated_at = %d,updated_uid = %s,updated_name = '%s',prefix = '%s'`,
		helper.GenId(), device, version, content, url, isForce, time.Now().Unix(), admin["id"], admin["name"], meta.Prefix,
		version, content, url, isForce, time.Now().Unix(), admin["id"], admin["name"], meta.Prefix)

	_, err := meta.MerchantDB.Exec(query)
	if err != nil {
		fmt.Println("AppUpgradeUpdate query = ", query)
		fmt.Println("AppUpgradeUpdate err = ", err)
		return pushLog(err, helper.DBErr)
	}

	// 更新缓存
	LoadAppUpgrades()

	return nil
}

func AppUpgradeList() ([]AppUpgrade, error) {

	var data []AppUpgrade
	query, _, _ := dialect.From("tbl_app_upgrade").Select(colsAppUpgrade...).Where(g.Ex{"prefix": meta.Prefix}).ToSQL()
	err := meta.MerchantDB.Select(&data, query)
	if err != nil && err != sql.ErrNoRows {
		return data, pushLog(err, helper.DBErr)
	}

	return data, nil
}

func LoadAppUpgrades() {

	var data []UpgradeInfo
	query, _, _ := dialect.From("tbl_app_upgrade").Select("platform", "version", "is_force", "content", "url").Where(g.Ex{"prefix": meta.Prefix}).ToSQL()
	query = "/* master */ " + query
	fmt.Println(query)
	err := meta.MerchantDB.Select(&data, query)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	pipe := meta.MerchantRedis.TxPipeline()
	defer pipe.Close()

	for _, v := range data {
		key := fmt.Sprintf("%s:upgrade:%s", meta.Prefix, v.Platform)
		bytes, err := helper.JsonMarshal(&v)
		if err != nil {
			fmt.Println(err)
			continue
		}

		pipe.Unlink(ctx, key)
		pipe.Set(ctx, key, string(bytes), 100*time.Hour)
		pipe.Persist(ctx, key)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
