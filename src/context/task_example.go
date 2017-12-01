package main

import (
	"context"
	"fmt"
	"time"
)

var begin = time.Now()

func log(text string) {
	fmt.Println(time.Since(begin), "\t", text)
}

func mark(task string) func() {
	log(" ==> " + task + "...")
	return func() {
		log("<==  " + task + "...")
	}
}

func task1(ctx context.Context) {
	defer mark("task_1")()

	task2(ctx)
}

func task2(ctx context.Context) {
	defer mark("task_2")()

	ctx2, _ := context.WithCancel(ctx)
	// defer cancel()

	task3(ctx2)
}

func task3(ctx context.Context) {
	defer mark("task_3")()

	ctx3, _ := context.WithTimeout(ctx, time.Second*5)
	// defer cancel()

	go task4(ctx3)
	go task5(ctx3)

	select {
	case <-ctx3.Done():
		log("task3: <-context.Done()")
	}
}

func task4(ctx context.Context) {
	defer mark("task_4")()
	//  Create a context with 3 seconds timeout
	ctx4, _ := context.WithTimeout(ctx, time.Second*3)
	// defer cancel()

	//  wait until be canceled
	select {
	case <-ctx4.Done():
		log("task4: <- context.Done()")
	}
}

func task5(ctx context.Context) {
	defer mark("task_5")()
	//  Create context with 5 seconds timeout
	ctx5, _ := context.WithTimeout(ctx, time.Second*6)
	// defer cancel()

	//  Call following tasks
	go task6(ctx5)

	//  wait until be canceled
	select {
	case <-ctx5.Done():
		log("task5: <- context.Done()")
	}
}

func task6(ctx context.Context) {
	defer mark("task_6")()
	//  Create a context with a value in it
	ctx6 := context.WithValue(ctx, "userID", 12)

	//  wait until be canceled
	select {
	case <-ctx6.Done():
		log("task6: <- context.Done()")
	}
}

func main() {
	task1(context.Background())
}

/*
2.547µs           ==> task_1...
291.085µs         ==> task_2...
318.259µs         ==> task_3...
391.607µs         ==> task_5...
461.169µs         ==> task_6...
510.938µs         ==> task_4...
3.000740662s     task4: <- context.Done()
3.000849859s     <==  task_4...
5.000677529s     task5: <- context.Done()
5.000800435s     <==  task_5...
5.000759311s     task6: <- context.Done()
5.000839787s     <==  task_6...
5.000732088s     task3: <-context.Done()
5.000890002s     <==  task_3...
5.000906884s     <==  task_2...
5.000921517s     <==  task_1...


 context.WithTimeout 达到超时时间，后代ctx退出
 比如 将ctx3, _ := context.WithTimeout(ctx, time.Second*1)
 ctx3 1秒超时退出，则导致ctx4 ctx5 ctx6 即使没有到超时时间，也触发退出
 因为他们是ctx3 的后代
*/
