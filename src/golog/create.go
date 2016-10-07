// package main
//
// import (
//     "fmt"
//     "os"
// )
//
// func main() {
//     f, _ := os.OpenFile("C:\\tmp\\11.txt", os.O_WRONLY|os.O_CREATE|os.O_SYNC,
//         0755)
//     os.Stdout = f
//     os.Stderr = f
//     fmt.Println("fmt")
//     fmt.Print(make(map[int]int)[0])
// }


package main

import (
  "fmt"
  "os"
  "syscall"
)

func main() {
  logFile, _ := os.OpenFile("./log_file.log", os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0755)
  syscall.Dup2(int(logFile.Fd()), 1)
  syscall.Dup2(int(logFile.Fd()), 2)
  defer logFile.Close()
  fmt.Printf("Hello from fmt\n")
  fmt.Printf("Hello from fmt\n")
  fmt.Printf("Hello from fmt\n")
  fmt.Printf("Hello from fmt\n")
  fmt.Printf("Hello from fmt\n")
  fmt.Printf("Hello from fmt\n")
  fmt.Printf("Hello from fmt\n")
  fmt.Printf("Hello from fmt\n")
  fmt.Printf("Hello from fmt\n")
  fmt.Printf("Hello from fmt\n")
  panic("Hello from panic\n")
}
