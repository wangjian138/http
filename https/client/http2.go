package main

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"learn/http/go/crypto/tls"
	"learn/http/go/net/http"
	"learn/http/go/net/http2"
	"log"
	"time"
)

func main() {
	clientCertFile := "/Users/wangjian01/Documents/go/learn/http/https/client/client.pem"
	clientKeyFile := "/Users/wangjian01/Documents/go/learn/http/https/client/client.key"
	caCertFile := "/Users/wangjian01/Documents/go/learn/http/https/CA/ca.pem"
	var cert tls.Certificate
	var err error
	if clientCertFile != "" && clientKeyFile != "" {
		cert, err = tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
		if err != nil {
			log.Fatalf("Error creating x509 keypair from client cert file %s and client key file %s", clientCertFile, clientKeyFile)
		}
	}
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		fmt.Printf("Error opening cert file %s, Error: %s", caCertFile, err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	/*http 1.1
	  t := &http.Transport{
	      TLSClientConfig: &tls.Config{
	          Certificates: []tls.Certificate{cert},
	          RootCAs:      caCertPool,
	      },
	  }
	*/
	t := &http2.Transport{
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: true,
		},
	}

	client := http.Client{Transport: t, Timeout: 15 * time.Second}
	resp, err := client.Get("https://localhost:8084/")
	if err != nil {
		fmt.Printf("Failed get: %s\r\n", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed reading response body: %s\r\n", err)
	}
	fmt.Printf("Client Got response %d: %s %s\r\n", resp.StatusCode, resp.Proto, string(body))
}

func getTLSConfig(host, caCertFile string, certOpt tls.ClientAuthType) *tls.Config {
	var caCert []byte
	var err error
	var caCertPool *x509.CertPool
	if certOpt > tls.RequestClientCert {
		caCert, err = ioutil.ReadFile(caCertFile)
		if err != nil {
			fmt.Printf("Error opening cert file %s error: %v", caCertFile, err)
		}
		caCertPool = x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
	}

	return &tls.Config{
		ServerName: host,
		ClientAuth: certOpt,
		ClientCAs:  caCertPool,
		MinVersion: tls.VersionTLS12, // TLS versions below 1.2 are considered insecure - see https://www.rfc-editor.org/rfc/rfc7525.txt for details
	}
}
