package view

import (
	"github.com/shopspring/decimal"
)

// swagger:parameters ChangePassword
type ChangePasswordReq struct {
	// 代理商(aid)
	// in:formData
	VendorID string `json:"vendorId"`
	// 代理商标识符
	// in:formData
	Signature string `json:"signature"`
	// 帐号
	// in:formData
	User string `json:"user"`
	// 新密码
	// in:formData
	NewPassword string `json:"newpassword"`
	// 时间戳
	// in:formData
	Timestamp int64 `json:"timestamp"`
	// 0:中文, 1:英文 (非必要)
	// in:formData
	SyslangStr string `json:"syslang" form:"syslang"`
	// swagger:ignore
	Syslang string
}

// swagger:model
type ChangePasswordResp struct {
	Result string `json:"result"`
}

// swagger:parameters EditLimit
type EditLimitReq struct {
	// 代理商(aid)
	// in:formData
	VendorID string `json:"vendorId"`
	// 代理商标识符
	// in:formData
	Signature string `json:"signature"`
	// 帐号
	// in:formData
	User string `json:"user"`
	// 限制类型 (非必要)例:2,5,9,14,48,107,111,131 (非必要)
	// in:formData
	LimitType string `json:"limitType"`
	// 最大赢（非必要）
	// in:formData
	Maxwin int64 `json:"maxwin"`
	// 最大输（非必要）
	// in:formData
	Maxlose int64 `json:"maxlose"`
	// 重置时间
	// in:formData
	Reset int `json:"reset"`
	// 时间戳
	// in:formData
	Timestamp int64 `json:"timestamp"`
	// 0:中文, 1:英文 (非必要)
	// in:formData
	SyslangStr string `json:"syslang" form:"syslang"`
	// swagger:ignore
	Syslang string
}

// swagger:model
type EditLimitResp struct {
	Result string `json:"result"`
}

type UserDtl struct {
	Type      int    `json:"type"`
	Bkwater   string `json:"bkwater"`
	Rate      string `json:"rate"`
	Betlimit  string `json:"betlimit"`
	Maxwin    string `json:"maxwin"`
	Maxlose   string `json:"maxlose"`
	Nbetlimit string `json:"nbetlimit"`
	Netwater  string `json:"netwater"`
	Netrate   string `json:"netrate"`
}

// swagger:parameters LogoutGame
type LogoutGameReq struct {
	// 代理商(aid)
	// in:formData
	VendorID string `json:"vendorId" form:"vendorId"`
	// 代理商标识符
	// in:formData
	Signature string `json:"signature" form:"signature"`
	// 帐号
	// in:formData
	User string `json:"user" form:"user"`
	// 时间戳
	// in:formData
	Timestamp int64 `json:"timestamp" form:"timestamp"`
	// 0:中文, 1:英文 (非必要)
	// in:formData
	SyslangStr string `json:"syslang" form:"syslang"`
	// swagger:ignore
	Syslang string
}

// swagger:model
type LogoutGameResp struct {
	Result string `json:"result"`
}

//swagger:parameters GetBalance
type GetBalanceReq struct {
	// 代理商(aid)
	// in:formData
	VendorID string `json:"vendorId" form:"vendorId"`
	// 代理商标识符
	// in:formData
	Signature string `json:"signature"`
	// 会员(user)
	// in:formData
	User string `json:"user" form:"user"`
	// 时间戳
	// in:formData
	Timestamp int64 `json:"timestamp"`
	// 0:中文, 1:英文 (非必要)
	// in:formData
	SyslangStr string `json:"syslang" form:"syslang"`
	// swagger:ignore
	Syslang string
}

// swagger:model
type GetBalanceResp struct {
	Result decimal.Decimal `json:"result"`
}

