package main

import (
	"fmt"
	"learn/http/gnet"
	"learn/http/gnet/pool/goroutine"
	"time"
)

type echoServer struct {
	*gnet.EventServer
	pool *goroutine.Pool
}

func (es *echoServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	data := append([]byte{}, frame...)

	fmt.Printf("React data:%v\n", string(data))
	_ = es.pool.Submit(func() {
		time.Sleep(1 * time.Second)
		c.AsyncWrite(data)
	})

	return
}

func main() {
	p := goroutine.Default()
	defer p.Release()

	echo := &echoServer{pool: p}
	err := gnet.Serve(echo, "tcp://:9000", gnet.WithMulticore(true))
	fmt.Printf("err:%v", err)
}
