package main

import (
    "log"
    "os/exec"
)

func main() {
    cmd := exec.Command("ls", "-l", "/var/log/")
    err := cmd.Run()
    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err)
    }
}
// 只执行命令 不获取结果