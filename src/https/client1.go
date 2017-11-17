package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	pool := x509.NewCertPool()
	caCertPath := "./ca_key/ca.crt"

	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return
	}

	pool.AppendCertsFromPEM(caCrt) // 客户端添加ca证书

	transport := &http.Transport{
		TLSClientConfig:    &tls.Config{RootCAs: pool}, // 客户端加载ca证书
		DisableCompression: true,
	}

	client := &http.Client{Transport: transport}
	resp, err := client.Get("https://wsecho.com:9099")

	if err != nil {
		fmt.Println("Get error:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("response: ", string(body))
}
