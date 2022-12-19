package ziface

// 接口定义
type IServer interface {
	//启动
	Start()
	//停止
	Stop()
	//运行状态
	Serve()
	//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
	AddRouter(router IRouter)
}
