package main

import (
	"sort"
	"fmt"
)

type Peak struct {
	Name      string
	Elevation int
}

func main() {
	peaks := []Peak{
		{"Aconcagua", 22838},
		{"Denali", 20322},
		{"Puncak Jaya", 16024},
		{"Mount Elbrus", 19431},
	}

	sort.Slice(peaks, func(i, j int) bool {
		return peaks[i].Elevation >= peaks[j].Elevation
	})
	fmt.Println(peaks)
}
