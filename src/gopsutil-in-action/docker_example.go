package main

import (
	"fmt"

	"github.com/shirou/gopsutil/docker"
)

func main() {
	s, _ := docker.GetDockerStat()

	list, _ := docker.GetDockerIDList()

	fmt.Println(s)
	fmt.Println(list)
}
