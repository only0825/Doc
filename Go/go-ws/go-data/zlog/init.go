package zlog

import (
	"io"
	"log"
	"os"
	"time"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

func init() {
	//日志输出文件
	path := "./log/"
	err := os.MkdirAll(path, 0777)
	if err != nil {
		panic(err)
	}
	file, err := os.OpenFile(path+time.Now().Format("20060102")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0766)
	if err != nil {
		log.Fatalln("Faild to open error zlog file:", err)
	}
	//自定义日志格式
	Info = log.New(io.MultiWriter(file, os.Stderr), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(file, os.Stderr), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}