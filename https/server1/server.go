package main

import (
	"fmt"
	"learn/http/go/net/http"
	"os"
	"runtime"
	"strings"
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
	_, fp, _, _ := runtime.Caller(0)
	dir := getParentDirectory(fp)

	mux.Handle("/", &indexHandler{})

	thWelcome := &textHandler{"TextHandler !"}
	mux.Handle("/text", thWelcome)

	//http.ListenAndServe(":8084", mux)
	err := http.ListenAndServeTLS(":8084", fmt.Sprintf("%s/%s", dir, "server.crt"), fmt.Sprintf("%s/%s", dir, "server.key"), mux)
	fmt.Printf("err:%v\n", err)
}

func getParentDirectory(directory string) string {
	return substr(directory, 0, strings.LastIndex(directory, string(os.PathSeparator)))
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}
