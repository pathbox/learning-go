package main

import (
	"fmt"
	"time"
)

const (
	//Second has 1 * 1e9 nanoseconds
	Second time.Duration = time.Second
	//Minute has 60 seconds
	Minute time.Duration = time.Minute
	//Hour has 60 minutes
	Hour time.Duration = time.Hour
	//Day has 24 hours
	Day time.Duration = time.Hour * 24
	//Week has 7 days
	Week time.Duration = Day * 7
)

func main() {
	fmt.Println(Second)
	fmt.Println(Minute)
	fmt.Println(Hour)
	fmt.Println(Day)
	fmt.Println(Week)
}

/*
1s
1m0s
1h0m0s
24h0m0s
168h0m0s
*/
