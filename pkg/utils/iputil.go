package utils

import (
	"net"

	"github.com/oschwald/geoip2-golang"
)

func GetCityByIP(ip string) (string, error) {
	// 打开GeoLite2数据库
	db, err := geoip2.Open("../GeoLite2-City.mmdb")
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
