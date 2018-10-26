package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/disk"
)

func main() {
	v, _ := disk.Partitions(true)

	fmt.Println(v)

	for {
		time.Sleep(1 * time.Second)
		r, _ := disk.IOCounters("disk1s1")
		fmt.Println(r)
	}
}
