package main

import (
	"log"
)

func main() {
	arr := []int{
		2,
		3,
	}

	log.Print("Print array ", arr, "\n")
	log.Println("Println array ", arr)
	log.Printf("Printf array with item [%d, %d]\n", arr[0], arr[1])
}
