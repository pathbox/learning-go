package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt) // 把关闭通知传到channel c

	for {
		select {
		case now := <-ticker.C:
			fmt.Printf("tick %s\n", now.UTC().Format("20060102-150405.000000000"))
		case <-c: // c中有值后，就是接收到了中断通知，这样就会走这一步
			fmt.Println("Receiverd C-c - shutting down")
			return
		}
	}
}
