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
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"

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

var ChipsMap = map[int]string{
	17: "10000,50000,100000,1000000,5000000",
	18: "1000,10000,100000,200000,1000000",
}

var ChipsCheck = []int{1, 5, 10, 20, 50, 100, 500, 1000, 5000, 10000, 20000, 50000, 100000, 200000, 1000000, 5000000, 10000000, 20000000, 50000000, 10000000, 20000000, 50000000, 100000000}

type PublicApiService interface {
	//用户信息
	GetUserInfo(ctx context.Context, userID int64) (*view.GetUserInfoResp, error)
	GetUserByAccountAndPwd(ctx context.Context, account, pwd string) (*view.GetUserInfoResp, error)
	GetUserByAccount(ctx context.Context, account string) (*view.MemberCache, error)
	SigninGame(ctx context.Context, req *view.SigninGameReq) (*view.SigninGameResp, error)
	MemberRegister(ctx context.Context, req *view.MemberRegisterReq) (*view.MemberRegisterResp, error)
	AgentVerify(ctx context.Context, req *view.AgentVerifyReq) (*view.AgentVerifyResp, error)
	EditLimit(ctx context.Context, req *view.EditLimitReq) (*view.EditLimitResp, error)
	LogoutGame(ctx context.Context, req *view.LogoutGameReq) (*view.LogoutGameResp, error)
	ChangePassword(ctx context.Context, req *view.ChangePasswordReq) (*view.ChangePasswordResp, error)
	GetAgentBalance(ctx context.Context, req *view.GetAgentBalanceReq) (*view.GetAgentBalanceResp, error)
	GetBalance(ctx context.Context, req *view.GetBalanceReq) (*view.GetBalanceResp, error)
	ChangeBalance(ctx context.Context, req *view.ChangeBalanceReq) (*view.ChangeBalanceResp, error)
	GetMemberTradeReport(ctx context.Context, req *view.GetMemberTradeReportReq) (*view.GetMemberTradeReportResp, error)
	EnableOrDisableMem(ctx context.Context, req *view.EnableOrDisableMemReq) (*view.EnableOrDisableMemResp, error)
}

type MemDtlDao interface {
	QueryByMemberID(tx *gorm.DB, memberID int64) ([]*db.MemberDtl, error)
}

type BetLimitDefaultDao interface {
	QueryAll(tx *gorm.DB) ([]*db.BetLimitDefault, error)
	QueryByGtype(tx *gorm.DB, gtype int) ([]*db.BetLimitDefault, error)
}

type publicApiService struct {
	userDao             db.UserDao
	apiurlDao           db.ApiurlDao
	wechatURLDao        db.WechatURLDao
	agentsLoginPassDao  db.AgentsLoginPassDao
	agentDao            db.AgentDao
	memLoginDao         db.MemLoginDao
	bet02Dao            db.Bet02Dao
	agentDtlDao         db.AgentDtlDao
	betLimitDefaultDao  db.BetLimitDefaultDao
	memDtlDao           db.MemberDtlDao
	gameTypeDao         db.GameTypeDao
	inOutMDao           db.InOutMDao
	logAgeCashChangeDao db.LogAgeCashChangeDao
	alertMessageDao     db.AlertMessageDao

	s3Client *s3.Client
	redisCli *redis.Client
	*service.Session
}

func NewPublicApiService(
	sess *service.Session,
	userDao db.UserDao,
	apiurlDao db.ApiurlDao,
	wechatURLDao db.WechatURLDao,
	agentsLoginPassDao db.AgentsLoginPassDao,
	agentDao db.AgentDao,
	memLoginDao db.MemLoginDao,
	bet02Dao db.Bet02Dao,
	agentDtlDao db.AgentDtlDao,
	betLimitDefaultDao db.BetLimitDefaultDao,
	memDtlDao db.MemberDtlDao,
	gameTypeDao db.GameTypeDao,
	inOutMDao db.InOutMDao,
	logAgeCashChangeDao db.LogAgeCashChangeDao,
	alertMessageDao db.AlertMessageDao,

	s3Client *s3.Client,
	redisCli *redis.Client,
) PublicApiService {
	srv := &publicApiService{
		userDao:             userDao,
		apiurlDao:           apiurlDao,
		wechatURLDao:        wechatURLDao,
		agentsLoginPassDao:  agentsLoginPassDao,
		agentDao:            agentDao,
		memLoginDao:         memLoginDao,
		bet02Dao:            bet02Dao,
		agentDtlDao:         agentDtlDao,
		betLimitDefaultDao:  betLimitDefaultDao,
		memDtlDao:           memDtlDao,
		gameTypeDao:         gameTypeDao,
		inOutMDao:           inOutMDao,
		logAgeCashChangeDao: logAgeCashChangeDao,
		alertMessageDao:     alertMessageDao,

		s3Client: s3Client,
		redisCli: redisCli,
	}
	srv.Session = sess
	return srv
}

