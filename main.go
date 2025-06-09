package main

import (
	"fmt"
	"go-zrbc/config"
	"go-zrbc/db"
	"go-zrbc/http"
	awsS3 "go-zrbc/pkg/oss"
	"go-zrbc/pkg/xlog"
	"go-zrbc/service"
	pService "go-zrbc/service/public"
	sService "go-zrbc/service/s3"
	wService "go-zrbc/service/web"
	"go-zrbc/wschannel"

	"os"
	"os/signal"

	_ "go-zrbc/docs"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
)

const (
	DATAID = "activity-api"
	GROUP  = "business-app"
)

func main() {
	Execute()
}

var rootCmd = &cobra.Command{
	Use:   "YsWebSystem",
	Short: "YsWebSystem is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		if configPath == "config.json" {
			fmt.Println("config path:", configPath)
			config.Init(configPath)
		} else {
			// config.InitConfigWithJson(DATAID, GROUP, config.Global)
			config.GetConfigFromApollo(config.Global)
		}
		xlog.LogLevel = config.Global.LogLevel
		xlog.LogFile = config.Global.LogFile
		xlog.Init()
		Start()
	},
}

var (
	configPath string
)

func init() {
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "path to config file")
	//rootCmd.MarkFlagRequired("config")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func Start() {
	dbh := db.NewDBHandler()
	sess := service.NewSession(dbh)

	redisCli := redis.NewClient(&redis.Options{
		Addr:     config.Global.Redis.Addr,
		Password: config.Global.Redis.Password, // no password set
		DB:       config.Global.Redis.DB,       // use default DB
	})
	s3Client := awsS3.S3Client()

	userDao := db.NewMemberDao()
	barrageDao := db.NewBarrageDao()
	apiurlDao := db.NewApiurlDao()
	wechatURLDao := db.NewWechatURLDao()
	agentsLoginPassDao := db.NewAgentsLoginPassDao()
	agentDao := db.NewAgentDao()
	memLoginDao := db.NewMemLoginDao()
	bet02Dao := db.NewBet02Dao()
	agentDtlDao := db.NewAgentDtlDao()
	betLimitDao := db.NewBetLimitDefaultDao()
	memberDtlDao := db.NewMemberDtlDao()
	gameTypeDao := db.NewGameTypeDao()
	inOutMDao := db.NewInOutMDao()
	logAgeCashChangeDao := db.NewLogAgeCashChangeDao()
	alertMessageDao := db.NewAlertMessageDao()

	userSrv := pService.NewPublicApiService(sess, userDao, apiurlDao, wechatURLDao, agentsLoginPassDao, agentDao, memLoginDao, bet02Dao, agentDtlDao, betLimitDao, memberDtlDao, gameTypeDao, inOutMDao, logAgeCashChangeDao, alertMessageDao, s3Client, redisCli)
	s3Srv := sService.NewS3Service()
	webSrv := wService.NewWebService(sess, barrageDao, s3Client, redisCli)

	httpSrv := http.NewServer(webSrv, userSrv, s3Srv)

	go httpSrv.RunMetric()
	go httpSrv.Run()

	wsSrv := wschannel.NewWsServer("0.0.0.0:8082", webSrv, userSrv)
	go wsSrv.Run()

	xlog.Info("server start success!")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
	case s := <-c:
		xlog.Info("receive interrupt signal", s)
	}
}
