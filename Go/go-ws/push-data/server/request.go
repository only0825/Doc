package server

import (
	"push-data/iserver"
)

type Request struct {
	//当前用户连接
	conn iserver.IConnection
	//消息封装
	msg iserver.IMessage
}

// 获取请求连接信息
func (r *Request) GetConnection() iserver.IConnection {
	return r.conn
}

// 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// 获取请求的消息的ID
func (r *Request) GetMsgType() int {
	return r.msg.GetMsgType()
}

// 获取请求的消息的ID
func (r *Request) GetMsgID() string {
	return r.msg.GetMsgID()
}