func (srv *publicApiService) GetUserInfo(ctx context.Context, userID int64) (*view.GetUserInfoResp, error) {
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

func (srv *publicApiService) GetUserByAccountAndPwd(ctx context.Context, account, pwd string) (*view.GetUserInfoResp, error) {
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

func (srv *publicApiService) GetUserByAccount(ctx context.Context, account string) (*view.MemberCache, error) {
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

func (srv *publicApiService) validateRequest(user, password string, isTest bool) error {
	if !isTest {
		if password == "" {
			return utils.ErrInvalidPasswordEmpty
		}
		if user == "" {
			return utils.ErrInvalidAccountEmpty
		}

		mem, err := srv.userDao.QueryByAccount(srv.DB(), user)
		if err == gorm.ErrRecordNotFound {
			return utils.ErrParamInvalidAccountNotExist
		}
		if err != nil {
			return err
		}

		if mem.Password != password {
			return utils.ErrParamInvalidAccountPasswordError
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

func (srv *publicApiService) getAPIURL(code int) (string, error) {
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

func (srv *publicApiService) getWechatURL() (string, error) {
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

func (srv *publicApiService) SigninGame(ctx context.Context, req *view.SigninGameReq) (*view.SigninGameResp, error) {
	if err := utils.CheckTimestamp(req.Timestamp); err != nil {
		xlog.Errorf("error to check timestamp, err:%+v", err)
		return nil, err
	}
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
			Result: gameURL,
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
		Result: gameURL,
	}, nil
}

func (srv *publicApiService) AgentVerify(ctx context.Context, req *view.AgentVerifyReq) (*view.AgentVerifyResp, error) {
	// Validate request parameters
	if req.VendorID == "" && req.Signature == "" {
		return nil, utils.ErrAgentIDAndSignatureFormatError
	}
	if req.VendorID == "" {
		return nil, utils.ErrAgentIDEmpty
	}
	if req.Signature == "" {
		return nil, utils.ErrAgentSignatureEmpty
	}
	var err error
	var agent *db.Agent
	var agentsLoginPass *db.AgentsLoginPass
	err = srv.Tx(func(tx *gorm.DB) error {
		agent, err = srv.agentDao.QueryByVendorID(srv.DB(), req.VendorID)
		if err != nil {
			xlog.Errorf("error to query agent by vendor id:%s, err:%+v", req.VendorID, err)
			return utils.ErrAgentIDNotExist
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
		err := utils.ErrAgentIDExistButSignatureError
		xlog.Error(err)
		return nil, err
	}
	return &view.AgentVerifyResp{
		Agent:           DBToViewAgent(agent),
		AgentsLoginPass: DBToViewAgentsLoginPass(agentsLoginPass),
	}, nil
}

func (srv *publicApiService) validateMemberRegisterInput(req *view.MemberRegisterReq) error {
	// Password validation
	if req.Password == "" {
		return utils.ErrParamInvalidPasswordEmpty
	}

	// Check for Chinese characters in password (Unicode range 4e00-9fa5)
	matched, _ := regexp.MatchString(`^[\x{4e00}-\x{9fa5}]{1,30}$`, req.Password)
	if matched {
		return utils.ErrInvalidPasswordChinese
	}

	if len(req.Password) <= 5 {
		return utils.ErrInvalidPasswordLengthShort
	}
	if len(req.Password) >= 65 {
		return utils.ErrInvalidPasswordLengthLong
	}

	// Account validation
	if req.User == "" {
		return utils.ErrParamInvalidAccountNameEmpty
	}

	matched, _ = regexp.MatchString(`^[A-Za-z0-9@_]+$`, req.User)
	if !matched {
		return utils.ErrInvalidAccountFormat
	}

	if strings.Contains(req.User, ".") || strings.Contains(req.Password, ".") {
		return utils.ErrParamInvalidAccountPasswordIllegal
	}

	if len(req.User) <= 4 {
		return utils.ErrInvalidAccountLengthShort
	}
	if len(req.User) > 31 {
		return utils.ErrInvalidAccountLength
	}

	// Username validation
	if req.Username == "" {
		return utils.ErrInvalidUsernameEmpty
	}

	if len(req.Username) > 31 {
		return utils.ErrInvalidUsernameLengthLong
	}

	// Mark validation
	if len(req.Mark) >= 21 {
		return utils.ErrInvalidMarkLengthLong
	}

	return nil
}

func validateChips(chips string) (string, error) {
	chips = strings.ReplaceAll(chips, " ", "")
	chipsList := strings.Split(chips, ",")

	// Validate each chip value
	chipMap := make(map[string]bool)
	for _, chip := range chipsList {
		// Add validation for allowed chip values
		chipInt, err := strconv.Atoi(chip)
		if err != nil {
			return "", utils.ErrInvalidChipsFormat
		}
		if !slices.Contains(ChipsCheck, chipInt) {
			return "", utils.ErrInvalidChipsType
		}
		chipMap[chip] = true
	}

	if len(chipMap) < 5 || len(chipMap) > 10 {
		return "", utils.ErrInvalidChipsCount
	}
	chipsStr := ""
	for k := range chipMap {
		chipsStr += fmt.Sprintf("%s,", k)
	}

	return chipsStr, nil
}

func (srv *publicApiService) MemberRegister(ctx context.Context, req *view.MemberRegisterReq) (*view.MemberRegisterResp, error) {
	if err := utils.CheckTimestamp(req.Timestamp); err != nil {
		xlog.Errorf("error to check timestamp, err:%+v", err)
		return nil, err
	}
	avResp, err := srv.AgentVerify(ctx, &view.AgentVerifyReq{VendorID: req.VendorID, Signature: req.Signature})
	if err != nil {
		xlog.Errorf("error to verify agent, err:%+v", err)
		return nil, err
	}
	if avResp.Agent.Age015 == "N" || avResp.Agent.Age016 == "N" {
		return nil, utils.ErrParamInvalidAgentDeactivated
	}

	// Validate input
	if err := srv.validateMemberRegisterInput(req); err != nil {
		return nil, err
	}

	// Check if account exists
	exists, err := srv.userDao.QueryByAccount(srv.DB(), req.User)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}
	if exists != nil {
		return nil, utils.ErrAccountExists
	}

	currentCount, err := srv.userDao.GetMemberCountByAgentID(srv.DB(), avResp.Agent.ID)
	if err != nil {
		return nil, err
	}

	// Check member count limit
	if avResp.Agent.Membermax > 0 && currentCount >= avResp.Agent.Membermax {
		return nil, utils.ErrInvalidMemberLimit
	}

	// Check profit max
	if avResp.Agent.Profitmax != "" {
		profitMax, err := decimal.NewFromString(avResp.Agent.Profitmax)
		if err != nil {
			return nil, err
		}
		winloss, err := srv.bet02Dao.GetAgentWinloss(srv.DB(), avResp.Agent.ID)
		if err != nil {
			return nil, err
		}
		if winloss.LessThan(profitMax) {
			return nil, utils.ErrInvalidAgentProfitLimit
		}
	}

	insinfo := make(map[string]interface{}, 0)
	// Use agent default limits
	agentlimitsMap, err := srv.agentDtlDao.GetDefaultLimits(srv.DB(), avResp.Agent.ID)
	if err != nil {
		return nil, err
	}
	// Get parent agent's limits
	agentLimits, err := srv.agentDtlDao.QueryByAgentID(srv.DB(), avResp.Agent.ID)
	if err != nil {
		return nil, err
	}
	if req.LimitType == "" {
		for k, v := range agentlimitsMap {
			insinfo[k] = v
		}
	} else {
		// Validate and use custom limits
		req.LimitType = strings.ReplaceAll(req.LimitType, " ", "")
		limitTypes := strings.Split(req.LimitType, ",")
		// Check if limit types exist in parent agent's limits
		validLimits := make([]string, 0)
		for _, agentLimit := range agentLimits {
			tmpAgentLimit := strings.ReplaceAll(agentLimit.Ag014, " ", "")
			if tmpAgentLimit != "" {
				parentLimits := strings.Split(tmpAgentLimit, ",")
				for _, requestedLimit := range limitTypes {
					for _, parentLimit := range parentLimits {
						if requestedLimit == parentLimit {
							validLimits = append(validLimits, requestedLimit)
						}
					}
				}
			}
		}

		if len(utils.ArrayDiffString(limitTypes, validLimits)) != 0 {
			return nil, utils.ErrInvalidLimitTypeNotOpen
		}

		// Get default bet limits
		defaultLimits := []string{}
		bLDs, err := srv.betLimitDefaultDao.QueryAll(srv.DB())
		if err != nil {
			return nil, err
		}
		for _, dl := range bLDs {
			defaultLimits = append(defaultLimits, strconv.Itoa(int(dl.ID)))
		}
		if len(utils.ArrayDiffString(limitTypes, defaultLimits)) != 0 {
			return nil, utils.ErrInvalidLimitType
		} else {
			limitTmpMap, err := srv.GetLimit(bLDs, limitTypes)
			if err != nil {
				return nil, err
			}
			for k, v := range limitTmpMap {
				if v == "" {
					limitTmpMap[k] = agentlimitsMap[k]
				}
			}
			for k, v := range limitTmpMap {
				insinfo[k] = v
			}
		}

	}

	// Set chips
	chipsStr := ""
	if avResp.Agent.Currency == 17 || avResp.Agent.Currency == 18 {
		chipsStr = ChipsMap[avResp.Agent.Currency]
	} else if req.Chips != "" {
		// Validate custom chips
		chipsStr, err = validateChips(req.Chips)
		if err != nil {
			return nil, err
		}
	}
	insinfo["chips"] = chipsStr

	kickperiod := 0
	if req.Kickperiod == "" {
		kickperiod = avResp.Agent.Kickperiod
	} else {
		reqKickperiod, err := strconv.Atoi(req.Kickperiod)
		if err != nil {
			return nil, err
		}
		if reqKickperiod >= 0 {
			if (reqKickperiod <= avResp.Agent.Kickperiod && reqKickperiod != 0) || avResp.Agent.Kickperiod == 0 {
				kickperiod = reqKickperiod
			} else {
				return nil, utils.ErrInvalidKickPeriodGreaterThanUpper
			}
		} else {
			return nil, utils.ErrInvalidKickPeriodNegative
		}
	}
	insinfo["kickperiod"] = kickperiod

	insinfo["account"] = req.User
	insinfo["password"] = req.Password
	insinfo["name"] = req.Username

	insinfo["remark"] = req.Mark
	insinfo["maxwin"] = req.Maxwin
	insinfo["maxlose"] = req.Maxlose
	insinfo["head"] = req.Profile
	insinfo["lv"] = 6
	insinfo["uid"] = avResp.Agent.ID
	insinfo["ulv"] = avResp.Agent.ULV
	insinfo["kind"] = "a"
	insinfo["type"] = 0
	insinfo["tel"] = " "
	insinfo["currency"] = avResp.Agent.Currency
	insinfo["set101_1"] = 0
	insinfo["set102_1"] = 0
	insinfo["set103_1"] = 0
	insinfo["set104_1"] = 0
	insinfo["set105_1"] = 0
	insinfo["set106_1"] = 0
	insinfo["set107_1"] = 0
	insinfo["set108_1"] = 0
	insinfo["set110_1"] = 0
	insinfo["set111_1"] = 0
	insinfo["set112_1"] = 0
	insinfo["set113_1"] = 0
	insinfo["set117_1"] = 0
	insinfo["set121_1"] = 0
	insinfo["set126_1"] = 0
	insinfo["opengame"] = "101,102,103,104,105,106,107,108,110,111,112,113,117,121,126"

	//舊版註解
	insinfo["set301_1"] = 0
	insinfo["set301_2"] = 0
	insinfo["set301_9"] = 7
	insinfo["set109_1"] = 0
	insinfo["set109_9_0"] = 1
	insinfo["set109_9_1"] = 2
	insinfo["set109_9_2"] = 4
	insinfo["set109_1"] = 0
	insinfo["set109_9_0"] = 1
	insinfo["set109_9_1"] = 2
	insinfo["set109_9_2"] = 4
	insinfo["set109_9_3"] = 0
	insinfo["set109_9_4"] = 0
	insinfo["set109_9_5"] = 0
	insinfo["set109_9_6"] = 0
	insinfo["set109_9_7"] = 0
	insinfo["set109_9_8"] = 0
	insinfo["set109_9_9"] = 0
	insinfo["set109_14"] = 0
	insinfo["shortCode"] = 0

	// set 301_4
	for _, agentLimit := range agentLimits {
		if agentLimit.Ag002 == 301 {
			insinfo["set301_4"] = agentLimit.Ag015
		}
	}

	userDtl, err := srv.GetOneUserDtl(avResp.Agent.ID, avResp.Agent.ULV)
	if err != nil {
		return nil, err
	}
	insinfo["tip"] = avResp.Agent.Tip
	for _, v := range userDtl {
		typeStr := v["type"]
		if req.Rakeback == 0 {
			insinfo["set"+typeStr+"_2"] = v["bkwater"]
		} else {
			insinfo["set"+typeStr+"_2"] = 0
		}
		insinfo["set"+typeStr+"_9"] = v["betlimit"]
	}
	insinfo["set_109_9"] = 7

	userBase, err := srv.GetOneUserBase("agent", avResp.Agent.ID)
	if err != nil {
		return nil, err
	}
	insinfo["info"] = userBase

	gameTypes, err := srv.gameTypeDao.QueryByStatus(srv.DB(), 1)
	if err != nil {
		return nil, err
	}

	member := db.Member{
		User:     insinfo["account"].(string),
		Password: insinfo["password"].(string),
		UserName: insinfo["name"].(string),
		Opengame: insinfo["opengame"].(string),
		Mem005:   time.Now(),
		Mem006:   7,
	}
	if avResp.Agent.ULV == 1 {
		member.Mem007 = avResp.Agent.ID
	} else {
		member.Mem007 = avResp.Agent.Age007
	}

	if avResp.Agent.ULV == 2 {
		member.Mem008 = avResp.Agent.ID
	} else {
		member.Mem008 = avResp.Agent.Age008
	}

	if avResp.Agent.ULV == 3 {
		member.Mem009 = avResp.Agent.ID
	} else {
		member.Mem009 = avResp.Agent.Age009
	}

	if avResp.Agent.ULV == 4 {
		member.Mem010 = avResp.Agent.ID
	} else {
		member.Mem010 = avResp.Agent.Age010
	}

	if avResp.Agent.ULV == 5 {
		member.Mem011 = avResp.Agent.ID
	} else {
		member.Mem011 = avResp.Agent.Age011
	}

	if avResp.Agent.ULV == 5 {
		if avResp.Agent.Tip != "" {
			member.Tip = avResp.Agent.Tip
		} else {
			if tmpTip, ok := insinfo["tip"]; ok {
				member.Tip = tmpTip.(string)
			}
		}
	} else {
		if tmpTip, ok := insinfo["tip"]; ok {
			member.Tip = tmpTip.(string)
		}
	}

	member.Mem012 = 0
	member.Mem013 = time.Time{}
	member.Mem016 = "Y"
	member.Mem017 = "Y"
	member.Mem018 = "N"
	member.Mem019 = "N"
	member.Mem022 = insinfo["tel"].(string)
	member.Kickperiod = insinfo["kickperiod"].(int)
	member.Identity = 1
	if shortCode, ok := insinfo["shortCode"]; ok {
		member.Mem022a = shortCode.(int)
	}
	if head, ok := insinfo["head"]; ok {
		member.Head = head.(int)
	}
	member.Mem028 = insinfo["remark"].(string)
	if wallet, ok := insinfo["wallet"]; ok {
		member.Wallet = wallet.(string)
	} else {
		member.Wallet = ""
	}
	if avResp.Agent.Type != 3 {
		member.Type = avResp.Agent.Type
	} else {
		member.Type = insinfo["type"].(int)
	}
	member.Currency = avResp.Agent.Currency
	if chips, ok := insinfo["chips"]; ok {
		member.Chips = chips.(string)
	}

	// Insert member
	var newMemID int64
	err = srv.Tx(func(tx *gorm.DB) error {
		var err error
		newMemID, err = srv.userDao.CreateUser(tx, &member)
		if err != nil {
			err = errors.New("新增会员资料错误")
			return err
		}
		xlog.Debugf("newMemID: %d", newMemID)

		var memberDtls []*db.MemberDtl
		for _, gameType := range gameTypes {
			memberDtl := &db.MemberDtl{
				Mem001: newMemID,
				Mem002: gameType.Code,
			}
			if tmpItem, ok := insinfo["set"+strconv.Itoa(int(gameType.Code))+"_2"]; ok {
				memberDtl.Mem003, err = decimal.NewFromString(tmpItem.(string))
				if err != nil {
					err = errors.New("退水信息错误")
					return err
				}
			}
			memberDtl.Mem004 = member.Mem007
			memberDtl.Mem005 = member.Mem008
			memberDtl.Mem006 = member.Mem009
			memberDtl.Mem007 = member.Mem010
			memberDtl.Mem008 = member.Mem011
			memberDtl.Mem013 = time.Now()
			if gameType.Code != 109 && gameType.Code != 301 {
				if tmpItem, ok := insinfo["set"+strconv.Itoa(int(gameType.Code))+"_14"]; ok {
					memberDtl.Mem015 = tmpItem.(string)
				}
			}
			if gameType.Code == 301 {
				if tmpItem, ok := insinfo["set"+strconv.Itoa(int(gameType.Code))+"_4"]; ok {
					xlog.Debugf("tmpItem: %v", tmpItem)
					if v, ok := tmpItem.(int); ok {
						memberDtl.Mem016 = decimal.NewFromInt(int64(v))
					} else if v, ok := tmpItem.(string); ok {
						memberDtl.Mem016, err = decimal.NewFromString(v)
						if err != nil {
							err = errors.New("电投退水信息错误")
							return err
						}
					} else if v, ok := tmpItem.(decimal.Decimal); ok {
						memberDtl.Mem016 = v
					} else {
						err = errors.New("电投退水信息错误")
						return err
					}
				}
			}
			if tmpItem, ok := insinfo["maxwin"]; ok {
				memberDtl.Mem012 = tmpItem.(int64)
			}
			if tmpItem, ok := insinfo["maxlose"]; ok {
				memberDtl.Mem014 = tmpItem.(int64)
			}
			memberDtls = append(memberDtls, memberDtl)
		}
		err = srv.memDtlDao.CreateMemberDtls(tx, memberDtls)
		if err != nil {
			err = errors.New("新增会员細項资料错误")
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &view.MemberRegisterResp{
		Result: "操作成功",
	}, nil
}

// GetLimit groups bet limits by game type and returns concatenated limit IDs
func (srv *publicApiService) GetLimit(bLDs []*db.BetLimitDefault, limitTypes []string) (map[string]string, error) {
	// Initialize limit arrays for each game type
	set101_14 := make([]string, 0)
	set102_14 := make([]string, 0)
	set103_14 := make([]string, 0)
	set104_14 := make([]string, 0)
	set105_14 := make([]string, 0)
	set106_14 := make([]string, 0)
	set107_14 := make([]string, 0)
	set108_14 := make([]string, 0)
	set110_14 := make([]string, 0)
	set111_14 := make([]string, 0)
	set112_14 := make([]string, 0)
	set113_14 := make([]string, 0)
	set117_14 := make([]string, 0)
	set121_14 := make([]string, 0)
	set126_14 := make([]string, 0)

	// Get all bet limit defaults
	for _, limit := range bLDs {
		// Check if this limit ID is in the requested limitType slice
		limitID := strconv.Itoa(int(limit.ID))
		if !slices.Contains(limitTypes, limitID) {
			continue
		}

		// Add limit ID to appropriate game type array
		switch limit.Gtype {
		case 101:
			set101_14 = append(set101_14, limitID)
		case 102:
			set102_14 = append(set102_14, limitID)
		case 103:
			set103_14 = append(set103_14, limitID)
		case 104:
			set104_14 = append(set104_14, limitID)
		case 105:
			set105_14 = append(set105_14, limitID)
		case 106:
			set106_14 = append(set106_14, limitID)
		case 107:
			set107_14 = append(set107_14, limitID)
		case 108:
			set108_14 = append(set108_14, limitID)
		case 110:
			set110_14 = append(set110_14, limitID)
		case 111:
			set111_14 = append(set111_14, limitID)
		case 112:
			set112_14 = append(set112_14, limitID)
		case 113:
			set113_14 = append(set113_14, limitID)
		case 121:
			set121_14 = append(set121_14, limitID)
		case 126:
			set126_14 = append(set126_14, limitID)
		}
	}

	// Create result map with concatenated limit IDs
	result := map[string]string{
		"set101_14": strings.Join(set101_14, ","),
		"set102_14": strings.Join(set102_14, ","),
		"set103_14": strings.Join(set103_14, ","),
		"set104_14": strings.Join(set104_14, ","),
		"set105_14": strings.Join(set105_14, ","),
		"set106_14": strings.Join(set106_14, ","),
		"set107_14": strings.Join(set107_14, ","),
		"set108_14": strings.Join(set108_14, ","),
		"set110_14": strings.Join(set110_14, ","),
		"set111_14": strings.Join(set111_14, ","),
		"set112_14": strings.Join(set112_14, ","),
		"set113_14": strings.Join(set113_14, ","),
		"set117_14": strings.Join(set117_14, ","),
		"set121_14": strings.Join(set121_14, ","),
		"set126_14": strings.Join(set126_14, ","),
	}

	return result, nil
}

func (srv *publicApiService) GetOneUserDtl(id int64, lv int) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)

	// Query based on level
	if lv != 7 {
		// Query agent_dtl for non-member levels
		var agentDtls []*db.AgentDtl
		err := srv.Tx(func(tx *gorm.DB) error {
			var err error
			agentDtls, err = srv.agentDtlDao.QueryByAgentID(tx, id)
			return err
		})
		if err != nil {
			return nil, err
		}

		// Process agent details
		for _, dtl := range agentDtls {
			typeStr := strconv.FormatInt(dtl.Ag002, 10)
			result[typeStr] = map[string]string{
				"type":      typeStr,
				"bkwater":   dtl.Ag003.String(),
				"rate":      dtl.Ag012.String(),
				"betlimit":  strconv.Itoa(int(dtl.Ag013)),
				"nbetlimit": dtl.Ag014,
				"netwater":  dtl.Ag015.String(),
				"netrate":   strconv.Itoa(int(dtl.Ag016)),
			}
			// Add limitary array
			if dtl.Ag014 != "" {
				result[typeStr]["limitary"] = dtl.Ag014
			}
		}
	} else {
		// Query member_dtl for member level
		var memberDtls []*db.MemberDtl
		err := srv.Tx(func(tx *gorm.DB) error {
			var err error
			memberDtls, err = srv.memDtlDao.QueryByMemberID(tx, id)
			return err
		})
		if err != nil {
			return nil, err
		}

		// Process member details
		for _, dtl := range memberDtls {
			typeStr := strconv.FormatInt(dtl.Mem002, 10)
			result[typeStr] = map[string]string{
				"type":      typeStr,
				"bkwater":   dtl.Mem003.String(),
				"rate":      "0",
				"betlimit":  strconv.FormatInt(dtl.Mem011, 10),
				"maxwin":    strconv.FormatInt(dtl.Mem012, 10),
				"maxlose":   strconv.FormatInt(dtl.Mem014, 10),
				"nbetlimit": dtl.Mem015,
				"netwater":  dtl.Mem016.String(),
			}
			// Add limitary array
			if dtl.Mem015 != "" {
				result[typeStr]["limitary"] = dtl.Mem015
			}
		}
	}

	// Add default type 301 if not exists
	if _, exists := result["301"]; !exists {
		result["301"] = map[string]string{
			"type":        "301",
			"bkwater":     "0",
			"rate":        "0",
			"betlimit":    "1023", // As per PHP code
			"maxwin":      "0",
			"maxlose":     "0",
			"nbetlimit":   "0",
			"betlimitary": "1,1,1,1,1,1,1,1,1,1",
		}
	}

	// Add default type 105 if not exists
	if _, exists := result["105"]; !exists {
		// Get bet limit defaults for type 105
		var betLimits []*db.BetLimitDefault
		err := srv.Tx(func(tx *gorm.DB) error {
			var err error
			betLimits, err = srv.betLimitDefaultDao.QueryByGtype(tx, 105)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		// Build nbetlimit string
		var nbetlimitParts []string
		for _, limit := range betLimits {
			nbetlimitParts = append(nbetlimitParts, strconv.FormatInt(limit.ID, 10))
		}
		nbetlimit := strings.Join(nbetlimitParts, ",")

		result["105"] = map[string]string{
			"type":      "105",
			"bkwater":   "2.00",
			"rate":      "90",
			"betlimit":  "0",
			"maxwin":    "0",
			"maxlose":   "0",
			"netwater":  "0.00",
			"netrate":   "0",
			"nbetlimit": nbetlimit,
		}
	}

	return result, nil
}

func (srv *publicApiService) GetOneUserBase(table string, id int64) (map[string]interface{}, error) {
	var query string
	switch table {
	case "agent":
		query = `SELECT 
			age001 as id,
			age002 as account, 
			age003 as password, 
			age004 as name, 
			age015 as useflag, 
			age016 as betflag,
			age024 as remark,
			type,
			currency,
			prefix_add as prefixadd,
			tip
		FROM agent WHERE age001 = ?`
	case "member":
		query = `SELECT 
			mem001 as id,
			mem002 as account, 
			mem003 as password, 
			mem004 as name, 
			mem016 as useflag, 
			mem017 as betflag,
			mem022 as tel,
			mem022a as shortCode,
			mem028 as remark,
			type,
			head,
			currency,
			tip
		FROM member WHERE mem001 = ?`
	default:
		return nil, fmt.Errorf("invalid table name: %s", table)
	}

	var result map[string]interface{}
	err := srv.DB().Raw(query, id).Scan(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying user base info: %w", err)
	}

	return result, nil
}

func (srv *publicApiService) EditLimit(ctx context.Context, req *view.EditLimitReq) (*view.EditLimitResp, error) {
	if err := utils.CheckTimestamp(req.Timestamp); err != nil {
		xlog.Errorf("error to check timestamp, err:%+v", err)
		return nil, err
	}

	// Verify agent
	avResp, err := srv.AgentVerify(ctx, &view.AgentVerifyReq{VendorID: req.VendorID, Signature: req.Signature})
	if err != nil {
		xlog.Errorf("error to verify agent, err:%+v", err)
		return nil, err
	}

	// Handle maxwin/maxlose update without account
	if req.User == "" && (req.Maxwin != -1 || req.Maxlose != -1) {
		err = srv.Tx(func(tx *gorm.DB) error {
			updates := make(map[string]interface{})
			if req.Maxwin != -1 {
				updates["mem012"] = req.Maxwin
			}
			if req.Maxlose != -1 {
				updates["mem014"] = req.Maxlose
			}
			if len(updates) > 0 {
				err := srv.memDtlDao.UpdatesByAgentID(tx, avResp.Agent.ID, updates)
				if err != nil {
					xlog.Errorf("error to update by agent, err:%+v", err)
					return err
				}
			}
			return nil
		})
		if err != nil {
			xlog.Errorf("error to tx operation, err:%+v", err)
			return nil, err
		}
		return &view.EditLimitResp{Result: "操作成功"}, nil
	}

	// Validate account
	if req.User == "" {
		xlog.Errorf("error to validate account, err:%+v", utils.ErrInvalidAccountEmpty)
		return nil, utils.ErrInvalidAccountEmpty
	}

	// Get member info
	member, err := srv.userDao.QueryByAccount(srv.DB(), req.User)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			xlog.Errorf("error to query member by account, err:%+v", utils.ErrParamInvalidAccountNotExist)
			return nil, utils.ErrParamInvalidAccountNotExist
		}
		return nil, err
	}

	// Verify member belongs to agent
	if member.Mem011 != avResp.Agent.ID {
		xlog.Errorf("error to verify member belongs to agent, err:%+v", utils.ErrParamInvalidAccountNotBelongToAgent)
		return nil, utils.ErrParamInvalidAccountNotBelongToAgent
	}

	insinfo := make(map[string]interface{})
	insinfo["account"] = req.User

	// Handle limit type
	if req.LimitType != "" {
		req.LimitType = strings.ReplaceAll(req.LimitType, " ", "")
		limitTypes := strings.Split(req.LimitType, ",")

		// Get parent agent's limits
		agentLimits, err := srv.agentDtlDao.QueryByAgentID(srv.DB(), avResp.Agent.ID)
		if err != nil {
			xlog.Errorf("error to query agent dtl, err:%+v", err)
			return nil, err
		}

		// Validate limit types against parent agent's limits
		validLimits := make([]string, 0)
		for _, agentLimit := range agentLimits {
			tmpAgentLimit := strings.ReplaceAll(agentLimit.Ag014, " ", "")
			if tmpAgentLimit != "" {
				parentLimits := strings.Split(tmpAgentLimit, ",")
				for _, requestedLimit := range limitTypes {
					for _, parentLimit := range parentLimits {
						if requestedLimit == parentLimit {
							validLimits = append(validLimits, requestedLimit)
						}
					}
				}
			}
		}
		xlog.Infof("validLimits: %+v", validLimits)

		if len(utils.ArrayDiffString(limitTypes, validLimits)) != 0 {
			xlog.Errorf("error to validate limit types, err:%+v", utils.ErrInvalidLimitTypeNotOpen)
			return nil, utils.ErrInvalidLimitTypeNotOpen
		}

		// Get default bet limits
		bLDs, err := srv.betLimitDefaultDao.QueryAll(srv.DB())
		if err != nil {
			xlog.Errorf("error to get default bet limits, err:%+v", err)
			return nil, err
		}

		defaultLimits := make([]string, 0)
		for _, dl := range bLDs {
			defaultLimits = append(defaultLimits, strconv.Itoa(int(dl.ID)))
		}

		if len(utils.ArrayDiffString(limitTypes, defaultLimits)) != 0 {
			xlog.Errorf("error to validate limit types, err:%+v", utils.ErrInvalidLimitType)
			return nil, utils.ErrInvalidLimitType
		}

		// Get limit map
		limitMap, err := srv.GetLimit(bLDs, limitTypes)
		if err != nil {
			xlog.Errorf("error to get limit map, err:%+v", err)
			return nil, err
		}

		for k, v := range limitMap {
			insinfo[k] = v
		}
	}

	// Handle maxwin/maxlose
	if req.Maxwin != -1 {
		insinfo["maxwin"] = req.Maxwin
	}
	if req.Maxlose != -1 {
		insinfo["maxlose"] = req.Maxlose
	}

	// Handle reset
	if req.Reset == 1 {
		insinfo["reset"] = time.Now()
	}

	// Update member details
	err = srv.Tx(func(tx *gorm.DB) error {
		memberDtls, err := srv.memDtlDao.QueryByMemberID(tx, member.ID)
		if err != nil {
			xlog.Errorf("error to query member dtl, err:%+v", err)
			return err
		}

		for _, dtl := range memberDtls {
			updates := make(map[string]interface{})

			if maxwin, ok := insinfo["maxwin"]; ok {
				updates["mem012"] = maxwin
			}
			if maxlose, ok := insinfo["maxlose"]; ok {
				updates["mem014"] = maxlose
			}
			if reset, ok := insinfo["reset"]; ok {
				updates["mem013"] = reset
			}
			if limitKey := fmt.Sprintf("set%d_14", dtl.Mem002); insinfo[limitKey] != nil {
				updates["mem015"] = insinfo[limitKey]
			}

			if len(updates) > 0 {
				err = srv.memDtlDao.UpdateMemberDtlByMemberID(tx, dtl, updates)
				if err != nil {
					xlog.Errorf("error to update member dtl, err:%+v", err)
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		xlog.Errorf("error to tx operation, err:%+v", err)
		return nil, err
	}

	return &view.EditLimitResp{Result: "操作成功"}, nil
}

func (srv *publicApiService) LogoutGame(ctx context.Context, req *view.LogoutGameReq) (*view.LogoutGameResp, error) {
	if err := utils.CheckTimestamp(req.Timestamp); err != nil {
		xlog.Errorf("error to check timestamp, err:%+v", err)
		return nil, err
	}

	// Verify agent
	avResp, err := srv.AgentVerify(ctx, &view.AgentVerifyReq{VendorID: req.VendorID, Signature: req.Signature})
	if err != nil {
		xlog.Errorf("error to verify agent, err:%+v", err)
		return nil, err
	}

	var memberIDs []int64
	if req.User != "" {
		// Get member by account
		member, err := srv.userDao.QueryByAccount(srv.DB(), req.User)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				xlog.Errorf("error to query member by account, err:%+v", utils.ErrParamInvalidAccountNotExist)
				return nil, utils.ErrParamInvalidAccountNotExist
			}
			return nil, err
		}

		// Verify member belongs to agent
		if member.Mem011 != avResp.Agent.ID {
			xlog.Errorf("error to verify member belongs to agent, err:%+v", utils.ErrParamInvalidAccountNotBelongToAgent)
			return nil, utils.ErrParamInvalidAccountNotBelongToAgent
		}

		memberIDs = append(memberIDs, member.ID)
	} else {
		// Get all members under this agent
		members, err := srv.userDao.QueryByAgentID(srv.DB(), avResp.Agent.ID)
		if err != nil {
			return nil, err
		}
		if len(members) == 0 {
			xlog.Errorf("error to query members by agent, err:%+v", utils.ErrParamInvalidAccountNotExist)
			return nil, utils.ErrParamInvalidAccountNotExist
		}

		for _, member := range members {
			memberIDs = append(memberIDs, member.ID)
		}
	}

	// Delete mem_login records for the members
	err = srv.Tx(func(tx *gorm.DB) error {
		if err := srv.memLoginDao.UpdateMemLoginByMemIDs(tx, memberIDs); err != nil {
			xlog.Errorf("error to delete mem_login, memberIDs:%+v, err:%+v", memberIDs, err)
			return err
		}
		return nil
	})
	if err != nil {
		xlog.Errorf("error to tx operation, err:%+v", err)
		return nil, err
	}

	return &view.LogoutGameResp{
		Result: "操作成功",
	}, nil
}

func (srv *publicApiService) ChangePassword(ctx context.Context, req *view.ChangePasswordReq) (*view.ChangePasswordResp, error) {
	if err := utils.CheckTimestamp(req.Timestamp); err != nil {
		xlog.Errorf("error to check timestamp, err:%+v", err)
		return nil, err
	}

	// Verify agent
	avResp, err := srv.AgentVerify(ctx, &view.AgentVerifyReq{VendorID: req.VendorID, Signature: req.Signature})
	if err != nil {
		xlog.Errorf("error to verify agent, err:%+v", err)
		return nil, err
	}

	// Validate input
	if req.User == "" {
		return nil, utils.ErrInvalidAccountEmpty
	}
	if req.NewPassword == "" {
		return nil, utils.ErrInvalidPasswordEmpty
	}

	// Check password format
	matched, _ := regexp.MatchString(`^[\x{4e00}-\x{9fa5}]{1,30}$`, req.NewPassword)
	if matched {
		return nil, utils.ErrInvalidPasswordChinese
	}
	if len(req.NewPassword) <= 5 {
		return nil, utils.ErrInvalidPasswordLengthShort
	}
	if len(req.NewPassword) >= 65 {
		return nil, utils.ErrInvalidPasswordLengthLong
	}

	// Get member info
	member, err := srv.userDao.QueryByAccount(srv.DB(), req.User)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrParamInvalidAccountNotExist
		}
		return nil, err
	}

	// Verify member belongs to agent
	if member.Mem011 != avResp.Agent.ID {
		xlog.Errorf("error to verify member belongs to agent, err:%+v", utils.ErrParamInvalidAccountNotBelongToAgent)
		return nil, utils.ErrParamInvalidAccountNotBelongToAgent
	}

	// Check if new password is same as old password
	if member.Password == req.NewPassword {
		return nil, utils.ErrParamInvalidPasswordSame
	}

	// Update password in database
	err = srv.Tx(func(tx *gorm.DB) error {
		err := srv.userDao.UpdatesMember(tx, member.ID, map[string]interface{}{
			"mem003": req.NewPassword,
		})
		if err != nil {
			xlog.Errorf("Failed to update member password, err:%+v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Get language-specific response
	var result string
	switch req.Syslang {
	case 0:
		result = fmt.Sprintf("密码:%s,修改完成", req.NewPassword)
	case 1:
		result = fmt.Sprintf("Password:%s,Change completed", req.NewPassword)
	default:
		result = fmt.Sprintf("密码:%s,修改完成", req.NewPassword)
	}

	return &view.ChangePasswordResp{
		Result: result,
	}, nil
}

func (srv *publicApiService) GetAgentBalance(ctx context.Context, req *view.GetAgentBalanceReq) (*view.GetAgentBalanceResp, error) {
	// Validate timestamp
	if err := utils.CheckTimestamp(req.Timestamp); err != nil {
		xlog.Errorf("error to check timestamp, err:%+v", err)
		return nil, err
	}

	// Verify agent
	avResp, err := srv.AgentVerify(ctx, &view.AgentVerifyReq{VendorID: req.VendorID, Signature: req.Signature})
	if err != nil {
		xlog.Errorf("error to verify agent, err:%+v", err)
		return nil, err
	}

	// Get agent's balance
	agent, err := srv.agentDao.QueryByID(srv.DB(), avResp.Agent.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrParamInvalidAgentNotExist
		}
		xlog.Errorf("error to get agent balance, err:%+v", err)
		return nil, err
	}

	return &view.GetAgentBalanceResp{
		Result: agent.Cash,
	}, nil
}

// GetBalance gets the balance for a member
func (srv *publicApiService) GetBalance(ctx context.Context, req *view.GetBalanceReq) (*view.GetBalanceResp, error) {
	// Validate timestamp
	if err := utils.CheckTimestamp(req.Timestamp); err != nil {
		xlog.Errorf("error to check timestamp, err:%+v", err)
		return nil, err
	}

	// Verify agent
	avResp, err := srv.AgentVerify(ctx, &view.AgentVerifyReq{VendorID: req.VendorID, Signature: req.Signature})
	if err != nil {
		xlog.Errorf("error to verify agent, err:%+v", err)
		return nil, err
	}

	// Get member by account
	member, err := srv.userDao.QueryByAccount(srv.DB(), req.User)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			xlog.Error("no data")
			return nil, utils.ErrParamInvalidAccountNotExist
		}
		xlog.Errorf("error to get member by account, err:%+v", err)
		return nil, err
	}

	// Verify member belongs to agent
	if member.Mem011 != avResp.Agent.ID {
		xlog.Errorf("error to verify member belongs to agent, err:%+v", utils.ErrParamInvalidAccountNotBelongToAgent)
		return nil, utils.ErrParamInvalidAccountNotBelongToAgent
	}

	return &view.GetBalanceResp{
		Result: member.Cash,
	}, nil
}

func (srv *publicApiService) ChangeBalance(ctx context.Context, req *view.ChangeBalanceReq) (*view.ChangeBalanceResp, error) {
	// Validate timestamp
	if err := utils.CheckTimestamp(req.Timestamp); err != nil {
		xlog.Errorf("error to check timestamp, err:%+v", err)
		return nil, err
	}

	// Verify agent
	avResp, err := srv.AgentVerify(ctx, &view.AgentVerifyReq{VendorID: req.VendorID, Signature: req.Signature})
	if err != nil {
		xlog.Errorf("error to verify agent, err:%+v", err)
		return nil, err
	}

	// Validate money parameter
	if req.Money == "" {
		return nil, utils.ErrWalletAddOrSubPointEmptyOrMoneyParamNotSet
	}

	// Check for Chinese characters in money parameter
	matched, _ := regexp.MatchString(`[\p{Han}]`, req.Money)
	if matched {
		return nil, utils.ErrWalletAddOrSubPointChinese
	}

	// Parse money value
	money, err := decimal.NewFromString(req.Money)
	if err != nil || money.IsZero() {
		return nil, utils.ErrWalletAddOrSubPointEmptyOrMoneyParamNotSet
	}

	// Get member by account
	member, err := srv.userDao.QueryByAccount(srv.DB(), req.User)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrParamInvalidAccountNotExist
		}
		return nil, err
	}

	// Verify member belongs to agent
	if member.Mem011 != avResp.Agent.ID {
		return nil, utils.ErrParamInvalidAccountNotBelongToAgent
	}

	// Check if member is locked
	if member.Mem020 == "Y" {
		return nil, utils.ErrInvalidTransactionError
	}

	// Check transaction timing
	const ChangeBalance_chkTime = "ChangeBalance_chkTime"
	currentTime := time.Now().Unix()

	// Get last transaction time from Redis hash
	lastTxTime, err := srv.redisCli.HGet(ctx, ChangeBalance_chkTime, strconv.FormatInt(member.ID, 10)).Int64()
	if err == redis.Nil {
		// No previous transaction time, set current time
		err = srv.redisCli.HSet(ctx, ChangeBalance_chkTime, strconv.FormatInt(member.ID, 10), currentTime).Err()
		if err != nil {
			xlog.Errorf("error setting transaction time in Redis: %v", err)
			return nil, utils.ErrWalletTransferLineError
		}
	} else if err != nil {
		xlog.Errorf("error getting transaction time from Redis: %v", err)
		return nil, utils.ErrWalletTransferLineError
	} else {
		// Calculate time difference
		timeDiff := currentTime - lastTxTime

		// Special handling for specific vendors/agents
		if req.VendorID == "igktwapi" || req.VendorID == "ocmsapi" || avResp.Agent.ID == 1717 {
			if timeDiff < 2 {
				return nil, utils.ErrWalletTransferRepeatIn2Error
			}
		} else {
			if timeDiff < 5 {
				return nil, utils.ErrWalletTransferRepeatIn5Error
			}
		}

		// Update transaction time
		err = srv.redisCli.HSet(ctx, ChangeBalance_chkTime, strconv.FormatInt(member.ID, 10), currentTime).Err()
		if err != nil {
			xlog.Errorf("error updating transaction time in Redis: %v", err)
			return nil, utils.ErrWalletTransferLineError
		}
	}

	// Check launder
	launder, err := srv.CheckLaunder(srv.DB(), avResp, member)
	if err != nil {
		xlog.Errorf("error to check launder, err:%+v", err)
		return nil, err
	}
	if launder {
		xlog.Errorf("error to check launder, err:%+v", utils.ErrWalletTransferLockError)
		return nil, utils.ErrWalletTransferLockError
	}

	// Check 5 seconds duplicate transactions (same amount)
	res, err := srv.CheckDoubleDeals(ctx, member.ID, money, req.Order)
	if err != nil {
		xlog.Errorf("error to check double deals, err:%+v", err)
		return nil, err
	}
	if res == 1 {
		xlog.Errorf("error to check double deals, err:%+v", utils.ErrWalletTransferRepeatIn5Error)
		return nil, utils.ErrWalletTransferRepeatIn5Error
	}
	if res == 3 {
		xlog.Errorf("error to check double deals, err:%+v", utils.ErrWalletTransferExist)
		return nil, utils.ErrWalletTransferExist
	}

	if money.IsPositive() {
		err := srv.ProDealAddValue(ctx, money, avResp, member, req.Order)
		if err != nil {
			xlog.Errorf("error to pro deal add value, err:%+v", err)
			return nil, err
		}
		xlog.Infof("pro deal add value result: %+v", res)
	}
	if money.IsNegative() {
		err := srv.ProDealDecValue(ctx, money, avResp, member, req.Order)
		if err != nil {
			xlog.Errorf("error to pro deal dec value, err:%+v", err)
			return nil, err
		}
		xlog.Infof("pro deal dec value result: %+v", res)
	}

	// Get language-specific response
	var result string
	switch req.Syslang {
	case 1:
		result = fmt.Sprintf("Balance change completed, amount: %s", req.Money)
	default:
		result = fmt.Sprintf("余额变更完成，金额: %s", req.Money)
	}

	return &view.ChangeBalanceResp{
		Result: result,
	}, nil
}

// CheckLaunder checks if a member has too many transactions in the last minute
func (srv *publicApiService) CheckLaunder(tx *gorm.DB, avResp *view.AgentVerifyResp, member *db.Member) (bool, error) {
	// Get transactions in last minute
	count, err := srv.inOutMDao.CountTransactionsInLastMinute(tx, member.ID)
	if err != nil {
		xlog.Errorf("error to count transactions in last minute, err:%+v", err)
		return false, err
	}

	if count > 9 {
		// Create alert message
		lastTime := time.Now().Add(-1 * time.Minute)
		nowTime := time.Now()
		message := fmt.Sprintf("請注意! 會員帳號 : %s從%s~%s已有%d次,轉帳紀錄!",
			member.User,
			lastTime.Format("2006-01-02 15:04:05"),
			nowTime.Format("2006-01-02 15:04:05"),
			count,
		)
		message += fmt.Sprintf(" 請與代理商 : %s, skype群組%s  確認",
			avResp.Agent.Name,
			avResp.AgentsLoginPass.Skyname,
		)

		// Insert alert message
		alertMsg := &db.AlertMessage{
			Mid:          member.ID,
			Message:      message,
			Status:       0,
			ErrorTime:    nowTime,
			UnierrorTime: nowTime.Unix(),
			Operator:     "", // 设置默认操作者
		}
		_, err = srv.alertMessageDao.Create(tx, alertMsg)
		if err != nil {
			xlog.Errorf("error to insert alert message, err:%+v", err)
			return false, err
		}

		// Update member status
		err = srv.userDao.UpdatesMember(tx, member.ID, map[string]interface{}{
			"mem020": "Y",
		})
		if err != nil {
			xlog.Errorf("error to update member status, err:%+v", err)
			return false, err
		}

		return true, nil
	}

	return false, nil
}

// CheckDoubleDeals checks for duplicate transactions within a short time window
func (srv *publicApiService) CheckDoubleDeals(ctx context.Context, memberID int64, amount decimal.Decimal, orderNum string) (int, error) {
	inOutM, err := srv.inOutMDao.GetLastTransaction(srv.DB(), memberID, orderNum)
	if err != nil {
		xlog.Errorf("error to get last transaction, err:%+v", err)
		return 0, err
	}
	if inOutM == nil {
		return 0, nil
	}

	// Check if order numbers match when orderNum is provided
	if orderNum != "" && inOutM.Iom008 == orderNum {
		return 3, nil
	}

	// Get current time
	now := time.Now()

	// Check if the transaction is within 5 seconds and has the same amount
	if now.Sub(inOutM.Iom002).Seconds() < 5 && amount.Equal(inOutM.Iom004) {
		return 1, nil
	}

	return 0, nil
}

func (srv *publicApiService) GetUpLv5(ctx context.Context, member *db.Member) int64 {
	if member.Mem011 != 0 {
		return member.Mem011
	}
	if member.Mem010 != 0 {
		return member.Mem010
	}
	if member.Mem009 != 0 {
		return member.Mem009
	}
	if member.Mem008 != 0 {
		return member.Mem008
	}
	return member.Mem007
}

func (srv *publicApiService) ProDealAddValue(ctx context.Context, money decimal.Decimal, avResp *view.AgentVerifyResp, member *db.Member, orderNum string) error {
	upid := srv.GetUpLv5(ctx, member)
	lv5Before := avResp.Agent.Cash
	if member.Type == 1 {
		lv5Before = avResp.Agent.Credit
	}
	pointtype := 0
	if member.Type == 1 {
		pointtype = 1
	}
	if lv5Before.Sub(money).IsNegative() {
		return utils.ErrWalletAgentOverLimit
	}
	var err error
	err = srv.Tx(func(tx *gorm.DB) error {
		logAgeCashChange := &db.LogAgeCashChange{
			Lacc02:    member.Mem006,
			Lacc03:    member.ID,
			Lacc06:    money,
			Lacc07:    time.Now(),
			Lacc08:    "",
			Lacc09:    member.Cash,
			Lacc10:    upid,
			Lacc11:    lv5Before,
			Pointtype: pointtype,
		}
		_, err = srv.logAgeCashChangeDao.Create(tx, logAgeCashChange)
		if err != nil {
			xlog.Errorf("error to create log age cash change, err:%+v", err)
			return err
		}
		if err = srv.userDao.UpdatesMember(tx, member.ID, map[string]interface{}{
			"cash": gorm.Expr("cash + ?", money),
		}); err != nil {
			xlog.Errorf("error to update member cash, err:%+v", err)
			return err
		}
		if member.Type != 1 {
			if err = srv.agentDao.UpdatesAgent(tx, avResp.Agent.ID, map[string]interface{}{
				"cash": gorm.Expr("cash - ?", money),
			}); err != nil {
				xlog.Errorf("error to update agent cash, err:%+v", err)
				return err
			}
		}
		_, err = srv.inOutMDao.DealInsRecord(tx, "121", 0, int64(avResp.Agent.ULV), avResp.Agent.ID, member.ID, money, orderNum, member.Cash)
		if err != nil {
			xlog.Errorf("error to deal ins record, err:%+v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return utils.ErrWalletTransferLineError
	}

	return nil
}

func (srv *publicApiService) ProDealDecValue(ctx context.Context, money decimal.Decimal, avResp *view.AgentVerifyResp, member *db.Member, orderNum string) error {
	upid := srv.GetUpLv5(ctx, member)
	lv5Before := avResp.Agent.Cash
	if member.Type == 1 {
		lv5Before = avResp.Agent.Credit
	}
	pointtype := 0
	if member.Type == 1 {
		pointtype = 1
	}
	if member.Cash.Add(money).IsNegative() {
		return utils.ErrWalletTransferBalanceNotEnough
	}
	var err error
	err = srv.Tx(func(tx *gorm.DB) error {
		logAgeCashChange := &db.LogAgeCashChange{
			Lacc02:    member.Mem006,
			Lacc03:    member.ID,
			Lacc06:    money,
			Lacc07:    time.Now(),
			Lacc08:    "",
			Lacc09:    member.Cash,
			Lacc10:    upid,
			Lacc11:    lv5Before,
			Pointtype: pointtype,
		}
		_, err = srv.logAgeCashChangeDao.Create(tx, logAgeCashChange)
		if err != nil {
			xlog.Errorf("error to create log age cash change, err:%+v", err)
			return err
		}
		if err = srv.userDao.UpdatesMember(tx, member.ID, map[string]interface{}{
			"cash": gorm.Expr("cash + ?", money),
		}); err != nil {
			xlog.Errorf("error to update member cash, err:%+v", err)
			return err
		}
		if member.Type != 1 {
			if err = srv.agentDao.UpdatesAgent(tx, avResp.Agent.ID, map[string]interface{}{
				"cash": gorm.Expr("cash - ?", money),
			}); err != nil {
				xlog.Errorf("error to update agent cash, err:%+v", err)
				return err
			}
		}
		_, err = srv.inOutMDao.DealInsRecord(tx, "122", 0, int64(avResp.Agent.ULV), avResp.Agent.ID, member.ID, money, orderNum, member.Cash)
		if err != nil {
			xlog.Errorf("error to deal ins record, err:%+v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return utils.ErrWalletTransferLineError
	}

	return nil
}

func (srv *publicApiService) GetMemberTradeReport(ctx context.Context, req *view.GetMemberTradeReportReq) (*view.GetMemberTradeReportResp, error) {
	// Validate timestamp
	if err := utils.CheckTimestamp(req.Timestamp); err != nil {
		xlog.Errorf("error to check timestamp, err:%+v", err)
		return nil, err
	}

	// Verify agent
	avgResp, err := srv.AgentVerify(ctx, &view.AgentVerifyReq{VendorID: req.VendorID, Signature: req.Signature})
	if err != nil {
		xlog.Errorf("error to verify agent, err:%+v", err)
		return nil, err
	}

	if req.User != "" {
		// Get member by account
		member, err := srv.userDao.QueryByAccount(srv.DB(), req.User)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, utils.ErrParamInvalidAccountNotExist
			}
			return nil, err
		}
		mIDs := []int64{member.ID}
		var result []*db.InOutM
		err = srv.Tx(func(tx *gorm.DB) error {
			result, err = srv.inOutMDao.GetInOutMs(tx, mIDs, req.OrderID, req.Order, req.StartTime, req.EndTime)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		var tradeItems []*view.TradeItem
		for _, v := range result {
			tradeItem := &view.TradeItem{
				MID:      v.Iom003,
				OrderID:  v.Iom001,
				OrderNum: v.Iom008,
				AddTime:  v.Iom002.Unix(),
				Money:    v.Iom004,
				OpCode:   v.Iom005,
				Subtotal: v.Iom010,
			}
			tradeItems = append(tradeItems, tradeItem)
		}

		return &view.GetMemberTradeReportResp{
			Result: tradeItems,
		}, nil
	} else {
		if req.StartTime != 0 && req.EndTime != 0 {
			if req.EndTime-req.StartTime > 86400 {
				return nil, utils.ErrFunctionOnlyQueryOneDayReport
			}
		}

		var mIDs []int64
		result := []*db.InOutM{}
		err = srv.Tx(func(tx *gorm.DB) error {
			mIDs, err = srv.bet02Dao.GetBet02s(tx, avgResp.Agent.ID, req.StartTime, req.EndTime)
			if err != nil {
				return err
			}
			result, err = srv.inOutMDao.GetInOutMs(tx, mIDs, req.OrderID, req.Order, req.StartTime, req.EndTime)
			if err != nil {
				return err
			}
			return nil
		})

		// Get member trade report key
		GetMemberTradeReport := "GetMemberTradeReport"
		vid := req.VendorID

		// Get last check time from Redis
		chkt, err := srv.redisCli.HGet(ctx, GetMemberTradeReport, vid).Result()
		if err != nil && err != redis.Nil {
			xlog.Warnf("error to get last check time from Redis: %v", err)
		}

		// Get current time
		chkTime := time.Now().Unix()
		// Check if it's igktwapi or no result
		if req.VendorID == "igktwapi" || len(result) == 0 {
			if chkt == "" {
				// First time access, set the time
				err = srv.redisCli.HSet(ctx, GetMemberTradeReport, vid, chkTime).Err()
				if err != nil {
					xlog.Errorf("error setting check time in Redis: %v", err)
					return nil, utils.ErrRedisError
				}
			} else {
				// Check time difference
				lastChkTime, err := strconv.ParseInt(chkt, 10, 64)
				if err != nil {
					xlog.Errorf("error parsing last check time: %v", err)
					return nil, utils.ErrRedisError
				}

				ct := chkTime - lastChkTime
				if ct < 10 {
					return nil, utils.ErrInvalidTransactionTimeoutRepeat
				}

				// Update check time
				err = srv.redisCli.HSet(ctx, GetMemberTradeReport, vid, chkTime).Err()
				if err != nil {
					xlog.Errorf("error updating check time in Redis: %v", err)
					return nil, utils.ErrRedisError
				}
			}
		} else {
			if chkt == "" {
				// First time access, set the time
				err = srv.redisCli.HSet(ctx, GetMemberTradeReport, vid, chkTime).Err()
				if err != nil {
					xlog.Errorf("error setting check time in Redis: %v", err)
					return nil, utils.ErrRedisError
				}
			} else {
				// Check time difference
				lastChkTime, err := strconv.ParseInt(chkt, 10, 64)
				if err != nil {
					xlog.Errorf("error parsing last check time: %v", err)
					return nil, utils.ErrRedisError
				}

				ct := chkTime - lastChkTime
				if ct < 30 {
					return nil, utils.ErrInvalidTransactionTimeout
				}

				// Update check time
				err = srv.redisCli.HSet(ctx, GetMemberTradeReport, vid, chkTime).Err()
				if err != nil {
					xlog.Errorf("error updating check time in Redis: %v", err)
					return nil, utils.ErrRedisError
				}
			}
		}
		var tradeItems []*view.TradeItem
		for _, v := range result {
			tradeItem := &view.TradeItem{
				MID:      v.Iom003,
				OrderID:  v.Iom001,
				OrderNum: v.Iom008,
				AddTime:  v.Iom002.Unix(),
				Money:    v.Iom004,
				OpCode:   v.Iom005,
				Subtotal: v.Iom010,
			}
			tradeItems = append(tradeItems, tradeItem)
		}
		return &view.GetMemberTradeReportResp{
			Result: tradeItems,
		}, nil
	}
}

func (srv *publicApiService) EnableOrDisableMem(ctx context.Context, req *view.EnableOrDisableMemReq) (*view.EnableOrDisableMemResp, error) {
	// Validate timestamp
	if err := utils.CheckTimestamp(req.Timestamp); err != nil {
		xlog.Errorf("error to check timestamp, err:%+v", err)
		return nil, err
	}

	// Verify agent
	avResp, err := srv.AgentVerify(ctx, &view.AgentVerifyReq{VendorID: req.VendorID, Signature: req.Signature})
	if err != nil {
		xlog.Errorf("error to verify agent, err:%+v", err)
		return nil, err
	}

	// Check agent status
	if avResp.Agent.Age015 == "N" {
		return nil, utils.ErrParamInvalidAgentDeactivated
	}

	// Validate status
	if req.Status != "Y" && req.Status != "N" {
		return nil, utils.ErrCommandSuccessButNoData
	}

	// Validate type and determine which field to update
	var columnToUpdate string
	var actionType string
	switch req.Type {
	case "login":
		columnToUpdate = "mem016"
		actionType = "登入"
		if req.Status == "Y" && avResp.Agent.Age015 == "N" {
			return nil, utils.ErrParamInvalidAgentDeactivated
		}
	case "bet":
		columnToUpdate = "mem017"
		actionType = "下注"
		if req.Status == "Y" && avResp.Agent.Age016 == "N" {
			return nil, utils.ErrParamInvalidAgentDeactivated
		}
	default:
		return nil, utils.ErrCommandSuccessButNoData
	}

	// Split and clean user accounts
	users := strings.Split(strings.ReplaceAll(req.User, " ", ""), ",")
	if len(users) == 0 {
		return nil, utils.ErrCommandSuccessButNoData
	}

	// Get status text based on language
	var statusText string
	if req.Status == "Y" {
		if req.Syslang == 1 {
			statusText = "enabled"
		} else {
			statusText = "已启用"
		}
	} else {
		if req.Syslang == 1 {
			statusText = "disabled"
		} else {
			statusText = "已停用"
		}
	}

	members, err := srv.userDao.QueryByAccounts(srv.DB(), users, avResp.Agent.ID)
	if err != nil {
		return nil, err
	}
	if len(members) == 0 {
		return nil, utils.ErrParamInvalidAccountNotExist
	}
	for _, member := range members {
		// Update member status
		err = srv.Tx(func(tx *gorm.DB) error {
			updates := map[string]interface{}{
				columnToUpdate: req.Status,
			}
			err := srv.userDao.UpdatesMember(tx, member.ID, updates)
			if err != nil {
				xlog.Errorf("Failed to update member status, err:%+v", err)
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	// Generate response message based on language
	var memberText, result string
	if req.Syslang == 1 {
		memberText = "Member"
		result = fmt.Sprintf("%s:%s %s %s", memberText, strings.Join(users, ","), statusText, actionType)
	} else {
		memberText = "会员"
		result = fmt.Sprintf("%s:%s %s %s", memberText, strings.Join(users, ","), statusText, actionType)
	}

	return &view.EnableOrDisableMemResp{
		Result: result,
	}, nil
}
