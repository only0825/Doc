package router

import (
	"fmt"
	"merchant/controller"
	"runtime/debug"
	"time"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

var (
	ApiTimeoutMsg = `{"status": "false","data":"服务器响应超时，请稍后重试"}`
	ApiTimeout    = time.Second * 30
	route         *router.Router
	buildInfo     BuildInfo
)

type BuildInfo struct {
	GitReversion   string
	BuildTime      string
	BuildGoVersion string
}

func apiServerPanic(ctx *fasthttp.RequestCtx, rcv interface{}) {

	err := rcv.(error)
	fmt.Println(err)
	debug.PrintStack()

	if r := recover(); r != nil {
		fmt.Println("recovered failed", r)
	}

	ctx.SetStatusCode(500)
	return
}

func Version(ctx *fasthttp.RequestCtx) {

	ctx.SetContentType("text/html; charset=utf-8")
	fmt.Fprintf(ctx, "merchant<br />Git reversion = %s<br />Build Time = %s<br />Go version = %s<br />System Time = %s<br />",
		buildInfo.GitReversion, buildInfo.BuildTime, buildInfo.BuildGoVersion, ctx.Time())

	//ctx.Request.Header.VisitAll(func (key, value []byte) {
	//	fmt.Fprintf(ctx, "%s: %s<br/>", string(key), string(value))
	//})
}

// SetupRouter 设置路由列表
func SetupRouter(b BuildInfo) *router.Router {

	route = router.New()
	route.PanicHandler = apiServerPanic

	buildInfo = b

	groupCtl := new(controller.GroupController)
	// 权限管理
	privCtl := new(controller.PrivController)
	// 用户管理
	adminCtl := new(controller.AdminController)
	// tree
	treeCtl := new(controller.TreeController)
	//ip区域
	areaCtl := new(controller.AreaController)
	//场馆列表
	platformCtl := new(controller.PlatformController)
	// 游戏列表
	slotsCtl := new(controller.SlotsController)
	//会员
	memberCtl := new(controller.MemberController)
	//会员转代
	memberTransferCtl := new(controller.MemberTransferController)
	//短链域名设置
	shortURLCtl := new(controller.ShortURLController)
	//会员等级设置
	levelCtl := new(controller.MemberLevelController)
	// 账户调整
	adjustCtl := new(controller.AdjustController)
	//银行卡管理
	bankCtl := new(controller.BankcardController)
	//红利
	dividendCtl := new(controller.DividendController)
	// 标签管理
	tagsCtl := new(controller.TagsController)
	// 广告
	bannerCtl := new(controller.BannerController)
	// app 版本管理
	appUpCtl := new(controller.AppUpgradeController)
	//黑名单管理
	blacklistCtl := new(controller.BlacklistController)
	//流水稽查管理
	inspectionCtl := new(controller.InspectionController)
	// 日志管理
	logCtl := new(controller.LogController)
	// 返水管理
	rebateCtl := new(controller.RebateController)
	//公告管理
	noticeCtl := new(controller.NoticeController)
	//转账，游戏，账变记录
	recordCtl := new(controller.RecordController)
	//内容管理
	msgCtl := new(controller.MessageController)
	//短信验证码管理
	smsCtl := new(controller.SmsRecordController)
	//邮箱验证码管理
	mailCtl := new(controller.EmailRecordController)
	// 短信通道管理
	smsChannelCtl := new(controller.SMSChannelController)
	// 短信通道管理
	linkCtl := new(controller.LinkController)

	get("/merchant/version", Version)

	// 代理管理-推广链接-查询
	get("/merchant/link/list", linkCtl.List)
	// 代理管理-推广链接-设置是否显示广告页
	post("/merchant/link/set", linkCtl.SetNoAd)
	// 代理管理-推广链接-删除
	get("/merchant/link/delete", linkCtl.Delete)

	// 权限管理-用户组管理-新增分组
	post("/merchant/group/insert", groupCtl.Insert)
	// 权限管理-用户组管理-修改分组
	post("/merchant/group/update", groupCtl.Update)
	// 权限管理-用户组管理列表
	get("/merchant/group/list", groupCtl.List)
	// 权限管理-获取分组列表
	get("/merchant/priv/list", privCtl.List)

	// 后台用户登录
	post("/merchant/admin/login", adminCtl.Login)
	// 权限管理-后台账号管理-后台账号列表
	get("/merchant/admin/list", adminCtl.List)
	// 权限管理-后台账号管理-新建账号
	post("/merchant/admin/insert", adminCtl.Insert)
	// 权限管理-后台账号管理-编辑更新账号
	post("/merchant/admin/update", adminCtl.Update)
	// 权限管理-后台账号管理-启用
	// 权限管理-后台账号管理-禁用
	get("/merchant/admin/update/state", adminCtl.UpdateState)

	//获取静态json
	get("/merchant/tree", treeCtl.List)
	//ip查询
	post("/merchant/area/view", areaCtl.View)

	// 站点管理-场馆管理
	get("/merchant/platform/list", platformCtl.List)
	// 站点管理-场馆开启/关闭 维护/解除维护 修改排序
	post("/merchant/platform/update", platformCtl.Update)

	// 站点管理-场馆管理-游戏列表
	post("/merchant/slots/list", slotsCtl.List)
	// 站点管理-场馆管理-游戏列表-编辑游戏
	post("/merchant/slots/update", slotsCtl.Update)
	// 站点管理-场馆管理-游戏列表-上线下线
	post("/merchant/slots/updatestate", slotsCtl.UpdateState)

	// 会员管理-会员列表-新增总代
	post("/merchant/member/insert", memberCtl.Insert)
	// 会员管理-会员列表-帐户信息
	get("/merchant/member/info", memberCtl.AccountInfo)
	// 会员管理-会员中心钱包余额
	get("/merchant/member/balance", memberCtl.Balance)
	//批量获取用户余额
	post("/merchant/member/balancebatch", memberCtl.BalanceBatch)
	//批量获取用户标签
	post("/merchant/member/tagsbatch", memberCtl.TagBatch)
	// 会员管理-会员列表-修改会员状态 (批量禁用启用)
	post("/merchant/member/updatestate", memberCtl.UpdateState)
	// 会员管理-会员列表-修改会员信息
	post("/merchant/member/update", memberCtl.Update)
	// 会员管理-会员列表
	post("/merchant/member/list", memberCtl.List)
	// 会员管理-会员列表-用户标签
	get("/merchant/member/tags", memberCtl.Tags)
	// 会员管理-会员列表-批量添加标签/编辑标签
	post("/merchant/member/settags", memberCtl.SetTags)
	//会员列表-修改会员密码
	post("/merchant/member/updatepwd", memberCtl.UpdatePwd)
	// 会员管理-会员列表-批量取消标签
	post("/merchant/member/tags/cancel", memberCtl.CancelTags)
	// 会员管理-会员列表-解除密码限制/接触短信限制/场馆钱包限制
	get("/merchant/member/retry/reset", memberCtl.RetryReset)
	// 会员管理-会员列表-数据概览
	get("/merchant/member/overview", memberCtl.Overview)
	// 会员管理-会员列表-会员详情
	get("/merchant/member/detail", memberCtl.Detail)

	// 会员管理-会员列表-添加备注
	post("/merchant/member/remark/insert", memberCtl.RemarkLogInsert)
	// 会员管理-会员列表-基本信息-备注信息
	get("/merchant/member/remark/list", memberCtl.RemarkLogList)
	// 会员管理-会员列表-修改备注
	post("/merchant/member/remark/update", memberCtl.RemarkLogUpdate)
	// 会员管理-会员列表-删除备注
	post("/merchant/member/remark/delete", memberCtl.RemarkLogDelete)

	// 查询用户真实姓名/邮箱/手机号/银行卡号修改历史
	post("/merchant/member/history", memberCtl.History)
	// 查询用户真实姓名/邮箱/手机号/银行卡号修改历史
	post("/merchant/member/history/{field}", memberCtl.HistoryField)
	// 查询用户真实姓名/邮箱/手机号/银行卡号明文信息
	post("/merchant/member/full", memberCtl.Full)
	// 查询用户真实姓名/邮箱/手机号/银行卡号明文信息
	post("/merchant/member/full/{field}", memberCtl.FullField)
	// 跳线转代
	post("/merchant/member/transfer", memberTransferCtl.Transfer)
	// 钱包余额冲正
	post("/merchant/member/clear", memberCtl.SetBalanceZero)

	// 代理管理-代理列表
	post("/merchant/agency/list", memberCtl.Agency)
	// 代理管理-代理修改密码/开启关闭
	post("/merchant/agency/update", memberCtl.UpdateTopMember)
	// 代理管理-代理修改返水比例
	post("/merchant/agency/updaterebate", memberCtl.UpdateMemberRebate)
	// 代理管理-代理编辑维护人
	post("/merchant/agency/updatemaintain", memberCtl.UpdateMaintainName)
	// 代理管理-下级成员
	post("/merchant/agency/memberlist", memberCtl.MemberList)

	// 代理管理-团队转代-列表
	post("/merchant/agency/transfer/list", memberTransferCtl.List)
	// 代理管理-团队转代-新增
	post("/merchant/agency/transfer/insert", memberTransferCtl.Insert)
	// 代理管理-团队转代-审核
	post("/merchant/agency/transfer/review", memberTransferCtl.Review)
	// 代理管理-团队转代-删除
	get("/merchant/agency/transfer/delete", memberTransferCtl.Delete)

	//代理管理-记录管理-游戏
	post("/merchant/agency/record/game", recordCtl.Game)
	//代理管理-记录管理-登录
	post("/merchant/agency/record/loginlog", recordCtl.LoginLog)
	//代理管理-记录管理-会员存款
	post("/merchant/agency/record/deposit", recordCtl.Deposit)
	//代理管理-记录管理-红利
	post("/merchant/agency/record/dividend", recordCtl.Dividend)
	//代理管理-记录管理-返水
	post("/merchant/agency/record/rebate", recordCtl.Rebate)
	//代理管理-记录管理-账户调整
	post("/merchant/agency/record/adjust", recordCtl.Adjust)
	//代理管理-记录管理-提款
	post("/merchant/agency/record/withdraw", recordCtl.Withdraw)
	//代理管理-团队管理-操作记录
	post("/merchant/agency/record/group", recordCtl.Group)
	//聊天室-随机计划-查看统计-明细
	get("/merchant/record/plan/list", recordCtl.Order)
	//用户VIP配置
	get("/merchant/vip/list", levelCtl.List)
	//用户VIP设置
	post("/merchant/vip/update", levelCtl.Update)
	//新增VIP
	post("/merchant/vip/insert", levelCtl.Insert)
	//用户VIP调整记录
	get("/merchant/vip/record", levelCtl.Record)

	// 会员列表-账户调整
	post("/merchant/adjust/insert", adjustCtl.Insert)
	// 会员管理-账户调整审核列表
	get("/merchant/adjust/list", adjustCtl.List)
	// 后台会员管理-账户调整审核
	post("/merchant/adjust/review", adjustCtl.Review)

	//新增银行卡
	post("/merchant/bankcard/insert", bankCtl.Insert)
	//查询银行卡
	post("/merchant/bankcard/list", bankCtl.List)
	//修改银行卡信息
	post("/merchant/bankcard/update", bankCtl.Update)
	//删除银行卡
	get("/merchant/bankcard/delete", bankCtl.Delete)
	// 后台管理-银行卡校验日志记录查询
	get("/merchant/bankcard/log/list", bankCtl.Log)

	// 运营管理-红利管理-单会员发放
	post("/merchant/dividend/insert", dividendCtl.Insert)
	// 运营管理-红利管理-审核列表
	post("/merchant/dividend/list", dividendCtl.List)
	// 运营管理-会员管理-基本信息-红利列表
	post("/merchant/dividend/member/list", dividendCtl.MemberList)
	// 运营管理-红利管理-更新
	post("/merchant/dividend/update", dividendCtl.Update)
	// 运营管理-红利管理-修改审核备注
	post("/merchant/dividend/review", dividendCtl.ReviewRemarkUpdate)

	// 会员管理-会员配置-标签管理-新增
	post("/merchant/tags/insert", tagsCtl.Insert)
	// 会员管理-会员配置-标签管理-列表
	post("/merchant/tags/list", tagsCtl.List)
	// 会员管理-会员配置-标签管理-修改
	post("/merchant/tags/update", tagsCtl.Update)
	// 会员管理-会员配置-标签管理-删除
	get("/merchant/tags/delete", tagsCtl.Delete)

	// 运营管理-广告管理-新增
	post("/merchant/banner/insert", bannerCtl.Insert)
	// 运营管理-广告管理-修改
	post("/merchant/banner/update", bannerCtl.Update)
	// 运营管理-广告管理-列表
	post("/merchant/banner/list", bannerCtl.List)
	// 运营管理-广告管理-删除
	get("/merchant/banner/delete", bannerCtl.Delete)
	// 运营管理-广告管理-状态(启用|停用)
	post("/merchant/banner/updatestate", bannerCtl.UpdateState)

	// App 升级配置更新
	post("/merchant/app/update", appUpCtl.Update)
	// App 升级配置列表
	get("/merchant/app/list", appUpCtl.List)

	// 短链接-域名设置
	post("/merchant/shorturl/set", shortURLCtl.Set)
	// 短链接-域名查询
	get("/merchant/shorturl/get", shortURLCtl.Get)
	//查询会员登录日志
	get("/merchant/blacklist/loginlog", blacklistCtl.LogList)
	//查询登陆日志会员信息
	get("/merchant/blacklist/assoclog", blacklistCtl.AssociateList)

	//风控管理-黑名单查询
	get("/merchant/blacklist/list", blacklistCtl.List)
	//风控管理-黑名单新增
	post("/merchant/blacklist/insert", blacklistCtl.Insert)
	//风控管理-黑名单修改
	post("/merchant/blacklist/update", blacklistCtl.Update)
	//风控管理-黑名单删除
	get("/merchant/blacklist/delete", blacklistCtl.Delete)
	//风控管理-黑名单-手机号码解锁
	post("/merchant/blacklist/clear/phone", blacklistCtl.ClearPhone)

	//风控管理-流水稽查查询
	get("/merchant/inspection/list", inspectionCtl.List)
	//风控管理-流水稽查审核
	post("/merchant/inspection/review", inspectionCtl.Review)
	//风控管理-流水稽查记录
	post("/merchant/inspection/history", inspectionCtl.History)

	// 系统管理-日志管理-登录日志
	get("/merchant/sys/log/login/list", logCtl.AdminLoginLog)
	// 系统管理-日志管理-系统日志
	post("/merchant/sys/log/system/list", logCtl.SystemLog)

	//运营管理-系统公告-增加
	post("/merchant/notice/insert", noticeCtl.Insert)
	//运营管理-系统公告-列表
	post("/merchant/notice/list", noticeCtl.List)
	//运营管理-系统公告-编辑
	post("/merchant/notice/update", noticeCtl.Update)
	//运营管理-系统公告-停用/启用
	post("/merchant/notice/updatestate", noticeCtl.UpdateState)
	//运营管理-系统公告-删除
	get("/merchant/notice/delete", noticeCtl.Delete)

	// 运营管理-短信通道-列表
	post("/merchant/sms/list", smsChannelCtl.List)
	// 运营管理-短信通道-增加
	post("/merchant/sms/insert", smsChannelCtl.Insert)
	// 运营管理-短信通道-编辑
	post("/merchant/sms/update", smsChannelCtl.Update)
	// 运营管理-短信通道-更新短信通道状态
	post("/merchant/sms/update/state", smsChannelCtl.UpdateState)
	// 运营管理-验证码查询-短信验证码列表
	get("/merchant/member/sms/list", smsCtl.List)
	// 运营管理-验证码查询-邮箱验证码列表
	get("/merchant/member/mail/list", mailCtl.List)

	// 获取返水上限
	get("/merchant/rebate/scale", rebateCtl.Scale)
	// 是否允许编辑返水比例
	get("/merchant/rebate/enablemod", rebateCtl.EnableMod)

	// 会员管理-会员列表-账变记录
	get("/merchant/record/transaction", recordCtl.Transaction)
	// 记录管理-平台转帐
	post("/merchant/record/transfer", recordCtl.Transfer)
	// 会员管理-会员列表-有效投注查询
	// 会员管理-投注管理
	// 会员管理-投注管理-会员游戏记录详情列表
	// 会员列表->输赢信息
	// 会员列表->输赢信息->场馆合并
	// 会员列表->输赢信息->日期合并
	post("/merchant/record/game", recordCtl.RecordGame)
	//会员管理-会员列表-详情-返水记录
	get("/merchant/member/rebate", rebateCtl.RebatePersonal)
	//运营管理-站内消息-添加
	post("/merchant/message/insert", msgCtl.Insert)
	//运营管理-站内消息-列表
	post("/merchant/message/list", msgCtl.List)
	//运营管理-站内消息-编辑
	post("/merchant/message/update", msgCtl.Update)
	//运营管理-站内消息-审核
	post("/merchant/message/review", msgCtl.Review)
	//运营管理-站内消息-已发站内信详情
	get("/merchant/message/detail", msgCtl.Detail)
	//运营管理-站内消息-已发系统站内信列表
	post("/merchant/message/system", msgCtl.System)
	//运营管理-站内消息-删除
	post("/merchant/message/delete", msgCtl.Delete)

	return route
}

// get is a shortcut for router.GET(path string, handle fasthttp.RequestHandler)
func get(path string, handle fasthttp.RequestHandler) {
	route.GET(path, fasthttp.TimeoutHandler(handle, ApiTimeout, ApiTimeoutMsg))
}

// head is a shortcut for router.HEAD(path string, handle fasthttp.RequestHandler)
func head(path string, handle fasthttp.RequestHandler) {
	route.HEAD(path, fasthttp.TimeoutHandler(handle, ApiTimeout, ApiTimeoutMsg))
}

// options is a shortcut for router.OPTIONS(path string, handle fasthttp.RequestHandler)
func options(path string, handle fasthttp.RequestHandler) {
	route.OPTIONS(path, fasthttp.TimeoutHandler(handle, ApiTimeout, ApiTimeoutMsg))
}

// post is a shortcut for router.POST(path string, handle fasthttp.RequestHandler)
func post(path string, handle fasthttp.RequestHandler) {
	route.POST(path, fasthttp.TimeoutHandler(handle, ApiTimeout, ApiTimeoutMsg))
}

// put is a shortcut for router.PUT(path string, handle fasthttp.RequestHandler)
func put(path string, handle fasthttp.RequestHandler) {
	route.PUT(path, fasthttp.TimeoutHandler(handle, ApiTimeout, ApiTimeoutMsg))
}

// patch is a shortcut for router.PATCH(path string, handle fasthttp.RequestHandler)
func patch(path string, handle fasthttp.RequestHandler) {
	route.PATCH(path, fasthttp.TimeoutHandler(handle, ApiTimeout, ApiTimeoutMsg))
}

// delete is a shortcut for router.DELETE(path string, handle fasthttp.RequestHandler)
func delete(path string, handle fasthttp.RequestHandler) {
	route.DELETE(path, fasthttp.TimeoutHandler(handle, ApiTimeout, ApiTimeoutMsg))
}
