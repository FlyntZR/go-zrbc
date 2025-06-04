package view

import (
	"github.com/shopspring/decimal"
)

// swagger:model
type Member struct {
	// 主键
	ID int64 `json:"id"`
	// 用户名
	User string `json:"user"`
	// 用户密码
	Mem003 string `json:"-"`
	// 用户名称
	UserName string `json:"username"`
	Mem005   int64  `json:"mem005"`
	Mem006   int    `json:"mem006"`
	Mem007   int    `json:"mem007"`
	Mem008   int    `json:"mem008"`
	Mem009   int    `json:"mem009"`
	Mem010   int    `json:"mem010"`
	Mem011   int    `json:"mem011"`
	Mem012   int    `json:"mem012"`
	Mem013   int64  `json:"mem013"`
	Mem014   string `json:"mem014"`
	// login_error
	Mem015 int `json:"mem015"`
	// enable
	Mem016 string `json:"mem016"`
	// canbet
	Mem017 string `json:"mem017"`
	// chg_pw
	Mem018 string `json:"mem018"`
	// is_test
	Mem019 string `json:"mem019"`
	// be_traded
	Mem020 string `json:"mem020"`
	Mem021 int    `json:"mem021"`
	// 電話
	Mem022 string `json:"mem022"`
	// 電話簡碼
	Mem022a int `json:"mem022a"`
	Mem023  int `json:"mem023"`
	Mem024  int `json:"mem024"`
	// mem_risk
	Mem026 int `json:"mem026"`
	// 備註
	Mem028 string `json:"mem028"`
	// 0:現金;1:信用;2:電投
	Type int `json:"type"`
	// 幣別
	Currency int `json:"currency"`
	// 現金
	Cash decimal.Decimal `json:"cash"`
	// 己出碼額度
	Money decimal.Decimal `json:"money"`
	// 鎖定金額
	Lockmoney decimal.Decimal `json:"lockmoney"`
	// 頭像ID
	Head int `json:"head"`
	// 籌碼選擇
	Chips string `json:"chips"`
	// 注關會員
	Follow1 string `json:"follow1"`
	// 關注荷官
	Follow2 string `json:"follow2"`
	Tip     string `json:"tip"`
	// 红包开关
	Red    string `json:"red"`
	Wallet string `json:"wallet"`
	// 輸入要開啟的種類
	Opengame string `json:"opengame"`
	Site     string `json:"site"`
	// Line ID
	Lineid string `json:"lineid"`
	// 踢出局數
	Kickperiod int `json:"kickperiod"`
	// 0:信用,1:api
	Identity int `json:"identity"`
	// 單注金額告警
	Singlebetprompt decimal.Decimal `json:"singlebetprompt"`
	// 連贏局數告警
	Conwinprompt int `json:"conwinprompt"`
	// 最大可輸可贏提示
	Winlossprompt string `json:"winlossprompt"`
	// 上線提示
	Onlineprompt string `json:"onlineprompt"`
	// 盈利告警
	Profitprompt int `json:"profitprompt"`
}

// swagger:parameters GetUserInfo
type GetUserInfoReq struct {
	// in:header
	Token string `json:"Authorization"`
}

// swagger:model
type GetUserInfoResp struct {
	User *Member `json:"data"`
}

// swagger:parameters SigninGame
type SigninGameReq struct {
	// in:formData
	VendorID string `form:"vendorId" json:"vendorId"`
	// in:formData
	Signature string `form:"signature" json:"signature"`
	// in:formData
	User string `form:"user" json:"user"`
	// in:formData
	Device string `form:"device" json:"device"`
	// in:formData
	Lang string `form:"lang" json:"lang"`
	// in:formData
	IsTest bool `form:"isTest" json:"isTest"`
	// in:formData
	Mode string `form:"mode" json:"mode"`
	// in:formData
	TableID string `form:"tableid" json:"tableid"`
	// in:formData
	Site string `form:"site" json:"site"`
	// in:formData
	Password string `form:"password" json:"password"`
	// in:formData
	GameType string `form:"gameType" json:"gameType"`
	// in:formData
	Width string `form:"width" json:"width"`
	// in:formData
	ReturnURL string `form:"returnurl" json:"returnurl"`
	// in:formData
	Size string `form:"size" json:"size"`
	// in:formData
	UI string `form:"ui" json:"ui"`
	// in:formData
	Mute string `form:"mute" json:"mute"`
	// in:formData
	Video string `form:"video" json:"video"`
}

// swagger:model
type Session struct {
	SID string
}

// swagger:model
type SigninGameResp struct {
	GameURL string `json:"gameURL"`
}
