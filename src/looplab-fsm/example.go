package main

import (
	"fmt"
	"github.com/looplab/fsm"
)

func main() {
	fsm := fsm.NewFSM(
		"closed",
		fsm.Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"}, // Name 事件， 由Src 状态转为 Dst状态
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		fsm.Callbacks{},
	)

	fmt.Println("current state: ", fsm.Current())

	err := fsm.Event("open")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("current state: ", fsm.Current())

	err = fsm.Event("close")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("current state: ", fsm.Current())
}
