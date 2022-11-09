const WebSocket = require('ws')
const wss = new WebSocket.Server({ port: 3000 })

// 事件
wss.on('connection', (ws) => {
  console.log('one people login');

  // 如果收到客户端发来的信息
  ws.on('message', (data) => {
    ws.send(data + ' hello world')
  })

  // 单一连接中进行操作
  ws.on('close', () => {
    console.log('one people logout');
  })
})