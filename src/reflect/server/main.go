package main

import (
	"./server"
	"fmt"
	"net/http"
)

type Hello struct {
}

func (this *Hello) Print(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	w.Write([]byte("print"))
	return nil
}
func (this *Hello) Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
func (this *Hello) JsonHello(r *http.Request) {
}
func main() {
	server := server.NewServer()
	fmt.Println(server.Register(new(Hello)))
	server.Start(":8080")
}
