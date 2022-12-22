package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"zinx_ws/configs"
	"zinx_ws/iserver"
	"zinx_ws/server"
	"zinx_ws/zlog"
)

var (
	configFile string
)

//	WebSocket服务端
//
// ping test 自定义路由
type PingRouter struct {
	server.BaseRouter
}

// Ping Handle
func (this *PingRouter) Handle(request iserver.IRequest) {
	fmt.Println(request.GetMsgID())
	if request.GetMsgID() == "1040" {
		for {
			err := request.GetConnection().SendMessage(request.GetMsgType(), request.GetData())
			if err != nil {
				logger.Error.Println("回调失败", err)
			}
			time.Sleep(time.Duration(3) * time.Second)
		}
	}
}

func init() {
	flag.StringVar(&configFile, "config", "./config.json", "where load config json")
	flag.Parse()

	if err := configs.LoadConfig(configFile); err != nil {
		fmt.Println("Load config json error:", err)
		return
	}
	//common.InitRedis()
}

func main() {
	//配置路由
	server.GWServer = server.NewServer()
	server.GWServer.AddRouter("ping", &PingRouter{})
	server.GWServer.AddRouter("1040", &PingRouter{})

	bindAddress := fmt.Sprintf("%s:%d", configs.GConf.Ip, configs.GConf.Port)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/score", server.WsHandler)
	r.Run(bindAddress)
}
