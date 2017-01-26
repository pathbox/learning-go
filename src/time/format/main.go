package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Now().Format(time.Kitchen))

	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(time.Now().Format("15:04:05"))

	fmt.Println(time.Now().Unix())

	timestamp := time.Now().Unix()
	fmt.Println(timestamp)

	tm := time.Unix(timestamp, 0)
	fmt.Println(tm.Format("2006-01-02 03:04:05 PM"))
	fmt.Println(tm.Format("02/01/2006 03:04:05 PM"))

	t2, _ := time.Parse("01/02/2006", "02/08/2015")
	fmt.Println(t2)
}
