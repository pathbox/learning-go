package main

import (
	"fmt"
	"sync/atomic"
)

type User struct {
	FirstName string
	LastName  string
}

var GlobalUser atomic.Value

func main() {
	user := User{"Ramsay", "Bolton"}
	GlobalUser.Store(user)           // atomic/thread-safe
	data := GlobalUser.Load().(User) // atomic/thread-safe
	fmt.Printf("%+v", data)
}
