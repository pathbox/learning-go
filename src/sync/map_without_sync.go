package main

import (
  "fmt"
)

var m = map[string]string{}

func main() {

  for i := 0; i < 100; i++ {
    fmt.Println(i)
    for j := 0; j < 100000; j++ {
      go func() {
        m["no"] = "1"
        fmt.Println("no", m["no"])
      }()
    }
  }
  // m["no"] = "1"
  // fmt.Println("no", m["no"])
}

/*
所有goroutine都崩
goroutine 8668 [runnable]:
main.main.func1()
  /home/user/temp_test/map_go_sync.go:14
created by main.main
  /home/user/temp_test/map_go_sync.go:14 +0x42

goroutine 8669 [runnable]:
main.main.func1()
  /home/user/temp_test/map_go_sync.go:14
created by main.main
  /home/user/temp_test/map_go_sync.go:14 +0x42

goroutine 8670 [runnable]:
main.main.func1()
  /home/user/temp_test/map_go_sync.go:14
created by main.main
  /home/user/temp_test/map_go_sync.go:14 +0x42

goroutine 8671 [runnable]:
main.main.func1()
  /home/user/temp_test/map_go_sync.go:14
created by main.main
  /home/user/temp_test/map_go_sync.go:14 +0x42

goroutine 8672 [runnable]:
main.main.func1()
  /home/user/temp_test/map_go_sync.go:14
created by main.main

*/


// package main

// import (
//  "fmt"
//  "sync"
//  "time"
// )

// var mutex sync.Mutex

// func main() {
//  c := make(map[string]int)

//  for i := 0; i < 100; i++ {
//    go func() {
//      for j := 0; j < 1000000; j++ {
//        // mutex.Lock()
//        // defer mutex.Unlock()
//        c[fmt.Sprintf("%d", j)] = j
//        fmt.Println(c[fmt.Sprintf("%d", j)])
//      }
//    }()
//  }

//  time.Sleep(100 * time.Second)
// }
