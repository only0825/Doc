package controller

import (
	"merchant/contrib/helper"
	"merchant/contrib/validator"
	"merchant/model"
	"strings"

	g "github.com/doug-martin/goqu/v9"

	"github.com/valyala/fasthttp"
)

type MemberLevelController struct{}

type VipUpdateParam struct {
	ID                string `rule:"digit" name:"id" json:"id"`
	RechargeNum       string `rule:"none" name:"recharge_num" json:"recharge_num"`
	LevelName         string `rule:"none" name:"level_name" json:"level_name"`
	UpgradeDeposit    string `rule:"none" name:"upgrade_deposit" json:"upgrade_deposit"`
	UpgradeRecord     string `rule:"none" name:"upgrade_record" json:"upgrade_record"`
	RelegationFlowing string `rule:"none" name:"relegation_flowing" json:"relegation_flowing"`
	UpgradeGift       string `rule:"none" name:"upgrade_gift" json:"upgrade_gift"`
	BirthGift         string `rule:"none" name:"birth_gift" json:"birth_gift"`
	WithdrawCount     string `rule:"none" name:"withdraw_count" json:"withdraw_count"`
	WithdrawMax       string `rule:"none" name:"withdraw_max" json:"withdraw_max"`
	EarlyMonthPacket  string `rule:"none" name:"early_month_packet" json:"early_month_packet"`
	LateMonthPacket   string `rule:"none" name:"late_month_packet" json:"late_month_packet"`
}

// 会员等级记录
type vipRecordListParam struct {
	UserName    string `rule:"none" name:"username"`     //会员账号
	Ty          string `rule:"none" name:"ty"`           //调整类型
	CreatedName string `rule:"none" name:"created_name"` // 操作人名
	Level       string `rule:"none" name:"level"`
	StartTime   string `rule:"none" name:"start_time" msg:"start_time error"` // 创建时间 开始时间
	EndTime     string `rule:"none" name:"end_time" msg:"end_time"`           // 结束时间
	Page        uint   `rule:"digit" name:"page"  default:"1" msg:"page error"`
	PageSize    uint   `rule:"digit" name:"page_size" default:"15" msg:"page_size error"`
}

func (that *MemberLevelController) List(ctx *fasthttp.RequestCtx) {

	data, err := model.MemberLevelList()
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}

func (that *MemberLevelController) Update(ctx *fasthttp.RequestCtx) {

	vip := VipUpdateParam{}
	err := validator.Bind(ctx, &vip)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	record := g.Record{
		"recharge_num":       vip.RechargeNum,
		"level_name":         vip.LevelName,
		"upgrade_deposit":    vip.UpgradeDeposit,
		"upgrade_record":     vip.UpgradeRecord,
		"relegation_flowing": vip.RelegationFlowing,
		"upgrade_gift":       vip.UpgradeGift,
		"birth_gift":         vip.BirthGift,
		"withdraw_count":     vip.WithdrawCount,
		"withdraw_max":       vip.WithdrawMax,
		"early_month_packet": vip.EarlyMonthPacket,
		"late_month_packet":  vip.LateMonthPacket,
	}
	// 更新VIP信息
	err = model.VIPUpdate(vip.ID, record)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

func (that *MemberLevelController) Insert(ctx *fasthttp.RequestCtx) {

	vip := model.MemberLevel{}
	err := validator.Bind(ctx, &vip)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	vip.CreateAt = uint32(ctx.Time().Unix())
	vip.UpdatedAt = uint32(ctx.Time().Unix())

	// 更新权限信息
	err = model.VIPInsert(vip)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, helper.Success)
}

// 会员vip等级历史记录
func (that *MemberLevelController) Record(ctx *fasthttp.RequestCtx) {

	param := vipRecordListParam{}
	err := validator.Bind(ctx, &param)
	if err != nil {
		helper.Print(ctx, false, helper.ParamErr)
		return
	}

	ex := g.Ex{}
	if param.UserName != "" {
		//字母数字组合，4-9，前2个字符必须为字母
		param.UserName = strings.ToLower(param.UserName)
		if !validator.CheckUName(param.UserName, 5, 14) {
			helper.Print(ctx, false, helper.UsernameErr)
			return
		}
		ex["username"] = param.UserName
	} else {
		if param.Ty != "" {

			if !validator.CheckIntScope(param.Ty, model.MemberLevelUpgrade, model.MemberLevelRecover) { // 全部类型
				helper.Print(ctx, false, helper.MemberLevelAlterTyErr)
				return
			}

			ex["ty"] = param.Ty
		}

		if param.CreatedName != "" {
			if !validator.CtypeAlnum(param.CreatedName) {
				helper.Print(ctx, false, helper.CreateNameErr)
				return
			}

			ex["created_name"] = param.CreatedName
		}

		if param.Level != "" {
			if !validator.CheckIntScope(param.Level, 1, 11) {
				helper.Print(ctx, false, helper.MemberLevelErr)
				return
			}

			ex["after_level"] = param.Level
		}
	}
	data, err := model.VipRecord(param.Page, param.PageSize, param.StartTime, param.EndTime, ex)
	if err != nil {
		helper.Print(ctx, false, err.Error())
		return
	}

	helper.Print(ctx, true, data)
}
