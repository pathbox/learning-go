package main

import (
    "sync"
    "fmt"
    "net/http"
)

func main() {
  var wg sync.WaitGroup
    var urls = []string{
        "http://www.baidu.com/",
        "http://www.baidu.com/",
    }

    // correct
    for _, url := range urls {
        wg.Add(1)
        go func(url string) {
            defer wg.Done()
            http.Get(url)
        }(url)
    }
    wg.Wait()

    // wg := sync.WaitGroup{}
    // wrong
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(wg sync.WaitGroup, i int) {
            fmt.Printf("i:%d", i)
            wg.Done()
        }(wg, i)
    }
    wg.Wait()
    // i:1i:3i:2i:0i:4fatal error: all goroutines are asleep - deadlock!

    // goroutine 1 [semacquire]:
    // sync.runtime_Semacquire(0xc000094018)
    //         /home/keke/soft/go/src/runtime/sema.go:56 +0x39
    // sync.(*WaitGroup).Wait(0xc000094010)
    //         /home/keke/soft/go/src/sync/waitgroup.go:130 +0x64
    // main.main()
    //         /home/keke/go/Test/wait.go:17 +0xab
    // exit status 2
    fmt.Println("exit")
}
