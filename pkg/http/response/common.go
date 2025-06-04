package response

import (
	"net/http"

	"go-zrbc/pkg/xlog"

	"github.com/gin-gonic/gin"
)

func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"code":   401,
		"msg":    http.StatusText(http.StatusUnauthorized),
		"err":    "user not found",
		"status": "fail",
	})
}

func BadRequestResp(c *gin.Context, err error) {
	xlog.Error(err)
	c.JSON(400, gin.H{
		"code":   400,
		"msg":    http.StatusText(400),
		"err":    err.Error(),
		"status": "fail",
	})
}

func ServerErrResp(c *gin.Context, err error) {
	c.JSON(500, gin.H{
		"code":   500,
		"msg":    http.StatusText(500),
		"err":    err.Error(),
		"status": "fail",
	})
}

func UnauthorizationResp(c *gin.Context, err error) {
	c.JSON(403, gin.H{
		"msg":    http.StatusText(http.StatusUnauthorized),
		"code":   403,
		"status": "fail",
	})
}

func ErrorPasswordResp(c *gin.Context, err error) {
	c.JSON(403, gin.H{
		"msg":/* "The account or password is incorrect"*/ "账户不存在",
		"code":   40301,
		"status": "fail",
	})
}

func ErrorMAResp(c *gin.Context, err error) {
	c.JSON(403, gin.H{
		"msg":/*"The message authentication code is incorrect"*/ "短信验证码错误",
		"code":   40302,
		"status": "fail",
	})
}

func ErrResp(c *gin.Context, err error) {
	xlog.Error(err)
	if customErr, ok := err.(*CustomError); ok {
		c.JSON(200, gin.H{
			"code":   customErr.Code,
			"msg":    customErr.Message,
			"status": "fail",
			"data":   nil,
		})
		return
	}
	// 默认系统错误
	c.JSON(200, gin.H{
		"code":   CodeSystemError,
		"msg":    err.Error(),
		"status": "fail",
		"data":   nil,
	})
}

func JsonResp(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "",
		"status": "success",
		"data":   data,
	})
}

func ForbiddenResp(c *gin.Context, err error) {
	xlog.Error(err)
	c.JSON(403, gin.H{
		"code": 40300,
		"msg":/*http.StatusText(403)*/ "账户被封禁",
		"err":    err.Error(),
		"status": "forbidden",
	})
}

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func AbortResp(c *gin.Context, code int, msg ...string) {
	respMsg := http.StatusText(code)
	if len(msg) > 0 {
		respMsg = msg[0]
	}
	c.AbortWithStatusJSON(code,
		&Response{
			Code: code,
			Msg:  respMsg,
		})
}
