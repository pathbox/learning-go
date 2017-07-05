package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(3 * time.Second) // 每隔三秒往C channel中传入当前时间值
	defer ticker.Stop()

	for {
		select {
		case now := <-ticker.C:
			fmt.Printf("tick %s\n", now.UTC().Format("20060102-150405.000000000"))
		}
	}
}
