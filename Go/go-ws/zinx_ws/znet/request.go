package znet

import "zinx_ws/ziface"

type Request struct {
	//当前用户连接
	conn ziface.IConnection
	//消息封装
	message string
}

// 得到当前连接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// 得到请求数据
func (r *Request) GetData() string {
	return r.message
}
