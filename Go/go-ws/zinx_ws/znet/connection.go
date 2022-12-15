package znet

import (
	"github.com/gorilla/websocket"
	"go-ws/zinx_ws/ziface"
)

// 初始化连接方法
func NewConnection(server ziface.IServer, conn *websocket.Conn, connID uint32, mh ziface.IMsgHandle) *Connection {
	c := &Connection{
		WsServer:    server,
		Conn:        conn,
		ConnID:      connID,
		MsgHandle:   mh,
		isClosed:    false,
		msgChan:     make(chan string, 1),
		msgBuffChan: make(chan string, utils.GlobalObject.MaxMsgChanLen),
		ExitChan:    make(chan bool, 1),
		property:    make(map[string]interface{}),
		messageType: websocket.TextMessage, //默认文本协议
	}

	//将当前连接放入connmgr
	c.WsServer.GetConnMgr().Add(c)

	return c
}
