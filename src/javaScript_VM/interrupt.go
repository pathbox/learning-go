package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/robertkrimen/otto"
)

var halt = errors.New("Stahp")

func main() {
	runUnsafe(`var abc = [];`)
	runUnsafe(`
			while(true) {
				// loop forever
			}
		`)
}

func runUnsafe(unsafe string) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		if caught := recover(); caught != nil {
			if caught == halt {
				fmt.Fprintf(os.Stderr, "Some code took to long! Stopping after: %v\n", duration)
				return
			}
			panic(caught) // Something else happened, repanic!
		}
		fmt.Fprintf(os.Stderr, "Ran code successfully: %v\n", duration)
	}()

	vm := otto.New()
	vm.Interrupt = make(chan func(), 1) // vm Interrupt 其实就是一个 chan

	go func() {
		select {
		case <-time.After(10 * time.Second):
			vm.Interrupt <- func() {
				fmt.Println("Interrupt!")
				panic(halt)
			}
		}
	}()
	vm.Run(unsafe)
}
