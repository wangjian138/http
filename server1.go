package main

import (
	"fmt"
	"net/http"
	"time"
)

func index1(w http.ResponseWriter, r *http.Request) {
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

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", index1)

	server := &http.Server{
		Addr:         ":8010",
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 2 * time.Second,
		Handler:      mux,
	}
	err := server.ListenAndServe()
	fmt.Printf("err:%v\n", err)
}
