package main

import (
	"encoding/json"
	"net/http"
  "fmt"
	"github.com/dpapathanasiou/go-api"
)

type Message struct {
	Text string
}

// The helloWorldJSON function accepts an http.ResponseWriter and
// http.Request object as input. For this simple example, which returns a
// greeting in JSON format regardless of the request parameters, neither
// input objects are used, but for more complex servers, the http.Request
// object has several attributes which help inform what the exact reply
// will be (see http://golang.org/pkg/net/http/#Request for the full list
// of attributes). Similarly, the http.ResponseWriter object can be used
// to write additional headers to the reply, beyond the Content-type and
// Content-length values provided automatically by the api package.

func helloWorldJSON(w http.ResponseWriter, r *http.Request) string {
  m := Message{"hello world"}
  b, err := json.Marshal(m)
  if err != nil {
    panic(err)
  }
  fmt.Println(b)
  fmt.Println(string(b))
  return string(b)
}

func main(){
  handlers := map[string]func(http.ResponseWriter, *http.Request){}
  handlers["/hello"] = func(w http.ResponseWriter, r *http.Request){
    api.Respond("application/json", "utf-8", helloWorldJSON)(w, r)
  }

  api.NewLocalServer(api.DefaultServerTransport, 9001, api.DefaultServerReadTimeout, false, handlers)
}

// To run the api server on a specific IP address, e.g., 192.168.1.1, use NewServer() instead:
	//api.NewServer("192.168.1.1", api.DefaultServerTransport, 9001, api.DefaultServerReadTimeout, false, handlers)

	// Another set of options are the transport layer (default is TCP) and FastCGI

	// To run the api server as a UDP server on a specific IP address or domain, change the transport:
	//api.NewServer("192.168.1.1", "udp", 9001, api.DefaultServerReadTimeout, false, handlers)
	// The same example, but with FastCGI:
	//api.NewServer("192.168.1.1", "udp", 9001, api.DefaultServerReadTimeout, true, handlers)
