package main

import (
	"fmt"
)

func main() {
	{
		payloads := [][]byte{}
		slice := []byte{0x11, 0x22, 0x33, 0x44}
		payloads = append(payloads, slice)
		fmt.Println(payloads[0])
	}
	{
		payloads := []byte{}
		payloads = append(payloads, 0x11, 0x22, 0x33, 0x44)
		fmt.Println(payloads)
	}
	{
		payloads := []byte{}
		slice := []byte{0x11, 0x22, 0x33, 0x44}
		payloads = append(payloads, slice...)
		fmt.Println(payloads)
	}

}
