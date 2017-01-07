package main

import (
	"./util/"
	"fmt"
)

func main() {
	fmt.Println("Let's try it!")
	util.Trace("Hello World~", "Hello Kitty")
	util.Info("Hello World~")
	util.Success("Hello World~")
	util.Warning("Hello World~")
	util.Error("Hello World~")
}
