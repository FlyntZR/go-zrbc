package response

import "fmt"

// ErrorCode 定义错误码类型
type ErrorCode int

// 定义系统级别错误码
const (
	// 成功
	CodeSuccess ErrorCode = 200
	// 系统错误
	CodeSystemError ErrorCode = 50001
	// 参数错误
	CodeParamError ErrorCode = 40000
	// 未授权
	CodeUnauthorized ErrorCode = 40100
	// 禁止访问
	CodeForbidden ErrorCode = 40300
	// 资源不存在
	CodeNotFound ErrorCode = 40400
	// 请求超时
	CodeTimeout ErrorCode = 40800
	// 用户不存在
	CodeUserNotExist ErrorCode = 10501
	// 用户名不得为空
	CodeUserEmpty ErrorCode = 10502
	// 密码不得为空
	CodeUserPWDEmpty ErrorCode = 10503
	// 用户密码错误
	CodeUserPWD ErrorCode = 10504
	// 用户已停用
	CodeUserDeactivated ErrorCode = 10505
)

// CustomError 自定义错误类型
type CustomError struct {
	Code    ErrorCode
	Message string
}

// Error 实现error接口
func (e *CustomError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// NewError 创建新的错误
func NewError(code ErrorCode, message string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}

// 预定义的错误
var (
	ErrSystem          = NewError(CodeSystemError, "系统错误")
	ErrParam           = NewError(CodeParamError, "参数错误")
	ErrUnauth          = NewError(CodeUnauthorized, "未授权访问")
	ErrForbidden       = NewError(CodeForbidden, "禁止访问")
	ErrNotFound        = NewError(CodeNotFound, "资源不存在")
	ErrTimeout         = NewError(CodeTimeout, "请求超时")
	ErrUserNotExist    = NewError(CodeUserNotExist, "用户不存在")
	ErrUserPWD         = NewError(CodeUserPWD, "此用户的密码错误")
	ErrUserEmpty       = NewError(CodeUserEmpty, "用户名不得为空")
	ErrUserPWDEmpty    = NewError(CodeUserPWDEmpty, "密码不得为空")
	ErrUserDeactivated = NewError(CodeUserDeactivated, "用户已停用")
)
