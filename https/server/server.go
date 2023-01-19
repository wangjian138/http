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
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Hello"))
}

func main() {
	mux := http.NewServeMux()
	_, fp, _, _ := runtime.Caller(0)
	dir := getParentDirectory(fp)

	mux.Handle("/", &indexHandler{})

	thWelcome := &textHandler{"TextHandler !"}
	mux.Handle("/text", thWelcome)

	err := http.ListenAndServeTLS(":8084", fmt.Sprintf("%s/%s", dir, "server.pem"), fmt.Sprintf("%s/%s", dir, "server.key"), mux)
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
