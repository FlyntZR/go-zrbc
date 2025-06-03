package main

import (
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestHttp_flynt(t *testing.T) {
	tURL := "https://wys.dev.zhanggao223.com/v1/ys_like_records"
	client := resty.New()
	tResp, err := client.R().Get(tURL)
	if err != nil {
		t.Errorf("http get failed, err:%v", err)
	}
	t.Logf("http get result, resp: %v", tResp)
}
