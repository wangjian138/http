// +build linux

package netpoll

import (
	"bytes"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

func TestEpollCreate(t *testing.T) {
	s, err := EpollCreate(epollConfig(t))
	if err != nil {
		t.Fatal(err)
	}
	if err = s.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestEpollAddClosed(t *testing.T) {
	s, err := EpollCreate(epollConfig(t))
	if err != nil {
		t.Fatal(err)
	}
	if err = s.Close(); err != nil {
		t.Fatal(err)
	}
	if err = s.Add(42, 0, nil); err != ErrClosed {
		t.Fatalf("Add() = %s; want %s", err, ErrClosed)
	}
}

func TestEpollDel(t *testing.T) {
	ln := RunEchoServer(t)
	defer ln.Close()

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	s, err := EpollCreate(epollConfig(t))
	if err != nil {
		t.Fatal(err)
	}

	f, err := conn.(filer).File()
	if err != nil {
		t.Fatal(err)
	}

	err = s.Add(int(f.Fd()), EPOLLIN, func(events EpollEvent) {})
	if err != nil {
		t.Fatal(err)
	}
	if err = s.Del(int(f.Fd())); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err = s.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestEpollServer(t *testing.T) {
	ep, err := EpollCreate(epollConfig(t))
	if err != nil {
		t.Fatal(err)
	}

	// Create listener on port 4444.
	ln, err := listen(4444)
	if err != nil {
		t.Fatal(err)
	}
	defer unix.Close(ln)
	println("listen(4444) ln:", ln)

	var received bytes.Buffer
	done := make(chan struct{})

	// Add listener fd to epoll instance to know when there are new incoming
	// connections.
	_ = ep.Add(ln, EPOLLIN, func(evt EpollEvent) {
		if evt&_EPOLLCLOSED != 0 {
			return
		}

		// Accept new incoming connection.
		conn, _, err := unix.Accept(ln)
		if err != nil {
			t.Fatalf("could not accept: %s", err)
		}

		println("TestEpollServer ln:", ln, " conn:", conn)

		// Socket must not block read() from it.
		_ = unix.SetNonblock(conn, true)

		// Add connection fd to epoll instance to get notifications about
		// available data.
		_ = ep.Add(conn, EPOLLIN|EPOLLET|EPOLLHUP|EPOLLRDHUP, func(evt EpollEvent) {
			// If EPOLLRDHUP is supported, it will be triggered after conn
			// close() or shutdown(). In older versions EPOLLHUP is triggered.
			if evt&_EPOLLCLOSED != 0 {
				return
			}

			var buf [128]byte
			for {
				n, _ := unix.Read(conn, buf[:])
				if n == 0 {
					close(done)
				}
				if n <= 0 {
					break
				}
				println("TestEpollServer buf[:n]:", string(buf[:n]))
				received.Write(buf[:n])
			}
		})
	})

	conn, err := dial(4444)
	println("dial(4444) conn:", conn)
	if err != nil {
		t.Fatal(err)
	}

	// Write some data bytes one by one to the conn.
	data := []byte("hello, epoll!")
	for i := 0; i < len(data); i++ {
		println("unix.Write data:", string(data[i:i+1]))
		if _, err := unix.Write(conn, data[i:i+1]); err != nil {
			t.Fatalf("could not make %d-th write (%v): %s", i, string(data[i]), err)
		}
		time.Sleep(time.Millisecond)
	}

	unix.Close(conn)
	<-done

	if err = ep.Close(); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(received.Bytes(), data) {
		t.Errorf("bytes not equal")
	}
}

func dial(port int) (conn int, err error) {
	conn, err = unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return
	}

	addr := &unix.SockaddrInet4{
		Port: port,
		Addr: [4]byte{0x7f, 0, 0, 1}, // 127.0.0.1
	}

	err = unix.Connect(conn, addr)
	if err == nil {
		err = unix.SetNonblock(conn, true)
	}

	return
}

func listen(port int) (ln int, err error) {
	ln, err = unix.Socket(unix.AF_INET, unix.O_NONBLOCK|unix.SOCK_STREAM, 0)
	if err != nil {
		return
	}

	// Need for avoid receiving EADDRINUSE error.
	// Closed listener could be in TIME_WAIT state some time.
	unix.SetsockoptInt(ln, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)

	addr := &unix.SockaddrInet4{
		Port: port,
		Addr: [4]byte{0x7f, 0, 0, 1}, // 127.0.0.1
	}

	if err = unix.Bind(ln, addr); err != nil {
		return
	}
	err = unix.Listen(ln, 4)

	return
}

func TestListen(t *testing.T) {
	ln, err := listen(4443)
	if err != nil {
		t.Fatalf("err:%v", err)
	}

	defer func() {
		err = unix.Close(ln)
		if err != nil {
			t.Errorf("unix.Close(ln) err:%v", err)
		}
	}()

	for {
		println("TestListen ", ln)
		conn, _, err := unix.Accept(ln)
		if err != nil {
			t.Errorf("TestListen could not accept: err:%v", err)
			time.Sleep(3 * time.Second)
			continue
		}
		var buf [128]byte
		for {
			n, _ := unix.Read(conn, buf[:])
			if n <= 0 {
				break
			}
			println("TestListen buf[:n]:", string(buf[:n]))
		}

		println("TestListen ln:", ln, " conn:", conn)
	}
}

func TestDial(t *testing.T) {
	conn, err := dial(4443)
	defer unix.Close(conn)
	println("dial(4444) conn:", conn)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDia1l1AndSend(t *testing.T) {
	conn, err := dial(4443)
	defer unix.Close(conn)
	println("dial(4444) conn:", conn)
	if err != nil {
		t.Fatal(err)
	}

	// Write some data bytes one by one to the conn.
	data := []byte("hello, epoll!")
	for i := 0; i < len(data); i++ {
		println("unix.Write data:", string(data[i:i+1]))
		if _, err := unix.Write(conn, data[i:i+1]); err != nil {
			t.Fatalf("could not make %d-th write (%v): %s", i, string(data[i]), err)
		}
		time.Sleep(time.Millisecond)
	}

}

// RunEchoServer starts tcp echo server.
func RunEchoServer(tb testing.TB) net.Listener {
	ln, err := net.Listen("tcp", "localhost:")
	if err != nil {
		tb.Fatal(err)
		return nil
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					// Server closed.
					return
				}

				tb.Fatal(err)
			}
			go func() {
				if _, err := io.Copy(conn, conn); err != nil && err != io.EOF {
					tb.Fatal(err)
				}
			}()
		}
	}()
	return ln
}

func epollConfig(tb testing.TB) *EpollConfig {
	return &EpollConfig{
		OnWaitError: func(err error) {
			tb.Fatal(err)
		},
	}
}
