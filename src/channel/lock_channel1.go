package main

import (
	"fmt"
)

func main() {
	c := make(chan int)
	c <- 60    // write channel Here is waiting the reader be build, but it is later, and it will not be build, so locking!
	val := <-c // read channl
	fmt.Println(val)
}

// 代码在次以单线程的方式运行，逐行运行。向channel写入的操作（c <- 60）会锁住整个程序的执行进程，
// 因为在同步channel中的写操作只有在读取器准备就绪后才能成功执行。然而在这里，我们在写操作的下一行才创建了读取器。
