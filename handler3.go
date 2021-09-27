package main

import (
	"fmt"
	"net/http"
)

type indexHandler1 struct {
	content string
}

func (ih *indexHandler1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, ih.content)
}
func main() {
	http.Handle("/", &indexHandler1{content: "hello world!"})
	http.ListenAndServe(":4443", nil)
}
