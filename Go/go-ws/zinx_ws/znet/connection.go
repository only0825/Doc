package znet

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"zinx_ws/ziface"
)

type Connection struct {
	//当前连接的socket TCP套接字
	Conn *websocket.Conn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnID uint32
	//当前连接的关闭状态
	isClosed bool
	//该连接的处理方法router
	Router ziface.IRouter
	//告知该链接已经退出/停止的channel
	ExitBuffChan chan bool
}

func (c *Connection) Close() {
	//TODO implement me
	panic("implement me")
}

func (c *Connection) GetConnID() uint64 {
	//TODO implement me
	panic("implement me")
}

// 创建连接的方法
func NewConnection(conn *websocket.Conn, connID uint32) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		ExitBuffChan: make(chan bool, 1),
	}

	return c
}

// 启动连接，让当前连接，开始工作
func (c *Connection) Start() {
	//启动读数据业务
	go c.StartReader()

	for {
		select {
		case <-c.ExitBuffChan:
			//得到退出消息，不再阻塞
			return
		}
	}
}

// 停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	//1. 如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//TODO Connection Stop() 如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用

	// 关闭socket链接
	c.Conn.Close()

	//通知从缓冲队列读数据的业务，该链接已经关闭
	c.ExitBuffChan <- true

	//关闭该链接全部管道
	close(c.ExitBuffChan)
}

/* 处理conn读数据的Goroutine */
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()

	for {
		//读取数据到内存中 messageType:TextMessage/BinaryMessage
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("connection startReader read err:", err)
			break
		}
		log.Println("connection StartReader recv from client1:", string(data))
		//得到request数据
		req := Request{
			conn:    c,
			message: string(data),
		}

		//从路由Routers 中找到注册绑定Conn的对应Handle
		//go func(request ziface.IRequest) {
		//	//执行注册的路由方法
		//	c.Router.PreHandle(request)
		//	c.Router.Handle(request)
		//	c.Router.PostHandle(request)
		//}(&req)
		fmt.Println(req.GetData())
	}
}

// 获取当前连接的websocket conn
func (c *Connection) GetWsConnection() *websocket.Conn {
	return c.Conn
}

// 获取当前连接ID
//func (c *Connection) GetConnID() uint32 {
//	return c.ConnID
//}

// 获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
