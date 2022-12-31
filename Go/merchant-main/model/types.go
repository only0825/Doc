package model

import (
	"database/sql"
	"github.com/shopspring/decimal"
)

const (
	TGDateFormat = "2006-01-02T15:04:05+07:00"
)

type Cate struct {
	ID     int64  `db:"id" json:"id"`
	Level  string `db:"level" json:"level"`
	Name   string `db:"name" json:"name"`
	Sort   int    `db:"sort" json:"sort"`
	Prefix string `db:"prefix" json:"prefix"`
}

type MemberData struct {
	Member
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Zalo     string `json:"zalo"`
}
type MemberPageData struct {
	T int64        `json:"t"`
	D []MemberData `json:"d"`
	S uint         `json:"s"`
}

type Member struct {
	UID                 string `db:"uid" json:"uid"`
	Username            string `db:"username" json:"username"`                           //会员名
	Avatar              string `db:"avatar" json:"avatar"`                               //头像
	Password            string `db:"password" json:"password"`                           //密码
	Birth               string `db:"birth" json:"birth"`                                 //生日日期
	BirthHash           string `db:"birth_hash" json:"birth_hash"`                       //生日日期哈希
	RealnameHash        string `db:"realname_hash" json:"realname_hash"`                 //真实姓名哈希
	EmailHash           string `db:"email_hash" json:"email_hash"`                       //邮件地址哈希
	PhoneHash           string `db:"phone_hash" json:"phone_hash"`                       //电话号码哈希
	ZaloHash            string `db:"zalo_hash" json:"zalo_hash"`                         //zalo哈希
	Prefix              string `db:"prefix" json:"prefix"`                               //站点前缀
	Tester              string `db:"tester" json:"tester"`                               //1正式 0测试
	WithdrawPwd         uint64 `db:"withdraw_pwd" json:"withdraw_pwd"`                   //取款密码哈希
	Regip               string `db:"regip" json:"regip"`                                 //注册IP
	RegUrl              string `db:"reg_url" json:"reg_url"`                             //注册域名
	RegDevice           string `db:"reg_device" json:"reg_device"`                       //注册设备号
	CreatedAt           uint32 `db:"created_at" json:"created_at"`                       //注册时间
	LastLoginIp         string `db:"last_login_ip" json:"last_login_ip"`                 //最后登陆ip
	LastLoginAt         uint32 `db:"last_login_at" json:"last_login_at"`                 //最后登陆时间
	SourceId            uint8  `db:"source_id" json:"source_id"`                         //注册来源 1 pc 2h5 3 app
	FirstDepositAt      uint32 `db:"first_deposit_at" json:"first_deposit_at"`           //首充时间
	FirstDepositAmount  string `db:"first_deposit_amount" json:"first_deposit_amount"`   //首充金额
	SecondDepositAt     uint32 `db:"second_deposit_at" json:"second_deposit_at"`         //二存时间
	SecondDepositAmount string `db:"second_deposit_amount" json:"second_deposit_amount"` //二充金额
	FirstBetAt          uint32 `db:"first_bet_at" json:"first_bet_at"`                   //首投时间
	FirstBetAmount      string `db:"first_bet_amount" json:"first_bet_amount"`           //首投金额
	TopUid              string `db:"top_uid" json:"top_uid"`                             //总代uid
	TopName             string `db:"top_name" json:"top_name"`                           //总代代理
	ParentUid           string `db:"parent_uid" json:"parent_uid"`                       //上级uid
	ParentName          string `db:"parent_name" json:"parent_name"`                     //上级代理
	BankcardTotal       uint8  `db:"bankcard_total" json:"bankcard_total"`               //用户绑定银行卡的数量
	LastLoginDevice     string `db:"last_login_device" json:"last_login_device"`         //最后登陆设备
	LastLoginSource     uint8  `db:"last_login_source" json:"last_login_source"`         //上次登录设备来源:1=pc,2=h5,3=ios,4=andriod
	Remarks             string `db:"remarks" json:"remarks"`                             //备注
	State               uint8  `db:"state" json:"state"`                                 //状态 1正常 2禁用
	Level               int    `db:"level" json:"level" redis:"level"`                   //等级
	Balance             string `db:"balance" json:"balance"`                             //余额
	LockAmount          string `db:"lock_amount" json:"lock_amount"`                     //锁定金额
	Commission          string `db:"commission" json:"commission"`                       //佣金
	MaintainName        string `db:"maintain_name" json:"maintain_name"`                 //维护人
	GroupName           string `db:"group_name" json:"group_name"`                       //团队名称
	AgencyType          int64  `db:"agency_type" json:"agency_type"`                     //391团队代理 393普通代理
	Address             string `db:"address" json:"address"`                             //收货地址
	LastWithdrawAt      uint32 `db:"last_withdraw_at" json:"last_withdraw_at"`           //上次提现时间

}

// MemberPlatform 会员场馆表
type MemberPlatform struct {
	ID                    string `db:"id" json:"id" redis:"id"`                                                                //
	Username              string `db:"username" json:"username" redis:"username"`                                              //用户名
	Pid                   string `db:"pid" json:"pid" redis:"pid"`                                                             //场馆ID
	Password              string `db:"password" json:"password" redis:"password"`                                              //平台密码
	Balance               string `db:"balance" json:"balance" redis:"balance"`                                                 //平台余额
	State                 int    `db:"state" json:"state" redis:"state"`                                                       //状态:1=正常,2=锁定
	CreatedAt             uint32 `db:"created_at" json:"created_at" redis:"created_at"`                                        //
	TransferIn            int    `db:"transfer_in" json:"transfer_in" redis:"transfer_in"`                                     //0:没有转入记录1:有
	TransferInProcessing  int    `db:"transfer_in_processing" json:"transfer_in_processing" redis:"transfer_in_processing"`    //0:没有转入等待记录1:有
	TransferOut           int    `db:"transfer_out" json:"transfer_out" redis:"transfer_out"`                                  //0:没有转出记录1:有
	TransferOutProcessing int    `db:"transfer_out_processing" json:"transfer_out_processing" redis:"transfer_out_processing"` //0:没有转出等待记录1:有
	Extend                uint64 `db:"extend" json:"extend" redis:"extend"`                                                    //兼容evo
}

type MBBalance struct {
	UID        string `db:"uid" json:"uid"`
	Balance    string `db:"balance" json:"balance"`         //余额
	LockAmount string `db:"lock_amount" json:"lock_amount"` //锁定额度
	Commission string `db:"commission" json:"commission"`   //代理余额
}

//账变表
type MemberTransaction struct {
	AfterAmount  string `db:"after_amount"`  //账变后的金额
	Amount       string `db:"amount"`        //用户填写的转换金额
	BeforeAmount string `db:"before_amount"` //账变前的金额
	BillNo       string `db:"bill_no"`       //转账|充值|提现ID
	CashType     int    `db:"cash_type"`     //0:转入1:转出2:转入失败补回3:转出失败扣除4:存款5:提现
	CreatedAt    int64  `db:"created_at"`    //
	ID           string `db:"id"`            //
	UID          string `db:"uid"`           //用户ID
	Username     string `db:"username"`      //用户名
	Prefix       string `db:"prefix"`        //站点前缀
	Remark       string `db:"remark"`        //备注
}

//场馆转账表
type MemberTransfer struct {
	AfterAmount  string `db:"after_amount"`  //转账后的金额
	Amount       string `db:"amount"`        //金额
	Automatic    int    `db:"automatic"`     //1:自动转账2:脚本确认3:人工确认
	BeforeAmount string `db:"before_amount"` //转账前的金额
	BillNo       string `db:"bill_no"`       //
	CreatedAt    int64  `db:"created_at"`    //
	ID           string `db:"id"`            //
	PlatformID   string `db:"platform_id"`   //三方场馆ID
	State        int    `db:"state"`         //0:失败1:成功2:处理中3:脚本确认中4:人工确认中
	TransferType int    `db:"transfer_type"` //0:转入1:转出
	UID          string `db:"uid"`           //用户ID
	Username     string `db:"username"`      //用户名
	ConfirmAt    int64  `db:"confirm_at"`    //确认时间
	ConfirmUid   uint64 `db:"confirm_uid"`   //确认人uid
	ConfirmName  string `db:"confirm_name"`  //确认人名
}

type CommissionTransferData struct {
	S int                  `json:"s"`
	D []CommissionTransfer `json:"d"`
	T int64                `json:"t"`
}

