package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
)

type bankOp struct { // bank operation: deposit or withdraw
	howMuch int      // amount
	confirm chan int // confirmation channel
}

var accountBalance = 0        // shared account
var bankRequests chan *bankOp // channel to banker

// For now a no-op, but could save balance to a file with a timestamp.
func logBalance(current int) {}

func reportAndExit(msg string) {
	fmt.Println(msg)
	os.Exit(-1) // all 1s in binary
}
func updateBalance(amt int) int {
	// fmt.Println(amt) // 1 or -1
	update := &bankOp{howMuch: amt, confirm: make(chan int)}
	bankRequests <- update
	newBalance := <-update.confirm // 这里阻塞，直到request.confirm <- accountBalance 有数据传入到update.confirm。这样返回的就是不断被操作的accountBalance值
	// fmt.Println(newBalance)
	return newBalance
}

func main() {
	iterations := 10000

	bankRequests = make(chan *bankOp, 8) // 8 is channel buffer size

	var wg sync.WaitGroup
	// The banker: handles all requests for deposits and withdrawals through a channel.
	go func() {
		for {
			/* The select construct is non-blocking:
			   -- if there's something to read from a channel, do so
			   -- otherwise, fall through to the next case, if any */
			select {
			case request := <-bankRequests:
				accountBalance += request.howMuch // update account
				request.confirm <- accountBalance // confirm with current balance
			}
		}
	}()

	// miser increments the balance
	wg.Add(1) // increment WaitGroup counter
	go func() {
		defer wg.Done() // invoke Done on the WaitGroup when finished
		for i := 0; i < iterations; i++ {
			newBalance := updateBalance(1)
			logBalance(newBalance)
			runtime.Gosched() // yield to another goroutine
		}
	}()

	//spendthrift decrements the balance
	wg.Add(1) // increment WaitGroup counter
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			newBalance := updateBalance(-1)
			logBalance(newBalance)
			runtime.Gosched() // be nice--yield
		}
	}()

	wg.Wait()                                      // await completion of miser and spendthrift
	fmt.Println("Final balance: ", accountBalance) // confirm the balance is zero
}

// 实际上用了两个channel: bankRequests 和 confirm

/*

go func() {
  for {
    // The select construct is non-blocking:
        -- if there's something to read from a channel, do so
        -- otherwise, fall through to the next case, if any
    select {
    case request := <-bankRequests:
      accountBalance += request.howMuch // update account
      request.confirm <- accountBalance // confirm with current balance
    }
  }
}()
这个goroutine在不断的非阻塞的 等待 request := <-bankRequests
然后对accountBalance进行add操作，之后将accountBalance 传给request.confirm

而其他的update goroutine(可以很多),是将 bankRequests <- update 传给bankRequests，整个结构相当于是：多个生产者将消息传入到队列bankRequests，一个消费者按顺序接收消息并对accountBalance进行处理，然后返回给消费者，真正修改accountBalance只有一个goroutine并且是顺序执行的，这样就实现了无锁的并发模式，对于生产者来说，他们是并发执行的

*/
