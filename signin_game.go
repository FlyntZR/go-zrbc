package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GameConfig contains static configuration
type GameConfig struct {
	UIModes   map[string]int
	SiteTypes map[string]string
}

// SigninGame represents the main game signin handler
type SigninGame struct {
	logger   *Logger
	db       *sql.DB
	dbHelper *DatabaseHelper
}

// Logger represents a simple logging interface
type Logger struct {
	logFile *os.File
}

// DatabaseHelper handles all database operations
type DatabaseHelper struct {
	db *sql.DB
}

// RequestParams contains all possible request parameters
type RequestParams struct {
	User      string `json:"user"`
	Device    string `json:"device"`
	Lang      string `json:"lang"`
	IsTest    bool   `json:"isTest"`
	Mode      string `json:"mode"`
	TableID   string `json:"tableid"`
	Site      string `json:"site"`
	Password  string `json:"password"`
	GameType  string `json:"gameType"`
	Width     string `json:"width"`
	ReturnURL string `json:"returnurl"`
	Size      string `json:"size"`
	UI        string `json:"ui"`
	Mute      string `json:"mute"`
	Video     string `json:"video"`
}

// Session represents a user session
type Session struct {
	SID string
}

// NewSigninGame creates a new instance of SigninGame
func NewSigninGame(db *sql.DB) *SigninGame {
	sg := &SigninGame{
		db:       db,
		dbHelper: NewDatabaseHelper(db),
	}
	sg.initializeLogger()
	return sg
}

func (sg *SigninGame) initializeLogger() {
	logDir := filepath.Join("log", "public", "SigninGame")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(err)
	}

	logFile := filepath.Join(logDir, fmt.Sprintf("new%s.log", time.Now().Format("2006-01-02 15:00")))
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	sg.logger = &Logger{logFile: file}
}

func (sg *SigninGame) validateRequest(user, password string, isTest bool) error {
	if !isTest {
		if password == "" {
			return errors.New("password is required")
		}
		if user == "" {
			return errors.New("username is required")
		}

		exists, err := sg.dbHelper.GetUserByUsername(user)
		if err != nil {
			return err
		}
		if !exists {
			return errors.New("user not found")
		}

		valid, err := sg.dbHelper.CheckUserCredentials(user, password)
		if err != nil {
			return err
		}
		if !valid {
			return errors.New("invalid credentials")
		}
	}
	return nil
}

func (sg *SigninGame) processUI(ui string) string {
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

func (sg *SigninGame) buildGameURL(baseURL string, params map[string]string) string {
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
	loginPass, err := sg.dbHelper.GetLoginPassByVendor(vendorID)
	co := "wm"
	if err == nil && loginPass["co"] != "" {
		co = loginPass["co"]
	}
	urlParams = append(urlParams, "co="+co)

	return fmt.Sprintf("%s?%s", baseURL, strings.Join(urlParams, "&"))
}

func (sg *SigninGame) Execute(params RequestParams) (string, error) {
	// Log request
	log := map[string]interface{}{
		"class":       "SigninGame",
		"request":     params,
		"client_ip":   os.Getenv("USER_IP"),
		"request_url": os.Getenv("SERVER_NAME") + os.Getenv("REQUEST_URI"),
	}

	// Validate request
	if err := sg.validateRequest(params.User, params.Password, params.IsTest); err != nil {
		sg.logAndReturn(err, nil, log, true)
		return "", err
	}

	// Process UI settings
	params.UI = sg.processUI(params.UI)

	// Get base URL
	baseURL, err := sg.getBaseURL(params.Site, os.Getenv("AGENT_CURRENCY"))
	if err != nil {
		sg.logAndReturn(err, nil, log, true)
		return "", err
	}

	// Handle test mode
	if params.IsTest {
		urlParams := map[string]string{"sid": "ANONYMOUS"}
		gameURL := sg.buildGameURL(baseURL, urlParams)
		sg.logAndReturn(nil, gameURL, log, false)
		return gameURL, nil
	}

	// Handle normal login
	session, err := sg.login(params.User, params.Password)
	if err != nil {
		sg.logAndReturn(err, nil, log, true)
		return "", err
	}

	// Build final game URL
	urlParams := map[string]string{
		"sid":       session.SID,
		"mode":      params.Mode,
		"tableid":   params.TableID,
		"ui":        params.UI,
		"mute":      params.Mute,
		"lang":      params.Lang,
		"width":     params.Width,
		"returnurl": params.ReturnURL,
		"size":      params.Size,
		"video":     params.Video,
	}

	gameURL := sg.buildGameURL(baseURL, urlParams)
	sg.logAndReturn(nil, gameURL, log, false)
	return gameURL, nil
}

func (sg *SigninGame) getBaseURL(site, currency string) (string, error) {
	switch site {
	case "WECHAT_VERTICAL", "WECHAT_HORIZONTAL", "WECHAT_HORIZONTAL_ALT":
		return sg.getWechatURL()
	default:
		return sg.getAPIURL(currency)
	}
}

func (sg *SigninGame) getAPIURL(code string) (string, error) {
	url, err := sg.dbHelper.GetAPIURLByCode(code)
	if err != nil {
		if code != "0" {
			return sg.getAPIURL("0")
		}
		return "", err
	}
	return url, nil
}

func (sg *SigninGame) getWechatURL() (string, error) {
	urlData, err := sg.dbHelper.GetRandomWechatURL()
	if err != nil {
		return "", err
	}

	if err := sg.dbHelper.UpdateWechatURLUseCount(urlData.ID); err != nil {
		// Log error but continue
		sg.logger.Error("Failed to update Wechat URL use count", err)
	}

	return urlData.URL, nil
}

func (sg *SigninGame) logAndReturn(err error, result interface{}, log map[string]interface{}, isError bool) {
	response := map[string]interface{}{
		"errorCode":   0,
		"errorRemark": "",
		"result":      result,
	}

	if err != nil {
		response["errorCode"] = 1
		response["errorRemark"] = err.Error()
	}

	log["response"] = response

	if isError {
		sg.logger.Warning(err.Error(), log)
	} else {
		sg.logger.Info("Operation successful", log)
	}
}

// Logger methods
func (l *Logger) Info(msg string, data map[string]interface{}) {
	l.log("INFO", msg, data)
}

func (l *Logger) Warning(msg string, data map[string]interface{}) {
	l.log("WARNING", msg, data)
}

func (l *Logger) Error(msg string, err error) {
	l.log("ERROR", msg, map[string]interface{}{"error": err.Error()})
}

func (l *Logger) log(level, msg string, data map[string]interface{}) {
	entry := map[string]interface{}{
		"level":     level,
		"message":   msg,
		"data":      data,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal log entry: %v\n", err)
		return
	}

	if _, err := l.logFile.Write(append(jsonData, '\n')); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write log entry: %v\n", err)
	}
}

// login handles user authentication and session creation
func (sg *SigninGame) login(username, password string) (*Session, error) {
	// In a real implementation, this would integrate with your authentication system
	// For now, we'll create a simple session
	sid := fmt.Sprintf("session_%d", time.Now().Unix())
	return &Session{SID: sid}, nil
}
