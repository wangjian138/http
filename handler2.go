package main

import (
	"fmt"
	"net/http"
	"sync"

	"go.uber.org/ratelimit"

	"github.com/golang/groupcache/lru"
)

func text(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hrllo world")
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	html := `<doctype html>
        <html>
        <head>
          <title>Hello World</title>
        </head>
        <body>
        <p>
          <a href="/welcome">Welcome</a> |  <a href="/message">Message</a>
        </p>
        </body>
</html>`
	fmt.Fprintln(w, html)
}

var (
	rateLimitAction ratelimit.Limiter
	lruTest         = lru.New(100)
	once            sync.Once
)

var pool = sync.Pool{}

func getLru(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	fmt.Println("key:", key)

	value, ok := lruTest.Get(key)
	if ok != true {
		fmt.Println("err:", ok)
		fmt.Fprintln(w, "err")
		//return
	}
	fmt.Println("value:", value)

	getRateLimit()

	times := rateLimitAction.Take()
	fmt.Println("times:", times.Unix())

}

func setLru(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	fmt.Println("key:", key)

	value := r.URL.Query().Get("value")

	//lruTest.Add(key, value)
	//
	//ratelimits := ratelimit.New(1000000000)

	fmt.Fprintln(w, value)
}

func getRateLimit() {
	once.Do(func() {
		fmt.Println("do once")
		rateLimitAction = ratelimit.New(1)
	})
}

func acquireByteBuffer() ratelimit.Limiter {

	return byteBufferPool.Get().(ratelimit.Limiter)
}

func releaseByteBuffer(b ratelimit.Limiter) {
	if b != nil {
		byteBufferPool.Put(b)
	}
}

var byteBufferPool = &sync.Pool{
	New: func() interface{} {
		fmt.Println("默认")
		return ratelimit.New(1)
	},
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(index))
	mux.HandleFunc("/text", text)
	mux.HandleFunc("/getLru", getLru)
	mux.HandleFunc("/setLru", setLru)
	mux.HandleFunc("//setLru", setLru)
	http.ListenAndServe(":8090", mux)
}
