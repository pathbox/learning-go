package main

import (
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

func main() {
	cmd := exec.Command("/bin/sh", "-c", "watch date > date.txt")
	// Go会将PGID设置成与PID相同的值
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	start := time.Now()
	time.AfterFunc(3*time.Second, func() { syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) })
	err := cmd.Run()
	fmt.Printf("pid=%d duration=%s err=%s\n", cmd.Process.Pid, time.Since(start), err)
}

// kill(2)不但支持向单个PID发送信号，还可以向进程组发信号，传递进程组PGID的时候要使用负数的形式。我们只要把sh进程及其所有子进程放到一个进程组里，就可以批量Kill了。关键是PGID的设置，默认情况下，子进程会把自己的PGID设置成与父进程相同，所以，我们只要设置了sh进程的PGID，所有子进程也就相应的有了PGID。

// [work@vm killproc]$ ps -jl
// F S   UID   PID  PPID  PGID   SID  C PRI  NI ADDR SZ WCHAN  TTY          TIME CMD
// 4 S  1000 17156 17155 17156 17136  0  80   0 - 28845 wait   pts/0    00:00:00 bash
// 0 R  1000 17364 17156 17364 17136  0  80   0 - 30319 -      pts/0    00:00:00 ps
// 如我们所愿，watch进程并没有残留，目标达成。
