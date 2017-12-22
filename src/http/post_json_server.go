//client-get
package main

import (
	"fmt"
	"net/url"
	"net/http"
	"io/ioutil"
	"log"
)

func main() {
	u, _ := url.Parse("http://localhost:9001/xiaoyue")
	q := u.Query()
	q.Set("username", "user")
	q.Set("password", "passwd")
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String());
	if err != nil {
				log.Fatal(err) return
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
				log.Fatal(err) return
	}
	fmt.Printf("%s", result)
}

// client-post
package main

import (
	"fmt"
	"net/url"
	"net/http"
	"io/ioutil"
	"log"
	"bytes"
	"encoding/json"
)

type Server struct {
	ServerName string
	ServerIP   string
}

type Serverslice struct {
	Servers []Server
	ServersID  string
}


func main() {

	var s Serverslice

	var newServer Server;
	newServer.ServerName = "Guangzhou_VPN";
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
	res,err := http.Post("http://localhost:9001/xiaoyue", "application/json;charset=utf-8", body)
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

// json Server
package main

import (
	"fmt"
	"net/http"
	"strings"
	"html"
	"io/ioutil"
	"encoding/json"
)

type Server struct {
	ServerName string
	ServerIP   string
}

type Serverslice struct {
	Servers []Server
	ServersID  string
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":9001", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数，默认是不会解析的
	fmt.Fprintf(w, "Hi, I love you %s", html.EscapeString(r.URL.Path[1:]))
	if r.Method == "GET" {
		fmt.Println("method:", r.Method) //获取请求的方法

		fmt.Println("username", r.Form["username"])
		fmt.Println("password", r.Form["password"])

		for k, v := range r.Form {
						fmt.Print("key:", k, "; ")
						fmt.Println("val:", strings.Join(v, ""))
		}
	} else if r.Method == "POST" {
		result, _:= ioutil.ReadAll(r.Body)
		r.Body.Close()
		fmt.Printf("%s\n", result)

		//未知类型的推荐处理方法
		var f interface{}
		json.Unmarshal(result, &f)
		m := f.(map[string]interface{})
		for k, v := range m {
			switch vv := v.(type) {
				case string:
					fmt.Println(k, "is string", vv)
				case int:
					fmt.Println(k, "is int", vv)
				case float64:
					fmt.Println(k,"is float64",vv)
				case []interface{}:
					fmt.Println(k, "is an array:")
					for i, u := range vv {
						fmt.Println(i, u)
					}
				default:
					fmt.Println(k, "is of a type I don't know how to handle")
				}
			}

			//结构已知，解析到结构体

			var s Serverslice;
			json.Unmarshal([]byte(result), &s)

			fmt.Println(s.ServersID);

			for i:=0; i<len(s.Servers); i++ {
			fmt.Println(s.Servers[i].ServerName)
			fmt.Println(s.Servers[i].ServerIP)
			}
	}
}
