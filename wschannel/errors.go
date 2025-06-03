package wschannel

import (
	"errors"
)

var (
	ErrWsChannelFull       = errors.New("ws channel is full")
	ErrUserOffline         = errors.New("user offline already")
	ErrOrderStatusError    = errors.New("order status error")
	ErrSymbolIdError       = errors.New("symbol id error")
	ErrDuplicateLogin      = errors.New("duplicate login")
	ErrDateFormatError     = errors.New("data format error")
	ErrInternalServerError = errors.New("internal server error")
	ErrSymbolNameNullError = errors.New("symbol name is empty")
	ErrDepthLevelError     = errors.New("depth level error")
	ErrUserAlreadyInScene  = errors.New("user already in the room")
	ErrRoomFull            = errors.New("room is full")
	ErrChatSessExpired     = errors.New("chat sess expired")
)

var (
	RespDataFormatError     = NewConnResp(ErrCodeTypeAssertInvalid)
	RespUserOffline         = NewConnResp(ErrCodeUserOffline)
	RespInternalServerError = NewConnResp(ErrCodeInternalServerError)
	RespDuplicateLogin      = NewConnResp(ErrCodeDuplicateLogin)
	RespSymbolNameNullError = NewConnResp(ErrCodeSymbolNameNull)
)

type WSRespCode int

func (wcode WSRespCode) String() string {
	switch wcode {
	case ErrCodeUserOffline:
		return "user is offline"

	case ErrCodeWsChannelFull:
		return "ws channel is full"

	case ErrCodeTypeAssertInvalid:
		return "data format err"

	case ErrCodeDuplicateLogin:
		return "duplicate login"

	case ErrCodeSymbolIdError:
		return "symbol id error"

	case ErrCodeOrderStatusError:
		return "order status error"

	case ErrCodeInternalServerError:
		return "internal server error"
	}

	return ""
}

var (
	ErrCodeMap = map[error]WSRespCode{
		ErrUserOffline:         ErrCodeUserOffline,
		ErrWsChannelFull:       ErrCodeWsChannelFull,
		ErrOrderStatusError:    ErrCodeOrderStatusError,
		ErrDuplicateLogin:      ErrCodeDuplicateLogin,
		ErrSymbolIdError:       ErrCodeSymbolIdError,
		ErrInternalServerError: ErrCodeInternalServerError,
		ErrDateFormatError:     ErrCodeTypeAssertInvalid,
		ErrSymbolNameNullError: ErrCodeSymbolNameNull,
	}
)

var (
	ErrCodeUserOffline       WSRespCode = 4001
	ErrCodeWsChannelFull     WSRespCode = 4002
	ErrCodeDuplicateLogin    WSRespCode = 4003
	ErrCodeSymbolIdError     WSRespCode = 4004
	ErrCodeOrderStatusError  WSRespCode = 4005
	ErrCodeTypeAssertInvalid WSRespCode = 4006
	ErrCodeSymbolNameNull    WSRespCode = 4007

	ErrCodeInternalServerError WSRespCode = 5001
)

func NewConnResp(code WSRespCode) *ConnMessageResp {
	return &ConnMessageResp{
		Code:   int(code),
		ErrMsg: code.String(),
	}
}

type ConnMessageResp struct {
	Code   int    `json:"code"`
	ErrMsg string `json:"err_msg"`
}
