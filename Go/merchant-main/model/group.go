package model

import (
	"errors"
	"fmt"
	g "github.com/doug-martin/goqu/v9"
	"merchant/contrib/helper"
	"strings"
	"time"
)

type Group struct {
	CreateAt   int32  `db:"create_at" rule:"none" json:"create_at"`                                                            //创建时间
	Gid        string `db:"gid" rule:"none" json:"gid"`                                                                        //
	Gname      string `db:"gname" name:"gname" rule:"chn" min:"2" max:"8" msg:"gname error[2-8]" json:"gname"`                 //组名
	Lft        int64  `db:"lft" rule:"none" json:"lft"`                                                                        //节点左值
	Lvl        int64  `db:"lvl" rule:"none" json:"lvl"`                                                                        //
	Noted      string `db:"noted" name:"noted" rule:"none" default:"" min:"0" max:"511" msg:"noted error[0-511]" json:"noted"` //备注信息
	Permission string `db:"permission" name:"permission" rule:"sDigit" min:"2" msg:"permission error[2-]" json:"permission"`   //权限模块ID
	Rgt        int64  `db:"rgt" rule:"none" json:"rgt"`                                                                        //节点右值
	Pid        string `db:"pid" rule:"none" json:"pid"`                                                                        //父节点
	State      int    `db:"state" json:"state" name:"state" rule:"digit" min:"0" max:"1" msg:"state error"`                    //0:关闭1:开启
	Prefix     string `db:"prefix" rule:"none" json:"prefix"`
}

func GroupUpdate(gid, adminGid string, data Group) error {

	// 所修改分组的权限map
	gPrivMap := make(map[string]bool)
	permissions := strings.Split(data.Permission, ",")
	// 检查新增分组的权限是否大于上级分组
	for _, v := range permissions {
		key := fmt.Sprintf("%s:priv:GM%s", meta.Prefix, data.Pid)
		exists := meta.MerchantRedis.HExists(ctx, key, v).Val()
		if !exists {
			return errors.New(helper.MethodNoPermission)
		}

		gPrivMap[v] = true
	}

	// 检查当前后台账号是否有权限增加当前分组
	ok, err := groupSubCheck(gid, adminGid)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New(helper.MethodNoPermission)
	}

	// 检查当前分组名是否已存在
	ex := g.Ex{
		"gid":   g.Op{"neq": gid},
		"gname": data.Gname,
	}
	ok, err = groupExistCheck(ex)
	if err != nil {
		return err
	}

	if ok {
		return errors.New(helper.RecordExistErr)
	}

	var parent Group
	ex = g.Ex{
		"gid":    data.Pid,
		"prefix": meta.Prefix,
	}
	query, _, _ := dialect.From("tbl_admin_group").Select(colsGroup...).Order(g.C("lvl").Asc()).Where(ex).ToSQL()
	fmt.Println(query)
	err = meta.MerchantDB.Get(&parent, query)
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	parentPrivs := strings.Split(parent.Permission, ",")
	fmt.Println(len(permissions), len(parentPrivs))
	fmt.Println(permissions, parentPrivs)
	if len(permissions) >= len(parentPrivs) {
		return errors.New(helper.SubPermissionEqualErr)
	}

	// 获取当前分组的下级分组
	subGids, err := groupSubs(gid)
	if err != nil {
		return err
	}

	var subs []Group
	if len(subGids) > 0 {
		ex := g.Ex{
			"prefix": meta.Prefix,
			"gid":    subGids,
		}
		query, _, _ := dialect.From("tbl_admin_group").
			Select(colsGroup...).Where(ex).Order(g.C("lvl").Asc()).ToSQL()
		fmt.Println(query)
		err = meta.MerchantDB.Select(&subs, query)
		if err != nil {
			return pushLog(err, helper.DBErr)
		}
	}

	tx, err := meta.MerchantDB.Begin()
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	if len(subs) > 0 {
		fmt.Println(subs)
		for _, v := range subs {
			privs := ""
			for _, vv := range strings.Split(v.Permission, ",") {
				// 下级权限在分组权限调整后的范围内保留，不在则删除
				if _, ok = gPrivMap[vv]; ok {
					if privs != "" {
						privs += ","
					}
					privs += vv
				}
			}

			record := g.Record{
				"permission": privs,
			}
			query, _, _ := dialect.Update("tbl_admin_group").Set(record).Where(g.Ex{"gid": v.Gid}).ToSQL()
			fmt.Println(query)
			_, err = tx.Exec(query)
			if err != nil {
				return pushLog(err, helper.DBErr)
			}
		}
	}

	record := g.Record{
		"gname":      data.Gname,
		"noted":      data.Noted,
		"state":      data.State,
		"permission": data.Permission,
	}
	query, _, _ = dialect.Update("tbl_admin_group").Set(record).Where(g.Ex{"gid": gid}).ToSQL()
	fmt.Println(query)
	_, err = tx.Exec(query)
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	_ = tx.Commit()

	return LoadGroups()
}

