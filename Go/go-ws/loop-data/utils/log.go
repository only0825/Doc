package utils

import (
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"os"
)

func LogInit(fileName string) {
	//日志输出文件
	path := "./log"
	_, err := os.Stat(path) // err为nil说明目录存在
	if err != nil {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			logrus.Error("创建日志目录错误", err)
		}
	}
	// 设置在输出日志中添加文件名和方法信息
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(&lumberjack.Logger{
		Filename:   path + fileName,
		MaxSize:    10, // megabytes
		MaxBackups: 20,
		MaxAge:     28,    //days
		Compress:   false, // disabled by default
	})
}
