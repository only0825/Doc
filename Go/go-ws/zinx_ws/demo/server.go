package demo

import "go-ws/zinx_ws/znet"

func main() {
	//创建一个实例
	s := znet.NewServer()
	s.Serve()
}
