package main

import (
	"fmt"
	"net"
)

// UDP服务端配置
func main() {
	//1:启动监听
	listener, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 8888,
	})
	if err != nil {
		fmt.Println("启动server失败，err：", err)
		return
	}
	defer listener.Close()
	//获取连接数据
	for {
		var buf [1024]byte
		n, addr, err := listener.ReadFromUDP(buf[:])
		if err != nil {
			fmt.Println("读取失败,err:", err)
			return
		}
		//信息输出
		fmt.Printf("来自%v的消息：%v\n", addr, string(buf[:n]))
		//信息回复
		_, err = listener.WriteToUDP([]byte("hi"), addr)
		if err != nil {
			fmt.Println("回复信息失败，err:", err)
			return
		}
	}
}
