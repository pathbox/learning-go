package main

import (
	"fmt"
	"unsafe"

	"github.com/tidwall/fastlane"
)

type MyType struct {
	Hiya string
}

type MyChan struct {
	base fastlane.ChanPointer
}

func main() {
	var ch MyChan
	go func() {
		ch.Send(&MyType{Hiya: "howdy!"})
	}()

	v := ch.Recv()
	fmt.Println(v.Hiya)
}

func (ch *MyChan) Send(value *MyType) {
	ch.base.Send(unsafe.Pointer(value))
}

func (ch *MyChan) Recv() *MyType {
	return (*MyType)(ch.base.Recv())
}