type CommissionsData struct {
	S int           `json:"s"`
	D []Commissions `json:"d"`
	T int64         `json:"t"`
}

type CommissionTransfer struct {
	ID           string `json:"id" db:"id"`
	UID          string `json:"uid" db:"uid"`                     //用户ID
	Username     string `json:"username" db:"username"`           //用户名
	ReceiveUID   string `json:"receive_uid" db:"receive_uid"`     //用户ID
	ReceiveName  string `json:"receive_name" db:"receive_name"`   //用户名
	TransferType int    `json:"transfer_type" db:"transfer_type"` //转账类型 2 佣金提取 3佣金下发
	Amount       string `json:"amount" db:"amount"`               //金额
	CreatedAt    int64  `json:"created_at" db:"created_at"`       //创建时间
	State        int    `json:"state" db:"state"`                 //1 审核中 2 审核通过 3 审核不通过
	Automatic    int    `json:"automatic" db:"automatic"`         // 1自动 2手动
	ReviewAt     int64  `json:"review_at" db:"review_at"`         //审核时间
	ReviewUid    string `json:"review_uid" db:"review_uid"`       //审核人uid
	ReviewName   string `json:"review_name" db:"review_name"`     //审核人名
	ReviewRemark string `json:"review_remark" db:"review_remark"` //审核备注
	Prefix       string `json:"prefix" db:"prefix"`
}

type MemberLoginLogData struct {
	S int              `json:"s"`
	D []MemberLoginLog `json:"d"`
	T int64            `json:"t"`
}

//MemberAssocLogData 会员最近登陆信息
type MemberAssocLogData struct {
	S int                    `json:"s"`
	D []MemberAssocLogMember `json:"d"`
	T int64                  `json:"t"`
}

//MemberAssocLogMember 会员最近登陆字段
type MemberAssocLogMember struct {
	Username    string `db:"username" json:"username"`
	Device      int    `db:"device" json:"device"`               //24,25,35,36
	TopUID      string `db:"top_uid" json:"top_uid"`             //总代uid
	TopName     string `db:"top_name" json:"top_name"`           //总代代理名
	ParentName  string `db:"parent_name" json:"parent_name"`     // 上级代理
	CreatedAt   uint32 `db:"created_at" json:"created_at"`       //会员注册时间
	State       uint8  `db:"state" json:"state"`                 //账号状态  1正常 2禁用
	Remarks     string `db:"remarks" json:"remarks"`             //备注 账号状态备注
	LastLoginAt uint32 `db:"last_login_at" json:"last_login_at"` //最后登陆时间
	GroupName   string `db:"group_name" json:"group_name"`       //团队名称
}

//IpUser 计数字段
type IpUser struct {
	IP       string `json:"ip" db:"ip"`
	Username string `json:"username" db:"username"`
}

type MemberLoginLog struct {
	Username   string `db:"username" json:"username"`
	IP         string `db:"ip" json:"ip"`
	Device     int    `db:"device" json:"device"`
	DeviceNo   string `db:"device_no" json:"device_no"`
	TopUID     string `db:"top_uid" json:"top_uid"`         // 总代uid
	TopName    string `db:"top_name" json:"top_name"`       // 总代代理
	ParentUID  string `db:"parent_uid" json:"parent_uid"`   // 上级uid
	ParentName string `db:"parent_name" json:"parent_name"` // 上级代理
	CreateAt   int    `db:"create_at" json:"create_at"`
	Prefix     string `db:"prefix" json:"prefix"`
	Ts         string `db:"ts" json:"ts"`
	CountName  int    `db:"count_name" json:"count_name"` //ip 对应 username 数
	GroupName  string `db:"group_name" json:"group_name"` //团队名称
}

type MemberRemarkLogData struct {
	S int                `json:"s"`
	D []MemberRemarksLog `json:"d"`
	T int64              `json:"t"`
}

// 用户备注日志
type MemberRemarksLog struct {
	ID          string `msg:"id" json:"id" db:"id"`
	TS          string `msg:"ts" json:"ts" db:"ts"`
	UID         string `msg:"uid" json:"uid" db:"uid"`
	Username    string `msg:"username" json:"username" db:"username"`
	Msg         string `msg:"msg" json:"msg" db:"msg"`
	File        string `msg:"file" json:"file" db:"file"`
	CreatedName string `msg:"created_name" json:"created_name" db:"created_name"`
	CreatedAt   int64  `msg:"created_at" json:"created_at" db:"created_at"`
	Prefix      string `msg:"prefix" json:"prefix" db:"prefix"`
}

// MemberAdjust db structure
type MemberAdjust struct {
	ID            string  `db:"id" json:"id"`
	UID           string  `db:"uid" json:"uid"` // 会员id
	Prefix        string  `db:"prefix" json:"prefix"`
	Ty            int     `db:"ty" json:"ty"`                         //来源
	Username      string  `db:"username" json:"username"`             // 会员username
	TopUid        string  `db:"top_uid" json:"top_uid"`               //总代uid
	TopName       string  `db:"top_name" json:"top_name"`             //总代代理
	ParentUid     string  `db:"parent_uid" json:"parent_uid"`         //上级uid
	ParentName    string  `db:"parent_name" json:"parent_name"`       //上级代理
	Amount        float64 `db:"amount" json:"amount"`                 // 调整金额
	AdjustType    int     `db:"adjust_type" json:"adjust_type"`       // 调整类型:1=系统调整,2=输赢调整,3=线下转卡充值
	AdjustMode    int     `db:"adjust_mode" json:"adjust_mode"`       // 调整方式:1=上分,2=下分
	IsTurnover    int     `db:"is_turnover" json:"is_turnover"`       // 是否需要流水限制:1=需要,0=不需要
	TurnoverMulti int     `db:"turnover_multi" json:"turnover_multi"` // 流水倍数
	ApplyRemark   string  `db:"apply_remark" json:"apply_remark"`     // 申请备注
	ReviewRemark  string  `db:"review_remark" json:"review_remark"`   // 审核备注
	State         int     `db:"state" json:"state"`                   // 状态:1=审核中,2=审核通过,3=审核未通过
	HandOutState  int     `db:"hand_out_state" json:"hand_out_state"` // 上下分状态 1 失败 2成功 3场馆上分处理中
	Images        string  `db:"images" json:"images"`                 // 图片地址
	ApplyAt       int64   `db:"apply_at" json:"apply_at"`             // 申请时间
	ApplyUid      string  `db:"apply_uid" json:"apply_uid"`           // 申请人uid
	ApplyName     string  `db:"apply_name" json:"apply_name"`         // 申请人
	ReviewAt      int64   `db:"review_at" json:"review_at"`           // 审核时间
	ReviewUid     string  `db:"review_uid" json:"review_uid"`         // 审核人uid
	ReviewName    string  `db:"review_name" json:"review_name"`       // 审核人
	IsRisk        int     `db:"-" json:"is_risk"`
	Tester        string  `db:"tester" json:"tester"`
}

type DividendData struct {
	T   int64             `json:"t"`
	D   []MemberDividend  `json:"d"`
	Agg map[string]string `json:"agg"`
}

type Dividend struct {
	ID            string  `db:"id" json:"id"`
	UID           string  `db:"uid" json:"uid"`
	Prefix        string  `db:"prefix" json:"prefix"`
	Ty            int     `db:"ty" json:"ty"`
	Username      string  `db:"username" json:"username"`
	TopUid        string  `db:"top_uid" json:"top_uid"`         //总代uid
	TopName       string  `db:"top_name" json:"top_name"`       //总代代理
	ParentUid     string  `db:"parent_uid" json:"parent_uid"`   //上级uid
	ParentName    string  `db:"parent_name" json:"parent_name"` //上级代理
	Amount        float64 `db:"amount" json:"amount"`
	ReviewAt      uint64  `db:"review_at" json:"review_at"`
	ReviewUid     string  `db:"review_uid" json:"review_uid"`
	ReviewName    string  `db:"review_name" json:"review_name"`
	WaterFlow     float64 `db:"water_flow" json:"water_flow"`
	WaterMultiple int     `db:"water_multiple" json:"water_multiple"`
}

