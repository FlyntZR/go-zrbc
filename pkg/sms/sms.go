package sms

const (
	SMS_GEE    = 1 // 极验(暂时没用)
	SMS_JXT    = 2 // 吉信通(中正)，中国地区
	SMS_CHUANX = 3 // chuanxsms(东南亚地区)
)

type SMSServer interface {
	SendMessage(phone, code string) error
}
