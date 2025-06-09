package http

import (
	"context"
	"fmt"
	"go-zrbc/pkg/http/middleware"
	commonresp "go-zrbc/pkg/http/response"
	"go-zrbc/pkg/xlog"
	service "go-zrbc/service/public"
	"go-zrbc/view"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PublicApiHandler struct {
	srv service.PublicApiService
}

func NewPublicApiHandler(srv service.PublicApiService) *PublicApiHandler {
	return &PublicApiHandler{
		srv: srv,
	}
}

func (h *PublicApiHandler) SetRouter(r *gin.Engine) {
	r.POST("/api/public/Gateway.php", h.handlePublicApi)

	r.GET("/v1/user_info", h.GetUserInfo)
	r.POST("/v1/signin_game", h.SigninGame)
	r.POST("/v1/member_register", h.MemberRegister)
	r.POST("/v1/edit_limit", h.EditLimit)
	r.POST("/v1/logout_game", h.LogoutGame)
	r.POST("/v1/change_password", h.ChangePassword)
	r.POST("/v1/get_agent_balance", h.GetAgentBalance)
}

func (h *PublicApiHandler) handlePublicApi(c *gin.Context) {
	// Get command from request
	cmd := c.PostForm("cmd")
	if cmd == "" {
		cmd = c.Query("cmd")
	}

	// Handle GetMemberReport command
	if cmd == "GetMemberReport" {
		commonresp.JsonResp(c, map[string]interface{}{
			"errorCode": 100,
		})
		return
	}

	// Check if command is in allowed list
	passCommands := map[string]bool{
		"GetAgentBalance":        true,
		"MemberRegister":         true,
		"EditLimit":              true,
		"Hello":                  true,
		"MemberLogin":            true,
		"GetDateTimeReport":      true,
		"GetTipReport":           true,
		"EnableorDisablemem":     true,
		"GetMemberTradeReport":   true,
		"LogoutGame":             true,
		"GetUnsettleReport":      true,
		"GetDateTimeReportOld":   true,
		"GetDateTimeCountReport": true,
	}

	// Handle command
	switch cmd {
	case "MemberRegister":
		h.MemberRegister(c)
	case "LogoutGame":
		h.LogoutGame(c)
	case "EditLimit":
		h.EditLimit(c)
	case "ChangePassword":
		h.ChangePassword(c)
	case "GetAgentBalance":
		h.GetAgentBalance(c)
	// Add other command handlers as needed
	default:
		if !passCommands[cmd] {
			xlog.Errorf("Invalid command: %s", cmd)
			commonresp.ErrResp(c, fmt.Errorf("invalid command"))
			return
		}
	}
}

// swagger:route GET /v1/user_info 用户接口 GetUserInfo
// 获取用户信息
// responses:
//
//	200: GetUserInfoResp
//	500: CommonError
func (h *PublicApiHandler) GetUserInfo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	resp, err := h.srv.GetUserInfo(context.TODO(), userID)
	if err != nil {
		commonresp.ErrResp(c, err)
		return
	}
	commonresp.JsonResp(c, resp)
}

// swagger:route POST /v1/signin_game api渠道接口 SigninGame
// 开游戏
// consumes:
//   - multipart/form-data
//
// responses:
//
//	200: SigninGameResp
//	500: CommonError
func (h *PublicApiHandler) SigninGame(c *gin.Context) {
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
	resp, err := h.srv.SigninGame(c, &req)
	if err != nil {
		commonresp.ErrResp(c, err)
		return
	}
	commonresp.JsonResp(c, resp)
}

// swagger:route POST /v1/member_register api渠道接口 MemberRegister
// 注册用户
// consumes:
//   - multipart/form-data
//
// responses:
//
//	200: MemberRegisterResp
//	500: CommonError
func (h *PublicApiHandler) MemberRegister(c *gin.Context) {
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
	resp, err := h.srv.MemberRegister(c, &req)
	if err != nil {
		commonresp.ErrResp(c, err)
		return
	}
	commonresp.JsonResp(c, resp)
}

// swagger:route POST /v1/edit_limit api渠道接口 EditLimit
// 修改限额
// consumes:
//   - multipart/form-data
//
// responses:
//
//	200: EditLimitResp
//	500: CommonError
func (h *PublicApiHandler) EditLimit(c *gin.Context) {
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
	resp, err := h.srv.EditLimit(c, &req)
	if err != nil {
		commonresp.ErrResp(c, err)
		return
	}
	commonresp.JsonResp(c, resp)
}

// swagger:route POST /v1/logout_game api渠道接口 LogoutGame
// 登出游戏
// consumes:
//   - multipart/form-data
//
// responses:
//
//	200: LogoutGameResp
//	500: CommonError
func (h *PublicApiHandler) LogoutGame(c *gin.Context) {
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
	resp, err := h.srv.LogoutGame(c, &req)
	if err != nil {
		commonresp.ErrResp(c, err)
		return
	}
	commonresp.JsonResp(c, resp)
}

// swagger:route POST /v1/change_password api渠道接口 ChangePassword
// 修改密码
// consumes:
//   - multipart/form-data
//
// responses:
//
//	200: ChangePasswordResp
//	500: CommonError
func (h *PublicApiHandler) ChangePassword(c *gin.Context) {
	var req view.ChangePasswordReq
	req.VendorID = c.PostForm("vendorId")
	req.Signature = c.PostForm("signature")
	req.User = c.PostForm("user")
	req.NewPassword = c.PostForm("newpassword")
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

	xlog.Debugf("ChangePassword req: %+v", &req)
	resp, err := h.srv.ChangePassword(c, &req)
	if err != nil {
		commonresp.ErrResp(c, err)
		return
	}
	commonresp.JsonResp(c, resp)
}

// swagger:route POST /v1/get_agent_balance api渠道接口 GetAgentBalance
// 获取代理商余额
// consumes:
//   - multipart/form-data
//
// responses:
//
//	200: GetAgentBalanceResp
//	500: CommonError
func (h *PublicApiHandler) GetAgentBalance(c *gin.Context) {
	var req view.GetAgentBalanceReq
	req.VendorID = c.PostForm("vendorId")
	req.Signature = c.PostForm("signature")
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

	xlog.Debugf("GetAgentBalance req: %+v", &req)
	resp, err := h.srv.GetAgentBalance(c, &req)
	if err != nil {
		commonresp.ErrResp(c, err)
		return
	}
	commonresp.JsonResp(c, resp)
}
