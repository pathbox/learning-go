package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {
	go PBar()
	time.Sleep(1 * time.Second)
}

func PBar() {
	go Bar()
}

func Bar() {
	fmt.Printf("我是 %s, %s 又在调用我!\n", printMyName(), printCallerName())
	MyGoID()
}
func printMyName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
func printCallerName() string {
	pc, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name()
}

func MyGoID() {
	id := GoID()
	fmt.Println("Here goroutine id is: ", id)
}

func GoID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
