package model

import "time"

var (
	gameRecordFields = []string{"parent_uid", "parent_name", "top_uid", "top_name", "settle_time", "start_time", "resettle", "presettle", "odds", "handicap", "handicap_type", "game_name", "main_bill_no", "api_bill_no", "api_name", "updated_at", "created_at", "result", "prefix", "play_type", "flag", "valid_bet_amount", "bet_amount", "rebate_amount", "game_type", "bet_time", "net_amount", "uid", "name", "player_name", "api_type", "bill_no", "row_id", "id"}
)

const (
	LOAD_PAGE = 100
)

//动态验证码的有效时间
const (
	otpTimeout = 60
)

var (
	LockTimeout = 20 * time.Second
)

const (
	GameOffline = 0 // 0下线
	GameOnline  = 1 // 1上线
)

const (
	BannerStateWait       = 1 //待发布
	BannerStateProcessing = 2 //进行中
	BannerStateEnd        = 3 //停止
)

//广告显示类型
const (
	BannerShowTypeForever = "1" // 广告显示永久有效
	BannerShowTypeSpecify = "2" // 广告显示指定时间
)

const (
	defaultRedisKeyPrefix = "rlock:"
)

const (
	PlatformFlagEdit   = 0 //场馆编辑
	PlatformFlagState  = 1 //场馆状态
	PlatformFlagWallet = 2 //钱包状态
)

// 中心钱包上下分
const (
	DownPointApply       = 0 //下分申请
	DownPointApplyPass   = 1 //下分申请通过
	DownPointApplyReject = 2 //下分申请拒绝
	UpPointApplyPass     = 3 //上分申请成功
)

//注单Ty类型
const (
	GameTyValid             = 1 //会员有效投注查询
	GameTyRecord            = 2 //投注管理
	GameTyRecordDetail      = 3 //投注管理会员游戏记录详情
	GameMemberWinOrLose     = 4 //会员输赢信息
	GameMemberTransferGroup = 5 //会员输赢信息场馆统计
	GameMemberDayGroup      = 6 //会员输赢信息按天统计
)

var (
	DeviceMap = map[int]bool{
		DeviceTypeWeb:            true, //web
		DeviceTypeH5:             true, //h5
		DeviceTypeAndroidFlutter: true, //android_flutter
		DeviceTypeIOSFlutter:     true, //ios_flutter
	}
)

// 设备端类型
const (
	DeviceTypeWeb            = 24 //web
	DeviceTypeH5             = 25 //h5
	DeviceTypeIOS            = 26 //ios
	DeviceTypeAndroid        = 27 //android
	DeviceTypeIOSSport       = 28 //ios_sport
	DeviceTypeAndroidSport   = 29 //android_sport
	DeviceTypeIOSAgency      = 30 //ios_agency
	DeviceTypeAndroidAgency  = 31 //android_agency
	DeviceTypeWebAgency      = 32 //web_agency
	DeviceTypeH5Agency       = 33 //h5_agency
	DeviceTypeH5Tg           = 34 //h5_tg
	DeviceTypeAndroidFlutter = 35 //android_flutter
	DeviceTypeIOSFlutter     = 36 //ios_flutter
)

// 场馆转账类型
const (
	TransferIn           = 181 //场馆转入
	TransferOut          = 182 //场馆转出
	TransferUpPoint      = 183 //后台场馆上分
	TransferResetBalance = 184 //场馆钱包清零
	TransferDividend     = 185 //场馆红利
)

// 场馆转账状态
const (
	TransferStateFailed        = 191 //场馆转账失败
	TransferStateSuccess       = 192 //场馆转账成功
	TransferStateDealing       = 193 //场馆转账处理中
	TransferStateScriptConfirm = 194 //场馆转账脚本确认中
	TransferStateManualConfirm = 195 //场馆转账人工确认中
)

// 会员等级调整类型
const (
	MemberLevelUpgrade    = 201 //会员升级
	MemberLevelRelegation = 202 //会员保级
	MemberLevelDowngrade  = 203 //会员降级
	MemberLevelRecover    = 204 //会员等级恢复
)

// 红利类型
const (
	DividendSite      = 211 //平台红利(站点)
	DividendUpgrade   = 212 //升级红利
	DividendBirthday  = 213 //生日红利
	DividendMonthly   = 214 //每月红利
	DividendRedPacket = 215 //红包红利
	DividendMaintain  = 216 //维护补偿
	DividendDeposit   = 217 //存款优惠
	DividendPromo     = 218 //活动红利
	DividendInvite    = 219 //推荐红利
	DividendAdjust    = 220 //红利调整
	DividendResetPlat = 221 //场馆余额负数清零
	DividendAgency    = 222 //代理红利
)

