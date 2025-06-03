package mq

const (
	// http
	MQMsgTypeBatchKickUser = 2000
	MQMsgTypeKline         = 2001
	MQMsgTypeOrder         = 2002
	MQMsgTypeTradeClearing = 2003
	MQMsgTypeAccountUpdate = 2004
	MQMsgTypeTrade         = 2005
	MQMsgTypeLastDay       = 2006
	MQMsgTypeDepth         = 2007
)

type MQMsg struct {
	// 1: purely broadcast, no server involve
	MsgType   int         `json:"msg_type"`
	AccountId int64       `json:"accountId"`
	Payload   interface{} `json:"payload"`
}

type BatchKickUser struct {
	UserIDs   []int64 `json:"u_ids"`
	Timestamp int64   `json:"ts"`
}
