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
	Syslang int `json:"syslang"`
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
	Syslang int `json:"syslang"`
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
	Syslang int `json:"syslang" form:"syslang"`
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
	Syslang int `json:"syslang"`
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
	Syslang int `json:"syslang" form:"syslang"`
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
	Syslang int `json:"syslang" form:"syslang"`
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
	Syslang int `json:"syslang" form:"syslang"`
}

// swagger:model
type EnableOrDisableMemResp struct {
	Result string `json:"result"`
}