type MemberDividend struct {
	ID            string  `db:"id" json:"id"`
	UID           string  `db:"uid" json:"uid"`
	Prefix        string  `db:"prefix" json:"prefix"`
	Ty            int     `db:"ty" json:"ty"`
	Level         int     `db:"level" json:"level"`
	WaterLimit    uint8   `db:"water_limit" json:"water_limit"`
	PID           string  `db:"pid" json:"pid"`
	PTitle        string  `db:"ptitle" json:"ptitle"`
	Username      string  `db:"username" json:"username"`
	TopUid        string  `db:"top_uid" json:"top_uid"`         //总代uid
	TopName       string  `db:"top_name" json:"top_name"`       //总代代理
	ParentUid     string  `db:"parent_uid" json:"parent_uid"`   //上级uid
	ParentName    string  `db:"parent_name" json:"parent_name"` //上级代理
	Amount        float64 `db:"amount" json:"amount"`
	WaterFlow     float64 `db:"water_flow" json:"water_flow"`
	WaterMultiple int     `db:"water_multiple" json:"water_multiple"`
	State         int     `db:"state" json:"state"`
	Automatic     int     `db:"automatic" json:"automatic"`
	Remark        string  `db:"remark" json:"remark"`
	ReviewRemark  string  `db:"review_remark" json:"review_remark"`
	ApplyAt       uint64  `db:"apply_at" json:"apply_at"`
	ApplyUid      string  `db:"apply_uid" json:"apply_uid"`
	ApplyName     string  `db:"apply_name" json:"apply_name"`
	ReviewAt      uint64  `db:"review_at" json:"review_at"`
	ReviewUid     string  `db:"review_uid" json:"review_uid"`
	ReviewName    string  `db:"review_name" json:"review_name"`
	IsRisk        int     `db:"-" json:"is_risk"`
	Tester        string  `db:"tester" json:"tester"`
}

type BannerData struct {
	T int64    `json:"t"`
	D []Banner `json:"d"`
	S uint     `json:"s"`
}

type Banner struct {
	ID          string `json:"id" db:"id"`                     //
	Title       string `json:"title" db:"title"`               //标题
	Device      string `json:"device" db:"device"`             //设备类型(1,2)
	RedirectURL string `json:"redirect_url" db:"redirect_url"` //跳转地址
	Images      string `json:"images" db:"images"`             //图片路径
	Seq         string `json:"seq" db:"seq"`                   //排序
	Flags       string `json:"flags" db:"flags"`               //广告类型
	ShowType    string `json:"show_type" db:"show_type"`       //1 永久有效 2 指定时间
	ShowAt      string `json:"show_at" db:"show_at"`           //开始展示时间
	HideAt      string `json:"hide_at" db:"hide_at"`           //结束展示时间
	URLType     string `json:"url_type" db:"url_type"`         //链接类型 1站内 2站外
	UpdatedName string `json:"updated_name" db:"updated_name"` //更新人name
	UpdatedUID  string `json:"updated_uid" db:"updated_uid"`   //更新人id
	UpdatedAt   string `json:"updated_at" db:"updated_at"`     //更新时间
	State       uint8  `json:"state" db:"state"`               //0:关闭1:开启
	Prefix      string `json:"prefix" db:"prefix"`
}

type BlacklistData struct {
	T int64       `json:"t"`
	D []Blacklist `json:"d"`
	S uint        `json:"s"`
}

type Blacklist struct {
	ID          string `json:"id" db:"id"`                                 //id
	Ty          int    `json:"ty" db:"ty"`                                 //黑名单类型
	Value       string `json:"value" db:"value"`                           //黑名单类型值
	Remark      string `json:"remark" db:"remark"`                         //备注
	CreatedAt   string `json:"created_at" db:"created_at" rule:"none"`     //添加时间
	CreatedUID  string `json:"created_uid" db:"created_uid" rule:"none"`   //添加人uid
	CreatedName string `json:"created_name" db:"created_name" rule:"none"` //添加人name
}

type MemberAssocLoginLogData struct {
	S int                   `json:"s"`
	D []MemberAssocLoginLog `json:"d"`
	T int64                 `json:"t"`
}

type MemberAssocLoginLog struct {
	Username string `json:"username"`
	IP       int64  `json:"ip"`
	IPS      string `json:"ips"`
	Device   string `json:"device"`
	DeviceNo string `json:"device_no"`
	Date     uint32 `json:"date"`
	Serial   string `json:"serial"`
	Agency   bool   `json:"agency"`
	Parents  string `json:"parents"`
	Tags     string `json:"tags"`
}

// 数据库 游戏字段
type GameLists struct {
	ID         string `db:"id" json:"id"`
	PlatformId string `db:"platform_id" json:"platform_id"`
	Name       string `db:"name" json:"name"`
	EnName     string `db:"en_name" json:"en_name"`
	VnAlias    string `db:"vn_alias" json:"vn_alias"`
	ClientType string `db:"client_type" json:"client_type"`
	GameType   int64  `db:"game_type" json:"game_type"`
	GameId     string `db:"game_id" json:"game_id"`
	ImgPhone   string `db:"img_phone" json:"img_phone"`
	ImgPc      string `db:"img_pc" json:"img_pc"`
	ImgCover   string `db:"img_cover" json:"img_cover"`
	OnLine     int64  `db:"online" json:"online"`
	IsHot      int    `db:"is_hot" json:"is_hot"`
	IsNew      int    `db:"is_new" json:"is_new"`
	IsFs       int    `db:"is_fs" json:"is_fs"`
	Sorting    int64  `db:"sorting" json:"sorting"`
	CreatedAt  int64  `db:"created_at" json:"created_at"`
}

// 游戏列表返回数据结构
type GamePageList struct {
	D []GameLists `json:"d"`
	T int64       `json:"t"`
	S uint        `json:"s"`
}

type showGameJson struct {
	ID         string `db:"id" json:"id"`
	PlatformID string `db:"platform_id" json:"platform_id"`
	EnName     string `db:"en_name" json:"en_name"`
	ClientType string `db:"client_type" json:"client_type"`
	GameType   string `db:"game_type" json:"game_type"`
	GameID     string `db:"game_id" json:"game_id"`
	ImgPhone   string `db:"img_phone" json:"img_phone"`
	ImgPc      string `db:"img_pc" json:"img_pc"`
	IsHot      int    `db:"is_hot" json:"is_hot"`
	IsNew      int    `db:"is_new" json:"is_new"`
	Name       string `db:"name" json:"name"`
	ImgCover   string `db:"img_cover" json:"img_cover"`
	Sort       int    `db:"sorting" json:"sorting"`
	VnAlias    string `db:"vn_alias" json:"vn_alias"`
}

type Priv struct {
	ID        string `db:"id" json:"id" redis:"id"`                      //
	Name      string `db:"name" json:"name" redis:"name"`                //权限名字
	Module    string `db:"module" json:"module" redis:"module"`          //模块
	Sortlevel string `db:"sortlevel" json:"sortlevel" redis:"sortlevel"` //
	State     int    `db:"state" json:"state" redis:"state"`             //0:关闭1:开启
	Pid       int64  `db:"pid" json:"pid" redis:"pid"`                   //父级ID
}

// 后台用户登录记录
type adminLoginLogBase struct {
	UID       string `msg:"uid" json:"uid"`
	Name      string `msg:"name" json:"name"`
	IP        string `msg:"ip" json:"ip"`
	Device    string `msg:"device" json:"device"`
	Flag      int    `msg:"flag" json:"flag"` // 1 登录 2 登出
	CreatedAt uint32 `msg:"created_at" json:"created_at"`
	Prefix    string `msg:"prefix" json:"prefix"`
}

//type adminLoginLog struct {
//	Id string `msg:"_id" json:"id"`
//	adminLoginLogBase
//}

type adminLoginLog struct {
	UserName string `db:"username" json:"username"`
	IP       string `db:"ip" json:"ip"`
	DeviceNo string `db:"device_no" json:"device_no"`
	CreateAt uint32 `db:"create_at" json:"create_at"`
	Prefix   string `db:"prefix" json:"prefix"`
}

// 后台用户登录记录
type AdminLoginLogData struct {
	D []adminLoginLog `json:"d"`
	T int64           `json:"t"`
	S int             `json:"s"`
}

// 系统日志
type systemLogBase struct {
	UID       string `msg:"uid" json:"uid"`
	Name      string `msg:"name" json:"name"`
	IP        string `msg:"ip" json:"ip"`
	Title     string `msg:"title" json:"title"`
	Content   string `msg:"content" json:"content"`
	CreatedAt uint32 `msg:"created_at" json:"created_at"`
	Prefix    string `msg:"prefix" json:"prefix"`
}

type systemLog struct {
	Id string `msg:"_id" json:"id"`
	systemLogBase
}

