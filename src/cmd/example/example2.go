package main

import (
"fmt"
"log"
"os/exec"
)

func main() {
    cmd := exec.Command("ls", "-l", "/var/log/")
    out, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Printf("combined out:\n%s\n", string(out))
        log.Fatalf("cmd.Run() failed with %s\n", err)
    }
    fmt.Printf("combined out:\n%s\n", string(out))
}
// 执行命令，并获取结果