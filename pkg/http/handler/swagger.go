package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
)

type DocsHandler struct{}

func NewDocsHandler() *DocsHandler {
	return &DocsHandler{}
}

func (qa *DocsHandler) SetRouter(r *gin.Engine) {
	// 兼容本地测试时需要使用swagger_local.yaml文件
	swaggerFile := "swagger.yaml"
	if gin.Mode() == gin.DebugMode {
		swaggerFile = "swagger_local.yaml"
	} else if gin.Mode() == gin.TestMode {
		swaggerFile = "swagger_test.yaml"
	}
	opts := middleware.SwaggerUIOpts{
		SpecURL: swaggerFile,
	}
	sh := middleware.SwaggerUI(opts, nil)

	r.GET("/docs", WrapH(sh))
	//r.Static("/static", "./static")
	//r.GET("/static/*name", WrapH(NoCache(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))))
	r.GET("/swagger.yaml", WrapH(http.FileServer(http.Dir("./"))))
	r.GET("/swagger_local.yaml", WrapH(http.FileServer(http.Dir("./")))) // 兼容本地测试时需要使用swagger_local.yaml文件
	r.GET("/swagger_test.yaml", WrapH(http.FileServer(http.Dir("./"))))  // 兼容测试环境需要使用swagger_test.yaml文件
}

func WrapH(h http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
