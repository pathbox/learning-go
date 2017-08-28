package main

import (
	"github.com/LK4D4/trylock"
)

type LockedStruct struct {
	mu trylock.Mutex
}

func main() {
	storage := &LockedStruct{}

	if storage.mu.TryLock() {
		// do something with storage
	} else {
		// return busy or use some logic for unavailable storage
	}
}