// swagger:parameters ChangeBalance
type ChangeBalanceReq struct {
	// 代理商(aid)
	// in:formData
	VendorID string `json:"vendorId" form:"vendorId"`
	// 代理商标识符
	// in:formData
	Signature string `json:"signature" form:"signature"`
	// 会员(user)
	// in:formData
	User string `json:"user" form:"user"`
	// 加扣点金额(money)
	// in:formData
	Money string `json:"money" form:"money"`
	// 贵公司产生的订单序号(非必要)最大值:32字符
	// in:formData
	Order string `json:"order" form:"order"`
	// 时间戳
	// in:formData
	Timestamp int64 `json:"timestamp" form:"timestamp"`
	// 0:中文, 1:英文 (非必要)
	// in:formData
	SyslangStr string `json:"syslang" form:"syslang"`
	// swagger:ignore
	Syslang string
}

// swagger:model
type ChangeBalanceResp struct {
	Result string `json:"result"`
}

// swagger:parameters GetMemberTradeReport
type GetMemberTradeReportReq struct {
	// 代理商(aid)
	// in:formData
	VendorID string `json:"vendorId" form:"vendorId"`
	// 代理商标识符
	// in:formData
	Signature string `json:"signature" form:"signature"`
	// 会员(user)
	// in:formData
	User string `json:"user" form:"user"`
	// 订单ID
	// in:formData
	OrderID string `json:"orderid" form:"orderid"`
	// 订单号
	// in:formData
	Order string `json:"order" form:"order"`
	// 开始时间戳
	// in:formData
	StartTimeStr string `json:"startTime" form:"startTime"`
	// swagger:ignore
	StartTime int64
	// 结束时间戳
	// in:formData
	EndTimeStr string `json:"endTime" form:"endTime"`
	// swagger:ignore
	EndTime int64
	// 时间戳
	// in:formData
	Timestamp int64 `json:"timestamp" form:"timestamp"`
	// 0:中文, 1:英文 (非必要)
	// in:formData
	SyslangStr string `json:"syslang" form:"syslang"`
	// swagger:ignore
	Syslang string
}

type TradeItem struct {
	MID      int64           `json:"mid"`
	OrderID  int64           `json:"orderid"`
	OrderNum string          `json:"ordernum"`
	AddTime  int64           `json:"addtime"`
	Money    decimal.Decimal `json:"money"`
	OpCode   string          `json:"op_code"`
	Subtotal decimal.Decimal `json:"subtotal"`
}

// swagger:model
type GetMemberTradeReportResp struct {
	Result []*TradeItem `json:"result"` // Result data
}

// swagger:parameters GetTipReport
type GetTipReportReq struct {
	// 代理商(aid)
	// in:formData
	VendorID string `json:"vendorId" form:"vendorId"`
	// 代理商标识符
	// in:formData
	Signature string `json:"signature" form:"signature"`
	// 会员(user)
	// in:formData
	User string `json:"user" form:"user"`
	// 开始时间戳
	// in:formData
	StartTimeStr string `json:"startTime" form:"startTime"`
	// swagger:ignore
	StartTime int64
	// 结束时间戳
	// in:formData
	EndTimeStr string `json:"endTime" form:"endTime"`
	// swagger:ignore
	EndTime int64
	// 时间戳
	// in:formData
	Timestamp int64 `json:"timestamp" form:"timestamp"`
	// 0:中文, 1:英文 (非必要)
	// in:formData
	SyslangStr string `json:"syslang" form:"syslang"`
	// swagger:ignore
	Syslang string
}

type TipReportItem struct {
	BetID      string          `json:"betId"`
	ID         int64           `json:"id"`
	BetTime    string          `json:"betTime"`
	Tip        decimal.Decimal `json:"tip"`
	BetResult  string          `json:"betResult"`
	WinLoss    decimal.Decimal `json:"winLoss"`
	IP         string          `json:"ip"`
	GID        string          `json:"gid"`
	Event      string          `json:"event"`
	Round      string          `json:"round"`
	Subround   string          `json:"subround"`
	EventChild string          `json:"eventChild"`
	TableID    string          `json:"tableId"`
	Username   string          `json:"username"`
	GName      string          `json:"gname"`
}

