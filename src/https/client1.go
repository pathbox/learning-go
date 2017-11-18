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

/*
   client与server进行通信时 client也要对server返回数字证书进行校验
   因为server自签证书是无效的 为了client与server正常通信
   通过设置客户端跳过证书校验
   TLSClientConfig:{&tls.Config{InsecureSkipVerify: true}
   true:跳过证书校验
*/
