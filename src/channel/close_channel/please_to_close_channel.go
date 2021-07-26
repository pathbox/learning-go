package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

func bug(_ http.ResponseWriter, _ *http.Request) {
	taskChan := make(chan int, 100)

	for i := 0; i < 100; i++ {
		taskChan <- i
	}

	consumer := func() {
		for task := range taskChan {
			fmt.Println(task)
		}
	}

	for i := 0; i < 100; i++ {
		go consumer()
	}
}

func bugfix(_ http.ResponseWriter, _ *http.Request) {
	taskChan := make(chan int, 100)

	for i := 0; i < 100; i++ {
		taskChan <- i
	}

	consumer := func() {
		for task := range taskChan {
			fmt.Println(task)
		}
	}

	for i := 0; i < 100; i++ {
		go consumer()
	}
	close(taskChan) // bugfix
}

func main() {
	http.HandleFunc("/bug", bug)
	http.HandleFunc("/bugfix", bugfix)

	http.ListenAndServe(":8000", nil)
}
