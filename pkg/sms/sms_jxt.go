package sms

import (
	"bytes"
	"fmt"
	"go-zrbc/config"
	"go-zrbc/pkg/xlog"
	"io/ioutil"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type jxtSMS struct {
}

func NewJxtSMS() SMSServer {
	return &jxtSMS{}
}

func (s *jxtSMS) SendMessage(phone, code string) error {
	defer func() {
		e := recover()
		if e != nil {
			xlog.Errorf("recover from panic error, failed to http auth to user service, err:%+v", e)
		}
	}()

	// context := url.QueryEscape(fmt.Sprintf("验证码%s，请勿告知他人！【广州玉美】", code))
	content := fmt.Sprintf("验证码%s，请勿告知他人！【广州玉美】", code)
	data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(content)), simplifiedchinese.GBK.NewEncoder()))
	sURL := fmt.Sprintf("https://service.winic.org/sys_port/gateway/index.asp?id=Gho202266&pwd=Gho2022&to=%s&content=%s&time=", phone, data)

	client := resty.New()
	if config.Global.Agent != "" {
		client.SetProxy(config.Global.Agent)
	}
	sResp, err := client.R().Get(sURL)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("failed to http send message, agent:%s", config.Global.Agent))
	}
	xlog.Infof("http send message sURL:%s, resp: %v", sURL, sResp)
	return nil
}
