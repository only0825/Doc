package model

import (
	"errors"
	"merchant/contrib/helper"
)

type tbl_members_t struct {
	Zalo                string `db:"zalo" json:"zalo"`         // 会员名
	RealName            string `db:"realname" json:"realname"` // 会员名
	Phone               string `db:"phone" json:"phone"`       // 会员名
	Email               string `db:"email" json:"email"`       // 会员名
	UID                 string `db:"uid" json:"uid" redis:"uid"`
	Username            string `db:"username" json:"username" redis:"username"`                                     //会员名
	Password            string `db:"password" json:"password" redis:"password"`                                     //密码
	Birth               string `db:"birth" json:"birth" redis:"birth"`                                              //生日日期
	BirthHash           string `db:"birth_hash" json:"birth_hash" redis:"birth_hash"`                               //生日日期哈希
	RealnameHash        string `db:"realname_hash" json:"realname_hash" redis:"realname_hash"`                      //真实姓名哈希
	EmailHash           string `db:"email_hash" json:"email_hash" redis:"email_hash"`                               //邮件地址哈希
	PhoneHash           string `db:"phone_hash" json:"phone_hash" redis:"phone_hash"`                               //电话号码哈希
	ZaloHash            string `db:"zalo_hash" json:"zalo_hash" redis:"zalo_hash"`                                  //zalo哈希
	Prefix              string `db:"prefix" json:"prefix" redis:"prefix"`                                           //站点前缀
	Tester              string `db:"tester" json:"tester" redis:"tester"`                                           //1正式 0测试
	WithdrawPwd         uint64 `db:"withdraw_pwd" json:"withdraw_pwd" redis:"withdraw_pwd"`                         //取款密码哈希
	Regip               string `db:"regip" json:"regip" redis:"regip"`                                              //注册IP
	RegDevice           string `db:"reg_device" json:"reg_device" redis:"reg_device"`                               //注册设备号
	RegUrl              string `db:"reg_url" json:"reg_url" redis:"reg_url"`                                        //注册链接
	CreatedAt           uint32 `db:"created_at" json:"created_at" redis:"created_at"`                               //注册时间
	LastLoginIp         string `db:"last_login_ip" json:"last_login_ip" redis:"last_login_ip"`                      //最后登陆ip
	LastLoginAt         uint32 `db:"last_login_at" json:"last_login_at" redis:"last_login_at"`                      //最后登陆时间
	SourceId            uint8  `db:"source_id" json:"source_id" redis:"source_id"`                                  //注册来源 1 pc 2h5 3 app
	FirstDepositAt      uint32 `db:"first_deposit_at" json:"first_deposit_at" redis:"first_deposit_at"`             //首充时间
	FirstDepositAmount  string `db:"first_deposit_amount" json:"first_deposit_amount" redis:"first_deposit_amount"` //首充金额
	FirstBetAt          uint32 `db:"first_bet_at" json:"first_bet_at" redis:"first_bet_at"`                         //首投时间
	FirstBetAmount      string `db:"first_bet_amount" json:"first_bet_amount" redis:"first_bet_amount"`             //首投金额
	SecondDepositAt     uint32 `db:"second_deposit_at" json:"second_deposit_at"`                                    //二存时间
	SecondDepositAmount string `db:"second_deposit_amount" json:"second_deposit_amount"`                            //二充金额
	TopUid              string `db:"top_uid" json:"top_uid" redis:"top_uid"`                                        //总代uid
	TopName             string `db:"top_name" json:"top_name" redis:"top_name"`                                     //总代代理
	ParentUid           string `db:"parent_uid" json:"parent_uid" redis:"parent_uid"`                               //上级uid
	ParentName          string `db:"parent_name" json:"parent_name" redis:"parent_name"`                            //上级代理
	BankcardTotal       uint8  `db:"bankcard_total" json:"bankcard_total" redis:"bankcard_total"`                   //用户绑定银行卡的数量
	LastLoginDevice     string `db:"last_login_device" json:"last_login_device" redis:"last_login_device"`          //最后登陆设备
	LastLoginSource     int    `db:"last_login_source" json:"last_login_source" redis:"last_login_source"`          //上次登录设备来源:1=pc,2=h5,3=ios,4=andriod
	Remarks             string `db:"remarks" json:"remarks" redis:"remarks"`                                        //备注
	State               uint8  `db:"state" json:"state" redis:"state"`                                              //状态 1正常 2禁用
	Level               int    `db:"level" json:"level" redis:"level" redis:"level"`                                //等级
	Balance             string `db:"balance" json:"balance" redis:"balance"`                                        //余额
	LockAmount          string `db:"lock_amount" json:"lock_amount" redis:"lock_amount"`                            //锁定金额
	Commission          string `db:"commission" json:"commission" redis:"commission"`                               //佣金
	GroupName           string `db:"group_name" json:"group_name" redis:"group_name"`                               //团队名称
	AgencyType          int64  `db:"agency_type" json:"agency_type" redis:"agency_type"`                            //391团队代理 393普通代理
	Address             string `db:"address" json:"address" redis:"address"`                                        //收货地址
	Avatar              string `db:"avatar" json:"avatar" redis:"avatar"`                                           //收货地址
}

func memberInfoCache(username string) (tbl_members_t, error) {

	m := tbl_members_t{}

	key := meta.Prefix + ":member:" + username

	pipe := meta.MerchantRedis.TxPipeline()
	defer pipe.Close()

	exist := pipe.Exists(ctx, key)
	rs := pipe.HMGet(ctx, key, "uid", "username", "password", "birth", "birth_hash", "realname_hash", "email_hash", "phone_hash", "zalo_hash", "prefix", "tester", "withdraw_pwd", "regip", "reg_device", "reg_url", "created_at", "last_login_ip", "last_login_at", "source_id", "first_deposit_at", "first_deposit_amount", "first_bet_at", "first_bet_amount", "", "", "top_uid", "top_name", "parent_uid", "parent_name", "bankcard_total", "last_login_device", "last_login_source", "remarks", "state", "level", "balance", "lock_amount", "commission", "group_name", "agency_type", "address", "avatar")

	_, err := pipe.Exec(ctx)
	if err != nil {
		return m, pushLog(err, helper.RedisErr)
	}

	if exist.Val() == 0 {
		return m, errors.New(helper.UsernameErr)
	}

	if err = rs.Scan(&m); err != nil {
		return m, pushLog(rs.Err(), helper.RedisErr)
	}

	return m, nil
}

func MemberInfo(username string) (tbl_members_t, error) {

	res, err := memberInfoCache(username)
	if err != nil {
		return res, err
	}

	encRes := []string{}
	if res.RealnameHash != "0" {

		encRes = append(encRes, "realname")
	}
	if res.PhoneHash != "0" {

		encRes = append(encRes, "phone")
	}
	if res.EmailHash != "0" {

		encRes = append(encRes, "email")
	}
	if res.ZaloHash != "0" {

		encRes = append(encRes, "zalo")
	}

	res.Zalo = ""
	res.RealName = ""
	res.Phone = ""
	res.Email = ""

	if len(encRes) > 0 {
		recs, err := grpc_t.Decrypt(res.UID, true, encRes)
		if err != nil {

			//fmt.Println("MemberInfo res.MemberInfos.UID = ", res.MemberInfos.UID)
			//fmt.Println("MemberInfo grpc_t.Decrypt err = ", err.Error())
			return res, errors.New(helper.UpdateRPCErr)
		}

		res.Zalo = recs["zalo"]
		res.RealName = recs["realname"]
		res.Phone = recs["phone"]
		res.Email = recs["email"]
	}

	return res, nil
}
