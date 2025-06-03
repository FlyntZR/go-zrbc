package main

import (
	"go-zrbc/pkg/sms"
	"net/url"
	"testing"
)

func TestSms(t *testing.T) {
	auth := sms.GenerateAuth()
	t.Logf("auth:(%s)", auth)
}

func TestChinees(t *testing.T) {
	chineseText := "你好，世界！"
	encodedText := url.QueryEscape(chineseText)
	t.Logf("Encoded:%s", encodedText)

	decodedText, err := url.QueryUnescape(encodedText)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Decoded:%s", decodedText)
}
