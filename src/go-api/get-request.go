package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	// "strings"

	"github.com/dpapathanasiou/go-api"
)

type Message struct {
  Status string
  Data []string
}

// The logPostData function accepts an http.ResponseWriter and http.Request
// object as input; it uses the http.Request object to confirm that the
// API request from the client is a GET, and then echoes back the variable
// name/value pairs in a JSON object (a more complex API server would
// actually do something with the GET data, of course).

func logPostData(w http.ResponseWriter, r *http.Request) string{
  m := Message{Status: "Do you know how to get there", Data: []string{}}

  if "GET" == r.Method{
    r.ParseForm()

		// iterate over the data sent via a client GET request:
		// k = the variable name
		// v = the list of values corresponding to k

		// for this example, we're just going to echo the data
		// back as a single string message within the json object,
		// just to prove we can get all names and variables correctly

    var buffer bytes.Buffer
    var postData []string

    // for k, v := range r.URL.Query(){
    //   buffer.WriteString(k)
    //   buffer.WriteString("=")
    //   buffer.WriteString(strings.Join(++v, ","))
    //   postData = append(postData, buffer.String())
    //   buffer.Reset()
    // }
    buffer.WriteString("Hello, nice to meet you")
    postData = append(postData, buffer.String())
    buffer.Reset()
    m = Message{Status: "OK", Data: postData}
  }
  b, err := json.Marshal(m)
  if err != nil {
    panic(err)
  }
  return string(b)
}

func main() {
  handlers := map[string]func(http.ResponseWriter, *http.Request){}
  handlers["/logger"] = func(w http.ResponseWriter, r *http.Request){
    api.Respond("application/json", "utf-8", logPostData)(w, r)
  }

  api.NewLocalServer(api.DefaultServerTransport, 9001, api.DefaultServerReadTimeout, false, handlers)
}