// 系统日志 分页展示数据
type SystemLogData struct {
	D []systemLog `json:"d"`
	T int64       `json:"t"`
	S int         `json:"s"`
}

type MemberRebate struct {
	UID              string `db:"uid" json:"uid"`
	ZR               string `db:"zr" json:"zr"`                                 //真人返水
	QP               string `db:"qp" json:"qp"`                                 //棋牌返水
	TY               string `db:"ty" json:"ty"`                                 //体育返水
	DJ               string `db:"dj" json:"dj"`                                 //电竞返水
	DZ               string `db:"dz" json:"dz"`                                 //电游返水
	CP               string `db:"cp" json:"cp"`                                 //彩票返水
	FC               string `db:"fc" json:"fc"`                                 //斗鸡返水
	BY               string `db:"by" json:"by"`                                 //捕鱼返水
	CgOfficialRebate string `db:"cg_official_rebate" json:"cg_official_rebate"` //CG官方彩返点
	CgHighRebate     string `db:"cg_high_rebate" json:"cg_high_rebate"`         //CG高频彩返点
	CreatedAt        uint32 `db:"created_at" json:"created_at"`
	ParentUID        string `db:"parent_uid" json:"parent_uid"`
	Prefix           string `db:"prefix" json:"prefix"`
}

type MemberMaxRebate struct {
	ZR               sql.NullFloat64 `db:"zr" json:"zr"`                                 //真人返水
	QP               sql.NullFloat64 `db:"qp" json:"qp"`                                 //棋牌返水
	TY               sql.NullFloat64 `db:"ty" json:"ty"`                                 //体育返水
	DJ               sql.NullFloat64 `db:"dj" json:"dj"`                                 //电竞返水
	DZ               sql.NullFloat64 `db:"dz" json:"dz"`                                 //电游返水
	CP               sql.NullFloat64 `db:"cp" json:"cp"`                                 //彩票返水
	FC               sql.NullFloat64 `db:"fc" json:"fc"`                                 //斗鸡返水
	BY               sql.NullFloat64 `db:"by" json:"by"`                                 //捕鱼返水
	CgHighRebate     sql.NullFloat64 `db:"cg_high_rebate" json:"cg_high_rebate"`         //高频彩返点
	CgOfficialRebate sql.NullFloat64 `db:"cg_official_rebate" json:"cg_official_rebate"` //官方彩返点
}

type NoticeData struct {
	D []Notice `json:"d"`
	T int64    `json:"t"`
	S uint     `json:"s"`
}

// 系统公告
type Notice struct {
	ID          string `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`               // 标题
	Content     string `db:"content" json:"content"`           // 内容
	Redirect    int    `db:"redirect" json:"redirect"`         // 是否跳转：1是 2否
	RedirectUrl string `db:"redirect_url" json:"redirect_url"` // 跳转url
	State       int    `db:"state" json:"state"`               // 1停用 2启用
	CreatedAt   int64  `db:"created_at" json:"created_at"`     // 创建时间
	CreatedUid  string `db:"created_uid" json:"created_uid"`
	CreatedName string `db:"created_name" json:"created_name"`
	Prefix      string `db:"prefix" json:"prefix"`
}

// 帐变数据
type TransactionData struct {
	T   int64         `json:"t"`
	D   []Transaction `json:"d"`
	Agg string        `db:"agg" json:"agg"`
}

type Transaction struct {
	ID           string `db:"id" json:"id"`
	PlatformId   string `db:"platform_id" json:"platform_id"`
	BillNo       string `db:"bill_no" json:"bill_no"`
	OperationNo  string `db:"operation_no" json:"operation_no"`
	Uid          string `db:"uid" json:"uid"`
	Username     string `db:"username" json:"username"`
	CashType     int    `db:"cash_type" json:"cash_type"`
	Amount       string `db:"amount" json:"amount"`
	BeforeAmount string `db:"before_amount" json:"before_amount"`
	AfterAmount  string `db:"after_amount" json:"after_amount"`
	CreatedAt    uint64 `db:"created_at" json:"created_at"`
	Remark       string `db:"remark" json:"remark"`
}

// 场馆转账数据
type TransferData struct {
	T   int64      `json:"t"`
	D   []Transfer `json:"d"`
	Agg string     `db:"agg" json:"agg"`
}

//转账记录
type Transfer struct {
	ID           string `json:"id" db:"id"`
	UID          string `json:"uid" db:"uid"`
	BillNo       string `json:"bill_no" db:"bill_no"`
	PlatformId   string `json:"platform_id" db:"platform_id"`
	Username     string `json:"username" db:"username"`
	TransferType int    `json:"transfer_type" db:"transfer_type"`
	Amount       string `json:"amount" db:"amount"`
	BeforeAmount string `json:"before_amount" db:"before_amount"`
	AfterAmount  string `json:"after_amount" db:"after_amount"`
	CreatedAt    uint64 `json:"created_at" db:"created_at"`
	State        int    `json:"state" db:"state"`
	Automatic    int    `json:"automatic" db:"automatic"`
	ConfirmName  string `json:"confirm_name" db:"confirm_name"`
}

// 游戏记录数据
type GameRecordData struct {
	T   int64        `json:"t"`
	D   []GameRecord `json:"d"`
	Agg GameRecord   `json:"agg"`
}

//游戏投注记录结构
type GameRecord struct {
	ID             string  `db:"id" json:"id" form:"id"`
	RowId          string  `db:"row_id" json:"row_id" form:"row_id"`
	BillNo         string  `db:"bill_no" json:"bill_no" form:"bill_no"`
	ApiType        string  `db:"api_type" json:"api_types" form:"api_type"`
	PlayerName     string  `db:"player_name" json:"player_name" form:"player_name"`
	Name           string  `db:"name" json:"name" form:"name"`
	Uid            string  `db:"uid" json:"uid" form:"uid"`
	NetAmount      float64 `db:"net_amount" json:"net_amount" form:"net_amount"`
	BetTime        int64   `db:"bet_time" json:"bet_time" form:"bet_time"`
	StartTime      int64   `db:"start_time" json:"start_time" form:"start_time"`
	Resettle       uint8   `db:"resettle" json:"resettle" form:"resettle"`
	Presettle      uint8   `db:"presettle" json:"presettle" form:"presettle"`
	GameType       string  `db:"game_type" json:"game_type" form:"game_type"`
	BetAmount      float64 `db:"bet_amount" json:"bet_amount" form:"bet_amount"`
	ValidBetAmount float64 `db:"valid_bet_amount" json:"valid_bet_amount" form:"valid_bet_amount"`
	RebateAmount   float64 `db:"rebate_amount" json:"rebate_amount" form:"rebate_amount"`
	Flag           int     `db:"flag" json:"flag" form:"flag"`
	PlayType       string  `db:"play_type" json:"play_type" form:"play_type"`
	Prefix         string  `db:"prefix" json:"prefix" form:"prefix"`
	Result         string  `db:"result" json:"result" form:"result"`
	CreatedAt      uint64  `db:"created_at" json:"created_at" form:"created_at"`
	UpdatedAt      uint64  `db:"updated_at" json:"updated_at" form:"updated_at"`
	ApiName        string  `db:"api_name" json:"api_name" form:"api_name"`
	ApiBillNo      string  `db:"api_bill_no" json:"api_bill_no" form:"api_bill_no"`
	MainBillNo     string  `db:"main_bill_no" json:"main_bill_no" form:"main_bill_no"`
	GameName       string  `db:"game_name" json:"game_name" form:"game_name"`
	HandicapType   string  `db:"handicap_type" json:"handicap_type" form:"handicap_type"`
	Handicap       string  `db:"handicap" json:"handicap" form:"handicap"`
	Odds           float64 `db:"odds" json:"odds" form:"odds"`
	SettleTime     int64   `db:"settle_time" json:"settle_time" form:"settle_time"`
	ApiBetTime     uint64  `db:"api_bet_time" json:"api_bet_time" form:"api_bet_time"`
	ApiSettleTime  uint64  `db:"api_settle_time" json:"api_settle_time" form:"api_settle_time"`
	IsRisk         int     `db:"-" json:"is_risk"`
	TopUid         string  `db:"top_uid" json:"top_uid"`         //总代uid
	TopName        string  `db:"top_name" json:"top_name"`       //总代代理
	ParentUid      string  `db:"parent_uid" json:"parent_uid"`   //上级uid
	ParentName     string  `db:"parent_name" json:"parent_name"` //上级代理
}