func groupSubs(pid string) ([]string, error) {

	var descendants []string
	ex := g.Ex{
		"prefix":     meta.Prefix,
		"ancestor":   pid,
		"descendant": g.Op{"neq": pid},
	}
	query, _, _ := dialect.From("tbl_admin_group_tree").
		Select("descendant").Where(ex).GroupBy("descendant").Order(g.C("lvl").Asc()).ToSQL()
	fmt.Println(query)
	err := meta.MerchantDB.Select(&descendants, query)
	if err != nil {
		return descendants, pushLog(err, helper.DBErr)
	}

	return descendants, nil
}

func GroupInsert(adminGid string, data Group) error {

	privs := strings.Split(data.Permission, ",")
	// 检查新增分组的权限是否大于上级分组
	for _, v := range privs {
		key := fmt.Sprintf("%s:priv:GM%s", meta.Prefix, data.Pid)
		exists := meta.MerchantRedis.HExists(ctx, key, v).Val()
		if !exists {
			return errors.New(helper.MethodNoPermission)
		}
	}

	// 检查当前后台账号是否有权限增加当前分组
	ok, err := groupSubCheck(data.Pid, adminGid)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New(helper.MethodNoPermission)
	}

	// 检查当前分组名是否已存在
	ex := g.Ex{
		"gname": data.Gname,
	}
	ok, err = groupExistCheck(ex)
	if err != nil {
		return err
	}

	if ok {
		return errors.New(helper.RecordExistErr)
	}

	var parent Group
	ex = g.Ex{
		"gid":    data.Pid,
		"prefix": meta.Prefix,
	}
	query, _, _ := dialect.From("tbl_admin_group").Select(colsGroup...).Order(g.C("lvl").Asc()).Where(ex).ToSQL()
	fmt.Println(query)
	err = meta.MerchantDB.Get(&parent, query)
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	parentPrivs := strings.Split(parent.Permission, ",")
	fmt.Println(len(privs), len(parentPrivs))
	fmt.Println(privs, parentPrivs)
	if len(privs) >= len(parentPrivs) {
		return errors.New(helper.SubPermissionEqualErr)
	}

	tx, err := meta.MerchantDB.Begin()
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	gid := helper.GenId()
	data.Lvl = parent.Lvl + 1
	data.Gid = gid
	data.Prefix = meta.Prefix
	query, _, _ = dialect.Insert("tbl_admin_group").Rows(data).ToSQL()
	fmt.Println(query)
	_, err = tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		return pushLog(err, helper.DBErr)
	}

	query = GroupClosureInsert(gid, data.Pid)
	fmt.Println(query)
	_, err = tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		return pushLog(err, helper.DBErr)
	}

	err = tx.Commit()
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	return LoadGroups()
}

