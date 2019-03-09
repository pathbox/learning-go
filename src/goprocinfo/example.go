package main

import (
	"log"

	linuxproc "github.com/c9s/goprocinfo/linux"
)

func main() {
	stat, err := linuxproc.ReadStat("/proc/pid/status")
	if err != nil {
		log.Fatal("stat read fail")
	}
	for _, info := range stat.CPUStats {
		log.Printf("User: %v\n", info.User)
		log.Printf("Nice: %v\n", info.Nice)
		log.Printf("System: %v\n", info.System)
		log.Printf("Idle: %v\n", info.Idle)
		log.Printf("IOWait: %v\n", info.IOWait)

	}
}
