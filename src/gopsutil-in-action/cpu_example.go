package main

import (
	"fmt"
	"time"

	cpu "github.com/shirou/gopsutil/cpu"
)

func main() {
	for {
		r, _ := cpu.Percent(1*time.Second, true)
		fmt.Println(r)
	}

}
