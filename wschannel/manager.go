package wschannel

import (
	"context"
	"encoding/json"
	"fmt"
	"go-zrbc/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"

	bhttp "go-zrbc/pkg/http/handler"

	"go-zrbc/pkg/http/middleware"

	"go-zrbc/pkg/xlog"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	pubSrv "go-zrbc/service/public"
	webSrv "go-zrbc/service/web"
)

const (
	UniversalToken = "UNIVERSAL_WS_TOKEN_123456" // Universal token for testing/development
)

type Server struct {
	addr string
	sync.RWMutex
	closeCh chan struct{}

	maxUserInWys int64
	totalConns   int64
	// all room
	rooms map[string]*Room

	clients map[string]*Client

	webService  webSrv.WebService
	userService pubSrv.PublicApiService
	redisCli    *redis.Client
}

func NewWsServer(addr string, webService webSrv.WebService,
	userService pubSrv.PublicApiService) *Server {
	redisAddr := config.Global.Redis.Addr
	redisDB := config.Global.Redis.DB
	redisCli := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: config.Global.Redis.Password, // no password set
		DB:       redisDB,                      // use default DB
	})
	srv := &Server{
		addr:         addr,
		maxUserInWys: 100000,
		// 暂时没用，看后续有没有用桌台号（groupID)建房间
		rooms:   make(map[string]*Room, 0),
		clients: make(map[string]*Client, 0),
		closeCh: make(chan struct{}),

		webService:  webService,
		userService: userService,
		redisCli:    redisCli,
	}
	return srv
}

func (srv *Server) Run() {
	md := middleware.NewMiddlewareHandler("ws-channel")

	metricGin := gin.New()
	metricGin.GET("/metrics", bhttp.WrapH(promhttp.Handler()))
	go metricGin.Run(fmt.Sprintf(":%d", config.Global.MetricPort))

	r := gin.New()
	r.Use(gin.Recovery(), md.RequestLog, md.Options)

	docHandler := bhttp.NewDocsHandler()
	docHandler.SetRouter(r)

	r.GET("/home", func(c *gin.Context) { srv.ServeHome(c) })
	// 大厅长连接服务
	r.GET("/15109", func(c *gin.Context) { srv.EnterLobby(c) })
	r.GET("/stats", func(c *gin.Context) { srv.LobbyStats(c) })

	// go bib.Init()
	go r.Run(srv.addr)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	select {
	case s := <-ch:
		xlog.Infof("receive interrupt signal(%v), close ws-channel manager", s)
		srv.Close()
		xlog.Infof("ws-channel manager has been closed successful")
	}
}

func (srv *Server) Close() {
	srv.Lock()
	defer srv.Unlock()
	for _, cli := range srv.clients {
		cli.Close("close")
	}

	close(srv.closeCh)
	for _, r := range srv.rooms {
		r.Close()
	}
}

func (srv *Server) ServeHome(c *gin.Context) {
	r := c.Request
	w := c.Writer
	log.Printf("server home, url(%v), path(%v)", r.URL, r.URL.Path)
	if r.URL.Path != "/home" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func (srv *Server) EnterLobby(c *gin.Context) {
	xlog.Infof("<====> enter wys. req header upgrade(%v)\n", c.Request.Header.Get("Upgrade"))

	cli := NewClient(c, srv)
	cli.Register()
}

type JsonResult struct {
	Code      int         `json:"code"`
	StatsData interface{} `json:"stats_data"`
}

func (srv *Server) LobbyStats(c *gin.Context) {
	rid := c.Query("rid")
	if rid == "" {
		srv.RLock()
		defer srv.RUnlock()
		statsMap := map[string]interface{}{"rooms": srv.rooms}

		msg, _ := json.Marshal(JsonResult{Code: 200, StatsData: statsMap})
		c.Writer.Write(msg)
	} else {
		srv.RLock()
		defer srv.RUnlock()
		if room, ok := srv.rooms[rid]; ok {
			room.RLock()
			defer room.RUnlock()
			statsMap := map[string]interface{}{"clients": room.clients}

			msg, _ := json.Marshal(JsonResult{Code: 200, StatsData: statsMap})
			c.Writer.Write(msg)
		} else {
			msg, _ := json.Marshal(JsonResult{Code: 500, StatsData: "room id err"})
			c.Writer.Write(msg)
		}
	}
}

func (srv *Server) GetOrCreateRoom(ctx context.Context, vid string, rType int) (*Room, error) {
	srv.RLock()
	if r, ok := srv.rooms[vid]; ok {
		srv.RUnlock()
		return r, nil
	}
	srv.RUnlock()

	r := NewRoom(vid, rType)
	if r == nil {
		return nil, errors.New("room cannot be created")
	}
	srv.Lock()
	defer srv.Unlock()
	srv.rooms[vid] = r
	return r, nil
}

func (srv *Server) AddClient(cli *Client) error {
	// todo 检查用户是否已经登录, 目前建立连接时没有用户信息，无法检查，auth的时候再检查，详细见MapUIDAndConn
	// local user number limit
	total := srv.Incr()
	if total > srv.maxUserInWys {
		return ErrWsChannelFull
	}

	srv.Lock()
	defer srv.Unlock()
	srv.clients[cli.connID] = cli
	return nil
}

func (srv *Server) Desc() int64 {
	return int64(atomic.AddInt64(&srv.totalConns, -1))
}

func (srv *Server) Incr() int64 {
	return int64(atomic.AddInt64(&srv.totalConns, 1))
}

func (srv *Server) RemoveClient(cli *Client) {
	srv.Lock()
	defer srv.Unlock()

	delete(srv.clients, cli.ConnID())
}
