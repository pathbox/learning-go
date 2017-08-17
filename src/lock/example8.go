// A WaitGroup must not be copied after first use.

package main

import (
	"log"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(wg sync.WaitGroup, i int) {
			log.Printf("i: %d", i)
			wg.Done()
		}(wg, i)
	}
	wg.Wait()
	log.Println("exit")
}

/*

Result:
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [semacquire]:
sync.runtime_Semacquire(0xc42000e33c)
  /usr/local/go/src/runtime/sema.go:47 +0x34
sync.(*WaitGroup).Wait(0xc42000e330)
  /usr/local/go/src/sync/waitgroup.go:131 +0x7a
main.main()
  /Users/pathbox/code/learning-go/src/lock/example8.go:19 +0xab
exit status 2

因为 wg 给拷贝传递到了 goroutine 中，导致只有 Add 操作，其实 Done操作是在 wg 的副本执行的。因此 Wait 就死锁了。
go 中四种引用类型有 slice， channel， function， map
*/
