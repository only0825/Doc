package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"zinx_ws/common"
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
	switch request.GetMsgID() {
	case "ping":
		err := request.GetConnection().SendMessage(request.GetMsgType(), request.GetData())
		if err != nil {
			zlog.Error.Println("ping 回调失败", err)
		}
	}
}

func main() {

	if err := configs.LoadConfig(); err != nil {
		fmt.Println(err)
		zlog.Error.Println("Load config json error:", err)
		return
	}

	cache := configs.Conf.Cache
	err := common.InitCache(cache)
	if err != nil {
		zlog.Error.Println("Redis初始化错误:", err)
		return
	}

	//配置路由
	server.GWServer = server.NewServer()
	server.GWServer.AddRouter("ping", &PingRouter{})

	bindAddress := fmt.Sprintf("%s:%d", configs.Conf.Server.Ip, configs.Conf.Server.Port)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/score", server.WsHandler)
	r.Run(bindAddress)

	//server.GWServer.GetConnMgr().

}
