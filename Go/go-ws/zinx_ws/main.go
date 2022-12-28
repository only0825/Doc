package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
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

	// 获取命令行参数
	argc := len(os.Args)
	if argc != 2 {
		zlog.Error.Println("运行格式错误，格式为 ./应用 <配置文件名称>")
		return
	}

	if err := configs.LoadConfig(os.Args[1]); err != nil {
		zlog.Error.Println("Load config json error:", err)
		return
	}

	msg := fmt.Sprintf("服务监听 %s:%d", configs.Conf.Server.Ip, configs.Conf.Server.Port)
	zlog.Info.Printf(msg)

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
	r.GET("/change", server.WsHandler)
	err = r.Run(bindAddress)
	if err != nil {
		zlog.Error.Println("启动服务失败:", err)
	}
	//server.GWServer.GetConnMgr().

}
