package main

import (
	"fmt"
	"sync"
	"time"
)

type Person struct {
	Name      string
	Age       string
	waitGroup *sync.WaitGroup // wait for all goroutines
}

func main() {
	person := Person{
		Name:      "Cary",
		Age:       "26",
		waitGroup: &sync.WaitGroup{},
	}
	person.waitGroup.Add(1)

	go foo(&person)

	fmt.Println("before wait")
	person.waitGroup.Wait()
	fmt.Println("After wait")

}

func foo(person *Person) {
	for i := 0; i < 10; i++ {
		fmt.Println("times: ", i, person.Name)
		time.Sleep(time.Second)
	}
	// 计数器-1
	person.waitGroup.Done()
}

// before wait
// times:  0 Cary
// times:  1 Cary
// times:  2 Cary
// times:  3 Cary
// times:  4 Cary
// times:  5 Cary
// times:  6 Cary
// times:  7 Cary
// times:  8 Cary
// times:  9 Cary
// After wait

// if don't Wait() result is:
// before wait
// times:  0 Cary
// After wait

// main (or any goroutinue: do person.waitGroup.Add(1) ) doesn't wait for the other goroutinue
