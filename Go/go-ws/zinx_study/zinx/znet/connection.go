package znet

import (
	"go-ws/zinx_study/zinx/ziface"
	"net"
)

type Connection struct {
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnID uint32
	//当前连接的关闭状态
	isClosed bool
	//该连接的处理方法api
	handleAPI ziface.HandFunc
	//告知该链接已经退出/停止的channel
	ExitBuffChan chan bool
}
