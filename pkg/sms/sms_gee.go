package sms

import (
	"fmt"
	"go-zrbc/config"
	"go-zrbc/pkg/common"
	util "go-zrbc/pkg/utils"
	"go-zrbc/pkg/xlog"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/lifei6671/gorand"
	"github.com/pkg/errors"
)

type geeSMS struct {
}

func NewGeeSMS() SMSServer {
	return &geeSMS{}
}

func GenerateAuth() string {
	gtID := "xp9mzzxttrrjheg8jtojwskqzz64zq3j"
	gtKey := "h9yldjrzxaeiabtad0kb4ty5ivj7ehr1"
	// nonce := "frxwel0nioxt92smrtn509majr5750lj"
	nonce := strings.ToLower(string(gorand.KRand(32, gorand.KC_RAND_KIND_ALL)))
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	// timestamp := "1717579165"

	//将gtID、timestamp、nonce三个参数进行字典序排序
	var tempArray = []string{gtID, timestamp, nonce}
	sort.Strings(tempArray)
	//将三个参数字符串拼接成一个字符串进行sha1加密
	var sha1String string = ""
	for _, v := range tempArray {
		sha1String += v
	}
	signature := util.HmacSha256Hex(sha1String, gtKey)
	return fmt.Sprintf("gt_id=%s,nonce=%s,signature=%s,timestamp=%s", gtID, nonce, signature, timestamp)
}

type SendMessageReq struct {
	Phone      string            `json:"phone"`
	ModeID     string            `json:"modeId"`
	Arguments  map[string]string `json:"arguments"`
	ExtendCode string            `json:"extendCode"`
}

type SendMessageResp struct {
	Status   int    `json:"status"`
	TradeID  string `json:"trade_id"`
	ErrorMsg string `json:"error_msg"`
}

func (s *geeSMS) SendMessage(phone, code string) error {
	defer func() {
		e := recover()
		if e != nil {
			xlog.Errorf("recover from panic error, failed to http auth to user service, err:%+v", e)
		}
	}()

	sURL := "https://tectapi.geetest.com/v2/message"
	auth := GenerateAuth()
	client := resty.New()
	if config.Global.Agent != "" {
		client.SetProxy(config.Global.Agent)
	}
	sResp, err := client.R().SetHeader("Authorization", auth).SetBody(&SendMessageReq{Phone: phone, ModeID: config.Global.SMSModeID}).Post(sURL)
	if err != nil {
		return errors.Wrapf(err, "failed to http send message")
	}
	xlog.Infof("http send message sURL:%s, resp: %v", sURL, sResp)
	var result struct {
		SendMessageResp
	}
	common.DecodeJSONFromBytes(sResp.Body(), &result)
	if result.Status != 200 {
		return errors.New("failed to send message")
	}

	return nil
}
