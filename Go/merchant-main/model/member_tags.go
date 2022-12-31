package model

import (
	"errors"
	"fmt"
	"merchant/contrib/helper"
	"strconv"
	"strings"

	g "github.com/doug-martin/goqu/v9"
)

// Tags 标签表 table structure
type Tags struct {
	ID          int64  `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	Color       string `db:"color" json:"color"`
	Sort        int64  `db:"sort" json:"sort"`
	Flags       int64  `db:"flags" json:"flags"`
	Members     int64  `db:"-" json:"members"`
	CreatedAt   int64  `db:"created_at" json:"created_at"`
	UpdatedAt   int64  `db:"updated_at" json:"updated_at"`
}

// TagsData 标签管理-标签列表 response structure
type TagsData struct {
	D []Tags `json:"d"`
	T int64  `json:"t"`
	S int    `json:"s"`
}

// MemberTags 用户标签表
type MemberTags struct {
	ID        string `db:"id" json:"id"`
	UID       string `db:"uid" json:"uid"`
	AdminID   string `db:"admin_id" json:"admin_id"`
	TagID     string `db:"tag_id" json:"tag_id"`
	TagName   string `db:"tag_name" json:"tags_name"`
	CreatedAt int64  `db:"created_at" json:"created_at"`
	UpdatedAt int64  `db:"updated_at" json:"updated_at"`
	Prefix    string `db:"prefix" json:"prefix"`
}

type memberTagSum struct {
	ID  int64 `db:"tag_id"`
	Sum int64 `db:"sum"`
}

type MemberTagsData struct {
	D []MemberTags `json:"d"`
	T int64        `json:"t"`
}

// MemberTagsList 获取用户所有标签
func MemberTagsList(uid string) (MemberTagsData, error) {

	var data MemberTagsData

	ex := g.Ex{"uid": uid, "prefix": meta.Prefix}
	query, _, _ := dialect.From("tbl_member_tags").Select(g.COUNT(1)).Where(ex).ToSQL()
	err := meta.MerchantDB.Get(&data.T, query)
	if err != nil {
		return data, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	query, _, _ = dialect.From("tbl_member_tags").Select(colsMemberTags...).Where(ex).ToSQL()
	err = meta.MerchantDB.Select(&data.D, query)
	if err != nil {
		return data, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	return data, nil
}

// MemberTagsSet 设置用户标签 batch=0时编辑单个用户标签,batch=1时批量添加标签
func MemberTagsSet(batch int, adminID string, uids []string, tags []string, ts int64) error {

	// 验证tags是否存在
	var tagData []Tags

	tagQuery, _, _ := dialect.From("tbl_tags").Select(colsTags...).Where(g.Ex{"id": tags, "prefix": meta.Prefix}).ToSQL()
	err := meta.MerchantDB.Select(&tagData, tagQuery)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), tagQuery), helper.DBErr)
	}

	if len(tags) != len(tagData) {
		return errors.New(helper.IDErr)
	}

	var total int
	// 验证uid
	totalQuery, _, _ := dialect.From("tbl_members").Select(g.COUNT("uid")).Where(g.Ex{"uid": g.Op{"in": uids}}).ToSQL()
	_ = meta.MerchantDB.Get(&total, totalQuery)
	if total != len(uids) {
		return errors.New(helper.UIDErr)
	}

	var data []MemberTags
	ex := g.Ex{"uid": g.Op{"in": uids}, "prefix": meta.Prefix}
	// 编辑多用户标签
	if batch == 1 {
		ex["tag_id"] = g.Op{"in": tags}
	}

	var tagls []interface{}
	for _, v := range tags {
		t, _ := strconv.ParseInt(v, 10, 64)
		tagls = append(tagls, t)
	}

	var ids []interface{}
	// 拼接现在的标签
	for _, uid := range uids {
		ids = append(ids, uid)
		for _, v := range tagData {
			data = append(data, MemberTags{
				ID:        helper.GenId(),
				UID:       uid,
				AdminID:   adminID,
				TagID:     fmt.Sprintf("%d", v.ID),
				TagName:   v.Name,
				CreatedAt: ts,
				UpdatedAt: ts,
				Prefix:    meta.Prefix,
			})
		}
	}

	tx, err := meta.MerchantDB.Begin()
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	// 删除以前的标签
	query, _, _ := dialect.Delete("tbl_member_tags").Where(ex).ToSQL()
	_, err = tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	// 添加现在的标签
	query, _, _ = dialect.Insert("tbl_member_tags").Rows(data).ToSQL()
	_, err = tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	err = tx.Commit()
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	return nil
}

// MemberTagsCancel 取消用户标签
func MemberTagsCancel(uidStr string, tagStr string) error {

	var tags []interface{}
	for _, v := range strings.Split(tagStr, ",") {
		t, _ := strconv.ParseInt(v, 10, 64)
		tags = append(tags, t)
	}

	var uids []interface{}
	for _, uid := range strings.Split(uidStr, ",") {
		uids = append(uids, uid)
	}

	query, _, _ := dialect.Delete("tbl_member_tags").
		Where(g.Ex{"uid": g.Op{"in": uids}, "tag_id": g.Op{"in": tags}, "prefix": meta.Prefix}).ToSQL()
	tx, err := meta.MerchantDB.Beginx()
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	_, err = tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	err = tx.Commit()
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	return nil
}

// TagList 标签管理-标签列表
func TagList(name string, flag, all, page, pageSize int) (TagsData, error) {

	var data TagsData

	ex := g.Ex{
		"prefix": meta.Prefix,
	}
	if name != "" {
		if strings.Contains(name, " ") {
			ex["name"] = []string{strings.ReplaceAll(name, " ", "&nbsp;"), name}
		} else {
			ex["name"] = name
		}
	}

	if flag > 0 {
		ex["flags"] = flag
	}

	if page == 1 {
		query, _, _ := dialect.From("tbl_tags").Select(g.COUNT(1)).Where(ex).ToSQL()
		fmt.Println(query)
		err := meta.MerchantDB.Get(&data.T, query)
		if err != nil {
			return data, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
		}

		if data.T == 0 {
			return data, nil
		}
	}

	offset := (page - 1) * pageSize
	dl := dialect.From("tbl_tags").Select(colsTags...).Where(ex).Order(g.C("updated_at").Desc())
	if all == 0 { // 是否取所有的标签（用户列表需要使用所有的标签）
		dl = dl.Offset(uint(offset)).Limit(uint(pageSize))
	}

	query, _, _ := dl.ToSQL()
	fmt.Println(query)
	err := meta.MerchantDB.Select(&data.D, query)
	if err != nil {
		return data, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	if len(data.D) == 0 {
		return data, nil
	}

	var ids []int64
	for _, v := range data.D {
		ids = append(ids, v.ID)
	}

	var mts []memberTagSum
	// 构造查询用户数量的sql
	query, _, _ = dialect.From("tbl_member_tags").Select([]interface{}{"tag_id", g.COUNT("tag_id").As("sum")}...).
		Where(g.Ex{"tag_id": g.Op{"in": ids}, "prefix": meta.Prefix}).GroupBy("tag_id").ToSQL()
	err = meta.MerchantDB.Select(&mts, query)
	if err != nil {
		return data, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	m := make(map[int64]int64)
	for _, v := range mts {
		m[v.ID] = v.Sum
	}

	for k, v := range data.D {
		data.D[k].Members = m[v.ID]
	}

	data.S = pageSize
	return data, nil
}

// TagInsert 标签管理-新增标签
func TagInsert(params map[string]interface{}) error {

	params["prefix"] = meta.Prefix
	query, _, _ := dialect.Insert("tbl_tags").Rows(params).ToSQL()
	_, err := meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	return nil
}

// TagUpdate 标签管理-修改标签
func TagUpdate(ex g.Ex, record g.Record) error {

	ex["prefix"] = meta.Prefix
	query, _, _ := dialect.Update("tbl_tags").Set(record).Where(ex).ToSQL()
	_, err := meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	return nil
}

// TagDelete 标签管理-删除标签
func TagDelete(id string) error {

	var ids []string
	query, _, _ := dialect.From("tbl_member_tags").Select("id").Where(g.Ex{"tag_id": id, "prefix": meta.Prefix}).ToSQL()
	err := meta.MerchantDB.Select(&ids, query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	if len(ids) > 0 {
		return fmt.Errorf(helper.MemberTagInUse)
	}

	query, _, _ = dialect.Delete("tbl_tags").Where(g.Ex{"id": id, "prefix": meta.Prefix}).ToSQL()
	_, err = meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	return nil
}

// TagByNameAndFlag 通过name和flag查找tag
func TagByNameAndFlag(name, flag string) (Tags, error) {

	tags := Tags{}

	ex := g.Ex{
		"name":   name,
		"flags":  flag,
		"prefix": meta.Prefix,
	}
	query, _, _ := dialect.From("tbl_tags").Select(colsTags...).Where(ex).ToSQL()
	err := meta.MerchantDB.Get(&tags, query)
	if err != nil {
		return tags, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	return tags, nil
}
