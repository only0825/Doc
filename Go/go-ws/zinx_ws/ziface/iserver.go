package ziface

// 接口定义
type IServer interface {
	//启动
	Start()
	//停止
	Stop()
	//运行状态
	Serve()
}
