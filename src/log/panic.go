package main

import (
	"fmt"
	"log"
)

func main() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("Just comming recover")
			fmt.Println("e from recover is :", e)
			fmt.Println("After recover")
		}
	}()
	arr := []int{2, 3}
	log.Panic("Print array ", arr, "\n")
}
