package main

import (
  "net"
  "log"

  "github.com/zzh20/timewheel"
)

// 定义心跳包，设置心跳超时时间，处理函数
var wheelHeartbeat = timewheel.New(time.Second*1, 30, func(data interface{}) {
	c := data.(net.Conn)
	log.Printf("timeout close conn:%v", c)
	c.Close()
})

func main() {

  // 启动心跳包检查
  wheelHeartbeat.Start()

}

// 客户端连接成功
func SessionConnected() {
    wheelHeartbeat.Add(conn)
}

// 客户端连接断开
func SessionClosed() {
    wheelHeartbeat.Remove(conn))
}

// 处理客户端的心跳包
func HeartbeatHandler() {
  wheelHeartbeat.Add(conn)
}