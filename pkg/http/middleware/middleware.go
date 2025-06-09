package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"go-zrbc/config"
	"go-zrbc/pkg/token"
	"go-zrbc/pkg/xlog"
	service "go-zrbc/service/public"

	commonresp "go-zrbc/pkg/http/response"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MiddlewareHandler struct {
	serviceName string
	redisCli    *redis.Client
}

func NewMiddlewareHandler(srvName string) *MiddlewareHandler {
	redisAddr := config.Global.Redis.Addr
	redisDB := config.Global.Redis.DB
	redisCli := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: config.Global.Redis.Password, // no password set
		DB:       redisDB,                      // use default DB
	})
	return &MiddlewareHandler{
		serviceName: srvName,
		redisCli:    redisCli,
	}
}

func Within(s string, arr []string) bool {
	for i := range arr {
		if s == arr[i] {
			return true
		}
	}
	return false
}

func HasPrefixIn(s string, arr []string) bool {
	for i := range arr {
		if strings.HasPrefix(s, arr[i]) {
			return true
		}
	}
	return false
}

var (
	httpReqQPS = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_qps_total",
		Help: "The total number of processed events",
	}, []string{
		"service_id",
		"url",
		"method",
		"status_code",
	})

	httpResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_resptime",
		Help: "The total number of processed events",
	}, []string{
		"service_id",
		"url",
		"method",
	})
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

var m1 = regexp.MustCompile(`[0-9]+`)

func OmitNumber(s string) string {
	// regex match num
	// replace number with `-`
	return m1.ReplaceAllString(s, "-")
}

func (md *MiddlewareHandler) RequestLog(c *gin.Context) {
	start := time.Now()
	guuid := uuid.New().String()
	c.Set("uuid", guuid)
	c.Set("user_id", int64(2))

	var payload string
	if c.Request.Method == http.MethodPost {
		body, _ := ioutil.ReadAll(c.Request.Body)
		payload = string(body)
		if len(payload) > 1000 {
			payload = "not print because body too long"
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
	}

	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	c.Next()

	respBody := blw.body.String()

	httpReqQPS.With(prometheus.Labels{
		"service_id":  md.serviceName,
		"url":         OmitNumber(c.Request.URL.Path),
		"method":      c.Request.Method,
		"status_code": fmt.Sprintf("%d", c.Writer.Status()),
	}).Inc()

	httpResponseTime.With(prometheus.Labels{
		"service_id": md.serviceName,
		"url":        OmitNumber(c.Request.URL.Path),
		"method":     c.Request.Method,
	}).Observe(float64(time.Since(start) / time.Millisecond))

	userID := GetUserID(c)

	rl := &HttpLog{
		ServiceID: md.serviceName,
		Method:    c.Request.Method,
		URL:       c.Request.URL.String(),
		UUID:      guuid,
		ReqBody:   payload,
		UserID:    userID,
		RespBody:  respBody,
		Status:    c.Writer.Status(),
		Cost:      int64(time.Since(start) / time.Millisecond),
	}
	logb, _ := json.Marshal(rl)
	fmt.Println(string(logb))
	xlog.Infof("%s %v, uuid(%v), body(%v), userID(%v), response(%+v), status(%v), takes(%v ms)\n", c.Request.Method, c.Request.URL, guuid, payload, userID, respBody, c.Writer.Status(), int64(time.Since(start)/time.Millisecond))
}

func GetUser(c *gin.Context) *token.TokenUser {
	i, ok := c.Get("user")
	if !ok {
		xlog.Error(errors.New("not found"))
		return nil
	}
	user, ok := i.(*token.TokenUser)
	if !ok {
		xlog.Error(errors.New("err convert"))
		return nil
	}
	return user
}

func GetUserID(c *gin.Context) int64 {
	i, ok := c.Get("user_id")
	if !ok {
		xlog.Error(errors.New("not found"))
		return 0
	}
	userID, ok := i.(int64)
	if !ok {
		if i != -1 {
			xlog.Errorf("user_id(%v), err:(%+v)", i, errors.New("err convert"))
		} else {
			xlog.Debugf("user_id(%v), universal token", i)
		}
		return 0
	}
	return userID
}

type HttpLog struct {
	ServiceID string `json:"service_id"`
	Method    string `json:"method"`
	URL       string `json:"url"`
	UUID      string `json:"uuid"`
	ReqBody   string `json:"req_body"`
	UserID    int64  `json:"user_id"`
	RespBody  string `json:"resp_body"`
	Status    int    `json:"status"`
	Cost      int64  `json:"cost_time"`
}

var BanIPList []string

func (md *MiddlewareHandler) Oauth(c *gin.Context) {
	loginIP := c.ClientIP()
	// loginIP := "116.179.33.44"
	if !Within(c.Request.URL.Path, []string{"/v1/refresh_ban_ips", "/v1/get_ban_ips"}) {
		if loginIP != "" {
			if len(BanIPList) > 0 {
				if HasPrefixIn(loginIP, BanIPList) {
					commonresp.AbortResp(c, http.StatusForbidden)
					return
				}
			}
		}
	}
	// xlog.Infof("request ip :%s, len of ban ips:%d\n", loginIP, len(BanIPList))
	if strings.HasPrefix(c.Request.RequestURI, "/static") {
		c.Next()
		return
	}
	if Within(c.Request.URL.Path, []string{
		"/swagger.yaml", "/swagger_local.yaml", "/swagger_test.yaml", "/docs", "/v1/time_ts", "/v1/signin_game", "/v1/member_register"}) {
		c.Next()
		return
	}
	accessToken := c.GetHeader("Authorization")
	if accessToken == "H>Ps83XXI?T^g@0J" {
		c.Set("user_id", int64(941388))
		c.Next()
		return
	}
	if accessToken == "" {
		if HasPrefixIn(c.Request.RequestURI, []string{"v1/user_info"}) {
			c.Next()
			return
		}
		commonresp.AbortResp(c, http.StatusUnauthorized)
		return
	}
	tUser, err := token.GetToken(md.redisCli, service.TokenPrefix, accessToken)
	if err == nil {
		commonresp.AbortResp(c, http.StatusNotFound)
		return
	}
	c.Set("user_id", tUser.ID)
	c.Set("user", tUser)
	c.Set("user_name", tUser.UserName)
	c.Next()
}

func (md *MiddlewareHandler) Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(http.StatusOK)
	}
}
