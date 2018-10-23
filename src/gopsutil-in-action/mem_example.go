package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/mem"
)

func main() {
	// v, _ := mem.VirtualMemory()
	// vs, _ := mem.SwapMemory()

	// fmt.Println(v)

	// almost every return value is a struct
	for {
		v, _ := mem.VirtualMemory()
		time.Sleep(1 * time.Second)
		fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
	}
}
