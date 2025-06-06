package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-zrbc/config"
	"go-zrbc/db"
	"go-zrbc/pkg/utils"
	"go-zrbc/pkg/xlog"
	"go-zrbc/service"
	"go-zrbc/view"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	commonresp "go-zrbc/pkg/http/response"

	"gorm.io/gorm"
)

const (
	TokenPrefix = "zrys:access_token"

	// Language codes mapping
	LangChinese    = "cn"
	LangEnglish    = "en"
	LangThai       = "th"
	LangVietnamese = "vi"
	LangJapanese   = "ja"
	LangKorean     = "ko"
	LangHindi      = "hi"
	LangMalay      = "ms"
	LangIndonesian = "in"
	LangTaiwan     = "tw"
	LangSpanish    = "es"
)

// LanguageMap maps numeric language codes to their corresponding language parameter strings
var LanguageMap = map[int]string{
	0:  "&lang=" + LangChinese,
	1:  "&lang=" + LangEnglish,
	2:  "&lang=" + LangThai,
	3:  "&lang=" + LangVietnamese,
	4:  "&lang=" + LangJapanese,
	5:  "&lang=" + LangKorean,
	6:  "&lang=" + LangHindi,
	7:  "&lang=" + LangMalay,
	8:  "&lang=" + LangIndonesian,
	9:  "&lang=" + LangTaiwan,
	10: "&lang=" + LangSpanish,
}

type UserService interface {
	//用户信息
	GetUserInfo(ctx context.Context, userID int64) (*view.GetUserInfoResp, error)
	GetUserByAccountAndPwd(ctx context.Context, account, pwd string) (*view.GetUserInfoResp, error)
	SigninGame(ctx context.Context, req *view.SigninGameReq) (*view.SigninGameResp, error)
	AgentVerify(ctx context.Context, req *view.AgentVerifyReq) (*view.AgentVerifyResp, error)
}

type userService struct {
	userDao            db.UserDao
	apiurlDao          db.ApiurlDao
	wechatURLDao       db.WechatURLDao
	agentsLoginPassDao db.AgentsLoginPassDao
	agentDao           db.AgentDao
	memLoginDao        db.MemLoginDao

	s3Client *s3.Client
	redisCli *redis.Client
	*service.Session
}

func NewUserService(
	sess *service.Session,
	userDao db.UserDao,
	apiurlDao db.ApiurlDao,
	wechatURLDao db.WechatURLDao,
	agentsLoginPassDao db.AgentsLoginPassDao,
	agentDao db.AgentDao,
	memLoginDao db.MemLoginDao,

	s3Client *s3.Client,
	redisCli *redis.Client,
) UserService {
	srv := &userService{
		userDao:            userDao,
		apiurlDao:          apiurlDao,
		wechatURLDao:       wechatURLDao,
		agentsLoginPassDao: agentsLoginPassDao,
		agentDao:           agentDao,
		memLoginDao:        memLoginDao,

		s3Client: s3Client,
		redisCli: redisCli,
	}
	srv.Session = sess
	return srv
}

