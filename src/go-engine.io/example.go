package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/googollee/go-engine.io.v1"
)

func main() {
	server, _ := engineio.NewServer(nil)

	go func() {
		for {
			conn, _ := server.Accept()
			go func() {
				defer conn.Close()
				for {
					t, r, _ := conn.NextReader()
					b, _ := ioutil.ReadAll(r)
					r.Close()

					w, _ := conn.NextWriter(t)
					w.Write(b)
					w.Close()
				}
			}()
		}

	}()

	http.Handle("/engine.io/", server)
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
