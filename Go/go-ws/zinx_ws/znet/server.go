package znet

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"zinx_ws/utils"
	"zinx_ws/ziface"
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
	Port string
	//协议
	Path string
	//当前Server由用户绑定的回调router,也就是Server注册的链接对应的处理业务
	Router ziface.IRouter
}

// 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router

	fmt.Println("Add Router succ! ")
}

// 连接信息
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024, //读取最大值
	WriteBufferSize: 1024, //写最大值
	//解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 全局conectionid 后续使用uuid生成
var cid uint32

// websocket回调
func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("server wsHandler upgrade err:", err)
		return
	}

	log.Println("server wsHandler a new client coming ip:", conn.RemoteAddr())
	//处理新连接业务方法
	dealConn := NewConnection(conn, cid)
	go dealConn.Start()
	cid++
}

// 启动
func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)

	//开启一个go去做服务端Linster业务
	http.HandleFunc("/"+s.Path, s.wsHandler)
	err := http.ListenAndServe(s.IP+":"+s.Port, nil)
	if err != nil {
		log.Println("server start listen error:", err)
	}
}

// 停止服务
func (s *Server) Stop() {
	log.Println("server stop name:", s.Name)

	//TODO  Server.Stop() 将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
}

// 运行服务
func (s *Server) Serve() {
	s.Start()

	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞,否则主Go退出， listenner的go将会退出
	select {}
}

/*
创建一个服务器句柄
*/
func NewServer() ziface.IServer {
	//先初始化全局配置文件
	utils.GlobalObject.Reload()

	s := &Server{
		Name:   utils.GlobalObject.Name, //从全局参数获取
		Scheme: "tcp4",
		IP:     utils.GlobalObject.Host, //从全局参数获取
		Port:   utils.GlobalObject.Port, //从全局参数获取
		Router: nil,
	}
	return s
}
