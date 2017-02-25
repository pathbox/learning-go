package main

import (
	"crypto/rand"
	"time"
	"github.com/grooveshark/golib/agg"
)

func main() {
  agg.CreateInterrupt(1)

  // Create 20 goroutines, each repeatedly measuring how long it takes to run
	// readRand() and aggregating that data

  for i := 0; i < 20; i++ {
    go func ()  {
      for{
        start := time.Now()
        readRand()
        agg.Agg("readRand", time.Since(start).Seconds())
      }
    }()
  }
  select{}
}

// Reads some random data, we're testing how fast this actually takes
func readRand() {
	b := make([]byte, 16)
	rand.Read(b)
}
