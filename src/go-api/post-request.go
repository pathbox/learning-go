package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dpapathanasiou/go-api"
)

type Message struct {
	Status string
	Data   []string
}

// The logPostData function accepts an http.ResponseWriter and http.Request
// object as input; it uses the http.Request object to confirm that the API
// request from the client is a POST, and then echoes back the variable
// name/value pairs in a JSON object (a more complex API server would
// actually do something with the POST data, of course).
func logPostData(w http.ResponseWriter, r *http.Request) string {
	// prepare the default response, in case the request is invalid
	m := Message{Status: "Sorry, there was a problem", Data: []string{}}

	// this function only responds to POST requests
	if "POST" == r.Method {
		r.ParseForm()

		// iterate over the data sent via a client POST request:
		// k = the variable name
		// v = the list of values corresponding to k

		// for this example, we're just going to echo the data
		// back as a single string message within the json object,
		// just to prove we can get all names and variables correctly

		var buffer bytes.Buffer // efficient way to concanenate strings
		var postData []string

		for k, v := range r.PostForm {
			buffer.WriteString(k)
			buffer.WriteString("=")
			buffer.WriteString(strings.Join(v, ","))

			postData = append(postData, buffer.String())
			buffer.Reset()
		}
		m = Message{Status: "ok", Data: postData}
	}

	b, err := json.Marshal(m)
	if err != nil {
		panic(err) // no, not really
	}

	return string(b)
}

func main() {
	handlers := map[string]func(http.ResponseWriter, *http.Request){}
	handlers["/logger"] = func(w http.ResponseWriter, r *http.Request) {
		api.Respond("application/json", "utf-8", logPostData)(w, r)
	}

	api.NewLocalServer(api.DefaultServerTransport, 9001, api.DefaultServerReadTimeout, false, handlers)
	// To run the api server on a specific IP address, e.g., 192.168.1.1, use NewServer() instead:
	//api.NewServer("192.168.1.1", api.DefaultServerTransport, 9001, api.DefaultServerReadTimeout, false, handlers)

	// Another set of options are the transport layer (default is TCP) and FastCGI

	// To run the api server as a UDP server on a specific IP address or domain, change the transport:
	//api.NewServer("192.168.1.1", "udp", 9001, api.DefaultServerReadTimeout, false, handlers)
	// The same example, but with FastCGI:
	//api.NewServer("192.168.1.1", "udp", 9001, api.DefaultServerReadTimeout, true, handlers)
}
