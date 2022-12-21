package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync/atomic"
	"time"
	"zinx_ws/configs"
	"zinx_ws/iserver"
)

// 服务器实现 ws://127.0.0.1:8080/echo
type Server struct {
	//服务器名称
	Name string
	//服务器协议 ws,wss
	Scheme string
	//服务器ip地址
	IP string
	//服务器端口
	Port int
	//Server的消息管理模块
	MsgHandler iserver.IMsgHandle
	//当前Server链接管理器
	ConnMgr iserver.IConnMgr
}

// 得到链接管理
func (s *Server) GetConnMgr() iserver.IConnMgr {
	return s.ConnMgr
}

var (
	GWServer iserver.IServer
)

/*
创建一个服务器句柄
*/
func NewServer() iserver.IServer {
	s := &Server{
		Name:       configs.GConf.Name, //从全局参数获取
		Scheme:     configs.GConf.Scheme,
		IP:         configs.GConf.Ip,
		Port:       configs.GConf.Port,
		ConnMgr:    NewConnManager(),
		MsgHandler: NewMsgHandle(),
	}
	return s
}

// 全局conectionid 后续使用uuid生成
var cid uint32

// 启动
func (s *Server) Start(c *gin.Context) {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)

	//开启一个go去做服务端Linster业务
	go func() {

		//TODO server.go 应该有一个自动生成ID的方法
		curConnId := uint64(time.Now().Unix())
		connId := atomic.AddUint64(&curConnId, 1)
		//3.1 阻塞等待客户端建立连接请求
		var (
			err        error
			wsSocket   *websocket.Conn
			wsUpgrader = websocket.Upgrader{
				// 允许所有CORS跨域请求
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			}
		)

		if wsSocket, err = wsUpgrader.Upgrade(c.Writer, c.Request, nil); err != nil {
			return
		}
		fmt.Println("Get conn remote addr = ", wsSocket.RemoteAddr().String())
		//3 启动server网络连接业务

		//3.2 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
		/*
			if s.ConnMgr.Len() >= configs.GConf.MaxConn {
				wsSocket.Close()
				continue
			}
			**/
		//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
		dealConn := NewConnection(s, wsSocket, connId, s.MsgHandler)

		fmt.Println("Current connId:", connId)
		//3.4 启动当前链接的处理业务
		go dealConn.Start()

	}()

	//开启一个go去做服务端Linster业务
	//http.HandleFunc("/"+s.Path, s.wsHandler)
	//err := http.ListenAndServe(s.IP+":"+s.Port, nil)
	//if err != nil {
	//	log.Println("server start listen error:", err)
	//}
}

// 停止服务
func (s *Server) Stop() {
	log.Println("server stop name:", s.Name)

	//TODO  Server.Stop() 将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
}

// 运行服务
func (s *Server) Serve(c *gin.Context) {
	s.Start(c)

	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞,否则主Go退出， listenner的go将会退出
	select {}
}

// 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(msgId string, router iserver.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
}
