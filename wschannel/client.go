package wschannel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"go-zrbc/pkg/utils"
	"go-zrbc/pkg/xlog"
	"go-zrbc/view"

	bhttpresp "go-zrbc/pkg/http/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    128,
	WriteBufferSize:   128,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	connID   string // 连接ID，每个WS连接都有一个唯一的ID
	status   int64
	ginCtx   *gin.Context
	Room     *Room
	User     *view.WsUser
	DeviceID string

	logger *Logger

	mgr  *Server
	conn *websocket.Conn

	// used to send json text msg
	bytesSend chan []byte

	LastActive       int64 // 最后活动时间
	HeartbeatRetried int64 // 已连续 N 次没有响应服务端的主动 ping 消息（服务端会在将要超时的时候主动 ping）
	CreatedAt        int64
}

type PreparedMessage struct {
	CreatedAt int64
	wsp       *websocket.PreparedMessage
}

func NewPreparedMessage(wsp *websocket.PreparedMessage) *PreparedMessage {
	return &PreparedMessage{
		wsp:       wsp,
		CreatedAt: time.Now().Unix(),
	}
}

func NewClient(ginCtx *gin.Context, mgr *Server) *Client {
	// todo urlParam后面公共推送有用
	cli := &Client{
		connID:    uuid.New().String(),
		ginCtx:    ginCtx,
		mgr:       mgr,
		bytesSend: make(chan []byte, 256),
		CreatedAt: time.Now().Unix(),
	}
	cli.logger = cli.NewLogger()
	return cli
}

func (cli *Client) Response(resp *ConnMessageResp) {
	b, err := json.Marshal(resp)
	if err != nil {
		xlog.Error(err)
		return
	}
	cli.bytesSend <- b
}

func (cli *Client) ConnID() string {
	return cli.connID
}

func (cli *Client) UserID() int64 {
	return cli.User.ID
}

func (cli *Client) Register() {
	defer cli.Unregister()

	err := cli.mgr.AddClient(cli)
	if errors.Is(err, ErrWsChannelFull) {
		cli.logger.Error(err)
		bhttpresp.BadRequestResp(cli.ginCtx, err)
		return
	}
	if err != nil {
		cli.logger.Error(err)
		bhttpresp.ServerErrResp(cli.ginCtx, err)
		return
	}
	conn, err := upgrader.Upgrade(cli.ginCtx.Writer, cli.ginCtx.Request, nil)
	if err != nil {
		cli.logger.Error(err)
		bhttpresp.BadRequestResp(cli.ginCtx, err)
		return
	}
	cli.conn = conn
	cli.logger.Debugf("client(%s) enter ws channel, let's welcome", cli)

	go cli.writePump()
	cli.readPump()
}

func (cli *Client) Unregister() {
	cli.mgr.Desc()
	xlog.Infof("receive unregister msg, client(%s) will leave\n", cli)
	cli.mgr.RemoveClient(cli)
	cli.Close("unregister")
}

var (
	ErrTypeAssert = errors.New("type assert error")
	ErrMsgType    = errors.New("ws msg type not found")
)

func (cli *Client) NewLogger() *Logger {
	return &Logger{
		ServiceID: "ws-channel",
		UUID:      cli.connID,
		// UserID:    strconv.FormatInt(cli.User.ID, 10),
	}
}

const (
	statusClosed = 1
)

func (cli *Client) Close(reason string) {
	if atomic.LoadInt64(&cli.status) == statusClosed {
		return
	}
	cli.logger.Infof("client(%s) conn will be closed, reason(%s)\n", cli, reason)
	close(cli.bytesSend)
	if cli.conn != nil {
		cli.conn.Close()
	}
	atomic.StoreInt64(&cli.status, statusClosed)
}

func (c *Client) String() string {
	return fmt.Sprintf("client: conn_id(%s), user_id(%d), ts(%d)", c.connID, c.User.ID, c.CreatedAt)
}

func (cli *Client) WriteControlFrame() {
	err := cli.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye bye!"))
	if err != nil {
		cli.logger.Errorf("write close, err: %v", err)
		return
	}
}

var wsMsgTypeText = map[int]string{
	websocket.TextMessage:   "text",
	websocket.BinaryMessage: "binary",
}

