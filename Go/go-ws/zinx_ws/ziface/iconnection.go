package ziface

import (
	"github.com/gorilla/websocket"
	"net"
)

type IConnection interface {
	//启动链接开始工作
	Start()
	//关闭链接停止工作
	Close()
	//获取websocket链接
	GetConnection() *websocket.Conn
	//获取当前连接ID
	GetConnID() uint64
	//获取远程客户端地址信息
	RemoteAddr() net.Addr
}

// 定义一个统一处理链接业务的接口
type HandFunc func(*websocket.Conn, []byte, int) error
