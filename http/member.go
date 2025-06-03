package http

import (
	"context"
	"go-zrbc/pkg/http/middleware"
	commonresp "go-zrbc/pkg/http/response"
	service "go-zrbc/service/user"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	srv service.UserService
}

func NewUserHandler(srv service.UserService) *UserHandler {
	return &UserHandler{
		srv: srv,
	}
}

func (h *UserHandler) SetRouter(r *gin.Engine) {
	r.GET("/v1/user_info", h.GetUserInfo)
}

// swagger:route GET /v1/user_info 用户接口 GetUserInfo
// 获取会员列表
// responses:
//
//	200: GetUserInfoResp
//	500: CommonError
func (handler *UserHandler) GetUserInfo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	resp, err := handler.srv.GetUserInfo(context.TODO(), userID)
	if err != nil {
		commonresp.ErrResp(c, err)
		return
	}
	commonresp.JsonResp(c, resp)
}
