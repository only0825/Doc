package server

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"push-data/configs"
	"push-data/iserver"
	"push-data/task"
	"sync/atomic"
	"time"
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
		Name:       configs.Conf.Server.Name, //从全局参数获取
		Scheme:     configs.Conf.Server.Scheme,
		IP:         configs.Conf.Server.Ip,
		Port:       configs.Conf.Server.Port,
		ConnMgr:    NewConnManager(),
		MsgHandler: NewMsgHandle(),
	}

	return s
}

// 全局conectionid 后续使用uuid生成
var cid uint32

// 启动
func (s *Server) Start(c *gin.Context) {
	logrus.Info("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)

	//开启一个go去做服务端Linster业务
	go func() {
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
			logrus.Info("wsUpgrader.Upgrade wrong", err)
			return
		}

		logrus.Info("Get conn remote addr = ", wsSocket.RemoteAddr().String())

		//// WS token 鉴权
		//token := c.Query("token")
		//if token == "" {
		//	return
		//}
		//// 解码出时间，如果时间超过2分钟，则不通过
		//fmt.Println(token)
		//decrypt, err := utils.CBCDecrypt(token, configs.Conf.WsKey)

		//3 启动server网络连接业务

		//3.2 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
		/*
			if s.ConnMgr.Len() >= configs.Conf.Server.MaxConn {
				wsSocket.Close()
				continue
			}
			**/
		//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
		dealConn := NewConnection(s, wsSocket, connId, s.MsgHandler)

		logrus.Info("Current connId:", connId)
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

	go func() {
		for {
			scoreData, err := task.Score()
			if err == redis.Nil {
				return
			}
			if err != nil {
				logrus.Error("足球-比分 推送失败！", err)
				return
			}

			s.GetConnMgr().PushAll(scoreData)
			logrus.Info("足球-比分 推送成功！！！")
		}
	}()
	go func() {
		for {
			oddsData, err := task.Odds()
			if err == redis.Nil {
				return
			}
			if err != nil {
				logrus.Error("足球-指数 推送失败", err)
				return
			}

			s.GetConnMgr().PushAll(oddsData)
			logrus.Info("足球-指数 推送成功！！！")
		}
	}()
	go func() {
		for {
			oddsData, err := task.Score()
			if err == redis.Nil {
				return
			}
			if err != nil {
				logrus.Error("足球-指数 推送失败", err)
				return
			}

			s.GetConnMgr().PushAll(oddsData)
			logrus.Info("足球-指数 推送成功！！！")
		}
	}()

	//阻塞,否则主Go退出， listenner的go将会退出
	select {}
}

// 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(msgId string, router iserver.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
}
