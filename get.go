package main

import (
	"fmt"
	"net/http"

	"io/ioutil"
)

func main() {
	httpGet()
}

func httpGet() {
	resp, err := http.Get("http://127.0.0.1:8010/")
	if err != nil {
		fmt.Printf("resp err:%v\n", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	fmt.Println(string(body))
}
