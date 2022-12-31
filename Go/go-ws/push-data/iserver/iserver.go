package iserver

import "github.com/gin-gonic/gin"

// 接口定义
type IServer interface {
	//启动
	Start(c *gin.Context)
	//停止
	Stop()
	//运行状态
	Serve(c *gin.Context)
	//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
	AddRouter(msgId string, router IRouter)
	//得到链接管理
	GetConnMgr() IConnMgr
}
