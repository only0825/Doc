package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var m = make(map[string]interface{})

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}

	defer c.Close()

	// The event loop
	for {
		// 接收消息
		mt, message, err := c.ReadMessage()
		if err != nil {
			loger.Println("read:", err)
			break
		}
		loger.Printf("recv ms: %s", message)

		// 将json消息转为数组
		err = json.Unmarshal(message, &m)
		if err != nil {
			loger.Println("err :", err)
			break
		}

		if m["type"] == "ping" {
			//loger.Printf("recv: %s", message)
			message = []byte("pong")
			fmt.Println(mt)
			err = c.WriteMessage(mt, message)
			if err != nil {
				loger.Println("write:", err)
				break
			}
		}

		time.Sleep(time.Duration(3) * time.Second)
		res, err := http.Get(url)
		if err != nil {
			loger.Println("Request failed:", err)
			break
		}
		defer res.Body.Close()
		body, err := os.ReadAll(res.Body)
		if err != nil {
			loger.Println("Read body failed:", err)
			break
		}
		err = c.WriteMessage(1, body)
		if err != nil {
			loger.Println("Error during message writing:", err)
			break
		}
		loger.Println("Send message success")
	}
}
