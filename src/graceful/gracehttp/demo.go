package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tabalt/gracehttp"
)

func main() {
	http.HandleFunc("/sleep/", func(w http.ResponseWriter, r *http.Request) {
		duration, err := time.ParseDuration(r.FormValue("duration"))
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		time.Sleep(duration)

		fmt.Fprintf(
			w,
			"started at %s slept for %d nanoseconds from pid %d. It is new version\n",
			time.Now(),
			duration.Nanoseconds(),
			os.Getpid(),
		)
	})

	pid := os.Getpid()
	address := ":8080"
	body, err := ioutil.ReadFile("/Users/pathbox/code/learning-go/src/graceful/README.md")
	if err != nil {
		panic(err)
	}
	fmt.Printf("=body: %s=\n", string(body))
	log.Printf("process with pid %d serving %s.\n It is new", pid, address)
	err = gracehttp.ListenAndServe(address, nil)
	log.Printf("process with pid %d stoped, error: %s.\n", pid, err)
}

/*
This will output something like:

 2015/09/14 20:01:08 Serving :8080 with pid 4388.
In a second terminal start a slow HTTP request

 curl 'http://localhost:8080/sleep/?duration=20s'
In a third terminal trigger a graceful server restart (using the pid from your output):

 kill -SIGUSR2 $pid
Trigger another shorter request that finishes before the earlier request:

 curl 'http://localhost:8080/sleep/?duration=0s'
*/