// swagger:model
type GetTipReportResp struct {
	Result []*TipReportItem `json:"result"` // Result data
}

// swagger:parameters EnableOrDisableMem
type EnableOrDisableMemReq struct {
	// 代理商(aid)
	// in:formData
	VendorID string `json:"vendorId" form:"vendorId"`
	// 代理商标识符
	// in:formData
	Signature string `json:"signature" form:"signature"`
	// 会员账号，多个账号用逗号分隔
	// in:formData
	User string `json:"user" form:"user"`
	// 类型：login(登入) 或 bet(下注)
	// in:formData
	Type string `json:"type" form:"type"`
	// 状态：Y(启用) 或 N(停用)
	// in:formData
	Status string `json:"status" form:"status"`
	// 时间戳
	// in:formData
	Timestamp int64 `json:"timestamp" form:"timestamp"`
	// 0:中文, 1:英文 (非必要)
	// in:formData
	SyslangStr string `json:"syslang" form:"syslang"`
	// swagger:ignore
	Syslang string
}

// swagger:model
type EnableOrDisableMemResp struct {
	Result string `json:"result"`
}

// swagger:parameters GetDateTimeReport
type GetDateTimeReportReq struct {
	// 代理商(aid)
	// in:formData
	VendorID string `json:"vendorId" form:"vendorId"`
	// 代理商标识符
	// in:formData
	Signature string `json:"signature" form:"signature"`
	// 会员(user)
	// in:formData
	User string `json:"user" form:"user"`
	// 开始时间戳
	// in:formData
	StartTimeStr string `json:"startTime" form:"startTime"`
	// swagger:ignore
	StartTime int64
	// 结束时间戳
	// in:formData
	EndTimeStr string `json:"endTime" form:"endTime"`
	// swagger:ignore
	EndTime int64
	// 时间类型 0:下注时间 1:结算时间
	// in:formData
	TimeType int `json:"timetype" form:"timetype"`
	// 数据类型 0:一般数据 1:小费数据
	// in:formData
	DataType int `json:"datatype" form:"datatype"`
	// 游戏编号1
	// in:formData
	GameNo1 string `json:"gameno1" form:"gameno1"`
	// 游戏编号2
	// in:formData
	GameNo2 string `json:"gameno2" form:"gameno2"`
	// 时间戳
	// in:formData
	Timestamp int64 `json:"timestamp" form:"timestamp"`
	// 0:中文, 1:英文 (非必要)
	// in:formData
	SyslangStr string `json:"syslang" form:"syslang"`
	// swagger:ignore
	Syslang string
}

type DateTimeReportItem struct {
	User       string          `json:"user"`
	BetID      string          `json:"betId"`
	BetTime    string          `json:"betTime"`
	BeforeCash decimal.Decimal `json:"beforeCash"`
	Bet        decimal.Decimal `json:"bet"`
	ValidBet   decimal.Decimal `json:"validbet"`
	Water      decimal.Decimal `json:"water"`
	Result     decimal.Decimal `json:"result"`
	BetCode    string          `json:"betCode"`
	BetResult  string          `json:"betResult"`
	WaterBet   decimal.Decimal `json:"waterbet"`
	WinLoss    decimal.Decimal `json:"winLoss"`
	IP         string          `json:"ip"`
	GID        string          `json:"gid"`
	Event      string          `json:"event"`
	EventChild string          `json:"eventChild"`
	Round      string          `json:"round"`
	Subround   string          `json:"subround"`
	TableID    string          `json:"tableId"`
	Commission decimal.Decimal `json:"commission"`
	Settime    string          `json:"settime"`
	Reset      string          `json:"reset"`
	GameResult string          `json:"gameResult"`
	GName      string          `json:"gname"`
}

// swagger:model
type GetDateTimeReportResp struct {
	Result []*DateTimeReportItem `json:"result"`
}
