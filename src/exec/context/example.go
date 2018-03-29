package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctx, "./cmd")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = os.Stdout
	cmd.Start()

	time.Sleep(10 * time.Second)
	fmt.Println("退出程序中...", cmd.Process.Pid)
	cancel()
	cmd.Wait()
}