type Commissions struct {
	Id               string  `json:"id" db:"id"`
	Uid              string  `json:"uid" db:"uid"`
	Username         string  `json:"username" db:"username"`
	CreatedAt        int64   `json:"created_at" db:"created_at"`
	CommissionMonth  int64   `json:"commission_month" db:"commission_month"`
	TeamNum          int     `json:"team_num" db:"team_num"`
	ActiveNum        int     `json:"active_num" db:"active_num"`
	DepositAmount    float64 `json:"deposit_amount" db:"deposit_amount"`
	WithdrawAmount   float64 `json:"withdraw_amount" db:"withdraw_amount"`
	WinAmount        float64 `json:"win_amount" db:"win_amount"`
	PlatformAmount   float64 `json:"platform_amount" db:"platform_amount"`
	RebateAmount     float64 `json:"rebate_amount" db:"rebate_amount"`
	DividendAmount   float64 `json:"dividend_amount" db:"dividend_amount"`
	AdjustAmount     float64 `json:"adjust_amount" db:"adjust_amount"`
	NetWin           float64 `json:"net_win" db:"net_win"`
	BalanceAmount    float64 `json:"balance_amount" db:"balance_amount"`
	AdjustCommission float64 `json:"adjust_commission" db:"adjust_commission"`
	AdjustWin        float64 `json:"adjust_win" db:"adjust_win"`
	Amount           float64 `json:"amount" db:"amount"`
	Remark           string  `json:"remark" db:"remark"`
	Note             string  `json:"note" db:"note"`
	State            int     `json:"state" db:"state"`
	HandOutAt        int64   `json:"hand_out_at" db:"hand_out_at"`
	HandOutUid       string  `json:"hand_out_uid" db:"hand_out_uid"`
	HandOutName      string  `json:"hand_out_name" db:"hand_out_name"`
	DividendAgAmount float64 `json:"dividend_ag_amount" db:"dividend_ag_amount"`
	LastMonthAmount  float64 `json:"last_month_amount" db:"last_month_amount"`
	Prefix           string  `json:"prefix" db:"prefix"`
	PlanId           string  `json:"plan_id" db:"plan_id"`
	PlanName         string  `json:"plan_name" db:"plan_name"`
}

type CommissionTransaction struct {
	Id           string `json:"id" db:"id"`
	BillNo       string `json:"bill_no" db:"bill_no"`
	Uid          string `json:"uid" db:"uid"`
	Username     string `json:"username" db:"username"`
	CashType     int    `json:"cash_type" db:"cash_type"`
	Amount       string `json:"amount" db:"amount"`
	BeforeAmount string `json:"before_amount" db:"before_amount"`
	AfterAmount  string `json:"after_amount" db:"after_amount"`
	CreatedAt    int64  `json:"created_at" db:"created_at"`
	Prefix       string `json:"prefix" db:"prefix"`
}

type MembersTree struct {
	Ancestor   string `db:"ancestor" json:"ancestor,omitempty"`
	Descendant string `db:"descendant" json:"descendant,omitempty"`
	Lvl        int    `db:"lvl" json:"lvl,omitempty"`
	Prefix     string `db:"prefix" json:"prefix,omitempty"`
}

type CommssionConf struct {
	ID     string `json:"id" db:"id"`
	UID    string `json:"uid" db:"uid"`
	PlanID string `json:"plan_id" db:"plan_id"`
}

// CommissionPlan 返佣方案具体比例
type CommissionPlan struct {
	ID              string `db:"id" json:"id"`
	Name            string `db:"name" json:"name"`                         //方案名称
	CommissionMonth int64  `db:"commission_month" json:"commission_month"` // 生效月份
	CreatedAt       int64  `db:"created_at" json:"created_at"`
	UpdatedUID      string `db:"updated_uid" json:"updated_uid"`
	UpdatedName     string `db:"updated_name" json:"updated_name"`
	UpdatedAt       int64  `db:"updated_at" json:"updated_at"`
	Prefix          string `db:"prefix" json:"prefix"`
}

// CommissionDetail 返佣方案具体返水
type CommissionDetail struct {
	ID     string  `db:"id" json:"id"`
	PlanID string  `db:"plan_id" json:"plan_id"` // 所属方案
	WinMax float64 `db:"win_max" json:"win_max"` //净输赢最大值
	WinMin float64 `db:"win_min" json:"win_min"` //净输赢最小值
	Rate   float64 `db:"rate" json:"rate"`       //返佣比例
	Prefix string  `db:"prefix" json:"prefix"`
}

type CommPlanPageData struct {
	T       int64                         `json:"t"`
	D       []CommissionPlan              `json:"d"`
	S       uint                          `json:"s"`
	Details map[string][]CommissionDetail `json:"details"`
}

//VIP的数据库结构
type MemberLevel struct {
	ID                string  `db:"id" json:"id" name:"id" rule:"none"`
	Level             int     `db:"level" json:"level" name:"level" rule:"digit" msg:"level error"`
	LevelName         string  `db:"level_name" json:"level_name" name:"level_name" rule:"chnAlnum" msg:"level_name error"`
	RechargeNum       int     `db:"recharge_num" json:"recharge_num" name:"recharge_num" rule:"digit" msg:"recharge_num error"`
	UpgradeDeposit    int     `db:"upgrade_deposit" json:"upgrade_deposit" name:"upgrade_deposit" rule:"digit" msg:"upgrade_deposit error"`
	UpgradeRecord     int     `db:"upgrade_record" json:"upgrade_record" name:"upgrade_record" rule:"digit" msg:"upgrade_record error"`
	RelegationFlowing int     `db:"relegation_flowing" json:"relegation_flowing" name:"relegation_flowing" rule:"digit" msg:"relegation_flowing error"`
	UpgradeGift       int     `db:"upgrade_gift" json:"upgrade_gift" name:"upgrade_gift" rule:"digit" msg:"upgrade_gift error"`
	BirthGift         int     `db:"birth_gift" json:"birth_gift" name:"birth_gift" rule:"digit" msg:"birth_gift error"`
	WithdrawCount     int     `db:"withdraw_count" json:"withdraw_count" name:"withdraw_count" rule:"digit" msg:"withdraw_count error"`
	WithdrawMax       float64 `db:"withdraw_max" json:"withdraw_max" name:"withdraw_max" rule:"float" msg:"withdraw_max error"`
	EarlyMonthPacket  int     `db:"early_month_packet" json:"early_month_packet" name:"early_month_packet" rule:"digit" msg:"early_month_packet error"`
	LateMonthPacket   int     `db:"late_month_packet" json:"late_month_packet" name:"late_month_packet" rule:"digit" msg:"late_month_packet error"`
	CreateAt          uint32  `db:"created_at" json:"created_at" name:"created_at" rule:"none"` // 创建时间
	UpdatedAt         uint32  `db:"updated_at" json:"updated_at" name:"updated_at" rule:"none"`
	UserCount         int     `db:"user_count" json:"user_count"`
}

type MemberLevelRecordData struct {
	T int64               `json:"t"`
	D []MemberLevelRecord `json:"d"`
}

// 会员等级调整记录
type MemberLevelRecord struct {
	ID                  string `db:"id" json:"id"`
	UID                 string `db:"uid" json:"uid"`                                     //会员id
	Username            string `db:"username" json:"username"`                           //会员账号
	BeforeLevel         int    `db:"before_level" json:"before_level"`                   //调整前会员等级
	AfterLevel          int    `db:"after_level" json:"after_level"`                     //调整后会员等级
	TotalDeposit        string `db:"total_deposit" json:"total_deposit"`                 //累计存款
	TotalWaterFlow      string `db:"total_water_flow" json:"total_water_flow"`           //累计流水
	RelegationWaterFlow string `db:"relegation_water_flow" json:"relegation_water_flow"` //累计保级流水
	Ty                  int    `db:"ty" json:"ty"`                                       //会员等级调整类型
	CreatedAt           uint64 `db:"created_at" json:"created_at"`                       //操作时间
	CreatedUid          string `db:"created_uid" json:"created_uid"`                     //操作人uid
	CreatedName         string `db:"created_name" json:"created_name"`                   //操作人名
}

type MemberListParam struct {
	ParentName string `rule:"none" name:"parent_name"`
	Username   string `rule:"none" name:"username"`
	State      int    `rule:"none" name:"state"`
	StartAt    string `rule:"none" name:"start_at"`
	EndAt      string `rule:"none" name:"end_at"`
	RegStart   string `rule:"none" name:"reg_start"`
	RegEnd     string `rule:"none" name:"reg_end"`
	Page       int    `rule:"digit" default:"1" min:"1" msg:"page error" name:"page"`
	PageSize   int    `rule:"digit" min:"10" max:"200" msg:"page_size error" name:"page_size"`
}

