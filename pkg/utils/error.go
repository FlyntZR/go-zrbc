package utils

import "fmt"

// ErrorCode 定义错误码类型
type ErrorCode int

// 定义系统级别错误码
const (
	// 成功
	CodeSuccess ErrorCode = 200

	// 系统错误
	CodeSystemError ErrorCode = 50001
	// 内部服务器错误
	CodeInternalServerError ErrorCode = 500

	// ws 错误
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

	// 操作失敗
	CodeError ErrorCode = 1
	// 网路连线錯誤，或未有此用户取得帐号资料失败
	CodeNetworkError ErrorCode = 11
	// 指令操作成功但查无此笔资料
	CodeCommandSuccessButNoData ErrorCode = 107
	// 指令操作成功，但未写入资料库
	CodeCommandSuccessButNoDataInsert ErrorCode = 10711
	// 查无此命令
	CodeNoCommand ErrorCode = 100
	// 身分认证错误!
	CodeAuthenticationError ErrorCode = 101
	// 时间格式错误
	CodeTimeFormatError ErrorCode = 102
	// 此功能仅能查询一天内的报表，您已超过上限
	CodeFunctionOnlyQueryOneDayReport ErrorCode = 10201
	// 时间戳异常,超过30秒
	CodeTimestampError ErrorCode = 10202
	// 代理商ID与识别码格式错误
	CodeAgentIDAndSignatureFormatError ErrorCode = 103
	// 代理商ID为空,请检查(vendorid)
	CodeAgentIDEmpty ErrorCode = 10301
	// 没有这个代理商ID
	CodeAgentIDNotExist ErrorCode = 10302
	// 有此代理商ID,但代理商代码(signature)错误
	CodeAgentIDExistButSignatureError ErrorCode = 10303
	// 代理商代码(signature)为空
	CodeAgentSignatureEmpty ErrorCode = 10304
	// 代理商已被停用登入或下注
	CodeAgentDisabled ErrorCode = 10305

	// 账号已存在
	CodeAccountExists ErrorCode = 104
	// 限制类型无效
	CodeInvalidLimitType ErrorCode = 106
	// 限額未開放，請檢查
	CodeInvalidLimitTypeNotOpen ErrorCode = 10601
	// 期数不得为空
	CodeInvalidPeriodEmpty ErrorCode = 10602
	// 期数资料不存在
	CodeInvalidPeriodNotExist ErrorCode = 10603

	// 新增会员资料错误 帐号(user)不能为空值
	CodeInvalidAccountEmpty ErrorCode = 10401
	// 新增会员资料错误 密码(password)不能为空值
	CodeInvalidPasswordEmpty ErrorCode = 10402
	// 新增会员资料错误 姓名(username)不能为空值
	CodeInvalidUsernameEmpty ErrorCode = 10403
	// 帐号长度过长,因介于5~17字元
	CodeInvalidAccountLength ErrorCode = 10404
	// 帐号长度过短,因介于5~17字元
	CodeInvalidAccountLengthShort ErrorCode = 10405
	// 密码长度过短
	CodeInvalidPasswordLengthShort ErrorCode = 10406
	// 密码长度过长
	CodeInvalidPasswordLengthLong ErrorCode = 10407
	// 备注长度过长
	CodeInvalidMarkLengthLong ErrorCode = 10408
	// 姓名长度过长
	CodeInvalidUsernameLengthLong ErrorCode = 10409
	// 会员上笔交易未成功，请联系客服人员解锁
	CodeInvalidTransactionError ErrorCode = 10410
	// 请于30秒後再试
	CodeInvalidTransactionTimeout ErrorCode = 10411
	// 超过可注册会员数限制
	CodeInvalidMemberLimit ErrorCode = 10412
	// 代理输赢低于设定上限暂停使用
	CodeInvalidAgentProfitLimit ErrorCode = 10413
	// 密码不得为中文
	CodeInvalidPasswordChinese ErrorCode = 10414
	// 帐号不得为中文
	CodeInvalidAccountChinese ErrorCode = 10415
	// 未下注踢出局数大於上層
	CodeInvalidKickPeriodGreaterThanUpper ErrorCode = 10416
	// 未下注踢出局数不能为负数
	CodeInvalidKickPeriodNegative ErrorCode = 10417
	// 10秒內禁止重複操作
	CodeInvalidTransactionTimeoutRepeat ErrorCode = 10418
	// 筹码格式错误(请用逗号隔开)
	CodeInvalidChipsFormat ErrorCode = 10419
	// 筹码个数错误(介于5-10个)
	CodeInvalidChipsCount ErrorCode = 10420
	// 筹码种类错误
	CodeInvalidChipsType ErrorCode = 10421
	// 帐号只接受英文、数字、下划线与@
	CodeInvalidAccountFormat ErrorCode = 10422

	// 代入参数error
	// 帐号密码格式错误
	CodeParamInvalidAccountPasswordFormat ErrorCode = 105
	// 查无此帐号,请检查
	CodeParamInvalidAccountNotExist ErrorCode = 10501
	// 帐号名不得为空
	CodeParamInvalidAccountNameEmpty ErrorCode = 10502
	// 密码不得为空
	CodeParamInvalidPasswordEmpty ErrorCode = 10503
	// 此帐号的密码错误
	CodeParamInvalidAccountPasswordError ErrorCode = 10504
	// 此帐号已被停用
	CodeParamInvalidAccountDeactivated ErrorCode = 10505
	// 此账号非此代理下线
	CodeParamInvalidAccountNotBelongToAgent ErrorCode = 10506
	// 此账号非此代理下线,不能使用此功能
	CodeParamInvalidAccountNotBelongToAgentFunction ErrorCode = 10507

	CodeParamInvalidUsernameEmpty ErrorCode = 10509
	// 此代理商尚未有下线会员
	CodeParamInvalidAgentNoMember ErrorCode = 10510
	// 修改密码与原密码相同
	CodeParamInvalidPasswordSame ErrorCode = 10511
	// 帳號或密碼有非法字元
	CodeParamInvalidAccountPasswordIllegal ErrorCode = 10512
	// 无效的代理推广码
	CodeParamInvalidAgentPromotionCode ErrorCode = 10513
	// 找不到代理帐号
	CodeParamInvalidAgentNotExist ErrorCode = 10514
	// 代理推广码或币别其中一个为必填
	CodeParamInvalidAgentPromotionCodeOrCurrency ErrorCode = 10515
	// 姓名格式错误
	CodeParamInvalidUsernameFormat ErrorCode = 10516
	// sid不得為空
	CodeParamInvalidSidEmpty ErrorCode = 10517
	// 旧帐号错误
	CodeParamInvalidOldAccountError ErrorCode = 10518
	// 密码只能使用英数混合
	CodeParamInvalidPasswordFormat ErrorCode = 10519
	// 上层代理停用或停押
	CodeParamInvalidAgentDeactivated ErrorCode = 10520

	// wallet 单一钱包
	// 运营商代码不得为空
	CodeWalletOperatorCodeEmpty ErrorCode = 10701
	// 运营商代码不正确
	CodeWalletOperatorCodeIncorrect ErrorCode = 10702
	// 流水号不得为空
	CodeWalletSerialNumberEmpty ErrorCode = 10703
	// 此流水号查不到资料
	CodeWalletSerialNumberNotExist ErrorCode = 10704
	// 回傳網址未填入
	CodeWalletReturnUrlEmpty ErrorCode = 10705

	// 加扣点不得为零或英文字
	CodeWalletAddOrSubPointZeroOrEnglish ErrorCode = 10801
	// 加扣点为空,或未设置(money)参数
	CodeWalletAddOrSubPointEmptyOrMoneyParamNotSet ErrorCode = 10802
	// 加扣点不得为汉字
	CodeWalletAddOrSubPointChinese ErrorCode = 10803
	// 不得5秒內重复转帐
	CodeWalletTransferRepeatIn5Error ErrorCode = 10804
	// 转帐失败，该帐号余额不足
	CodeWalletTransferBalanceNotEnough ErrorCode = 10805
	// 此帐号代理商已超过代理可开点数
	CodeWalletAgentOverLimit ErrorCode = 10806
	// 转帐失败,该笔单号已存在
	CodeWalletTransferExist ErrorCode = 10807
	// 转帐失败,一分钟内转帐次数超过10次,帐号已锁定
	CodeWalletTransferLockError ErrorCode = 10808
	// 不得2秒內重复转帐
	CodeWalletTransferRepeatIn2Error ErrorCode = 10809
	// 连线异常，交易未成功
	CodeWalletTransferLineError ErrorCode = 10810

	// 注单编号不可为空
	CodeWalletBetNumberEmpty ErrorCode = 10910
	// 無此注单资料
	CodeWalletBetNumberNotExist ErrorCode = 10911
	// 参数格式错误
	CodeWalletParamFormatError ErrorCode = 10912

	// 识别码验证失败
	CodeWalletCodeVerifyError ErrorCode = 201
	// 识别码不得为空
	CodeWalletCodeEmpty ErrorCode = 202

	// system error
	// 查无此函数
	CodeSystemFunctionNotExist ErrorCode = 900
	// 交易未成功
	CodeSystemTransactionError ErrorCode = 901
	// 维护中
	CodeSystemMaintenance ErrorCode = 911
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
	ErrSuccess                                     = NewError(CodeSuccess, "操作成功")
	ErrError                                       = NewError(CodeError, "操作失败")
	ErrNetworkError                                = NewError(CodeNetworkError, "网路连线錯誤，或未有此用户取得帐号资料失败")
	ErrCommandSuccessButNoData                     = NewError(CodeCommandSuccessButNoData, "指令操作成功但查无此笔资料")
	ErrCommandSuccessButNoDataInsert               = NewError(CodeCommandSuccessButNoDataInsert, "指令操作成功，但未写入资料库")
	ErrNoCommand                                   = NewError(CodeNoCommand, "查无此命令")
	ErrAuthenticationError                         = NewError(CodeAuthenticationError, "身分认证错误!")
	ErrTimeFormatError                             = NewError(CodeTimeFormatError, "时间格式错误")
	ErrTimestampError                              = NewError(CodeTimestampError, "时间戳异常,超过30秒")
	ErrAgentIDAndSignatureFormatError              = NewError(CodeAgentIDAndSignatureFormatError, "代理商ID与识别码格式错误")
	ErrAgentIDEmpty                                = NewError(CodeAgentIDEmpty, "代理商ID为空,请检查(vendorid)")
	ErrAgentIDNotExist                             = NewError(CodeAgentIDNotExist, "没有这个代理商ID")
	ErrAgentIDExistButSignatureError               = NewError(CodeAgentIDExistButSignatureError, "有此代理商ID,但代理商代码(signature)错误")
	ErrAgentSignatureEmpty                         = NewError(CodeAgentSignatureEmpty, "代理商代码(signature)为空")
	ErrAgentDisabled                               = NewError(CodeAgentDisabled, "代理商已被停用登入或下注")
	ErrAccountExists                               = NewError(CodeAccountExists, "账号已存在")
	ErrInvalidLimitType                            = NewError(CodeInvalidLimitType, "新增资料错误")
	ErrInvalidAccountEmpty                         = NewError(CodeInvalidAccountEmpty, "新增会员资料错误 帐号(user)不能为空值")
	ErrInvalidPasswordEmpty                        = NewError(CodeInvalidPasswordEmpty, "新增会员资料错误 密码(password)不能为空值")
	ErrInvalidUsernameEmpty                        = NewError(CodeInvalidUsernameEmpty, "新增会员资料错误 姓名(username)不能为空值")
	ErrInvalidAccountLength                        = NewError(CodeInvalidAccountLength, "帐号长度过长,因介于5~17字元")
	ErrInvalidAccountLengthShort                   = NewError(CodeInvalidAccountLengthShort, "帐号长度过短,因介于5~17字元")
	ErrInvalidPasswordLengthShort                  = NewError(CodeInvalidPasswordLengthShort, "密码长度过短")
	ErrInvalidPasswordLengthLong                   = NewError(CodeInvalidPasswordLengthLong, "密码长度过长")
	ErrInvalidMarkLengthLong                       = NewError(CodeInvalidMarkLengthLong, "备注长度过长")
	ErrInvalidTransactionError                     = NewError(CodeInvalidTransactionError, "会员上笔交易未成功，请联系客服人员解锁")
	ErrInvalidTransactionTimeout                   = NewError(CodeInvalidTransactionTimeout, "请于30秒後再试")
	ErrInvalidMemberLimit                          = NewError(CodeInvalidMemberLimit, "超过可注册会员数限制")
	ErrInvalidAgentProfitLimit                     = NewError(CodeInvalidAgentProfitLimit, "代理输赢低于设定上限暂停使用")
	ErrInvalidPasswordChinese                      = NewError(CodeInvalidPasswordChinese, "密码不得为中文")
	ErrInvalidAccountChinese                       = NewError(CodeInvalidAccountChinese, "帐号不得为中文")
	ErrInvalidKickPeriodGreaterThanUpper           = NewError(CodeInvalidKickPeriodGreaterThanUpper, "未下注踢出局数大於上層")
	ErrInvalidKickPeriodNegative                   = NewError(CodeInvalidKickPeriodNegative, "未下注踢出局数不能为负数")
	ErrInvalidTransactionTimeoutRepeat             = NewError(CodeInvalidTransactionTimeoutRepeat, "10秒內禁止重複操作")
	ErrInvalidChipsFormat                          = NewError(CodeInvalidChipsFormat, "筹码格式错误(请用逗号隔开)")
	ErrInvalidChipsCount                           = NewError(CodeInvalidChipsCount, "筹码个数错误(介于5-10个)")
	ErrInvalidChipsType                            = NewError(CodeInvalidChipsType, "筹码种类错误")
	ErrInvalidAccountFormat                        = NewError(CodeInvalidAccountFormat, "帐号只接受英文、数字、下划线与@")
	ErrParamInvalidAccountPasswordFormat           = NewError(CodeParamInvalidAccountPasswordFormat, "帐号密码格式错误")
	ErrParamInvalidAccountNotExist                 = NewError(CodeParamInvalidAccountNotExist, "查无此帐号,请检查")
	ErrParamInvalidAccountNameEmpty                = NewError(CodeParamInvalidAccountNameEmpty, "帐号名不得为空")
	ErrParamInvalidPasswordEmpty                   = NewError(CodeParamInvalidPasswordEmpty, "密码不得为空")
	ErrParamInvalidAccountPasswordError            = NewError(CodeParamInvalidAccountPasswordError, "此帐号的密码错误")
	ErrParamInvalidAccountDeactivated              = NewError(CodeParamInvalidAccountDeactivated, "此帐号已被停用")
	ErrParamInvalidAccountNotBelongToAgent         = NewError(CodeParamInvalidAccountNotBelongToAgent, "此账号非此代理下线")
	ErrParamInvalidAccountNotBelongToAgentFunction = NewError(CodeParamInvalidAccountNotBelongToAgentFunction, "此账号非此代理下线,不能使用此功能")
	ErrParamInvalidAgentNoMember                   = NewError(CodeParamInvalidAgentNoMember, "此代理商尚未有下线会员")
	ErrParamInvalidPasswordSame                    = NewError(CodeParamInvalidPasswordSame, "修改密码与原密码相同")
	ErrParamInvalidAccountPasswordIllegal          = NewError(CodeParamInvalidAccountPasswordIllegal, "帳號或密碼有非法字元")
	ErrParamInvalidAgentPromotionCode              = NewError(CodeParamInvalidAgentPromotionCode, "无效的代理推广码")
	ErrParamInvalidAgentNotExist                   = NewError(CodeParamInvalidAgentNotExist, "找不到代理帐号")
	ErrParamInvalidAgentPromotionCodeOrCurrency    = NewError(CodeParamInvalidAgentPromotionCodeOrCurrency, "代理推广码或币别其中一个为必填")
	ErrParamInvalidUsernameFormat                  = NewError(CodeParamInvalidUsernameFormat, "姓名格式错误")
	ErrParamInvalidSidEmpty                        = NewError(CodeParamInvalidSidEmpty, "sid不得為空")
	ErrParamInvalidOldAccountError                 = NewError(CodeParamInvalidOldAccountError, "旧帐号错误")
	ErrParamInvalidPasswordFormat                  = NewError(CodeParamInvalidPasswordFormat, "密码只能使用英数混合")
	ErrParamInvalidAgentDeactivated                = NewError(CodeParamInvalidAgentDeactivated, "上层代理停用或停押")
	ErrWalletOperatorCodeEmpty                     = NewError(CodeWalletOperatorCodeEmpty, "运营商代码不得为空")
	ErrWalletOperatorCodeIncorrect                 = NewError(CodeWalletOperatorCodeIncorrect, "运营商代码不正确")
	ErrWalletSerialNumberEmpty                     = NewError(CodeWalletSerialNumberEmpty, "流水号不得为空")
	ErrWalletSerialNumberNotExist                  = NewError(CodeWalletSerialNumberNotExist, "此流水号查不到资料")
	ErrWalletReturnUrlEmpty                        = NewError(CodeWalletReturnUrlEmpty, "回傳網址未填入")
	ErrWalletAddOrSubPointZeroOrEnglish            = NewError(CodeWalletAddOrSubPointZeroOrEnglish, "加扣点不得为零或英文字")
	ErrWalletAddOrSubPointEmptyOrMoneyParamNotSet  = NewError(CodeWalletAddOrSubPointEmptyOrMoneyParamNotSet, "加扣点为空,或未设置(money)参数")
	ErrWalletAddOrSubPointChinese                  = NewError(CodeWalletAddOrSubPointChinese, "加扣点不得为汉字")
	ErrWalletTransferRepeatIn5Error                = NewError(CodeWalletTransferRepeatIn5Error, "不得5秒內重复转帐")
	ErrWalletTransferBalanceNotEnough              = NewError(CodeWalletTransferBalanceNotEnough, "转帐失败，该帐号余额不足")
	ErrWalletAgentOverLimit                        = NewError(CodeWalletAgentOverLimit, "此帐号代理商已超过代理可开点数")
	ErrWalletTransferExist                         = NewError(CodeWalletTransferExist, "转帐失败,该笔单号已存在")
	ErrWalletTransferLockError                     = NewError(CodeWalletTransferLockError, "转帐失败,一分钟内转帐次数超过10次,帐号已锁定")
	ErrWalletTransferRepeatIn2Error                = NewError(CodeWalletTransferRepeatIn2Error, "不得2秒內重复转帐")
	ErrWalletTransferLineError                     = NewError(CodeWalletTransferLineError, "连线异常，交易未成功")
	ErrWalletBetNumberEmpty                        = NewError(CodeWalletBetNumberEmpty, "注单编号不可为空")
	ErrWalletBetNumberNotExist                     = NewError(CodeWalletBetNumberNotExist, "無此注单资料")
	ErrWalletParamFormatError                      = NewError(CodeWalletParamFormatError, "参数格式错误")
	ErrWalletCodeVerifyError                       = NewError(CodeWalletCodeVerifyError, "识别码验证失败")
	ErrWalletCodeEmpty                             = NewError(CodeWalletCodeEmpty, "识别码不得为空")
	ErrSystemFunctionNotExist                      = NewError(CodeSystemFunctionNotExist, "查无此函数")
	ErrSystemTransactionError                      = NewError(CodeSystemTransactionError, "交易未成功")
	ErrSystemMaintenance                           = NewError(CodeSystemMaintenance, "维护中")
	ErrParamError                                  = NewError(CodeParamError, "参数错误")
	ErrUnauthorized                                = NewError(CodeUnauthorized, "未授权")
	ErrForbidden                                   = NewError(CodeForbidden, "禁止访问")
	ErrNotFound                                    = NewError(CodeNotFound, "资源不存在")
	ErrTimeout                                     = NewError(CodeTimeout, "请求超时")
	ErrSystemError                                 = NewError(CodeSystemError, "系统错误")
	ErrInternalServerError                         = NewError(CodeInternalServerError, "内部服务器错误")
	ErrInvalidUsernameLengthLong                   = NewError(CodeInvalidUsernameLengthLong, "姓名长度过长")
	ErrInvalidLimitTypeNotOpen                     = NewError(CodeInvalidLimitTypeNotOpen, "限額未開放，請檢查")
	ErrInvalidPeriodEmpty                          = NewError(CodeInvalidPeriodEmpty, "期数不得为空")
	ErrInvalidPeriodNotExist                       = NewError(CodeInvalidPeriodNotExist, "期数资料不存在")
)
