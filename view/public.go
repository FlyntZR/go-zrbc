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

// GetBalanceResp represents the response for GetBalance
type GetBalanceResp struct {
	Result decimal.Decimal `json:"result"`
}