func groupExistCheck(ex g.Ex) (bool, error) {

	var count int
	ex["prefix"] = meta.Prefix
	query, _, _ := dialect.From("tbl_admin_group").Select(g.COUNT("gid")).Where(ex).ToSQL()
	err := meta.MerchantDB.Get(&count, query)
	fmt.Println(query)
	if err != nil {
		return false, pushLog(err, helper.DBErr)
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

// 检查当前后台账号所属分组是不是所操作分组的上级
func groupSubCheck(gid, parentGid string) (bool, error) {

	var count int
	ex := g.Ex{
		"prefix":     meta.Prefix,
		"ancestor":   parentGid,
		"descendant": gid,
	}
	query, _, _ := dialect.From("tbl_admin_group_tree").
		Select(g.COUNT("descendant")).Where(ex).GroupBy("descendant").ToSQL()
	fmt.Println(query)
	err := meta.MerchantDB.Get(&count, query)
	if err != nil {
		return false, pushLog(err, helper.DBErr)
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

// 查询当前代理的分组和下级分组gid
func groupSubList(gid string) ([]string, map[string]bool, error) {

	var gids []string
	gidMap := make(map[string]bool)
	ex := g.Ex{
		"ancestor": gid,
		"prefix":   meta.Prefix,
	}
	query, _, _ := dialect.From("tbl_admin_group_tree").
		Select("descendant").Where(ex).GroupBy("descendant").ToSQL()
	fmt.Println(query)
	err := meta.MerchantDB.Select(&gids, query)
	if err != nil {
		return nil, nil, pushLog(err, helper.DBErr)
	}

	for _, v := range gids {
		gidMap[v] = true
	}

	return gids, gidMap, nil
}

func GroupList(gid, adminGid string) (string, error) {

	if gid != "" {
		ok, err := groupSubCheck(gid, adminGid)
		if err != nil {
			return "[]", err
		}

		if !ok {
			return "[]", errors.New(helper.MethodNoPermission)
		}
	} else {
		gid = adminGid
	}

	gids, _, err := groupSubList(gid)
	if err != nil {
		return "[]", err
	}

	if len(gids) == 0 {
		return "[]", nil
	}

	var groups []Group
	ex := g.Ex{
		"gid":    gids,
		"prefix": meta.Prefix,
	}
	query, _, _ := dialect.From("tbl_admin_group").Select(colsGroup...).Where(ex).Order(g.C("create_at").Asc()).ToSQL()
	fmt.Println(query)
	err = meta.MerchantDB.Select(&groups, query)
	if err != nil {
		return "[]", pushLog(err, helper.DBErr)
	}

	recs, err := helper.JsonMarshal(groups)
	if err != nil {
		return "[]", errors.New(helper.FormatErr)
	}

	//key := fmt.Sprintf("%s:priv:GroupAll", meta.Prefix)
	//val, err := meta.MerchantRedis.Get(ctx, key).Result()
	//if err != nil && err != redis.Nil {
	//	return val, pushLog(err, helper.RedisErr)
	//}

	return string(recs), nil
}

func GroupClosureInsert(nodeID, targetID string) string {

	t := "SELECT ancestor, " + nodeID + ",prefix, lvl+1 FROM tbl_admin_group_tree WHERE prefix='" + meta.Prefix + "' and descendant = " + targetID + " UNION SELECT " + nodeID + "," + nodeID + "," + "'" + meta.Prefix + "'" + ",0"
	query := "INSERT INTO tbl_admin_group_tree (ancestor, descendant,prefix,lvl) (" + t + ")"

	return query
}

func LoadGroups() error {

	var (
		groups []Group
		privs  []Priv
	)
	cols := []interface{}{"noted", "gid", "gname", "permission", "create_at", "state", "lft", "rgt", "lvl", "pid"}
	ex := g.Ex{
		"prefix": meta.Prefix,
	}
	query, _, _ := dialect.From("tbl_admin_group").Select(cols...).Where(ex).ToSQL()
	query = "/* master */ " + query
	fmt.Println(query)
	err := meta.MerchantDB.Select(&groups, query)
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	query, _, _ = dialect.From("tbl_admin_priv").
		Select("pid", "state", "id", "name", "sortlevel", "module").Where(g.Ex{"prefix": meta.Prefix}).Order(g.C("sortlevel").Asc()).ToSQL()
	query = "/* master */ " + query
	fmt.Println(query)
	err = meta.MerchantDB.Select(&privs, query)
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	if len(groups) == 0 || len(privs) == 0 {
		return nil
	}

	privMap := make(map[string]Priv)
	permission := ""
	for _, v := range privs {
		privMap[v.ID] = v
		if permission != "" {
			permission += ","
		}
		permission += v.ID
	}

	record := g.Ex{
		"permission": permission,
	}
	query, _, _ = dialect.Update("tbl_admin_group").Set(record).Where(g.Ex{"gid": 2, "prefix": meta.Prefix}).ToSQL()
	query = "/* master */ " + query
	fmt.Println(query)
	_, err = meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	recs, err := helper.JsonMarshal(groups)
	if err != nil {
		return errors.New(helper.FormatErr)
	}

	pipe := meta.MerchantRedis.TxPipeline()
	defer pipe.Close()

	key := fmt.Sprintf("%s:priv:GroupAll", meta.Prefix)
	pipe.Unlink(ctx, key)
	pipe.Set(ctx, key, string(recs), 100*time.Hour)
	pipe.Persist(ctx, key)

	for _, val := range groups {

		id := fmt.Sprintf("%s:priv:GM%s", meta.Prefix, val.Gid)
		pipe.Unlink(ctx, id)
		// 只保存开启状态的分组
		if val.State == 1 {
			gKey := fmt.Sprintf("%s:priv:list:GM%s", meta.Prefix, val.Gid)
			pipe.Unlink(ctx, gKey)
			if val.Gid != "2" {
				var gPrivs []Priv
				for _, v := range strings.Split(val.Permission, ",") {
					pipe.HSet(ctx, id, v, "1")
					gPrivs = append(gPrivs, privMap[v])
				}
				gRecs, _ := helper.JsonMarshal(gPrivs)
				pipe.Set(ctx, gKey, string(gRecs), 100*time.Hour)
			} else {
				for _, v := range privs {
					pipe.HSet(ctx, id, v.ID, "1")
				}
				gRecs, _ := helper.JsonMarshal(privs)
				pipe.Set(ctx, gKey, string(gRecs), 100*time.Hour)
			}
			pipe.Persist(ctx, id)
			pipe.Persist(ctx, gKey)
		}
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return pushLog(err, helper.RedisErr)
	}

	return nil
}
