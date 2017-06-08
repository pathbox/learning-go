package main

import (
	"fmt"
	"os/exec"
	"time"
)

func main() {
	cmd := exec.Command("/bin/sh", "-c", "watch date > date.txt")
	start := time.Now()
	time.AfterFunc(3*time.Second, func() { cmd.Process.Kill() })
	err := cmd.Run()
	fmt.Printf("pid=%d duration=%s err=%s\n", cmd.Process.Pid, time.Since(start), err)
}

// Go是使用kill(2)向sh进程的PID发了一个KILL信号，但没有发给watch进程，sh进程被kill之后，导致watch进程变成孤儿进程。实际这是unix编程语言的一个非常正常的行为，只是...在很多场景下确实不适用。
// ps -jl
// F S   UID   PID  PPID  PGID   SID  C PRI  NI ADDR SZ WCHAN  TTY          TIME CMD
// 4 S  1000  3570  3569  3570  2519  0  80   0 - 28870 wait   pts/0    00:00:00 bash
// 0 S  1000  8754     1  8730  2519  0  80   0 - 31323 hrtime pts/0    00:00:00 watch
// 0 R  1000  8767  3570  8767  2519  0  80   0 - 30319 -      pts/0    00:00:00 ps
// 程序仍然是3s退出，/bin/sh被kill，但是残留了watch这个子进程，该子进程的PPID已经是1，即被init进程接管了
