package server

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type Server struct {
	name    string
	rcvr    reflect.Value
	typ     reflect.Type
	methods map[string]*Method
}

type Method struct {
	method reflect.Method
	json   bool
}

func NewServer() *Server {
	server := new(Server)
	server.methods = make(map[string]*Method)
	return server
}

func (this *Server) Start(port string) error {
	return http.ListenAndServe(port, this)
}

func (this *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for mname, mmethod := range this.methods {
		if strings.ToLower("/"+this.name+"."+mname) == r.URL.Path {
			if method.json {
				returnValues := mmethod.method.Func.Call(
					[]reflect.Value{this.rcvr, reflect.ValueOf(w), reflect.ValueOf(r)})
				content := returnValues[0].Interface()
				if content != nil {
					w.WriteHeader(500)
					//...
				}
			} else {
				mmethod.method.Func.Call([]reflect.Value{this.rcvr, reflect.ValueOf(w), reflect.ValueOf(r)})
			}
		}
	}
}

/*
  func(this *Hello) JsonHello(r *http.Request){}
  func(this *Hello) Hello(w http.ResponseWriter, r *http.Request){}
*/

func (this *Server) Register(rcvr interface{}) error {
	this.tpy = reflect.TypeOf(rcvr)
	this.rcvr = reflect.ValueOf(rcvr)
	this.name = reflect.Indirect(this.rcvr).Type().Name()
	if this.name == "" {
		return fmt.Errorf("no service name for type ", this.typ.String())
	}
	for m := 0; m < this.typ.NumMethod(); m++ {
		method := this.typ.Method(m)
		mtype := method.Type
		mname := method.Name
		if strings.HasPrefix(mname, "Json") {
			if mtype.NumIn() != 2 {
				return fmt.Errorf("method %s has wrong number of ins: %d", mname, mtype.NumIn())
			}
			arg := mtype.In(1)
			if arg.String() != "*http.Request" {
				return fmt.Errorf("%s argument type not exported: %s", mname, arg)
			}
			this.methods[mname] = &Method{method, true}
		} else {
			if mtype.NumIn() != 3 {
				return fmt.Errorf("method %s has wrong number of ins: %d", mname, mtype.NumIn())
			}
			reply := mtype.In(1)
			if reply.String() != "http.ResponseWriter" {
				return fmt.Errorf("%s argument type not exported: %s", mname, reply)
			}
			arg := mtype.In(2)
			if arg.String() != "*http.Request" {
				return fmt.Errorf("%s argument type not exported: %s", mname, arg)
			}
			this.methods[mname] = &Method{method, false}

		}
	}
	return nil
}
