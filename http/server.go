package http

import (
	"fmt"
	"go-zrbc/config"
	. "go-zrbc/pkg/http/handler"

	pService "go-zrbc/service/public"
	s3Service "go-zrbc/service/s3"
	wService "go-zrbc/service/web"

	. "go-zrbc/pkg/http/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	webService    wService.WebService
	pubApiService pService.PublicApiService
	s3Service     s3Service.S3Service
}

func NewServer(
	webService wService.WebService,
	userService pService.PublicApiService,
	s3Service s3Service.S3Service,
) *Server {
	return &Server{
		webService:    webService,
		pubApiService: userService,
		s3Service:     s3Service,
	}
}

var (
	ServiceID = "zrbc-web-api"
)

func (s *Server) RunMetric() {
	r := gin.New()

	r.GET("/metrics", WrapH(promhttp.Handler()))
	r.Run(fmt.Sprintf(":%d", config.Global.MetricPort))
}

func (s *Server) Run() {
	if config.Global.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else if config.Global.GinMode == "test" {
		gin.SetMode(gin.TestMode)
	}

	r := gin.New()
	configa := cors.DefaultConfig()
	configa.AllowHeaders = append(configa.AllowHeaders, "Authorization")
	configa.AllowAllOrigins = true
	md := NewMiddlewareHandler(ServiceID)
	r.Use(gin.Recovery(), cors.New(configa), md.RequestLog, md.Options)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	docHandler := NewDocsHandler()
	docHandler.SetRouter(r)

	WebHandler := NewWebHandler(s.webService)
	WebHandler.SetRouter(r)

	PublicApiHandler := NewPublicApiHandler(s.pubApiService)
	PublicApiHandler.SetRouter(r)

	S3Handler := NewOssHandler(s.s3Service)
	S3Handler.SetRouter(r)

	r.Run(fmt.Sprintf(":%d", config.Global.HttpServerPort))
}
