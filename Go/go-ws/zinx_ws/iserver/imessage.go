package iserver

/*
将请求的一个消息封装到message中，定义抽象层接口
*/
type IMessage interface {
	GetDataLen() uint32 //获取消息数据段长度
	GetMsgType() int    //获取消息类型
	GetData() []byte    //获取消息内容
	GetMsgID() string   //获取消息ID
}
