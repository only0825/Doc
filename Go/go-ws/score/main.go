package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var loger *log.Logger
var url string

func main() {
	loger.Println("scoreServer running")
	fmt.Println("scoreServer running", 6031)

	http.HandleFunc("/score", socketHandler)

	err := http.ListenAndServe(":6031", nil)
	if err != nil {
		loger.Println("boot wrong", err)
	}
}

func init() {
	argc := len(os.Args)
	if argc != 2 {
		fmt.Println("Boot failed:", "没有请求地址")
		loger.Println("Boot failed:", "没有请求地址") // http://api.wuhaicj.com/api/live/live
		return
	}

	url = os.Args[1]

	file := "./" + time.Now().Format("20060102") + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	loger = log.New(logFile, "[MESSAGE]", log.LstdFlags|log.Lshortfile|log.LUTC) // 将文件设置为loger作为输出
	return
}
