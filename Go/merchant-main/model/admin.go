package model

import (
	"database/sql"
	"errors"
	"fmt"
	g "github.com/doug-martin/goqu/v9"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
	"merchant/contrib/helper"
	"merchant/contrib/session"
	"strconv"
)

// Admin 用户数据库结构体
type Admin struct {
	ID            string `name:"id" db:"id" json:"id" rule:"none"` // 主键ID
	Prefix        string `name:"prefix" db:"prefix" json:"prefix" rule:"none"`
	GroupID       string `name:"group_id" db:"group_id" json:"group_id" rule:"digit" min:"1" msg:"group_id error"`         // 用户组ID
	LastLoginIP   string `name:"last_login_ip" db:"last_login_ip" rule:"none"`                                             // 最后登录IP
	LastLoginTime uint32 `name:"last_login_time" db:"last_login_time" rule:"none"`                                         // 最后登录时间
	Pwd           string `name:"password" db:"password" json:"password" rule:"apwd" min:"5" max:"20" msg:"password error"` // 密码
	State         int    `name:"state" db:"state" json:"state" rule:"digit" min:"0" max:"1" msg:"state error"`             // 状态
	Seamo         string `name:"seamo" db:"seamo" json:"seamo" rule:"none"`
	Name          string `name:"name" db:"name" json:"name" rule:"aname" min:"5" max:"20" msg:"name error"` // 用户名
	CreateAt      uint32 `name:"create_at" db:"create_at" json:"create_at" rule:"none"`                     // 创建时间
	CreatedUid    string `name:"created_uid" db:"created_uid" json:"created_uid" rule:"none"`               //创建人uid
	CreatedName   string `name:"created_name" db:"created_name" json:"created_name" rule:"none"`            //创建人名
	UpdatedAt     uint32 `name:"updated_at" db:"updated_at" json:"updated_at" rule:"none"`                  // 修改时间
	UpdatedUid    string `name:"updated_uid" db:"updated_uid" json:"updated_uid" rule:"none"`               //修改人uid
	UpdatedName   string `name:"updated_name" db:"updated_name" json:"updated_name" rule:"none"`            //修改人名
}

type AdminData struct {
	D []Admin `json:"d"`
	T int64   `json:"t"`
	S uint    `json:"s"`
}

// AdminLoginResp 用户登录返回数据结构体
type AdminLoginResp struct {
	Token     string `json:"token"`     // 用户token
	Allows    string `json:"allows"`    // 用户权限列表
	AdminName string `json:"user_name"` // 用户名
	Domain    string `json:"domain"`    // 用户名
	Prefix    string `json:"prefix"`    //站点前缀
}

func AdminInsert(data Admin) error {

	data.Prefix = meta.Prefix
	v := MurmurHash(data.Pwd, data.CreateAt)

	data.Pwd = fmt.Sprintf("%d", v)
	data.ID = helper.GenId()
	data.LastLoginIP = ""
	data.LastLoginTime = 0
	data.Prefix = meta.Prefix

	query, _, _ := dialect.Insert("tbl_admins").Rows(data).ToSQL()
	_, err := meta.MerchantDB.Exec(query)
	if err != nil {
		body := fmt.Errorf("%s,[%s]", err.Error(), query)
		return pushLog(body, helper.DBErr)
	}

	return nil
}

func AdminUpdateState(ctx *fasthttp.RequestCtx, id, state string) error {

	admin, err := AdminToken(ctx)
	if err != nil {
		return errors.New(helper.AccessTokenExpires)
	}

	record := g.Record{
		"state":        state,
		"updated_at":   ctx.Time().Unix(),
		"updated_uid":  admin["id"],
		"updated_name": admin["name"],
	}
	t := dialect.Update("tbl_admins")
	query, _, _ := t.Set(record).Where(g.Ex{"id": id}).ToSQL()
	_, err = meta.MerchantDB.Exec(query)
	if err != nil {
		body := fmt.Errorf("%s,[%s]", err.Error(), query)
		return pushLog(body, helper.DBErr)
	}

	return nil
}

func AdminUpdate(id, pwd string, recs g.Record) error {

	var createAt uint32
	ex := g.Ex{
		"id": id,
	}
	query, _, _ := dialect.From("tbl_admins").Select("create_at").Where(ex).ToSQL()
	err := meta.MerchantDB.Get(&createAt, query)
	if err != nil {
		body := fmt.Errorf("%s,[%s]", err.Error(), query)
		return pushLog(body, helper.DBErr)
	}

	if pwd != "" {
		v := MurmurHash(pwd, createAt)
		recs["password"] = fmt.Sprintf("%d", v)
	}

	query, _, _ = dialect.Update("tbl_admins").Set(recs).Where(g.Ex{"id": id}).ToSQL()
	_, err = meta.MerchantDB.Exec(query)
	if err != nil {
		body := fmt.Errorf("%s,[%s]", err.Error(), query)
		return pushLog(body, helper.DBErr)
	}

	return nil
}

