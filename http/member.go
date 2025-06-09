package http

import (
	"context"
	"go-zrbc/pkg/http/middleware"
	commonresp "go-zrbc/pkg/http/response"
	"go-zrbc/pkg/xlog"
	service "go-zrbc/service/user"
	"go-zrbc/view"
	"strconv"
	"time"

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
	r.POST("/v1/member_register", h.MemberRegister)
	r.POST("/v1/edit_limit", h.EditLimit)
	r.POST("/v1/logout_game", h.LogoutGame)
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
	timestamp, err := strconv.ParseInt(c.PostForm("timestamp"), 10, 64)
	if err != nil {
		xlog.Warnf("timestamp is not a number, use default value 0")
		// 方便测试自动时间戳
		timestamp = time.Now().Unix()
	}
	req.Timestamp = timestamp
	ui, err := strconv.Atoi(c.PostForm("ui"))
	if err != nil {
		xlog.Warnf("ui is not a number, use default value 0")
		ui = 0
	}
	req.UI = ui
	req.Mute = c.PostForm("mute")
	req.Video = c.PostForm("video")
	syslang, err := strconv.Atoi(c.PostForm("syslang"))
	if err != nil {
		xlog.Warnf("syslang is not a number, use default value 0")
		syslang = 0
	}
	req.Syslang = syslang

	xlog.Debugf("SigninGame req: %+v", &req)
	resp, err := handler.srv.SigninGame(c, &req)
	if err != nil {
		commonresp.ErrResp(c, err)
		return
	}
	commonresp.JsonResp(c, resp)
}

// swagger:route POST /v1/member_register 用户接口 MemberRegister
// 注册用户
// consumes:
//   - multipart/form-data
//
// responses:
//
//	200: MemberRegisterResp
//	500: CommonError
func (handler *UserHandler) MemberRegister(c *gin.Context) {
	var req view.MemberRegisterReq
	req.VendorID = c.PostForm("vendorId")
	req.Signature = c.PostForm("signature")
	req.User = c.PostForm("user")
	req.Password = c.PostForm("password")
	req.Username = c.PostForm("username")
	profile, err := strconv.Atoi(c.PostForm("profile"))
	if err != nil {
		xlog.Warnf("profile is not a number, use default value 0")
		profile = 0
	}
	req.Profile = profile
	maxwin, err := strconv.ParseInt(c.PostForm("maxwin"), 10, 64)
	if err != nil {
		xlog.Warnf("maxwin is not a number, use default value 0")
		maxwin = 0
	}
	req.Maxwin = maxwin
	maxlose, err := strconv.ParseInt(c.PostForm("maxlose"), 10, 64)
	if err != nil {
		xlog.Warnf("maxlose is not a number, use default value 0")
		maxlose = 0
	}
	req.Maxlose = maxlose
	req.Mark = c.PostForm("mark")
	rakeback, err := strconv.Atoi(c.PostForm("rakeback"))
	if err != nil {
		xlog.Warnf("rakeback is not a number, use default value 0")
		rakeback = 0
	}
	req.Rakeback = rakeback
	req.LimitType = c.PostForm("limitType")
	req.Chips = c.PostForm("chip")
	req.Kickperiod = c.PostForm("kickperiod")
	timestamp, err := strconv.ParseInt(c.PostForm("timestamp"), 10, 64)
	if err != nil {
		xlog.Warnf("timestamp is not a number, use default value 0")
		// 方便测试自动时间戳
		timestamp = time.Now().Unix()
	}
	req.Timestamp = timestamp
	syslang, err := strconv.Atoi(c.PostForm("syslang"))
	if err != nil {
		xlog.Warnf("syslang is not a number, use default value 0")
		syslang = 0
	}
	req.Syslang = syslang

	xlog.Debugf("MemberRegister req: %+v", &req)
	resp, err := handler.srv.MemberRegister(c, &req)
	if err != nil {
		commonresp.ErrResp(c, err)
		return
	}
	commonresp.JsonResp(c, resp)
}

// swagger:route POST /v1/edit_limit 用户接口 EditLimit
// 修改限额
// consumes:
//   - multipart/form-data
//
// responses:
//
//	200: EditLimitResp
//	500: CommonError
func (handler *UserHandler) EditLimit(c *gin.Context) {
	var req view.EditLimitReq
	req.VendorID = c.PostForm("vendorId")
	req.Signature = c.PostForm("signature")
	req.User = c.PostForm("user")
	req.LimitType = c.PostForm("limitType")

	maxwin, err := strconv.ParseInt(c.PostForm("maxwin"), 10, 64)
	if err != nil {
		xlog.Warnf("maxwin is not a number, use default value 0")
		maxwin = -1
	}
	req.Maxwin = maxwin

	maxlose, err := strconv.ParseInt(c.PostForm("maxlose"), 10, 64)
	if err != nil {
		xlog.Warnf("maxlose is not a number, use default value 0")
		maxlose = -1
	}
	req.Maxlose = maxlose

	reset, err := strconv.Atoi(c.PostForm("reset"))
	if err != nil {
		xlog.Warnf("reset is not a number, use default value 0")
		reset = 0
	}
	req.Reset = reset

	timestamp, err := strconv.ParseInt(c.PostForm("timestamp"), 10, 64)
	if err != nil {
		xlog.Warnf("timestamp is not a number, use default value 0")
		// Use current timestamp for testing
		timestamp = time.Now().Unix()
	}
	req.Timestamp = timestamp
	syslang, err := strconv.Atoi(c.PostForm("syslang"))
	if err != nil {
		xlog.Warnf("syslang is not a number, use default value 0")
		syslang = 0
	}
	req.Syslang = syslang

	xlog.Debugf("EditLimit req: %+v", &req)
	resp, err := handler.srv.EditLimit(c, &req)
	if err != nil {
		commonresp.ErrResp(c, err)
		return
	}
	commonresp.JsonResp(c, resp)
}

// swagger:route POST /v1/logout_game 用户接口 LogoutGame
// 登出游戏
// consumes:
//   - multipart/form-data
//
// responses:
//
//	200: LogoutGameResp
//	500: CommonError
func (handler *UserHandler) LogoutGame(c *gin.Context) {
	var req view.LogoutGameReq
	req.VendorID = c.PostForm("vendorId")
	req.Signature = c.PostForm("signature")
	req.User = c.PostForm("user")
	timestamp, err := strconv.ParseInt(c.PostForm("timestamp"), 10, 64)
	if err != nil {
		xlog.Warnf("timestamp is not a number, use default value 0")
		// 方便测试自动时间戳
		timestamp = time.Now().Unix()
	}
	req.Timestamp = timestamp
	syslang, err := strconv.Atoi(c.PostForm("syslang"))
	if err != nil {
		xlog.Warnf("syslang is not a number, use default value 0")
		syslang = 0
	}
	req.Syslang = syslang

	xlog.Debugf("LogoutGame req: %+v", &req)
	resp, err := handler.srv.LogoutGame(c, &req)
	if err != nil {
		commonresp.ErrResp(c, err)
		return
	}
	commonresp.JsonResp(c, resp)
}
