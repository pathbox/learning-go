package main

import (
	"fmt"
	"net/http"
)

type Context int

type ServiceA struct {
	ctx Context
	a   int
}

func (sa *ServiceA) Foo(w http.ResponseWriter, r *http.Request) {
	sa.a = 1
	fmt.Println(sa.a)
}

func (sa *ServiceA) Bar(w http.ResponseWriter, r *http.Request) {
	sa.a = 2
	fmt.Println(sa.a)
}

func (sa *ServiceA) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Println(sa.a)
}

func main() {
	a := &ServiceA{
		ctx: 11,
		a:   0,
	}

	http.HandleFunc("/foo", a.Foo)
	http.HandleFunc("/bar", a.Bar)
	http.HandleFunc("/get", a.Get)

	http.ListenAndServe(":9090", nil)
}
