package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var logger *log.Logger

func main() {
	//后台启动处理器
	go hub.run()
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
	conn, err := up.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("获取链接失败：", err)
		return
	}
	// 连接成功后注册用户
	user := &User{
		conn: conn,
		msg:  make(chan []byte),
	}
	hub.register <- user
	defer func() {
		hub.unregister <- user
	}()
	// 得到连接后，就可以开始读写数据了
	go read(user)
	go oddsChange(user)
	write(user)
}

func read(user *User) {
	// 从连接中循环读取信息
	for {
		_, msg, err := user.conn.ReadMessage()
		if err != nil {
			fmt.Println("用户退出:", user.conn.RemoteAddr().String())
			hub.unregister <- user
			break
		}
		// 将读取到的信息传入websocket处理器中的broadcast中，
		hub.broadcast <- msg
	}
}

func write(user *User) {
	for data := range user.msg {
		err := user.conn.WriteMessage(1, data)
		if err != nil {
			fmt.Println("写入错误")
			break
		}
	}
}

func oddsChange(user *User) {
	for {
		time.Sleep(time.Duration(3) * time.Second)
		res, err := http.Get("http://api.wuhaicj.com/api/live/live")
		if err != nil {
			logger.Println("URL Request failed:", err)
			break
		}
		msg, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.Println("Read body failed:", err)
			break
		}
		defer res.Body.Close()
		// 将读取到的信息传入websocket处理器中的broadcast中，

		err = user.conn.WriteMessage(1, msg)
		if err != nil {
			fmt.Println("写入错误", err)
			break
		}
		//hub.broadcast <- msg
	}
}
