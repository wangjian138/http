package main

import (
	"fmt"
	"learn/http/gnet"
)

type echoServer struct {
	*gnet.EventServer
}

func (es *echoServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	out = frame
	out = append(out, []byte("aaa")...)
	fmt.Printf("React out:%v", string(out))
	return
}

func main() {
	echo := new(echoServer)
	err := gnet.Serve(echo, "tcp://:9000", gnet.WithMulticore(true))
	fmt.Printf("err:%v", err)
}
