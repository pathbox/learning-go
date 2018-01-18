package main

import (
	"fmt"

	workers "github.com/Scalingo/go-workers"
)

func main() {
	workers.Configure(map[string]string{
		"process": "worker1",
		"server":  "localhost:6379",
	})

	workers.Process("myqueue", MyGoWorker, 10)
	workers.Run()
}

func MyGoWorker(msg *workers.Msg) {
	fmt.Println("running task", msg)
}
