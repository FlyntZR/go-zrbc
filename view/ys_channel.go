package view

import "github.com/gorilla/websocket"

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

type WsBettingInfoItem struct {
	BetArea     int `json:"betArea"`
	AddBetMoney int `json:"addBetMoney"`
}

type WsBettingData struct {
	BetSerialNumber int                 `json:"betSerialNumber"`
	GameNo          int                 `json:"gameNo"`
	GameNoRound     int                 `json:"gameNoRound"`
	BetArr          []WsBettingInfoItem `json:"betArr"`
	Commission      int                 `json:"commission"`
}

type TableDtExtend struct {
	NetGroupName   string `json:"netGroupName"`
	PhoneGroupName string `json:"phoneGroupName"`
	TableName      string `json:"tableName"`
	NetType        string `json:"netType"`
	PhoneType      string `json:"phoneType"`
}
type WsTableData struct {
	GameID                   int           `json:"gameID"`
	GroupID                  int           `json:"groupID"`
	GroupType                int           `json:"groupType"`
	GameNo                   int           `json:"gameNo"`
	GameNoRound              int           `json:"gameNoRound"`
	DealerID                 int           `json:"dealerID"`
	DealerName               string        `json:"dealerName"`
	DealerImage              string        `json:"dealerImage"`
	DealerImage2             string        `json:"dealerImage2"`
	Dealer2ID                int           `json:"dealer2ID"`
	Dealer2Name              string        `json:"dealer2Name"`
	Dealer2Image             string        `json:"dealer2Image"`
	Dealer2Image2            string        `json:"dealer2Image2"`
	BetMilliSecond           int           `json:"betMilliSecond"`
	BWantToShuffle           bool          `json:"bWantToShuffle"`
	BWantToEnd               bool          `json:"bWantToEnd"`
	KeyStatus                int           `json:"keyStatus"`
	GameMode                 int           `json:"gameMode"`
	SingleLimit              int           `json:"singleLimit"`
	TableMinBet              int           `json:"tableMinBet"`
	TableMaxBet              int           `json:"tableMaxBet"`
	TableTieMinBet           int           `json:"tableTieMinBet"`
	TableTieMaxBet           int           `json:"tableTieMaxBet"`
	TablePairMinBet          int           `json:"tablePairMinBet"`
	TablePairMaxBet          int           `json:"tablePairMaxBet"`
	TableType                int           `json:"tableType"`
	TableStatus              int           `json:"tableStatus"`
	TableSort                int           `json:"tableSort"`
	TableSort2               int           `json:"tableSort2"`
	ReservedTable            int           `json:"reservedTable"`
	ReservedTableParentIDArr []int         `json:"reservedTableParentIDArr"`
	ReservedTableMemberIDArr []int         `json:"reservedTableMemberIDArr"`
	TableDtExtend            TableDtExtend `json:"tableDtExtend"`
}

type WsTableData21 struct {
	Protocol int         `json:"protocol"`
	Data     WsTableData `json:"data"`
}

type WsJoinTableData struct {
	DtBetLimitSelectID map[string]int `json:"dtBetLimitSelectID"`
	GroupID            int            `json:"groupID"`
}

type WsBettingCh struct {
	Conn  *websocket.Conn
	BetCh WsBettingData
}

type WsBetTimeData struct {
	GameID          int                    `json:"gameID"`
	GroupID         int                    `json:"groupID"`
	GroupType       int                    `json:"groupType"`
	GameNo          int                    `json:"gameNo"`
	GameNoRound     int                    `json:"gameNoRound"`
	BetTimeCount    int                    `json:"betTimeCount"`
	BetTimeContent  map[string]interface{} `json:"betTimeContent"`
	TimeMillisecond int                    `json:"timeMillisecond"`
}

type WsBetTimeResp struct {
	Protocol int           `json:"protocol"`
	Data     WsBetTimeData `json:"data"`
}
