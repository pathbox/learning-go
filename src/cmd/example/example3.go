package main

import (
    "bytes"
    "fmt"
    "log"
    "os/exec"
)

func main() {
    cmd := exec.Command("ls", "-l", "/var/log/*.log")
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout  // 标准输出
    cmd.Stderr = &stderr  // 标准错误
    err := cmd.Run()
    outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
    fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err)
    }
}

// 执行命令，并区分stdout 和 stderr