func AdminList(adminGid string, page, pageSize uint, ex g.Ex) (AdminData, error) {

	data := AdminData{}
	gids, gidMap, err := groupSubList(adminGid)
	if err != nil {
		return data, err
	}

	if len(gids) == 0 {
		return data, nil
	}

	ex["prefix"] = meta.Prefix
	if gid, ok := ex["group_id"].(string); ok {
		if _, ok = gidMap[gid]; !ok {
			return data, errors.New(helper.MethodNoPermission)
		}
	} else {
		ex["group_id"] = gids
	}

	t := dialect.From("tbl_admins")
	if page == 1 {
		query, _, _ := t.Select(g.COUNT(1)).Where(ex).ToSQL()
		err := meta.MerchantDB.Get(&data.T, query)
		if err != nil {
			body := fmt.Errorf("%s,[%s]", err.Error(), query)
			return data, pushLog(body, helper.DBErr)
		}

		if data.T == 0 {
			return data, nil
		}
	}

	offset := (page - 1) * pageSize
	query, _, _ := t.Select(colsAdmin...).Where(ex).Offset(offset).Limit(pageSize).Order(g.C("create_at").Desc()).ToSQL()
	fmt.Println(query)
	err = meta.MerchantDB.Select(&data.D, query)
	if err != nil {
		body := fmt.Errorf("%s,[%s]", err.Error(), query)
		return data, pushLog(body, helper.DBErr)
	}

	data.S = pageSize
	return data, nil

}

func AdminToken(ctx *fasthttp.RequestCtx) (map[string]string, error) {

	b := ctx.UserValue("token").([]byte)

	var p fastjson.Parser

	data := map[string]string{}
	v, err := p.ParseBytes(b)
	if err != nil {
		return data, err
	}

	o, err := v.Object()
	if err != nil {
		return data, err
	}

	o.Visit(func(k []byte, v *fastjson.Value) {
		key := string(k)
		val, err := v.StringBytes()
		if err == nil {
			data[key] = string(val)
		}
	})

	return data, nil
}

func AdminLogin(deviceNo, username, password, seamo, ip string, lastLoginTime uint32) (AdminLoginResp, error) {

	rsp := AdminLoginResp{
		Prefix: meta.Prefix,
	}
	data := Admin{}
	t := dialect.From("tbl_admins")
	query, _, _ := t.Select(colsAdmin...).Where(g.Ex{"name": username, "prefix": meta.Prefix}).Limit(1).ToSQL()
	err := meta.MerchantDB.Get(&data, query)
	if err != nil && err != sql.ErrNoRows {
		body := fmt.Errorf("%s,[%s]", err.Error(), query)
		return rsp, pushLog(body, helper.DBErr)
	}

	// 账号不存在提示
	if err == sql.ErrNoRows {
		return rsp, errors.New(helper.UserNotExist)
	}

	if data.State == 0 {
		return rsp, errors.New(helper.Blocked)
	}

	if !meta.IsDev {
		slat := helper.TOTP(data.Seamo, otpTimeout)
		if s, err := strconv.Atoi(seamo); err != nil || s != slat {
			fmt.Println(err, s, slat)
			return rsp, errors.New(helper.DynamicVerificationCodeErr)
		}
	}

	pwd := fmt.Sprintf("%d", MurmurHash(password, data.CreateAt))
	if pwd != data.Pwd {
		return rsp, errors.New(helper.UsernameOrPasswordErr)
	}

	record := g.Record{
		"last_login_ip":   ip,
		"last_login_time": lastLoginTime,
	}
	query, _, _ = dialect.Update("tbl_admins").Set(record).Where(g.Ex{"id": data.ID}).ToSQL()
	_, err = meta.MerchantDB.Exec(query)
	if err != nil {
		body := fmt.Errorf("%s,[%s]", err.Error(), query)
		return rsp, pushLog(body, helper.DBErr)
	}

	permission := "111"
	token := fmt.Sprintf("%s:a:token:%s", meta.Prefix, data.ID)
	b, _ := helper.JsonMarshal(data)
	sid, err := session.AdminSet(b, token, deviceNo)
	if err != nil {
		fmt.Println("AdminSet = ", err)
		return rsp, errors.New(helper.SessionErr)
	}

	rsp.Token = sid
	rsp.AdminName = username
	rsp.Allows = permission
	rsp.Domain = meta.GcsDoamin

	/*
		adminLoginLog := map[string]string{
			"uid":      data.ID,
			"username": username,
			"ip":       ip,
		}
		tdlog.Login(adminLoginLog)
		// 写入登录日志
		log := adminLoginLogBase{
			UID:       data.ID,
			Name:      username,
			IP:        ip,
			Flag:      1,
			CreatedAt: lastLoginTime,
		}

		err = meta.Zlog.Post(esPrefixIndex("admin_login_log"), log)
		if err != nil {
			fmt.Printf("user %s login_log wirte err = %s", username, err.Error())
		}
	*/
	return rsp, nil
}

// 检测管理员账号是否已存在
func AdminExist(ex g.Ex) bool {

	var id string
	ex["prefix"] = meta.Prefix
	t := dialect.From("tbl_admins")
	query, _, _ := t.Select("id").Where(ex).Limit(1).ToSQL()
	err := meta.MerchantDB.Get(&id, query)
	return err != sql.ErrNoRows
}

// 后台管理员退出登录
func AdminLogout(ctx *fasthttp.RequestCtx) {
	/*
		admin, err := AdminToken(ctx)
		if err != nil {
			fmt.Println("admin_logout_err: not fount admin:", err.Error())
			return
		}

			// 写入登出日志
			log := adminLoginLogBase{
				UID:       admin["id"],
				Name:      admin["name"],
				IP:        helper.FromRequest(ctx),
				Flag:      2,
				CreatedAt: uint32(time.Now().Unix()),
			}

			err = meta.Zlog.Post(esPrefixIndex("admin_login_log"), log)
			if err != nil {
				fmt.Printf("user %s login_log wirte err = %s", admin["name"], err.Error())
				return
			}
	*/
}
