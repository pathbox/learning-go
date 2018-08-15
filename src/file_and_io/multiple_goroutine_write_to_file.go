package main

import (
	"os"
)

func main() {
	f, _ := os.OpenFile("aaa.txt", os.O_WRONLY, 7555)
	ch := make(chan struct{}, 2)
	go writeFile(f, "A", ch)
	go writeFile(f, "B", ch)
	<-ch
	<-ch
	defer f.Close()
}

func writeFile(f *os.File, c string, ch chan struct{}) {
	for i := 0; i < 5; i++ {
		_, err := f.WriteString(c)
		if err != nil {
			panic(err)
		}
	}

	ch <- struct{}{}

}

// 利用chan， 来作为FIFO bytes队列，是最好的选择

/*  会产生竞争，相当于打开了多个文件描述符进行写数据，有并发竞争问题，不会得到正确的数据
func main() {
    ch := make(chan int, 2)
    go writeFile("aaa.txt", "A", ch)
    go writeFile("aaa.txt", "B", ch)
    <- ch
    <- ch

}

func writeFile(fn string, c string, ch chan int) {
    f, _ := os.OpenFile(fn, os.O_WRONLY, 7555)
    defer f.Close()
    for i := 0; i < 5; i++ {
        _, err := f.WriteString(c)
        if err != nil {
            panic(err)
        }
    }
    ch <- 1
}
*/
