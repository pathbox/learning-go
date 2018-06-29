package main

import (
	"fmt"
	"syscall/js"
)

var (
	no             int
	beforeUnloadCh = make(chan struct{})
)

func main() {
	callback := js.NewCallback(printMessage)
	defer callback.Close()

	setPrintMessage := js.Gloabl().Get("setPrintMessage")
	setPrintMessage.Invoke(callback)

	beforeUnloadCb := js.NewEventCallback(0, beforeUnload)
	defer beforeUnloadCb.Close()
	addEventListener := js.Global().Get("addEventListener")
	addEventListener.Invoke("beforeunload", beforeUnloadCb)

	<-beforeUnloadCh
	fmt.Println("Bye Wasm !")

}

func printMessage(args []js.Value) {
	message := args[0].String()
	no++
	fmt.Printf("Message no %d: %s\n", no, message)
}

func beforeUnload(event js.Value) {
	beforeUnloadCh <- struct{}{}
}
