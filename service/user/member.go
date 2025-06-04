package service

import (
	"context"
	"errors"
	"fmt"
	"go-zrbc/db"
	"go-zrbc/pkg/xlog"
	"go-zrbc/service"
	"go-zrbc/view"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-redis/redis/v8"

	"gorm.io/gorm"
)

const (
	TokenPrefix = "zrys:access_token"
)

type UserService interface {
	//用户信息
	GetUserInfo(ctx context.Context, userID int64) (*view.GetUserInfoResp, error)
	GetUserByAccountAndPwd(ctx context.Context, account, pwd string) (*view.GetUserInfoResp, error)
	SigninGame(ctx context.Context, req *view.SigninGameReq) (*view.SigninGameResp, error)
}

type userService struct {
	userDao            db.UserDao
	apiurlDao          db.ApiurlDao
	wechatURLDao       db.WechatURLDao
	agentsLoginPassDao db.AgentsLoginPassDao

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

	s3Client *s3.Client,
	redisCli *redis.Client,
) UserService {
	srv := &userService{
		userDao:            userDao,
		apiurlDao:          apiurlDao,
		wechatURLDao:       wechatURLDao,
		agentsLoginPassDao: agentsLoginPassDao,

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

func (srv *userService) validateRequest(user, password string, isTest bool) error {
	if !isTest {
		if password == "" {
			return errors.New("password is required")
		}
		if user == "" {
			return errors.New("username is required")
		}

		mem, err := srv.userDao.QueryByAccount(srv.DB(), user)
		if err == gorm.ErrRecordNotFound {
			return errors.New("user not found")
		}
		if err != err {
			return err
		}

		if mem.Password != password {
			return errors.New("invalid credentials")
		}
	}
	return nil
}

func processUI(ui string) string {
	switch ui {
	case "1", "2":
		return "MODERN"
	case "4":
		return "CUSTOM_4"
	case "5":
		return "CUSTOM_5"
	default:
		return "DEFAULT"
	}
}

func (srv *userService) getBaseURL(site, currency string) (string, error) {
	switch site {
	case "WECHAT_VERTICAL", "WECHAT_HORIZONTAL", "WECHAT_HORIZONTAL_ALT":
		return srv.getWechatURL()
	default:
		return srv.getAPIURL(currency)
	}
}

func (srv *userService) getAPIURL(code string) (string, error) {
	codeInt, err := strconv.Atoi(code)
	if err != nil {
		xlog.Warnf("warn to strconv atoi, code:%s", code)
		codeInt = 0
	}
	apiURL, err := srv.apiurlDao.QueryByCode(srv.DB(), codeInt)
	if err != nil {
		if code != "0" {
			defaultApiURL, err := srv.apiurlDao.QueryByCode(srv.DB(), 0)
			if err != nil {
				return "", err
			}
			return defaultApiURL.URL, nil
		}
		return "", err
	}
	return apiURL.URL, nil
}

func (srv *userService) getWechatURL() (string, error) {
	urlData, err := srv.wechatURLDao.GetRandomWechatURL(srv.DB())
	if err != nil {
		return "", err
	}

	if err := srv.wechatURLDao.UpdateWechatURLUseCount(srv.DB(), urlData.ID); err != nil {
		// Log error but continue
		xlog.Errorf("Failed to update Wechat URL use count, err:%+v", err)
	}

	return urlData.URL, nil
}

func (srv *userService) buildGameURL(baseURL string, params map[string]string) string {
	vendorID := os.Getenv("AGENT_AGE002")
	urlParams := make([]string, 0)

	// Add SID
	if vendorID == "bbinapi" || vendorID == "bbtest" || vendorID == "bbintwapi" {
		urlParams = append(urlParams, "#sid="+params["sid"])
	} else {
		urlParams = append(urlParams, "sid="+params["sid"])
	}

	// Add other parameters
	if params["mode"] != "" {
		urlParams = append(urlParams, "mode="+params["mode"])
		if params["tableid"] != "" {
			urlParams = append(urlParams, "tableid="+params["tableid"])
		}
	}

	// Add optional parameters
	optionalParams := []string{"ui", "mute", "lang", "returnurl"}
	for _, param := range optionalParams {
		if params[param] != "" {
			urlParams = append(urlParams, param+"="+params[param])
		}
	}

	if params["width"] != "" {
		urlParams = append(urlParams, "checkWidth=deviceWidth")
	}
	if params["size"] != "" {
		urlParams = append(urlParams, "size=1")
	}
	if params["video"] == "off" {
		urlParams = append(urlParams, "video=off")
	}

	// Add company code
	loginPass, err := srv.agentsLoginPassDao.GetLoginPassByVendor(srv.DB(), vendorID)
	co := "wm"
	if err == nil && loginPass["co"] != "" {
		co = loginPass["co"]
	}
	urlParams = append(urlParams, "co="+co)

	return fmt.Sprintf("%s?%s", baseURL, strings.Join(urlParams, "&"))
}

// login handles user authentication and session creation
func (srv *userService) login(username, password string) (*view.Session, error) {
	// In a real implementation, this would integrate with your authentication system
	// For now, we'll create a simple session
	sid := fmt.Sprintf("session_%d", time.Now().Unix())
	return &view.Session{SID: sid}, nil
}

func (srv *userService) SigninGame(ctx context.Context, req *view.SigninGameReq) (*view.SigninGameResp, error) {
	// Validate request
	if err := srv.validateRequest(req.User, req.Password, req.IsTest); err != nil {
		xlog.Errorf("error to get user, err:%+v", err)
		return nil, err
	}

	// Process UI settings
	req.UI = processUI(req.UI)

	// Get base URL
	baseURL, err := srv.getBaseURL(req.Site, os.Getenv("AGENT_CURRENCY"))
	if err != nil {
		xlog.Errorf("error to get base url, err:%+v", err)
		return nil, err
	}

	// Handle test mode
	if req.IsTest {
		urlParams := map[string]string{"sid": "ANONYMOUS"}
		gameURL := srv.buildGameURL(baseURL, urlParams)
		resp := view.SigninGameResp{
			GameURL: gameURL,
		}
		return &resp, nil
	}

	// Handle normal login
	session, err := srv.login(req.User, req.Password)
	if err != nil {
		xlog.Errorf("error to get user by account and pwd, err:%+v", err)
		return nil, err
	}

	// Build final game URL
	urlParams := map[string]string{
		"sid":       session.SID,
		"mode":      req.Mode,
		"tableid":   req.TableID,
		"ui":        req.UI,
		"mute":      req.Mute,
		"lang":      req.Lang,
		"width":     req.Width,
		"returnurl": req.ReturnURL,
		"size":      req.Size,
		"video":     req.Video,
	}

	gameURL := srv.buildGameURL(baseURL, urlParams)
	resp := view.SigninGameResp{
		GameURL: gameURL,
	}
	return &resp, nil
}
