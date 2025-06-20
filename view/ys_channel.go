package view

import (
	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
)

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

// WsYsChannalRespData represents the data field in the channel response
type WsBettingRespData struct {
	GameID   int  `json:"gameID"`
	GroupID  int  `json:"groupID"`
	AreaID   int  `json:"areaID"`
	AreaType int  `json:"areaType"`
	BOk      bool `json:"bOk"`
}

// WsYsChannalRespData represents the data field in the channel response
type WsBettingResp struct {
	Protocol int               `json:"protocol"`
	Data     WsBettingRespData `json:"data"`
}

// WsGameResultRespData represents the data field for game result response
// 根据 {"protocol":25,"data":{"gameID":101,"groupID":115,"result":193,"playerScore":4,"bankerScore":8,"dtCard":{"1":45,"2":32,"3":9,"4":48,"5":15,"6":15},"winBetAreaArr":[1,7,8,13]}} 结构定义
// swagger:model
type WsGameResultRespData struct {
	GameID        int            `json:"gameID"`
	GroupID       int            `json:"groupID"`
	Result        int            `json:"result"`
	PlayerScore   int            `json:"playerScore"`
	BankerScore   int            `json:"bankerScore"`
	DtCard        map[string]int `json:"dtCard"`
	WinBetAreaArr []int          `json:"winBetAreaArr"`
}

// WsGameResultResp is the response struct for game result
// swagger:model
type WsGameResultResp struct {
	Protocol int                  `json:"protocol"`
	Data     WsGameResultRespData `json:"data"`
}

// WsPayoutResultData represents the data field for payout result response
type WsPayoutResultData struct {
	GameID       int                        `json:"gameID"`
	GroupID      int                        `json:"groupID"`
	AreaID       int                        `json:"areaID"`
	AreaType     int                        `json:"areaType"`
	MemberID     int64                      `json:"memberID"`
	MoneyWin     decimal.Decimal            `json:"moneyWin"`
	MoneyWinLoss decimal.Decimal            `json:"moneyWinLoss"`
	DtMoneyWin   map[string]decimal.Decimal `json:"dtMoneyWin"`
}

// WsPayoutResp is the response struct for payout result
type WsPayoutResp struct {
	Protocol int                `json:"protocol"`
	Data     WsPayoutResultData `json:"data"`
}

// WsTableEntryResp represents the response when successfully entering a table
type WsTableEntryResp struct {
	Protocol int              `json:"protocol"`
	Data     WsTableEntryData `json:"data"`
}

type WsTableEntryData struct {
	BOk            bool                         `json:"bOk"`
	GameID         int                          `json:"gameID"`
	GroupID        int                          `json:"groupID"`
	GroupType      int                          `json:"groupType"`
	AreaID         int                          `json:"areaID"`
	AreaType       int                          `json:"areaType"`
	MemberID       int64                        `json:"memberID"`
	Account        string                       `json:"account"`
	UserName       string                       `json:"userName"`
	HeadID         int                          `json:"headID"`
	SeatIDArr      []int                        `json:"seatIDArr"`
	CurrencyCode   string                       `json:"currencyCode"`
	CurrencyName   string                       `json:"currencyName"`
	CurrencyRate   float64                      `json:"currencyRate"`
	Balance        float64                      `json:"balance"`
	Chips          float64                      `json:"chips"`
	BetMilliSecond int                          `json:"betMilliSecond"`
	DtOdds         map[string][]decimal.Decimal `json:"dtOdds"`
	UserCount      int                          `json:"userCount"`
	DtCard         map[string]interface{}       `json:"dtCard"`
	MinBet01       int                          `json:"minBet01"`
	MaxBet01       int                          `json:"maxBet01"`
	MinBet02       int                          `json:"minBet02"`
	MaxBet02       int                          `json:"maxBet02"`
	MinBet03       int                          `json:"minBet03"`
	MaxBet03       int                          `json:"maxBet03"`
	MinBet04       int                          `json:"minBet04"`
	MaxBet04       int                          `json:"maxBet04"`
	MinBet05       int                          `json:"minBet05"`
	MaxBet05       int                          `json:"maxBet05"`
	MinBet06       int                          `json:"minBet06"`
	MaxBet06       int                          `json:"maxBet06"`
	MinBet07       int                          `json:"minBet07"`
	MaxBet07       int                          `json:"maxBet07"`
	MinBet08       int                          `json:"minBet08"`
	MaxBet08       int                          `json:"maxBet08"`
	MinBet09       int                          `json:"minBet09"`
	MaxBet09       int                          `json:"maxBet09"`
	MinBet10       int                          `json:"minBet10"`
	MaxBet10       int                          `json:"maxBet10"`
	MinBet11       int                          `json:"minBet11"`
	MaxBet11       int                          `json:"maxBet11"`
	MinBet12       int                          `json:"minBet12"`
	MaxBet12       int                          `json:"maxBet12"`
	MinBet13       int                          `json:"minBet13"`
	MaxBet13       int                          `json:"maxBet13"`
	MinBet14       int                          `json:"minBet14"`
	MaxBet14       int                          `json:"maxBet14"`
}

