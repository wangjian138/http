package main

import (
	"fmt"
	"net"
)

// UDP客户端配置
func main() {
	//1:连接服务器
	coon, err := net.Dial("udp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("连接失败，err:", err)
		return
	}
	defer coon.Close()
	//发送数据
	_, err = coon.Write([]byte("hello"))
	if err != nil {
		fmt.Println("发送信息失败，err:", err)
		return
	}
	//接收信息
	var buf [1024]byte
	n, err := coon.Read(buf[:])
	if err != nil {
		fmt.Println("接收信息失败，err:", err)
	}
	fmt.Println("接收消息：", string(buf[:n]))
}
