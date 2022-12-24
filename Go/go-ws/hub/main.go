package main

import (
	"fmt"
	"hub/server"
	"log"
	"net/http"
	"os"
	"time"
)

var logger *log.Logger

func main() {
	//后台启动处理器
	go server.Hubb.Run()
	http.HandleFunc("/score", wsHandle) // 将chat请求交给wsHandle处理
	http.ListenAndServe(":6031", nil)   // 开始监听
}

func init() {
	file := "./" + time.Now().Format("20060102") + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	logger = log.New(logFile, "[MESSAGE]", log.LstdFlags|log.Lshortfile|log.LUTC) // 将文件设置为loger作为输出
	return
}

func wsHandle(w http.ResponseWriter, r *http.Request) {
	// 通过升级后的升级起的到链接
	conn, err := server.Up.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("获取链接失败：", err)
		return
	}
	// 连接成功后注册用户
	user := &server.User{
		Conn: conn,
		Msg:  make(chan []byte),
	}
	server.Hubb.Register <- user
	defer func() {
		server.Hubb.Unregister <- user
	}()
	// 得到连接后，就可以开始读写数据了
	go read()
	write(user)
}

func read() {
	// 从连接中循环读取信息
	for {
		//_, msg, err := user.conn.ReadMessage()
		//if err != nil {
		//	fmt.Println("用户退出:", user.conn.RemoteAddr().String())
		//	hub.unregister <- user
		//	break
		//}
		msg := "hahaah"
		// 将读取到的信息传入websocket处理器中的broadcast中，
		server.Hubb.Broadcast <- []byte(msg)
		fmt.Println(msg)
		time.Sleep(time.Duration(5) * time.Second)
	}
}

func write(user *server.User) {
	for data := range user.Msg {
		err := user.Conn.WriteMessage(1, data)
		if err != nil {
			fmt.Println("写入错误")
			break
		}
	}
}
