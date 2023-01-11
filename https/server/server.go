package main

import (
	"fmt"
	"net/http"
)

type textHandler struct {
	responseText string
}

func (th *textHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, th.responseText)
}

type indexHandler struct{}

func (ih *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	mux.Handle("/", &indexHandler{})

	thWelcome := &textHandler{"TextHandler !"}
	mux.Handle("/text", thWelcome)

	http.ListenAndServe(":8084", mux)
	//http.ListenAndServeTLS(":8084")
}
