package main

import (
	"fmt"

	"io"

	"net/http"

	"net/http/httputil"

	"net/url"

	"os"
)

func main() {

	if len(os.Args) != 3 {

		fmt.Println("Usage: ", os.Args[0], "http://proxy-host:port http://host:port/page")

		os.Exit(1)

	}

	proxyString := os.Args[1]

	proxyURL, err := url.Parse(proxyString)

	checkError(err)

	rawURL := os.Args[2]

	url, err := url.Parse(rawURL)

	checkError(err)

	transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}

	client := &http.Client{Transport: transport}

	request, err := http.NewRequest("GET", url.String(), nil)

	dump, _ := httputil.DumpRequest(request, false)

	fmt.Println(string(dump))

	response, err := client.Do(request)

	checkError(err)

	fmt.Println("Read ok")

	if response.Status != "200 OK" {

		fmt.Println(response.Status)

		os.Exit(2)

	}

	fmt.Println("Reponse ok")

	var buf [512]byte

	reader := response.Body

	for {

		n, err := reader.Read(buf[0:])

		if err != nil {

			os.Exit(0)

		}

		fmt.Print(string(buf[0:n]))

	}

	os.Exit(0)

}

func checkError(err error) {

	if err != nil {

		if err == io.EOF {

			return

		}

		fmt.Println("Fatal error ", err.Error())

		os.Exit(1)

	}

}

/*
发起者
代理服务
真正服务

发起者请求代理服务，代理服务发起请求到真正服务，代理服务将返回信息再返回给发起者。真正服务不知道有发起者，只知道代理服务的存在
*/
