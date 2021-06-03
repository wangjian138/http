package main

import (
	"fmt"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("listen error: ", err)
		return
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept error: ", err)
			break
		}

		// start a new goroutine to handle the new connection
		go HandleConn(conn)
	}
}
func HandleConn(conn net.Conn) {
	defer conn.Close()
	packet := make([]byte, 1024)
	for {
		// 如果没有可读数据，也就是读 buffer 为空，则阻塞
		_, _ = conn.Read(packet)
		// 同理，不可写则阻塞
		_, _ = conn.Write(packet)
	}
}
