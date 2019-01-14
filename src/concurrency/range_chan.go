package main

import "fmt"

func main() {
	chanOwner := func() <-chan int {
		results := make(chan int, 5)
		go func() {
			defer close(results) // for 都循环完了才会执行，要不会阻塞在results <-i
			for i := 0; i <= 5; i++ {
				results <- i
			}
		}()
		return results
	}

	consumer := func(results <-chan int) {
		for result := range results {
			fmt.Printf("Received: %d\n", result)
		}
		fmt.Println("Done receiving!")
	}

	results := chanOwner()
	consumer(results)
}
