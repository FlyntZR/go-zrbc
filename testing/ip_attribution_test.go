package main

import (
	"go-zrbc/pkg/utils"
	"log"
	"net"
	"testing"

	"github.com/oschwald/geoip2-golang"
)

func getCityByIP(ip string) (string, error) {
	// 打开GeoLite2数据库
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		return "", err
	}
	defer db.Close()

	// 查询IP地址
	ipAddress := net.ParseIP(ip)
	record, err := db.City(ipAddress)
	if err != nil {
		return "", err
	}

	// 返回城市名
	return record.City.Names["zh-CN"], nil
}

func Test_IPAttribution(t *testing.T) {
	ip := "219.136.162.58" // 示例IP地址
	city, err := utils.GetCityByIP(ip)
	if err != nil {
		log.Fatal(err)
	}
	t.Logf("IP '%s' 所在城市: %s\n", ip, city)
}
