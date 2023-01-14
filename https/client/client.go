package main

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"learn/http/go/crypto/tls"
	"learn/http/go/net/http"
	"os"
	"runtime"
	"strings"
)

func main() {
	pool := x509.NewCertPool()
	_, fp, _, _ := runtime.Caller(0)
	dir := getParentDirectory(fp)
	caCertPath := fmt.Sprintf("%s/%s", dir, "client.pem")
	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return
	}

	pool.AppendCertsFromPEM(caCrt) //客户端添加ca证书

	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{RootCAs: pool, InsecureSkipVerify: true}, //客户端加载ca证书
		DisableCompression: true,
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Get("https://localhost:8084/")
	if err != nil {
		fmt.Println("client Get err:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
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