func (srv *userService) GetUserInfo(ctx context.Context, userID int64) (*view.GetUserInfoResp, error) {
	var ret *db.Member
	var err error
	err = srv.Tx(func(tx *gorm.DB) error {
		ret, err = srv.userDao.QueryByID(tx, userID)
		if err != nil {
			xlog.Errorf("error to get user info, err:%+v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	resp := view.GetUserInfoResp{
		User: DBToViewUser(ret),
	}
	return &resp, nil
}

func (srv *userService) GetUserByAccountAndPwd(ctx context.Context, account, pwd string) (*view.GetUserInfoResp, error) {
	var ret *db.Member
	var err error

	err = srv.Tx(func(tx *gorm.DB) error {
		ret, err = srv.userDao.QueryByAccountAndPwd(tx, account, pwd)
		if err != nil {
			xlog.Errorf("error to get user by account and pwd, err:%+v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	resp := view.GetUserInfoResp{
		User: DBToViewUser(ret),
	}
	return &resp, nil
}

func (srv *userService) GetUserByAccount(ctx context.Context, account string) (*view.MemberCache, error) {
	var ret *db.Member
	var err error

	// First try to get user info from Redis
	userInfoJSON, err := srv.redisCli.HGet(ctx, "user", account).Result()
	if err == nil && userInfoJSON != "" {
		// User info found in Redis
		var mem view.MemberCache
		if err := json.Unmarshal([]byte(userInfoJSON), &mem); err == nil {
			xlog.Debugf("debug to unmarshal user info from Redis, userInfoJSON:%s, mem:%+v", userInfoJSON, mem)
			return &mem, nil
		}
		// If unmarshal fails, continue with DB query
		xlog.Warnf("Failed to unmarshal user info from Redis: %v", err)
	}

	// If not found in Redis or error occurred, query from database
	err = srv.Tx(func(tx *gorm.DB) error {
		ret, err = srv.userDao.QueryByAccount(tx, account)
		if err != nil {
			xlog.Errorf("error to get user by account, err:%+v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	mem := DBToViewUserCache(ret)
	// Store user information in Redis
	var userRedisData []byte
	var marshalErr error
	userRedisData, marshalErr = json.Marshal(mem)
	if marshalErr != nil {
		xlog.Warnf("Failed to marshal user info for Redis: %v", marshalErr)
	} else {
		redisErr := srv.redisCli.HSet(ctx, "user", account, string(userRedisData)).Err()
		if redisErr != nil {
			xlog.Warnf("Failed to store user info in Redis: %v", redisErr)
		}
	}

	return mem, nil
}

func (srv *userService) validateRequest(user, password string, isTest bool) error {
	if !isTest {
		if password == "" {
			return commonresp.ErrUserPWDEmpty
		}
		if user == "" {
			return commonresp.ErrUserEmpty
		}

		mem, err := srv.userDao.QueryByAccount(srv.DB(), user)
		if err == gorm.ErrRecordNotFound {
			return commonresp.ErrUserNotExist
		}
		if err != err {
			return err
		}

		if mem.Password != password {
			return commonresp.ErrUserPWD
		}
	}
	return nil
}

func processUI(ui int) int {
	switch ui {
	case 1, 2:
		return 2
	case 4:
		return 4
	case 5:
		return 5
	default:
		return 0
	}
}

func (srv *userService) getAPIURL(code int) (string, error) {
	apiURL, err := srv.apiurlDao.QueryByCode(srv.DB(), code)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			defaultURL, err := srv.apiurlDao.QueryByCode(srv.DB(), 0)
			if err != nil {
				return "", err
			}
			return defaultURL.URL, nil
		}
		return "", err
	}
	return apiURL.URL, nil
}

func (srv *userService) getWechatURL() (string, error) {
	var urlData *db.WechatURL
	var err error
	err = srv.Tx(func(tx *gorm.DB) error {
		urlData, err = srv.wechatURLDao.GetRandomWechatURL(tx)
		if err != nil {
			return err
		}

		if err := srv.wechatURLDao.UpdateWechatURLUseCount(tx, urlData.ID); err != nil {
			// Log error but continue
			xlog.Errorf("Failed to update Wechat URL use count, err:%+v", err)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return urlData.URL, nil
}

// GetClientIP gets the client IP from the gin context
func GetClientIP(c *gin.Context) string {
	// Try to get IP from X-Real-IP header
	clientIP := c.GetHeader("X-Real-IP")
	if clientIP != "" {
		return clientIP
	}

	// Try to get IP from X-Forwarded-For header
	clientIP = c.GetHeader("X-Forwarded-For")
	if clientIP != "" {
		return strings.Split(clientIP, ",")[0]
	}

	// Get IP from RemoteAddr
	return c.ClientIP()
}

func (srv *userService) SigninGame(ctx context.Context, req *view.SigninGameReq) (*view.SigninGameResp, error) {
	avResp, err := srv.AgentVerify(ctx, &view.AgentVerifyReq{VendorID: req.VendorID, Signature: req.Signature})
	if err != nil {
		xlog.Errorf("error to verify agent, err:%+v", err)
		return nil, err
	}

	// Validate request
	if !req.IsTest {
		if err := srv.validateRequest(req.User, req.Password, req.IsTest); err != nil {
			xlog.Errorf("error to get user, err:%+v", err)
			return nil, err
		}
	}

	// Define valid game modes
	validModes := []string{
		"onlybac", "onlydgtg", "onlyrou", "onlysicbo", "onlyniuniu",
		"onlysamgong", "onlyfantan", "onlysedie", "onlyfishshrimpcrab",
		"onlygoldenflower", "onlymultiple", "onlypaigow", "onlythisbar",
		"onlybactable", "onlydgtgtable", "onlyroutable", "onlysicbotable",
		"onlyniuniutable", "onlysamgongtable", "onlyfantantable", "onlysedietable",
		"onlyfishshrimpcrabtable", "onlygoldenflowertable", "onlypaigowtable",
		"onlythisbartable",
	}

	// Process mode parameter
	modeParam := ""
	if req.Mode != "" {
		isValidMode := false
		for _, validMode := range validModes {
			if req.Mode == validMode {
				isValidMode = true
				break
			}
		}
		if isValidMode {
			if req.TableID != "" {
				modeParam = "&mode=" + req.Mode + "&tableid=" + req.TableID
			} else {
				modeParam = "&mode=" + req.Mode
			}
		}
	}

	// Process mute parameter
	muteParam := ""
	if req.Mute == "true" {
		muteParam = "&mute=" + req.Mute
	}

	// Process UI settings
	ui := strconv.Itoa(processUI(req.UI))

	langInt, err := strconv.Atoi(req.Lang)
	if err != nil {
		xlog.Warnf("Invalid language code %s, using default (0)", req.Lang)
		langInt = 0
	}

	var baseURL, apiURL, wechatURL, originURL string
	apiURL, err = srv.getAPIURL(avResp.Agent.Currency)
	if err != nil {
		xlog.Errorf("error to get base url, err:%+v", err)
		return nil, err
	}
	wechatURL, err = srv.getWechatURL()
	if err != nil {
		xlog.Errorf("error to get wechat url, err:%+v", err)
		return nil, err
	}

	// Check if vendor is in other link vendor list
	otherLinkVendors := []string{"keaoapi", "qianhuapi", "zunlongapi", "lsjrmbapi", "lsjthbapi", "lsjmyrapi"}
	isOtherLinkVendor := false
	for _, v := range otherLinkVendors {
		if avResp.Agent.VendorID == v {
			isOtherLinkVendor = true
			break
		}
	}
	serverName := utils.GetServerName(ctx.(*gin.Context))
	if isOtherLinkVendor {
		baseURL = wechatURL
		originURL = baseURL
	} else if serverName == "api-live01.wmexpo.net" {
		originURL = "https://a45.me/"
	} else {
		originURL = apiURL
	}

	// Handle site type specific URL selection
	siteType := req.Site
	switch siteType {
	case "6", "9", "99":
		baseURL = wechatURL
	default:
		baseURL = originURL
	}

	var gameURL string
	if req.IsTest {
		if req.Mode != "" {
			if len(ui) != 0 {
				if avResp.Agent.VendorID == "bbinapi" || avResp.Agent.VendorID == "bbtest" || avResp.Agent.VendorID == "bbintwapi" {
					gameURL = fmt.Sprintf("%s?#sid=ANONYMOUS%s%s&ui=%s%s", baseURL, modeParam, muteParam, ui, LanguageMap[langInt])
				} else {
					gameURL = fmt.Sprintf("%s?sid=ANONYMOUS%s%s&ui=%s%s", baseURL, modeParam, muteParam, ui, LanguageMap[langInt])
				}
			} else {
				if avResp.Agent.VendorID == "bbinapi" || avResp.Agent.VendorID == "bbtest" || avResp.Agent.VendorID == "bbintwapi" {
					gameURL = fmt.Sprintf("%s?#sid=ANONYMOUS%s%s%s", baseURL, modeParam, muteParam, LanguageMap[langInt])
				} else {
					gameURL = fmt.Sprintf("%s?sid=ANONYMOUS%s%s%s", baseURL, modeParam, muteParam, LanguageMap[langInt])
				}
			}
		} else {
			if len(ui) != 0 {
				if avResp.Agent.VendorID == "bbinapi" || avResp.Agent.VendorID == "bbtest" || avResp.Agent.VendorID == "bbintwapi" {
					gameURL = fmt.Sprintf("%s?#sid=ANONYMOUS&ui=%s%s%s", baseURL, muteParam, ui, LanguageMap[langInt])
				} else {
					gameURL = fmt.Sprintf("%s?sid=ANONYMOUS&ui=%s%s%s", baseURL, muteParam, ui, LanguageMap[langInt])
				}
			} else {
				if avResp.Agent.VendorID == "bbinapi" || avResp.Agent.VendorID == "bbtest" || avResp.Agent.VendorID == "bbintwapi" {
					gameURL = fmt.Sprintf("%s?#sid=ANONYMOUS%s%s", baseURL, muteParam, LanguageMap[langInt])
				} else {
					gameURL = fmt.Sprintf("%s?sid=ANONYMOUS%s%s", baseURL, muteParam, LanguageMap[langInt])
				}
			}
		}
		if req.ReturnURL != "" {
			gameURL += "&returnurl=" + req.ReturnURL
		}
		return &view.SigninGameResp{
			GameURL: gameURL,
		}, nil
	}

	mem, err := srv.GetUserByAccount(ctx, req.User)
	if err != nil {
		xlog.Errorf("error to get user by account, account:%s, err:%+v", req.User, err)
		return nil, err
	}
	xlog.Debugf("mem:%+v", mem)
	ulv, utp := strconv.Itoa(mem.ULV), mem.UTP
	sid := utils.ProSIDCreate(config.Global.Wcode, ulv, utp, mem.UID, config.Global.SidLen)

	// Set cookie with 18 hour expiration
	cookieTime := time.Now().Add(18 * time.Hour)
	ginCtx := ctx.(*gin.Context)
	ginCtx.SetCookie(strings.ToUpper(config.Global.Wcode)+"[1]", sid, int(time.Until(cookieTime).Seconds()), "/", ginCtx.Request.Host, false, true)

	var newMemLogin *db.MemLogin
	switch ulvInt, _ := strconv.Atoi(ulv); ulvInt {
	case 0:
		// No action needed
	case 1, 2, 3, 4, 5:
		// No action needed
	case 7:
		now := time.Now()
		err = srv.Tx(func(tx *gorm.DB) error {
			newMemLogin, err = srv.memLoginDao.CreateOrUpdateMemLogin(tx, mem.UID, 0, sid, GetClientIP(ctx.(*gin.Context)), now)
			if err != nil {
				return err
			}
			if err := srv.userDao.UpdatesMember(tx, mem.UID, map[string]interface{}{
				"mem013": now,
				"mem014": GetClientIP(ctx.(*gin.Context)),
			}); err != nil {
				xlog.Errorf("Failed to update member, err:%+v", err)
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	if newMemLogin != nil {
		xlog.Debug("newMemLogin is not nil")
		memLoginSID := ""
		if newMemLogin != nil && newMemLogin.Mlg003 != "" {
			memLoginSID = newMemLogin.Mlg003
		} else {
			memLoginSID = sid
		}

		if len(ui) != 0 && req.Mode != "" {
			if avResp.Agent.VendorID == "bbinapi" || avResp.Agent.VendorID == "bbtest" || avResp.Agent.VendorID == "bbintwapi" {
				gameURL = fmt.Sprintf("%s?#sid=%s%s%s&ui=%s%s", baseURL, memLoginSID, modeParam, muteParam, ui, LanguageMap[langInt])
			} else {
				gameURL = fmt.Sprintf("%s?sid=%s%s%s&ui=%s%s", baseURL, memLoginSID, modeParam, muteParam, ui, LanguageMap[langInt])
			}
		} else if len(ui) != 0 && req.Mode == "" {
			if avResp.Agent.VendorID == "bbinapi" || avResp.Agent.VendorID == "bbtest" || avResp.Agent.VendorID == "bbintwapi" {
				gameURL = fmt.Sprintf("%s?#sid=%s%s&ui=%s%s", baseURL, memLoginSID, muteParam, ui, LanguageMap[langInt])
			} else {
				gameURL = fmt.Sprintf("%s?sid=%s%s&ui=%s%s", baseURL, memLoginSID, muteParam, ui, LanguageMap[langInt])
			}
		} else if req.Mode != "" {
			if avResp.Agent.VendorID == "bbinapi" || avResp.Agent.VendorID == "bbtest" || avResp.Agent.VendorID == "bbintwapi" {
				gameURL = fmt.Sprintf("%s?#sid=%s%s%s%s", baseURL, memLoginSID, modeParam, muteParam, LanguageMap[langInt])
			} else {
				gameURL = fmt.Sprintf("%s?sid=%s%s%s%s", baseURL, memLoginSID, modeParam, muteParam, LanguageMap[langInt])
			}
		} else {
			if avResp.Agent.VendorID == "bbinapi" || avResp.Agent.VendorID == "bbtest" || avResp.Agent.VendorID == "bbintwapi" {
				gameURL = fmt.Sprintf("%s?#sid=%s%s%s", baseURL, memLoginSID, muteParam, LanguageMap[langInt])
			} else {
				gameURL = fmt.Sprintf("%s?sid=%s%s%s", baseURL, memLoginSID, muteParam, LanguageMap[langInt])
			}
		}
	} else {
		xlog.Debug("newMemLogin is nil")
		if req.Mode != "" {
			if avResp.Agent.VendorID == "bbinapi" || avResp.Agent.VendorID == "bbtest" || avResp.Agent.VendorID == "bbintwapi" {
				gameURL = fmt.Sprintf("%s?#sid=%s%s%s%s", baseURL, sid, modeParam, muteParam, LanguageMap[langInt])
			} else {
				gameURL = fmt.Sprintf("%s?sid=%s%s%s%s", baseURL, sid, modeParam, muteParam, LanguageMap[langInt])
			}
		} else {
			if avResp.Agent.VendorID == "bbinapi" || avResp.Agent.VendorID == "bbtest" || avResp.Agent.VendorID == "bbintwapi" {
				gameURL = fmt.Sprintf("%s?#sid=%s%s%s", baseURL, sid, muteParam, LanguageMap[langInt])
			} else {
				gameURL = fmt.Sprintf("%s?sid=%s%s%s", baseURL, sid, muteParam, LanguageMap[langInt])
			}
		}
	}

	if req.Width == "1" {
		gameURL += "&checkWidth=deviceWidth"
	}
	if req.ReturnURL != "" {
		gameURL += "&returnurl=" + req.ReturnURL
	}
	if req.Size == "1" {
		gameURL += "&size=1"
	}
	if req.Video == "off" {
		gameURL += "&video=" + req.Video
	}
	if avResp.AgentsLoginPass.Co != "" {
		gameURL += "&co=" + avResp.AgentsLoginPass.Co
	} else {
		gameURL += "&co=wm"
	}

	return &view.SigninGameResp{
		GameURL: gameURL,
	}, nil
}

func (srv *userService) AgentVerify(ctx context.Context, req *view.AgentVerifyReq) (*view.AgentVerifyResp, error) {
	// Validate request parameters
	if req.VendorID == "" && req.Signature == "" {
		err := errors.New("vendor id or signature is empty")
		xlog.Error(err)
		return nil, err
	}
	var err error
	var agent *db.Agent
	var agentsLoginPass *db.AgentsLoginPass
	err = srv.Tx(func(tx *gorm.DB) error {
		agent, err = srv.agentDao.QueryByVendorID(srv.DB(), req.VendorID)
		if err != nil {
			xlog.Errorf("error to query agent by vendor id:%s, err:%+v", req.VendorID, err)
			return err
		}
		agentsLoginPass, err = srv.agentsLoginPassDao.QueryByAidAndVendorID(tx, agent.Age001, req.VendorID)
		if err != nil {
			xlog.Errorf("error to query agents login pass by aid:%d and vendor id:%s, err:%+v", agent.Age001, req.VendorID, err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if agentsLoginPass.Signature != req.Signature {
		err := errors.New("invalid signature")
		xlog.Error(err)
		return nil, err
	}
	return &view.AgentVerifyResp{
		Agent:           DBToViewAgent(agent),
		AgentsLoginPass: DBToViewAgentsLoginPass(agentsLoginPass),
	}, nil
}
