package kaca

import (
	"log"
	"net/http"
)

func ServeWs(addr string, checkOrigin bool){
	go disp.run()
	if checkOrigin{
		http.HandleFunc("/ws", serveWscheckOrigin)
	} else {
		http.HandleFunc("/ws", serveWs)
	}
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

