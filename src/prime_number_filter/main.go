// package main

// import (
// 	"fmt"
// )

// func generate(ch chan<- int, n int) {
// 	for i := 2; i <= n; i++ {
// 		ch <- i
// 	}
// }

// func filter(src <-chan int, dst chan<- int, prime int) {
// 	for i := range src { // loop over values received from src
// 		fmt.Println("i ===: ", i)
// 		if i%prime != 0 {
// 			fmt.Println(i, "---", prime)
// 			dst <- i // send i to channel dst
// 		}
// 	}
// }

// func main() {
// 	ch := make(chan int)
// 	n := 100
// 	go generate(ch, n)
// 	for {
// 		prime := <-ch
// 		fmt.Println("prime: ", prime)
// 		ch1 := make(chan int)
// 		go filter(ch, ch1, prime)
// 		ch = ch1
// 	}
// }

/* 埃拉托斯特尼素数筛算法
算法描述：先用最小的素数2去筛，把2的倍数剔除掉；下一个未筛除的数就是素数(这里是3)。再用这个素数3去筛，筛除掉3的倍数... 这样不断重复下去，直到筛完为止。
把2-n的整数放入一个数组。从2开始，筛选所有二的倍数的数，并移除数组。然后取2的后一位的数，这时候是3，取3的所有倍数并移出数组
然后取3的下一位数组，应该是5，再筛选5的所有倍数并移出数组一直到数组的最后一位数字。在数组中保留的就是2-n中的素数
*/

// 从2开始每找到一个素数就标记所有能被该素数整除的所有数。直到没有可标记的数，剩下的就都是素数。下面以找出10以内所有素数为例，借用 CSP 方式解决这个问题。

package main

import "fmt"

func Processor(seq <-chan int, wait chan struct{}, level int) {
	go func() {
		prime, ok := <-seq
		if !ok {
			close(wait)
			return
		}
		fmt.Printf("[%d]: %d\n", level, prime)
		out := make(chan int)
		Processor(out, wait, level+1)
		for num := range seq {
			if num%prime != 0 {
				out <- num
			}
		}
		close(out)
	}()
}

func main() {
	origin, wait := make(chan int), make(chan struct{})
	Processor(origin, wait, 1)
	for num := 2; num < 10; num++ {
		origin <- num
	}
	close(origin)
	<-wait
}