type memberListShow struct {
	UID            string `db:"uid" json:"uid" redis:"uid"`               //
	Name           string `db:"username" json:"name" redis:"name"`        //
	State          int    `db:"state" json:"state" redis:"state"`         //状态 1正常 2禁用
	ParentUID      string `db:"parent_uid" json:"parent_uid"`             //
	ParentName     string `db:"parent_name" json:"parent_name"`           //
	TopUID         string `db:"top_uid" json:"top_uid"`                   //
	TopName        string `db:"top_name" json:"top_name"`                 //
	CreatedAt      uint32 `db:"created_at" json:"created_at"`             //
	FirstDepositAt uint32 `db:"first_deposit_at" json:"first_deposit_at"` //
	LastLoginIP    string `db:"last_login_ip" json:"last_login_ip"`       //
	LastLoginAt    uint32 `db:"last_login_at" json:"last_login_at"`       //成为代理时间
}

type AgencyBaseSumField struct {
	DepositAmount    float64 `json:"deposit_amount"`     // 存款
	WithdrawAmount   float64 `json:"withdraw_amount"`    // 提款
	ValidBetAmount   float64 `json:"valid_bet_amount"`   // 有效流水
	CompanyNetAmount float64 `json:"company_net_amount"` // 输赢
	DividendAmount   float64 `json:"dividend_amount"`    // 红利
	DividendAgency   float64 `json:"dividend_agency"`    // 代理红利
	RebateAmount     float64 `json:"rebate_amount"`      // 返水
	AdjustAmount     float64 `json:"adjust_amount"`      // 调整
}

type MemReport struct {
	Id               string  `json:"id" db:"id"`
	ReportTime       int64   `json:"report_time" db:"report_time"`
	DepositAmount    float64 `json:"deposit_amount" db:"deposit_amount"`
	WithdrawalAmount float64 `json:"withdrawal_amount" db:"withdrawal_amount"`
	AdjustAmount     float64 `json:"adjust_amount" db:"adjust_amount"`
	ValidBetAmount   float64 `json:"valid_bet_amount" db:"valid_bet_amount"`
	CompanyNetAmount float64 `json:"company_net_amount" db:"company_net_amount"`
	DividendAmount   float64 `json:"dividend_amount" db:"dividend_amount"`
	RebateAmount     float64 `json:"rebate_amount" db:"rebate_amount"`
	Prefix           string  `json:"prefix" db:"prefix"`
	Uid              string  `json:"uid" db:"uid"`
	Username         string  `json:"username" db:"username"`
}

type memberListData struct {
	memberListShow
	WithdrawalAmount float64 `json:"withdrawal_amount"`
	DepositAmount    float64 `json:"deposit_amount"`
	CompanyNetAmount float64 `json:"company_net_amount"`
	ValidBetAmount   float64 `json:"valid_bet_amount"`
	Balance          string  `json:"balance"`
	RebateAmount     float64 `json:"rebate_amount"`
	DividendAmount   float64 `json:"dividend_amount"`
	DividendAgency   float64 `json:"dividend_agency"`
	IsRisk           int     `json:"is_risk"`
}

type AgencyMemberData struct {
	D []memberListData `json:"d"`
	T int              `json:"t"`
}

// Withdraw 出款
type Withdraw struct {
	ID                string  `db:"id" json:"id"`                                   //
	Prefix            string  `db:"prefix" json:"prefix"`                           //转账后的金额
	BID               string  `db:"bid" json:"bid"`                                 //转账前的金额
	Flag              int     `db:"flag" json:"flag"`                               //
	OID               string  `db:"oid" json:"oid"`                                 //转账前的金额
	Level             int     `db:"level" json:"level"`                             //
	UID               string  `db:"uid" json:"uid"`                                 //用户ID
	ParentUID         string  `db:"parent_uid" json:"parent_uid"`                   //上级代理ID
	ParentName        string  `db:"parent_name" json:"parent_name"`                 //上级代理
	TopUID            string  `db:"top_uid" json:"top_uid"`                         //总代uid
	TopName           string  `db:"top_name" json:"top_name"`                       //总代
	Username          string  `db:"username" json:"username"`                       //用户名
	PID               string  `db:"pid" json:"pid"`                                 //用户ID
	Amount            float64 `db:"amount" json:"amount"`                           //金额
	State             int     `db:"state" json:"state"`                             //0:待确认:1存款成功2:已取消
	Automatic         int     `db:"automatic" json:"automatic"`                     //1:自动转账2:脚本确认3:人工确认
	BankName          string  `db:"bank_name" json:"bank_name"`                     //银行名
	RealName          string  `db:"real_name" json:"real_name"`                     //持卡人姓名
	CardNO            string  `db:"card_no" json:"card_no"`                         //银行卡号
	CreatedAt         int64   `db:"created_at" json:"created_at"`                   //
	ConfirmAt         int64   `db:"confirm_at" json:"confirm_at"`                   //确认时间
	ConfirmUID        string  `db:"confirm_uid" json:"confirm_uid"`                 //确认人uid
	ConfirmName       string  `db:"confirm_name" json:"confirm_name"`               //确认人名
	ReviewRemark      string  `db:"review_remark" json:"review_remark"`             //确认人名
	WithdrawAt        int64   `db:"withdraw_at" json:"withdraw_at"`                 //三方场馆ID
	WithdrawRemark    string  `db:"withdraw_remark" json:"withdraw_remark"`         //确认人名
	WithdrawUID       string  `db:"withdraw_uid" json:"withdraw_uid"`               //确认人uid
	WithdrawName      string  `db:"withdraw_name" json:"withdraw_name"`             //确认人名
	FinanceType       int     `db:"finance_type" json:"finance_type"`               // 财务类型 442=提款 444=代客提款 446=代理提款
	LastDepositAmount float64 `db:"last_deposit_amount" json:"last_deposit_amount"` // 上笔成功存款金额
	RealNameHash      string  `db:"real_name_hash" json:"real_name_hash"`           //真实姓名哈希
	HangUpUID         string  `db:"hang_up_uid" json:"hang_up_uid"`                 // 挂起人uid
	HangUpRemark      string  `db:"hang_up_remark" json:"hang_up_remark"`           // 挂起备注
	HangUpName        string  `db:"hang_up_name" json:"hang_up_name"`               //  挂起人名字
	RemarkID          int     `db:"remark_id" json:"remark_id"`                     // 挂起原因ID
	HangUpAt          int     `db:"hang_up_at" json:"hang_up_at"`                   //  挂起时间
	ReceiveAt         int64   `db:"receive_at" json:"receive_at"`                   //领取时间
	WalletFlag        int     `db:"wallet_flag" json:"wallet_flag"`                 //钱包类型:1=中心钱包,2=佣金钱包
}

