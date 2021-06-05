package main

import (
	"io"
	"io/ioutil"
	"learn/http/bytebufferpool"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"syscall"
)

var epoller *epoll

func main() {
	setLimit()

	ln, err := net.Listen("tcp", ":8972")
	if err != nil {
		panic(err)
	}

	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatalf("pprof failed: %v", err)
		}
	}()

	epoller, err = MkEpoll()
	if err != nil {
		panic(err)
	}

	go start()

	for {
		conn, e := ln.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				log.Printf("accept temp err: %v", ne)
				continue
			}

			log.Printf("accept err: %v", e)
			return
		}
		log.Printf("进行处理")
		handleConn(conn)

		if err := epoller.Add(conn); err != nil {
			log.Printf("failed to add connection %v", err)
			conn.Close()
		}
	}
}

func handleConn(conn net.Conn) {
	byteBufferPool := bytebufferpool.Get()
	for {
		n, err := byteBufferPool.ReadFrom(conn)
		log.Printf("handleConn v:%v n:%v err:%v", string(byteBufferPool.B), n, err)
		if err != nil || n == 0 {
			break
		}
	}
	bytebufferpool.Put(byteBufferPool)
	io.Copy(ioutil.Discard, conn)
}

func start() {
	var buf = make([]byte, 8)
	for {
		connections, err := epoller.Wait()
		if err != nil {
			log.Printf("failed to epoll wait %v", err)
			continue
		}
		for _, conn := range connections {
			if conn == nil {
				break
			}
			if _, err := conn.Read(buf); err != nil {
				if err := epoller.Remove(conn); err != nil {
					log.Printf("failed to remove %v", err)
				}
				conn.Close()
			}
		}
	}
}

func setLimit() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	log.Printf("set cur limit: %d", rLimit.Cur)
}
