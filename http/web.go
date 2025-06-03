package http

import (
	"context"
	webSrv "go-zrbc/service/web"

	"github.com/gin-gonic/gin"

	commonresp "go-zrbc/pkg/http/response"
)

type WebHandler struct {
	srv webSrv.WebService
}

func NewWebHandler(srv webSrv.WebService) *WebHandler {
	return &WebHandler{
		srv: srv,
	}
}

func (handler *WebHandler) SetRouter(r *gin.Engine) {
	// 获取服务器系统时间
	r.GET("/v1/time_ts", handler.GetTimeTs)
}

// swagger:route GET /v1/time_ts 基础接口 GetTimeTs
// 获取当前时间戳
// responses:
//
//	200: GetTimeTsResp
//	500: CommonError
func (handler *WebHandler) GetTimeTs(c *gin.Context) {
	resp, err := handler.srv.GetTimeTs(context.TODO())
	if err != nil {
		commonresp.ErrResp(c, err)
		return
	}
	commonresp.JsonResp(c, resp)
}
