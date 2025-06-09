package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"go-zrbc/pkg/xlog"
)

type Config struct {
	GinMode        string `json:"gin_mode"` // debug: 使用swagger_local.yaml；release: 使用swagger.yaml
	ServiceID      string `json:"service_id"`
	HttpServerPort int    `json:"http_server_port"`
	MetricPort     int    `json:"metric_port"`
	Mysql          Mysql  `json:"mysql"`
	Redis          Redis  `json:"redis"`
	AwsKey         string `json:"aws_key"`
	AwsSecret      string `json:"aws_secret"`
	Wcode          string `json:"wcode"`
	SidLen         int    `json:"sid_len"`
	AlertUrl       string `json:"alert_url"`
	WMAlertUrl     string `json:"wm_alert_url"`
	LogLevel       string `json:"log_level"`
	LogFile        string `json:"log_file"`
	Agent          string `json:"agent"`        // 网络代理，如果非空，则需要使用代理
	SMSSupplier    int    `json:"sms_supplier"` // 短信验证码供应商
	SMSModeID      string `json:"sms_mode_id"`  // 短信验证码模板id
}

type Redis struct {
	Addr     string `json:"addr"`
	Password string `json:"passwd"`
	DB       int    `json:"db"`
}

type Mysql struct {
	UserName string `json:"user_name"`
	Passwd   string `json:"passwd"`
	Database string `json:"database"`
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
}

func (mysql *Mysql) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysql.UserName,
		mysql.Passwd,
		mysql.Addr,
		mysql.Port,
		mysql.Database,
	)
}

var Global = &Config{}

func Init(path string) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		xlog.Error(err)
		os.Exit(1)
	}
	err = json.Unmarshal(b, Global)
	if err != nil {
		xlog.Error(err)
		os.Exit(1)
	}
}