// 红利审核状态
const (
	DividendReviewing    = 231 //红利审核中
	DividendReviewPass   = 232 //红利审核通过
	DividendReviewReject = 233 //红利审核不通过
)

// 红利发放状态
const (
	DividendFailed      = 236 //红利发放失败
	DividendSuccess     = 237 //红利发放成功
	DividendPlatDealing = 238 //红利发放场馆处理中
)

// 后台调整类型
const (
	AdjustUpMode   = 251 // 上分
	AdjustDownMode = 252 // 下分
)

// 后台上下分审核状态
const (
	AdjustReviewing    = 256 //后台调整审核中
	AdjustReviewPass   = 257 //后台调整审核通过
	AdjustReviewReject = 258 //后台调整审核不通过
)

// 后台上下分状态
const (
	AdjustFailed      = 261 //上下分失败
	AdjustSuccess     = 262 //上下分成功
	AdjustPlatDealing = 263 //上分场馆处理中
)

// 存款状态
const (
	DepositConfirming = 361 //确认中
	DepositSuccess    = 362 //存款成功
	DepositCancelled  = 363 //存款已取消
	DepositReviewing  = 364 //存款审核中
)

// 取款状态
const (
	WithdrawReviewing     = 371 //审核中
	WithdrawReviewReject  = 372 //审核拒绝
	WithdrawDealing       = 373 //出款中
	WithdrawSuccess       = 374 //提款成功
	WithdrawFailed        = 375 //出款失败
	WithdrawAbnormal      = 376 //异常订单
	WithdrawAutoPayFailed = 377 // 代付失败
	WithdrawHangup        = 378 // 挂起
	WithdrawDispatched    = 379 // 已派单
)

// 活动状态变更类型
const (
	PromoAddWaterFlow = 461 //增加流水
	PromoUnlock       = 462 //解锁
)

// 活动解锁类型
const (
	PromoUnlockReachWaterFlow = 466 //流水达标解锁
	PromoUnlockBalance        = 467 //余额解锁
	PromoUnlockCancelled      = 468 //取消解锁
	PromoUnlockAdmin          = 469 //管理员解锁
)

var (
	BlackTy = map[int]bool{
		TyDevice:         true, //设备号
		TyIP:             true, //ip地址
		TyEmail:          true, //邮箱地址
		TyPhone:          true, //电话号码
		TyBankcard:       true, //银行卡
		TyVirtualAccount: true, //虚拟币地址
		TyRebate:         true, //返水
		TyCGRebate:       true, //返点
		TyPromoteLink:    true, //推广链接
		TyWhiteIP:        true, //后台ip白名单
	}
)

// 黑名单类型
const (
	TyDevice         = 1   //设备号
	TyIP             = 2   //ip地址
	TyEmail          = 3   //邮箱地址
	TyPhone          = 4   //电话号码
	TyBankcard       = 5   //银行卡
	TyVirtualAccount = 6   //虚拟币地址
	TyRebate         = 7   //场馆返水
	TyCGRebate       = 8   //cg彩票返点
	TyPromoteLink    = 9   //推广链接
	TyWhiteIP        = 101 //后台访问ip白名单
)

const (
	// 后台调整审核
	adjustReviewFmt = `{
  "cn": {
    "title": "账户调整审核",
    "content": "用户 %s，发起账户调整，%s 调整 %s KVND，请尽快审核。",
    "url": "/vip/acc"
  },
  "en": {
    "title": "Xét duyệt điều chỉnh tài khoản",
    "content": "Người dùng %s, phát điều chỉnh tài khoản, %s điều chỉnh %s KVND, vui lòng nhanh chóng xét duyệt.",
    "url": "/vip/acc"
  },
  "vn": {
     "title": "Xét duyệt điều chỉnh tài khoản",
    "content": "Người dùng %s, phát điều chỉnh tài khoản, %s điều chỉnh %s KVND, vui lòng nhanh chóng xét duyệt.",
    "url": "/vip/acc"
  }
}`
	// 红利审核
	dividendReviewFmt = `{
  "cn": {
    "title": "红利调整审核",
    "content": "用户 %s，发起红利调整，%s 调整 %s KVND，请尽快审核。",
    "url": "/operation/bonusManager?name=bonusReview"
  },
  "en": {
    "title": "Xét duyệt điều chỉnh hoa hồng",
    "content": "Người dùng %s， phát điều chỉnh hoa hồng, %s điều chỉnh %s KVND, vui lòng nhanh chóng xét duyệt.",
    "url": "/operation/bonusManager?name=bonusReview"
  },
  "vn": {
     "title": "Xét duyệt điều chỉnh hoa hồng",
    "content": "Người dùng %s， phát điều chỉnh hoa hồng, %s điều chỉnh %s KVND, vui lòng nhanh chóng xét duyệt.",
    "url": "/operation/bonusManager?name=bonusReview"
  }
}`
)