// Deposit 存款
type Deposit struct {
	ID              string  `db:"id" json:"id"`                               //
	Prefix          string  `db:"prefix" json:"prefix"`                       //转账后的金额
	OID             string  `db:"oid" json:"oid"`                             //转账前的金额
	Level           int     `db:"level" json:"level"`                         //
	UID             string  `db:"uid" json:"uid"`                             //用户ID
	ParentUID       string  `db:"parent_uid" json:"parent_uid"`               //上级代理ID
	ParentName      string  `db:"parent_name" json:"parent_name"`             //上级代理
	TopUID          string  `db:"top_uid" json:"top_uid"`                     //总代uid
	TopName         string  `db:"top_name" json:"top_name"`                   //总代
	Username        string  `db:"username" json:"username"`                   //用户名
	ChannelID       string  `db:"channel_id" json:"channel_id"`               //
	CID             string  `db:"cid" json:"cid"`                             //分类ID
	PID             string  `db:"pid" json:"pid"`                             //用户ID
	FinanceType     int     `db:"finance_type" json:"finance_type"`           //
	Amount          float64 `db:"amount" json:"amount"`                       //金额
	USDTFinalAmount float64 `db:"usdt_final_amount" json:"usdt_final_amount"` // 到账金额
	USDTApplyAmount float64 `db:"usdt_apply_amount" json:"usdt_apply_amount"` // 提单金额
	Rate            float64 `db:"rate" json:"rate"`                           // 汇率
	State           int     `db:"state" json:"state"`                         //0:待确认:1存款成功2:已取消
	Automatic       int     `db:"automatic" json:"automatic"`                 //1:自动转账2:脚本确认3:人工确认
	CreatedAt       int64   `db:"created_at" json:"created_at"`               //
	CreatedUID      string  `db:"created_uid" json:"created_uid"`             //三方场馆ID
	CreatedName     string  `db:"created_name" json:"created_name"`           //确认人名
	ReviewRemark    string  `db:"review_remark" json:"review_remark"`         //确认人名
	ConfirmAt       int64   `db:"confirm_at" json:"confirm_at"`               //确认时间
	ConfirmUID      string  `db:"confirm_uid" json:"confirm_uid"`             //确认人uid
	ConfirmName     string  `db:"confirm_name" json:"confirm_name"`           //确认人名
	IsRisk          int     `db:"-" json:"is_risk"`                           //是否风控
	ProtocolType    string  `db:"protocol_type" json:"protocol_type"`         //地址类型 trc20 erc20
	Address         string  `db:"address" json:"address"`                     //收款地址
	HashID          string  `db:"hash_id" json:"hash_id"`                     //区块链订单号
	Flag            int     `db:"flag" json:"flag"`                           // 1 三方订单 2 三方usdt订单 3 线下转卡订单 4 线下转usdt订单
	BankcardID      string  `db:"bankcard_id" json:"bankcard_id"`             // 线下转卡 收款银行卡id
	ManualRemark    string  `db:"manual_remark" json:"manual_remark"`         // 线下转卡订单附言
	BankCode        string  `db:"bank_code" json:"bank_code"`                 // 银行编号
	BankNo          string  `db:"bank_no" json:"bank_no"`                     // 银行卡号
}

// 取款数据
type FWithdrawData struct {
	T   int64             `json:"t"`
	D   []Withdraw        `json:"d"`
	Agg map[string]string `json:"agg"`
}

// 存款数据
type FDepositData struct {
	T   int64             `json:"t"`
	D   []Deposit         `json:"d"`
	Agg map[string]string `json:"agg"`
}

type withdrawCols struct {
	Withdraw
	MemberBankID       string `json:"member_bank_id"`
	MemberBankNo       string `json:"member_bank_no"`
	MemberBankRealName string `json:"member_bank_real_name"`
	MemberBankAddress  string `json:"member_bank_address"`
	MemberRealName     string `json:"member_real_name"`
	IsRisk             int    `json:"is_risk"`
}

type WithdrawListData struct {
	T   int64             `json:"t"`
	D   []withdrawCols    `json:"d"`
	Agg map[string]string `json:"agg"`
}

// 返水数据
type RebateData struct {
	T   int64             `json:"t"`
	D   []Transaction     `json:"d"`
	Agg map[string]string `json:"agg"`
}

// 代理团队转代
type AgencyTransfer struct {
	ID           string `json:"id" db:"id"`
	Prefix       string `json:"prefix" db:"prefix"`
	UID          string `json:"uid" db:"uid"`
	Username     string `json:"username" db:"username"`
	BeforeUid    string `json:"before_uid" db:"before_uid"`
	BeforeName   string `json:"before_name" db:"before_name"`
	AfterUid     string `json:"after_uid" db:"after_uid"`
	AfterName    string `json:"after_name" db:"after_name"`
	Status       int    `json:"status" db:"status"`
	ApplyAt      uint32 `json:"apply_at" db:"apply_at"`
	ApplyUid     string `json:"apply_uid" db:"apply_uid"`
	ApplyName    string `json:"apply_name" db:"apply_name"`
	ReviewAt     uint32 `json:"review_at" db:"review_at"`
	ReviewUid    string `json:"review_uid" db:"review_uid"`
	ReviewName   string `json:"review_name" db:"review_name"`
	Remark       string `json:"remark" db:"remark"`
	ReviewRemark string `json:"review_remark" db:"review_remark"`
}

type AgencyTransferData struct {
	T int64            `json:"t"`
	D []AgencyTransfer `json:"d"`
}

type AgencyTransferRecord struct {
	Id            string `json:"id" db:"id"`
	Flag          int    `json:"flag" db:"flag"`
	Uid           string `json:"uid" db:"uid"`
	Username      string `json:"username" db:"username"`
	Type          int64  `json:"type" db:"type"`
	BeforeUid     string `json:"before_uid" db:"before_uid"`
	BeforeName    string `json:"before_name" db:"before_name"`
	AfterUid      string `json:"after_uid" db:"after_uid"`
	AfterName     string `json:"after_name" db:"after_name"`
	Remark        string `json:"remark" db:"remark"`
	UpdatedAt     int64  `json:"updated_at" db:"updated_at"`
	UpdatedUid    string `json:"updated_uid" db:"updated_uid"`
	UpdatedName   string `json:"updated_name" db:"updated_name"`
	BeforeTopUid  string `json:"before_top_uid" db:"before_top_uid"`
	BeforeTopName string `json:"before_top_name" db:"before_top_name"`
	AfterTopUid   string `json:"after_top_uid" db:"after_top_uid"`
	AfterTopName  string `json:"after_top_name" db:"after_top_name"`
	Prefix        string `json:"prefix" db:"prefix"`
}

type AgencyTransferRecordData struct {
	T int64                  `json:"t"`
	D []AgencyTransferRecord `json:"d"`
}

type SmsRecord struct {
	Username  string `json:"username"`
	IP        string `json:"ip"`
	CreateAt  uint64 `json:"create_at"`
	Code      string `json:"code"`
	Phone     string `json:"phone"`
	PhoneHash string `json:"phone_hash"`
}

type SmsRecordData struct {
	T   int64             `json:"t"`
	D   []SmsRecord       `json:"d"`
	Agg map[string]string `json:"agg"`
}

type Promote struct {
	UID                string  `json:"uid"`                  //代理id
	Username           string  `json:"username"`             //代理名
	URL                string  `json:"url"`                  //域名
	CallNum            int64   `json:"call_num"`             //访问数
	IpNum              int64   `json:"ip_num"`               //独立ip数
	RegNum             int64   `json:"reg_num"`              //注册数
	RegRatio           float64 `json:"reg_ratio"`            //注册/访问百分比（独立ip）
	FirstDepositNum    int64   `json:"first_deposit_num"`    //首存人数
	FirstDepositRatio  float64 `json:"first_deposit_ratio"`  //首存人数/注册人数百分比
	FirstDepositAmount float64 `json:"first_deposit_amount"` //首存金额
	DepositNum         int64   `json:"deposit_num"`          //存款人数
	DepositAmount      float64 `json:"deposit_amount"`       //存款金额
}

type PromoteData struct {
	T   int64             `json:"t"`
	D   []Promote         `json:"d"`
	Agg map[string]string `json:"agg"`
}

type LinkAgency struct {
	UID  string `json:"uid" redis:"uid"`
	Name string `json:"name" redis:"name"`
}

type depositData struct {
	ParentUID string          `json:"parent_uid" db:"parent_uid"`
	Num       sql.NullInt64   `json:"num" db:"num"`
	Amount    sql.NullFloat64 `json:"amount" db:"amount"`
}

type UrlAgencyCount struct {
	ParentUID          string          `json:"parent_uid" db:"parent_uid"`
	Num                sql.NullInt64   `json:"num" db:"num"`
	FirstDepositAmount sql.NullFloat64 `json:"first_deposit_amount" db:"first_deposit_amount"`
	FirstDepositNum    sql.NullInt64   `json:"first_deposit_num" db:"first_deposit_num"`
}

type TgIp struct {
	Id         string `json:"id"`
	UID        string `json:"uid"`
	Username   string `json:"username"`
	RemoteAddr string `json:"remote_addr"`
	HttpHost   string `json:"hostdomain"`
	TimeIso    string `json:"time_iso8601"`
	RequestUri string `json:"request_uri"`
}

type TgMember struct {
	Uid        string `json:"uid" db:"uid"`
	UserName   string `json:"username" db:"username"`
	ParentUID  string `json:"parent_uid" db:"parent_uid"`
	ParentName string `json:"parent_name" db:"parent_name"`
	RegIP      string `json:"regip" db:"regip"`
	CreatedAt  string `json:"created_at" db:"created_at"`
}

type TgMemberData struct {
	T   int64             `json:"t"`
	S   uint              `json:"s"`
	D   []TgMember        `json:"d"`
	Agg map[string]string `json:"agg"`
}

