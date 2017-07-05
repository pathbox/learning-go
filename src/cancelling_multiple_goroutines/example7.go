package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go tick(ctx, &wg)

	wg.Add(1)
	go tock(ctx, &wg)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("main: received C-c - shutting down")

	// tell the goroutines to stop
	fmt.Println("main: telling goroutines to stop")
	cancel()

	// and wait for them both to reply back
	wg.Wait()
	fmt.Println("main: all goroutines have told us they've finished")
}

func tick(ctx context.Context, wg *sync.WaitGroup) {
	// tell the caller we've stopped
	defer wg.Done()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case now := <-ticker.C:
			fmt.Printf("tick: tick %s\n", now.UTC().Format("20060102-150405.000000000"))
		case <-ctx.Done():
			fmt.Println("tick: caller has told us to stop")
			return
		}
	}
}

func tock(ctx context.Context, wg *sync.WaitGroup) {
	// tell the caller we've stopped
	defer wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case now := <-ticker.C:
			fmt.Printf("tock: tock %s\n", now.UTC().Format("20060102-150405.000000000"))
		case <-ctx.Done():
			fmt.Println("tock: caller has told us to stop")
			return
		}
	}
}