func (cli *Client) readPump() {
	cli.conn.SetReadLimit(maxMessageSize)
	cli.conn.SetReadDeadline(time.Now().Add(pongWait))
	cli.conn.SetPongHandler(func(string) error { cli.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		mt, message, err := cli.conn.ReadMessage()
		cli.logger.Debugf("receive message from wsclient(%v): ws_msg_type(%v), message(%v), size(%d bytes), err(%v)\n", cli, wsMsgTypeText[mt], string(message), len(message), err)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				cli.logger.Errorf("unexpected error: cli(%v), %v\n", cli, err)
			} else {
				cli.logger.Errorf("expected error: cli(%v) %v\n", cli, err)
			}
			break
		}
		var wsReq view.WsReq
		err = json.Unmarshal(message, &wsReq)
		if err != nil {
			cli.logger.Error(err)
			cli.Response(RespDataFormatError)
			continue
		}
		xlog.Debugf("receive message from ws client(%v): wsReq(%+v)\n", cli, wsReq)
		err = cli.HandlerWsReq(&wsReq)
		if err != nil {
			cli.logger.Error(err)
			cli.Response(RespDataFormatError)
			continue
		}
	}
}

func (cli *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	periodicMsgTicker := time.NewTicker(6 * time.Second)
	defer func() {
		ticker.Stop()
		periodicMsgTicker.Stop()
		cli.conn.Close()
	}()

	periodicMsg := view.WsResp{
		Protocol: 999,
		Data:     struct{}{},
	}
	periodicMsgBytes, _ := json.Marshal(periodicMsg)

	for {
		select {
		case bs, ok := <-cli.bytesSend:
			if !ok {
				cli.logger.Error("bytesSend pipe closed")
				cli.WriteControlFrame()
				cli.conn.Close()
				return
			}
			// todo test BinaryMessage改为TextMessage
			// cli.logger.Debugf("bytesSend get data :%v", string(bs))
			if err := cli.conn.WriteMessage(websocket.TextMessage, bs); err != nil {
				cli.logger.Error(err)
				return
			}

		case <-ticker.C:
			cli.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := cli.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
			cli.conn.SetWriteDeadline(time.Time{})

		case <-periodicMsgTicker.C:
			if err := cli.conn.WriteMessage(websocket.TextMessage, periodicMsgBytes); err != nil {
				cli.logger.Error(err)
				return
			}
		}
	}
}

func (cli *Client) HandlerWsReq(wsReq *view.WsReq) error {
	switch wsReq.Protocol {
	case 0: // 登录验证
		return cli.HandlerAuthReq(wsReq)
	case 115: // 不接受指定游戏资料
		return cli.Handler115Req(wsReq)
	default:
		cli.logger.Errorf("HandlerReqReq protocol err, wsReq:%+v, err:(%+v)", wsReq, errors.New("req protocol err"))
		return errors.New("req protocol err")
	}
}

func (cli *Client) HandlerAuthReq(wsReq *view.WsReq) error {
	resp := view.WsResp{
		Protocol: wsReq.Protocol,
	}

	var ad view.AuthData
	bb, _ := json.Marshal(wsReq.Data)
	if err := json.Unmarshal(bb, &ad); err != nil {
		cli.logger.Errorf("HandlerAuthReq data err, wsReq:%+v, err:(%+v)", wsReq, err)
		return errors.New("auth data err")
	}
	if ad.Account == "" {
		err := errors.New("account is empty")
		cli.logger.Errorf("HandlerAuthReq data err, wsReq:%+v, err:(%+v)", wsReq, err)
		resp.Data = view.AuthResp{
			BOk: false,
		}
		respBin, _ := json.Marshal(resp)
		cli.bytesSend <- respBin
		return err
	}
	if ad.Password == "" {
		err := errors.New("password is empty")
		cli.logger.Errorf("HandlerAuthReq data err, wsReq:%+v, err:(%+v)", wsReq, err)
		resp.Data = view.AuthResp{
			BOk:            false,
			BValidPassword: false,
		}
		respBin, _ := json.Marshal(resp)
		cli.bytesSend <- respBin
		return err
	}
	userResp, err := cli.mgr.userService.GetUserByAccountAndPwd(context.TODO(), ad.Account, ad.Password)
	if err != nil {
		resp.Data = view.AuthResp{
			BOk: false,
		}
		respBin, _ := json.Marshal(resp)
		cli.bytesSend <- respBin
		return err
	}
	resp.Data = view.AuthResp{
		MemberID:       userResp.User.ID,
		Account:        userResp.User.User,
		UserName:       userResp.User.UserName,
		BOk:            true,
		Sid:            utils.GenerateUserCode(ad.Account, time.Now().Unix()),
		BValidPassword: true,
	}

	respBin, _ := json.Marshal(resp)
	cli.bytesSend <- respBin
	return nil
}

func (cli *Client) Handler115Req(wsReq *view.WsReq) error {
	resp := view.WsResp{
		Protocol: wsReq.Protocol,
	}
	respBin, _ := json.Marshal(resp)
	cli.bytesSend <- respBin
	return nil
}