const (
	UrlTyOfficial = 1
	UrlTyAgency   = 2
	UrlTyGeneral  = 3
	UrlTyTg       = 4
)

type TgIpData struct {
	T   int64             `json:"t"`
	S   uint              `json:"s"`
	D   []TgIp            `json:"d"`
	Agg map[string]string `json:"agg"`
}

type MessageTD struct {
	Ts        string `json:"ts" db:"ts"`                 //会员站内信id
	MessageID string `json:"message_id" db:"message_id"` //站内信id
	Username  string `json:"username" db:"username"`     //会员名
	Title     string `json:"title" db:"title"`           //标题
	Content   string `json:"content" db:"content"`       //内容
	IsTop     int    `json:"is_top" db:"is_top"`         //0不置顶 1置顶
	IsVip     int    `json:"is_vip" db:"is_vip"`         //0非vip站内信 1vip站内信
	Ty        int    `json:"ty" db:"ty"`                 //1站内消息 2活动消息
	IsRead    int    `json:"is_read" db:"is_read"`       //是否已读 0未读 1已读
	SendName  string `json:"send_name" db:"send_name"`   //发送人名
	SendAt    int64  `json:"send_at" db:"send_at"`       //发送时间
	Prefix    string `json:"prefix" db:"prefix"`         //商户前缀
}

type MessageTDData struct {
	T int64       `json:"t"`
	S int         `json:"s"`
	D []MessageTD `json:"d"`
}

// 站内信
type Message struct {
	ID         string ` db:"id" json:"id"`
	Title      string `db:"title" json:"title"`             //标题
	Content    string `db:"content" json:"content"`         //内容
	IsTop      int    `db:"is_top" json:"is_top"`           //0不置顶 1置顶
	IsPush     int    `db:"is_push" json:"is_push"`         //0不推送 1推送
	IsVip      int    `db:"is_vip" json:"is_vip"`           //0非vip站内信 1vip站内信
	Ty         int    `db:"ty" json:"ty"`                   //1站内消息 2活动消息
	Level      string `db:"level" json:"level"`             //会员等级
	Usernames  string `db:"usernames" json:"usernames"`     //会员账号，多个用逗号分隔
	State      int    `db:"state" json:"state"`             //1审核中 2审核通过 3审核拒绝 4已删除
	SendState  int    `db:"send_state" json:"send_state"`   //1未发送 2已发送
	SendName   string `db:"send_name" json:"send_name"`     //发送人名
	SendAt     int64  `db:"send_at" json:"send_at"`         //发送时间
	ApplyAt    uint32 `db:"apply_at" json:"apply_at"`       //创建时间
	ApplyUid   string `db:"apply_uid" json:"apply_uid"`     //创建人uid
	ApplyName  string `db:"apply_name" json:"apply_name"`   //创建人名
	ReviewAt   uint32 `db:"review_at" json:"review_at"`     //创建时间
	ReviewUid  string `db:"review_uid" json:"review_uid"`   //创建人uid
	ReviewName string `db:"review_name" json:"review_name"` //创建人名
	Prefix     string `db:"prefix" json:"prefix"`           //商户前缀
}

type MessageData struct {
	T int64     `json:"t"`
	S uint      `json:"s"`
	D []Message `json:"d"`
}

type smsLog struct {
	Username  string `json:"username"`
	IP        string `json:"ip"`
	CreateAt  uint64 `json:"create_at"`
	Code      string `json:"code"`
	Phone     string `json:"phone"`
	PhoneHash string `json:"phone_hash"`
}

type smsData struct {
	T int64    `json:"t"`
	D []smsLog `json:"d"`
}

type VncpOrder struct {
	Username  string `json:"username" db:"username"`     //用户名
	PayAmount string `json:"pay_amount" db:"pay_amount"` //跟注金额
	Bonus     string `json:"bonus" db:"bonus"`           //奖金
	NetAmount string `json:"net_amount" db:"net_amount"` //输赢
}

// 注单数据
type OrderData struct {
	T int64       `json:"t"`
	D []VncpOrder `json:"d"`
}

type WithdrawRecord struct {
	ID        string  `db:"id"                  json:"id"                 redis:"id"`
	Prefix    string  `db:"prefix"              json:"prefix"             redis:"prefix"`
	UID       string  `db:"uid"                 json:"uid"                redis:"uid"`        //
	Username  string  `db:"username"            json:"username"           redis:"username"`   //
	Amount    float64 `db:"amount"              json:"amount"             redis:"amount"`     // 提款金额
	State     int     `db:"state"               json:"state"              redis:"state"`      // 371:审核中 372:审核拒绝 373:出款中 374:提款成功 375:出款失败 376:异常订单 377:代付失败
	CreatedAt int64   `db:"created_at"          json:"created_at"         redis:"created_at"` //
}

// SMSChannel 短信通道
type SMSChannel struct {
	Id          string `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Alias       string `db:"alias" json:"alias"`
	Txt         string `db:"txt" json:"txt"`               // 0:没有1:开启2:关闭
	Voice       string `db:"voice" json:"voice"`           // 0:没有1:开启2:关闭
	CreatedAt   int64  `db:"created_at" json:"created_at"` // 创建时间
	CreatedUid  string `db:"created_uid" json:"created_uid"`
	CreatedName string `db:"created_name" json:"created_name"`
	Prefix      string `db:"prefix" json:"prefix"` // 前缀
	Remark      string `db:"remark" json:"remark"` // 备注
}

type LinkData struct {
	T int64    `json:"t"`
	S uint     `json:"s"`
	D []Link_t `json:"d"`
}

type Link_t struct {
	ID               string `db:"id" json:"id"`
	UID              string `db:"uid" json:"uid"`
	Username         string `db:"username" json:"username"`
	ShortURL         string `db:"short_url" json:"short_url"`
	Prefix           string `db:"prefix" json:"prefix"`
	NoAd             int    `db:"no_ad" json:"no_ad"`                           //0展示广告页，1不展示广告页
	ZR               string `db:"zr" json:"zr"`                                 //真人返水
	QP               string `db:"qp" json:"qp"`                                 //棋牌返水
	TY               string `db:"ty" json:"ty"`                                 //体育返水
	DJ               string `db:"dj" json:"dj"`                                 //电竞返水
	DZ               string `db:"dz" json:"dz"`                                 //电子返水
	CP               string `db:"cp" json:"cp"`                                 //彩票返水
	FC               string `db:"fc" json:"fc"`                                 //斗鸡返水
	BY               string `db:"by" json:"by"`                                 //捕鱼返水
	CGHighRebate     string `db:"cg_high_rebate" json:"cg_high_rebate"`         //cg高频彩返点
	CGOfficialRebate string `db:"cg_official_rebate" json:"cg_official_rebate"` //cg高频彩返点
	CreatedAt        string `db:"created_at" json:"created_at"`
}

type PersonalRebateReportData struct {
	D []MemberPersonalRebate `json:"d"`
	T int64                  `json:"t"`
	S uint16                 `json:"s"`
}

type MemberPersonalRebate struct {
	RebateDate       string
	ZR               decimal.Decimal //2真人返水
	QP               decimal.Decimal //5棋牌返水
	TY               decimal.Decimal //3体育返水
	DJ               decimal.Decimal //8电竞返水
	DZ               decimal.Decimal //6电游返水
	CP               decimal.Decimal //1彩票返水
	FC               decimal.Decimal //4斗鸡返水
	BY               decimal.Decimal //7捕鱼返水
	CGHighRebate     decimal.Decimal //CG高频彩返点
	CGOfficialRebate decimal.Decimal //CG官方彩返点
	TotalRebate      decimal.Decimal //总计
	Prefix           string
}

type RebateReportItem struct {
	Username       string `json:"username" db:"username"`
	Uid            string `json:"uid" db:"uid"`
	CashType       int    `json:"cash_type" db:"cash_type"`
	Rate           string `json:"rate"`
	ValidBetAmount string `json:"valid_bet_amount" db:"valid_bet_amount"`
	RebateAmount   string `json:"rebate_amount" db:"rebate_amount"`
	ReportTime     string `json:"report_time" db:"report_time"`
	Level          int    `json:"level"`
	ParentName     string `json:"parent_name"`
	TopName        string `json:"top_name"`
}

type GameResult_t struct {
	NetAmount      sql.NullFloat64 `db:"net_amount" json:"net_amount"`
	ValidBetAmount sql.NullFloat64 `db:"valid_bet_amount" json:"valid_bet_amount"`
}
