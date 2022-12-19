package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"zinx_ws/ziface"
)

// 存储一些数据供全局使用
type GlobalObj struct {
	//当前zinx实例对象
	WsServer ziface.IServer
	//连接地址
	Host string
	//端口
	Port string
	//服务器名字
	Name string
	//当前Zinx版本号
	Version string
	//当前数据包最大值
	MaxPacketSize uint32
	//允许最大连接人数
	MaxConn int
}

var GlobalObject *GlobalObj

// 重装加载配置
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		log.Println("globalobj reload ReadFile err:", err)
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		log.Println("globalobj reload Unmarshal err:", err)
		return
	}
}

// 初始化
func init() {
	GlobalObject = &GlobalObj{
		Name:          "zinx websocket",
		Version:       "V0.4",
		Host:          "0.0.0.0",
		Port:          "7777",
		MaxConn:       12000,
		MaxPacketSize: 4096,
	}

	//从conf加载数据
	GlobalObject.Reload()
}
