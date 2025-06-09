package view

import (
	"github.com/shopspring/decimal"
)

// MemberInfo represents the member information from JSON data
type MemberCache struct {
	UID      int64  `json:"uid"`
	Account  string `json:"acco"`
	Password string `json:"pswd"`
	Enable   string `json:"enable"`
	Mem007   int64  `json:"mem007"`
	Mem008   int64  `json:"mem008"`
	Mem009   int64  `json:"mem009"`
	Mem010   int64  `json:"mem010"`
	Mem011   int64  `json:"mem011"`
	SN       int64  `json:"sn"`
	Name     string `json:"name"`
	ULV      int    `json:"ulv"`
	ENS      string `json:"ens"`
	LogFail  int    `json:"logfail"`
	UTP      string `json:"utp"`
}

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
	Mem007   int64  `json:"mem007"`
	Mem008   int64  `json:"mem008"`
	Mem009   int64  `json:"mem009"`
	Mem010   int64  `json:"mem010"`
	Mem011   int64  `json:"mem011"`
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
	// 代理商(aid)
	// in:formData
	VendorID string `form:"vendorId" json:"vendorId"`
	// 代理商标识符
	// in:formData
	Signature string `form:"signature" json:"signature"`
	// 帐号
	// in:formData
	User string `form:"user" json:"user"`
	// 设备
	// in:formData
	Device string `form:"device" json:"device"`
	// 语言 使用语言 0或空值 为简体中文
	// 1为英文
	// 2为泰文
	// 3为越文
	// 4为日文
	// 5为韩文
	// 6为印度文
	// 7为马来西亚文
	// 8为印尼文
	// 9为繁体中文
	// 10为西文
	// 11为俄文
	// in:formData
	Lang string `form:"lang" json:"lang"`
	// 试玩（非必要）
	// in:formData
	IsTest bool `form:"isTest" json:"isTest"`
	// 游戏模式（非必要）onlybac 略过大厅直接进入百家乐；onlydgtg 略过大厅直接进入龙虎；onlyrou 略过大厅直接进入轮盘；onlysicbo 略过大厅直接进入骰宝；onlyniuniu 略过大厅直接进入牛牛；onlyfantan 略过大厅直接进入番摊；onlysedie 略过大厅直接进入色碟；下方數值務必帶入參數tableid
	// onlybactable 略过大厅直接进入百家乐；onlydgtgtable 略过大厅直接进入龙虎；onlyroutable 略过大厅直接进入轮盘；onlysicbotable 略过大厅直接进入骰宝；onlyniuniutable 略过大厅直接进入牛牛；onlyfantantable 略过大厅直接进入番摊；onlysedietable 略过大厅直接进入色碟
	// in:formData
	Mode string `form:"mode" json:"mode"`
	// 桌子IDmode参数带入以下数值为必填
	// onlybactable
	// onlydgtgtable
	// onlyroutable
	// onlysicbotable
	// onlyniuniutable
	// onlyfantantable
	// onlysedietable
	// 否则非必要
	// in:formData
	TableID string `form:"tableid" json:"tableid"`
	// 站点
	// in:formData
	Site string `form:"site" json:"site"`
	// in:formData
	Password string `form:"password" json:"password"`
	// 游戏类型
	// in:formData
	GameType string `form:"gameType" json:"gameType"`
	// 宽度
	// in:formData
	Width string `form:"width" json:"width"`
	// 返回URL
	// in:formData
	ReturnURL string `form:"returnurl" json:"returnurl"`
	// 1:iframe嵌入 (非必要)
	// in:formData
	Size string `form:"size" json:"size"`
	// 用户ID
	// in:formData
	UI int `form:"ui" json:"ui"`
	// 静音
	// in:formData
	Mute string `form:"mute" json:"mute"`
	// 视频
	// in:formData
	Video string `form:"video" json:"video"`
	// 0:中文, 1:英文 (非必要)
	// in:formData
	Syslang int `form:"video" json:"syslang"`
	// 籌碼選擇
	Chips string `form:"chips" json:"chip"`
	// 时间戳
	// in:formData
	Timestamp int64 `form:"timestamp" json:"timestamp"`
}

// swagger:model
type Session struct {
	SID string
}

// swagger:model
type SigninGameResp struct {
	Result string `json:"result"`
}

// swagger:parameters MemberRegister
type MemberRegisterReq struct {
	// 代理商(aid)
	// in:formData
	VendorID string `json:"vendorId"`
	// 代理商标识符
	// in:formData
	Signature string `json:"signature"`
	// 帐号
	// in:formData
	User string `json:"user"`
	// in:formData
	Password string `json:"password"`
	// 用户名称
	// in:formData
	Username string `json:"username"`
	// 用户类型（非必要）
	// in:formData
	Profile int `json:"profile"`
	// 最大赢（非必要）
	// in:formData
	Maxwin int64 `json:"maxwin"`
	// 最大输（非必要）
	// in:formData
	Maxlose int64 `json:"maxlose"`
	// 备注（非必要）
	// in:formData
	Mark string `json:"mark"`
	// 會員退水是否歸零   0為: 不歸零 1為: 歸零  (非必要)
	// in:formData
	Rakeback int `json:"rakeback"`
	// 限制类型 (非必要)例:2,5,9,14,48,107,111,131 (非必要)
	// (此为风险控管的一环，请经营者多善利用)
	// http://wmapi.a45.me/Limit.html
	// in:formData
	LimitType string `json:"limitType"`
	// 籌碼選擇会员筹码 使用逗号隔开，可填入5-10组 (非必要)
	// 若没带入，会员筹码预设为 10,20,50,100,500,1000,5000,10000,20000,50000
	// 可用筹码种类：1, 5, 10, 20, 50, 100, 500, 1000, 5000, 10000, 20000, 50000, 100000, 200000, 1000000, 5000000, 10000000, 20000000, 50000000, 10000000, 20000000, 50000000, 100000000
	// in:formData
	Chips string `json:"chip"`
	// 提出局数（非必要）
	// in:formData
	Kickperiod string `json:"kickperiod"`
	// 时间戳
	// in:formData
	Timestamp int64 `json:"timestamp"`
	// 0:中文, 1:英文 (非必要)
	// in:formData
	Syslang int `json:"syslang"`
}

// swagger:model
type MemberRegisterResp struct {
	Result string `json:"result"`
}

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
