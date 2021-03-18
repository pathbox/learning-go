package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"timeloop/timer"
)

var (
	timerCtl    *timer.TimerHeapHandler
	asyncWorker int = 10
)

func init() {
	timerCtl = timer.New(10, 1000)
}

type timerHandler struct {
}

func AddTimerTask(dueInterval int, taskId string) {
	timerCtl.AddFuncWithId(time.Duration(dueInterval)*time.Second, taskId, func() {
		fmt.Printf("taskid is %v, time Duration is %v", taskId, dueInterval)
	})
}

func (t *timerHandler) StartLoop() {
	timerCtl.StartTimerLoop(timer.MIN_TIMER) // 扫描的间隔时间 eq cpu hz/tick
}

func (t *timerHandler) Exit(sigs chan os.Signal) {
	// graceful exit
	<-sigs
	fmt.Println("timer will exit")
	timerCtl.Exit()
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	timerEntry := timerHandler{}
	timerEntry.StartLoop()
	timerEntry.Exit(sigs)
}
