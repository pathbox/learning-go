package main

import workers "github.com/Scalingo/go-workers"

func main() {
	workers.Configure(map[string]string{
		"process": "client1",
		"server":  "localhost:6379",
	})
	workers.Enqueue("myqueue", "MyRubyWorker", []string{"hello"})
}
