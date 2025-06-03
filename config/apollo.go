package config

import (
	"go-zrbc/pkg/xlog"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env"
	"github.com/apolloconfig/agollo/v4/env/config"
)

func GetConfigFromApollo(gConfig *Config) {
	// c := &config.AppConfig{
	// 	AppID:          "testApplication_yang",
	// 	Cluster:        "dev",
	// 	IP:             "http://106.54.227.205:8080",
	// 	NamespaceName:  "application",
	// 	IsBackupConfig: true,
	// 	Secret:         "",
	// }

	c, err := env.InitConfig(nil)
	if err != nil {
		panic(err)
	}

	client, _ := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	xlog.Info("start apollo config success")

	gConfig.GinMode = client.GetStringValue("go.gin.mode", "release")
	gConfig.HttpServerPort = client.GetIntValue("go.http_server_port", 0)
	gConfig.MetricPort = client.GetIntValue("go.metric_port", 0)
	gConfig.Mysql.UserName = client.GetStringValue("go.mysql.user_name", "")
	gConfig.Mysql.Passwd = client.GetStringValue("go.mysql.passwd", "")
	gConfig.Mysql.Database = client.GetStringValue("go.mysql.database", "")
	gConfig.Mysql.Addr = client.GetStringValue("go.mysql.addr", "")
	gConfig.Mysql.Port = client.GetIntValue("go.mysql.port", 0)
	gConfig.Redis.Addr = client.GetStringValue("go.redis.addr", "")
	gConfig.Redis.Password = client.GetStringValue("go.redis.password", "")
	gConfig.Redis.DB = client.GetIntValue("go.redis.db", 0)
	gConfig.LogLevel = client.GetStringValue("go.log_level", "")
	gConfig.LogFile = client.GetStringValue("go.log_file", "")
	gConfig.Agent = client.GetStringValue("go.agent", "")
	xlog.Info("load apollo config end")
}
