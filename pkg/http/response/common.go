package response

import (
	"net/http"

	"go-zrbc/pkg/utils"
	"go-zrbc/pkg/xlog"

	"github.com/gin-gonic/gin"
)

func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"errorCode":   401,
		"errorRemark": http.StatusText(http.StatusUnauthorized),
	})
}

func BadRequestResp(c *gin.Context, err error) {
	xlog.Error(err)
	c.JSON(400, gin.H{
		"errorCode":   400,
		"errorRemark": http.StatusText(400),
	})
}

func ServerErrResp(c *gin.Context, err error) {
	c.JSON(500, gin.H{
		"errorCode":   500,
		"errorRemark": http.StatusText(500),
	})
}

func UnauthorizationResp(c *gin.Context, err error) {
	c.JSON(403, gin.H{
		"errorRemark": http.StatusText(http.StatusUnauthorized),
		"errorCode":   403,
	})
}

func ErrorPasswordResp(c *gin.Context, err error) {
	c.JSON(403, gin.H{
		"errorRemark": "账户不存在",
		"errorCode":   40301,
	})
}

func ErrorMAResp(c *gin.Context, err error) {
	c.JSON(403, gin.H{
		"errorRemark": "短信验证码错误",
		"errorCode":   40302,
	})
}

func ErrResp(c *gin.Context, err error) {
	xlog.Error(err)
	if customErr, ok := err.(*utils.CustomError); ok {
		c.JSON(200, gin.H{
			"errorCode":   customErr.Code,
			"errorRemark": customErr.Message,
			"data":        nil,
		})
		return
	}
	// 默认系统错误
	c.JSON(200, gin.H{
		"errorCode":   utils.CodeSystemError,
		"errorRemark": err.Error(),
		"data":        nil,
	})
}

func JsonResp(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"errorCode":   200,
		"errorRemark": "",
		"data":        data,
	})
}

func ForbiddenResp(c *gin.Context, err error) {
	xlog.Error(err)
	c.JSON(403, gin.H{
		"errorCode":   40300,
		"errorRemark": "账户被封禁",
	})
}

type Response struct {
	ErrorCode   int    `json:"errorCode"`
	ErrorRemark string `json:"errorRemark"`
}

func AbortResp(c *gin.Context, code int, msg ...string) {
	respMsg := http.StatusText(code)
	if len(msg) > 0 {
		respMsg = msg[0]
	}
	c.AbortWithStatusJSON(code,
		&Response{
			ErrorCode:   code,
			ErrorRemark: respMsg,
		})
}
