package main

import (
	"flag"
	"fmt"
)

func main() {
	dataPath := flag.String("D", "/home/manu/sample", "DB data path")
	logFile := flag.String("LOG", "/home/manu/sample.log", "log file")
	nowaitFlag := flag.Bool("W", false, "do not wait untile operation completes")

	flag.Parse()
	var cmd string = flag.Arg(0)

	fmt.Printf("action   : %s\n", cmd)
	fmt.Printf("data path: %s\n", *dataPath)
	fmt.Printf("log file : %s\n", *logFile)
	fmt.Printf("nowait     : %v\n", *nowaitFlag)

	fmt.Printf("-------------------------------------------------------\n")

	fmt.Printf("there are %d non-flag input param\n", flag.NArg())
	for i, param := range flag.Args() {
		fmt.Printf("#%d :%s\n", i, param)
	}
}
