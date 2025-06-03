package main

import (
	"fmt"
	"testing"

	gConfig "go-zrbc/config"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
)

func TestAgollo_configTest_1(t *testing.T) {
	c := &config.AppConfig{
		AppID:          "go-zrbc",
		Cluster:        "dev",
		IP:             "http://103.84.45.139:10021",
		NamespaceName:  "application",
		IsBackupConfig: false,
		Secret:         "",
	}

	client, _ := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	t.Logf("初始化Apollo配置成功")

	//Use your apollo key to test
	cache := client.GetConfigCache(c.NamespaceName)
	value, _ := cache.Get("go.http_server_port")
	fmt.Println(value)
	t.Logf("value:%+v", value)
}

var Conf = gConfig.Config{}

func TestAgollo_configTest_2(t *testing.T) {
	gConfig.GetConfigFromApollo(&Conf)
	t.Logf("Conf:%+v", Conf)
}