// WsBetLimitModifyData represents the data structure for bet limit modification
type WsBetLimitModifyData struct {
	GameID             int            `json:"gameID"`
	MemberID           int64          `json:"memberID"`
	DtBetLimitSelectID map[string]int `json:"dtBetLimitSelectID"`
	MinBet01           int            `json:"minBet01"`
	MaxBet01           int            `json:"maxBet01"`
	MinBet02           int            `json:"minBet02"`
	MaxBet02           int            `json:"maxBet02"`
	MinBet03           int            `json:"minBet03"`
	MaxBet03           int            `json:"maxBet03"`
	MinBet04           int            `json:"minBet04"`
	MaxBet04           int            `json:"maxBet04"`
	MinBet05           int            `json:"minBet05"`
	MaxBet05           int            `json:"maxBet05"`
	MinBet06           int            `json:"minBet06"`
	MaxBet06           int            `json:"maxBet06"`
	MinBet07           int            `json:"minBet07"`
	MaxBet07           int            `json:"maxBet07"`
	MinBet08           int            `json:"minBet08"`
	MaxBet08           int            `json:"maxBet08"`
	MinBet09           int            `json:"minBet09"`
	MaxBet09           int            `json:"maxBet09"`
	MinBet10           int            `json:"minBet10"`
	MaxBet10           int            `json:"maxBet10"`
	MinBet11           int            `json:"minBet11"`
	MaxBet11           int            `json:"maxBet11"`
	MinBet12           int            `json:"minBet12"`
	MaxBet12           int            `json:"maxBet12"`
	MinBet13           int            `json:"minBet13"`
	MaxBet13           int            `json:"maxBet13"`
	MinBet14           int            `json:"minBet14"`
	MaxBet14           int            `json:"maxBet14"`
}

// WsBetLimitModifyResp represents the response structure for bet limit modification
type WsBetLimitModifyResp struct {
	Protocol int                  `json:"protocol"`
	Data     WsBetLimitModifyData `json:"data"`
}

// WsGroupInfo represents a single group item in the group array
type WsGroupInfo struct {
	GroupID         int `json:"groupID"`
	GroupType       int `json:"groupType"`
	SingleLimit     int `json:"singleLimit"`
	TableMinBet     int `json:"tableMinBet"`
	TableMaxBet     int `json:"tableMaxBet"`
	TableTieMinBet  int `json:"tableTieMinBet"`
	TableTieMaxBet  int `json:"tableTieMaxBet"`
	TablePairMinBet int `json:"tablePairMinBet"`
	TablePairMaxBet int `json:"tablePairMaxBet"`
	TableStatus     int `json:"tableStatus"`
}

// WsGroupListData represents the data field for group list response
type WsGroupListData struct {
	GroupArr []WsGroupInfo `json:"groupArr"`
}

// WsGroupListResp represents the response structure for group list (protocol 35)
type WsGroupListResp struct {
	Protocol int             `json:"protocol"`
	Data     WsGroupListData `json:"data"`
}
