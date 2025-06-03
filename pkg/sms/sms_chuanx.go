package sms

import (
	"fmt"
	"go-zrbc/config"
	"go-zrbc/pkg/xlog"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type chuanxSMS struct {
}

func NewChuanxSMS() SMSServer {
	return &chuanxSMS{}
}

func (s *chuanxSMS) SendMessage(phone, code string) error {
	defer func() {
		e := recover()
		if e != nil {
			xlog.Errorf("recover from panic error, failed to http auth to user service, err:%+v", e)
		}
	}()

	context := url.QueryEscape(fmt.Sprintf("您的验证码为%s", code))
	sURL := fmt.Sprintf("http://47.242.85.7:9090/sms/batch/v2?appcode=1000&appkey=eerv9n&appsecret=WAPdLY&phone=%s&msg=%s", phone, context)
	client := resty.New()
	if config.Global.Agent != "" {
		client.SetProxy(config.Global.Agent)
	}
	sResp, err := client.R().Get(sURL)
	if err != nil {
		return errors.Wrapf(err, "failed to http send message")
	}
	xlog.Infof("http send message sURL:%s, resp: %v", sURL, sResp)
	return nil
}
