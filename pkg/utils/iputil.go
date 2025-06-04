package utils

import (
	"go-zrbc/pkg/xlog"
	"net"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
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

// ExtractMainDomain 提取主域名
func extractMainDomain(host string) (string, error) {
	// 解析主机名
	u, err := url.Parse("http://" + host)
	if err != nil {
		return "", err
	}

	// 分割主机名
	parts := strings.Split(u.Hostname(), ".")
	if len(parts) < 2 {
		return "", nil // 无效域名
	}

	// 提取主域名（最后两部分，例如 example.com）
	return strings.Join(parts[len(parts)-2:], "."), nil
}

// login handles user authentication and session creation
func GetServerName(ctx *gin.Context) string {
	host := ctx.Request.Host // 获取 Host 头，例如 "sub.example.com:8080"
	// 移除端口号（如果存在）
	host = strings.Split(host, ":")[0]
	// 解析主域名
	mainDomain, err := extractMainDomain(host)
	if err != nil {
		xlog.Errorf("error to extract main domain, err:%+v", err)
		return ""
	}
	return mainDomain
}
