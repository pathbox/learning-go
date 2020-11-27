package main

import (
	"os"
	"os/exec"
)

func main() {
	c1 := exec.Command("grep", "ERROR", "/var/log/messages")
	c2 := exec.Command("wc", "-l")
	c2.Stdin, _ = c1.StdoutPipe()
	c2.Stdout = os.Stdout
	_ = c2.Start()
	_ = c1.Run()
	_ = c2.Wait()
}

/*
多条命令组合，请使用管道
将上一条命令的执行输出结果，做为下一条命令的参数。在 Shell 中可以使用管道符 | 来实现。

比如下面这条命令，统计了 message 日志中 ERROR 日志的数量。

$ grep ERROR /var/log/messages | wc -l
19
*/
