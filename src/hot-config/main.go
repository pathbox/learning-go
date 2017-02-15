package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var gConfig string

func main() {
	quit := make(chan bool)
	readConfig()
	go signals(quit)
	go displayConfig(quit)
EXIT:
	for {
		select {
		case <-quit:
			break EXIT
		default:
		}
	}
	fmt.Println("[main()]  exit")
}

func signals(q chan bool) bool {
	sigs := make(chan os.Signal)
	defer close(sigs)
EXIT:
	for {
		signal.Notify(sigs, syscall.SIGQUIT,
			syscall.SIGTERM,
			syscall.SIGINT,
			syscall.SIGUSR1,
			syscall.SIGUSR2)

		sig := <-sigs

		switch sig {
		case syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT:
			fmt.Println("[signals()] Interrupt...")
			break EXIT
		case syscall.SIGUSR1:
			fmt.Println("[signals()] syscall.SIGUSR1...")
			updateConfig()
		case syscall.SIGUSR2:
			fmt.Println("[signals()] syscall.SIGUSR2...")
			//updateVersion()
		default:
			break EXIT
		}
	}
	q <- true
	return true
}

func readConfig() {
	gConfig = "init"
	fmt.Println("[readConfig()] ", gConfig)
}

func updateConfig() {
	gConfig = "update"
	fmt.Println("[updateConfig()] ", gConfig)
}

func displayConfig(quit chan bool) {
	for {
		select {
		case <-quit:
			fmt.Println("[displayConfig()] exit")
			return
		default:
		}
		fmt.Println("[displayConfig()] Config:", gConfig)
		time.Sleep(time.Second * 2)
	}
}
