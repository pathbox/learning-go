package main

import (
	"sync"
	"fmt"
	"reflect"
)
type worker struct {
	Func interface{}
	Args []reflect.Value
}

func main() {
	var wg sync.WaitGroup

	channels := make(chan worker, 10)

	for i := 0; i < 5; i++{
		wg.Add(1)
		go func(){
			defer wg.Done()

			for ch := range channels {
				reflect.ValueOf(ch.Func).Call(ch.Args)
			}
		}()
	}

	for i := 0; i < 100; i++{
		wk := worker {
			Func: func(x,y int){
				fmt.Println(x+y)
			},
			Args: []reflect.Value{reflect.ValueOf(i), reflect.ValueOf(i)},
		}
		channels <- wk
	}
	close(channels)
	wg.Wait()
}