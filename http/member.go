package http

import (
	"context"
	"go-zrbc/pkg/http/middleware"
	commonresp "go-zrbc/pkg/http/response"
	"go-zrbc/pkg/xlog"
	service "go-zrbc/service/user"
	"go-zrbc/view"
	"strconv"

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
	r.POST("/v1/signin_game", h.SigninGame)
}

// swagger:route GET /v1/user_info 用户接口 GetUserInfo
// 获取用户信息
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

// swagger:route POST /v1/signin_game 用户接口 SigninGame
// 开游戏
// consumes:
//   - multipart/form-data
//
// responses:
//
//	200: SigninGameResp
//	500: CommonError
func (handler *UserHandler) SigninGame(c *gin.Context) {
	var req view.SigninGameReq
	req.VendorID = c.PostForm("vendorId")
	req.Signature = c.PostForm("signature")
	req.User = c.PostForm("user")
	req.Password = c.PostForm("password")
	isTest := c.PostForm("isTest")
	req.IsTest = isTest == "true"
	req.Mode = c.PostForm("mode")
	req.TableID = c.PostForm("tableid")
	req.Site = c.PostForm("site")
	req.GameType = c.PostForm("gameType")
	req.Width = c.PostForm("width")
	req.ReturnURL = c.PostForm("returnurl")
	req.Size = c.PostForm("size")
	ui, err := strconv.Atoi(c.PostForm("ui"))
	if err != nil {
		xlog.Warnf("ui is not a number, use default value 0")
		ui = 0
	}
	req.UI = ui
	req.Mute = c.PostForm("mute")
	req.Video = c.PostForm("video")

	xlog.Debugf("SigninGame req: %+v", &req)
	resp, err := handler.srv.SigninGame(c, &req)
	if err != nil {
		commonresp.ErrResp(c, err)
		return
	}
	commonresp.JsonResp(c, resp)
}
