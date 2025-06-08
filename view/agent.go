package view

//swagger:model
type Agent struct {
	// sn
	ID int64 `json:"id"`
	// 代理商ID
	VendorID string `json:"vendorId"`
	// 代理商名称
	Name string `json:"name"`
	// 代理商密码
	Password string `json:"-"`
	// 创建时间
	CreatedAt int64 `json:"createdAt"`
	// 等级
	ULV int `json:"ulv"`
	// 上级代理商ID
	Age007 int64 `json:"age007"`
	// 上级代理商ID
	Age008 int64 `json:"age008"`
	// 上级代理商ID
	Age009 int64 `json:"age009"`
	// 上级代理商ID
	Age010 int64 `json:"age010"`
	// 上级代理商ID
	Age011 int64 `json:"age011"`
	// 风险重置开关
	RiskReset int `json:"riskReset"`
	// 启停用状态
	Status string `json:"status"`
	// 启停押状态
	BetStatus string `json:"betStatus"`
	// 结算报表开关
	ReportSwitch int `json:"reportSwitch"`
	// 结算报表格式
	ReportFormat int `json:"reportFormat"`
	// 结算报表语系
	ReportLang int `json:"reportLang"`
	// 测试线判别
	TestLine int `json:"testLine"`
	// 备注
	Remark string `json:"remark"`
	// 信用额度
	Credit string `json:"credit"`
	// 现金
	Cash string `json:"cash"`
	// 类型
	Type int `json:"type"`
	// 币别
	Currency int `json:"currency"`
	// 小费开关
	Tip string `json:"tip"`
	// 红包开关
	Red string `json:"red"`
	// 前缀码开关
	PrefixAdd string `json:"prefixAdd"`
	// 前缀码
	PrefixAcc string `json:"prefixAcc"`
	// 收款账号
	ReceiptAcc int `json:"receiptAcc"`
	// 通知账号
	Notification int `json:"notification"`
	// 显示账号
	Sacc string `json:"sacc"`
	// 开启游戏
	Opengame string `json:"opengame"`
	// 站点
	Site string `json:"site"`
	// 会员总数
	Membermax int `json:"membermax"`
	// 最大利润
	Profitmax string `json:"profitmax"`
	// 推广代码
	Promotecode string `json:"promotecode"`
	// 踢出局数
	Kickperiod int `json:"kickperiod"`
	// 0:信用,1:API
	Identity int `json:"identity"`
	// 检查密钥
	ChkKey string `json:"chkKey"`
	// 是否锁定
	ChkLock string `json:"chkLock"`
	// 上次修改密码时间
	LastChpsw int64 `json:"lastChpsw"`
	// 启停用状态
	Age015 string `json:"age015"`
	// 启停押状态
	Age016 string `json:"age016"`
}

//swagger:model
type AgentsLoginPass struct {
	// 主键ID
	ID int64 `json:"id"`
	// 代理商ID
	Aid int64 `json:"aid"`
	// lv5 代理商
	VendorID string `json:"vendorId"`
	// 密钥
	Signature string `json:"signature"`
	// 客户密钥
	Signature2 string `json:"signature2"`
	// 密码
	Password string `json:"-"`
	// 添加时间
	Addtime int64 `json:"addtime"`
	// 单一钱包回传网址
	URL string `json:"url"`
	// 天空名称
	Skyname string `json:"skyname"`
	// c:一般,w:单一
	Type string `json:"type"`
	// 0:中文 ,1:英文
	Lang int `json:"lang"`
	// 异常锁定,0:N;1:Y
	Betfeedback int `json:"betfeedback"`
	// 正机白名单使用的url
	GatewayURL string `json:"gatewayUrl"`
	// 白名单
	WhiteList string `json:"whiteList"`
	// 操作员
	Operator string `json:"operator"`
	// 修改时间
	ModifyTime int64 `json:"modifyTime"`
	// 0:呼叫php,1:呼叫客户
	Object int `json:"object"`
	// 单一钱包结算
	Settle int `json:"settle"`
	// 公司代码
	Co string `json:"co"`
	// 前缀码开关
	PrefixSwitch string `json:"prefixSwitch"`
	// 指定网址
	OpenGameURL string `json:"openGameUrl"`
	// 子域名
	Subdomain string `json:"subdomain"`
}

//swagger:model
type AgentVerifyReq struct {
	VendorID  string `json:"vendorId"`
	Signature string `json:"signature"`
}

//swagger:model
type AgentVerifyResp struct {
	Agent           *Agent           `json:"agent"`
	AgentsLoginPass *AgentsLoginPass `json:"agentsLoginPass"`
}
