package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Server struct {
	ServerName string
	ServerIP   string
}

type Serverslice struct {
	Servers   []Server
	ServersID string
}

func main() {

	var s Serverslice

	var newServer Server
	newServer.ServerName = "Guangzhou_VPN"
	newServer.ServerIP = "127.0.0.1"
	s.Servers = append(s.Servers, newServer)

	s.Servers = append(s.Servers, Server{ServerName: "Shanghai_VPN", ServerIP: "127.0.0.2"})
	s.Servers = append(s.Servers, Server{ServerName: "Beijing_VPN", ServerIP: "127.0.0.3"})

	s.ServersID = "team1"

	b, err := json.Marshal(s)
	if err != nil {
		fmt.Println("json err:", err)
	}

	body := bytes.NewBuffer([]byte(b))
	res, err := http.Post("http://localhost:9001/xiaoyue", "application/json;charset=utf-8", body)
	if err != nil {
		log.Fatal(err)
		return
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("%s", result)
}
