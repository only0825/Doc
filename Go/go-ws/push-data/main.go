package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"push-data/configs"
	"push-data/iserver"
	"push-data/model"
	"push-data/server"
	"push-data/utils"
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
			logrus.Error("ping 回调失败", err)
		}
	}
}

func main() {
	fmt.Println("staring push ...")

	// 获取命令行参数
	argc := len(os.Args)
	if argc != 2 {
		logrus.Error("运行格式错误，格式为 ./应用 <配置文件名称>")
		return
	}

	// 加载配置
	if err := configs.LoadConfig(os.Args[1]); err != nil {
		logrus.Error("Load config json error:", err)
		return
	}

	// 初始化日志
	utils.LogInit("/push-data.log")

	msg := fmt.Sprintf("服务监听 %s:%d", configs.Conf.Server.Ip, configs.Conf.Server.Port)
	logrus.Info(msg)

	rdb, err := utils.InitRedis()
	model.Rdb = rdb
	//rdb, err := utils.InitRedisCluster()
	//model.Rdb = rdb
	if err != nil {
		logrus.Error("Redis初始化错误:", err)
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
		logrus.Error("启动服务失败:", err)
	}
	//server.GWServer.GetConnMgr().

}
