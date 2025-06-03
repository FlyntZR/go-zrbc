package main

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestHttp_agent(t *testing.T) {
	context := url.QueryEscape(fmt.Sprintf("您的验证码为%s", "999999"))
	sURL := fmt.Sprintf("https://service.winic.org/sys_port/gateway/index.asp?id=Gho202266&pwd=Gho2022&to=%s&content=%s&time=", "18666666666", context)

	client := resty.New().SetProxy("http://18.167.194.10:443")
	sResp, err := client.R().Get(sURL)
	if err != nil {
		t.Errorf("http get failed, err:%v", err)
	}
	t.Logf("http send message sURL:%s, resp: %v", sURL, sResp)
}
