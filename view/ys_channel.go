package view

type WsReq struct {
	Protocol int         `json:"protocol"`
	Data     interface{} `json:"data"`
}

type WsResp struct {
	Protocol int         `json:"protocol"`
	Data     interface{} `json:"data"`
}

type AuthData struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type AuthResp struct {
	MemberID       int64  `json:"memberID"`
	Account        string `json:"account"`
	UserName       string `json:"userName"`
	Sid            string `json:"sid"`
	BOk            bool   `json:"bOk"`
	BValidPassword bool   `json:"bValidPassword"`
}

// swagger:model
type WsUser struct {
	// id
	ID int64 `json:"id" description:"会员id"`
}
