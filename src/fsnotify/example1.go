package main

import (
	"log"

	"github.com/howeyc/fsnotify"
)

func main() {
	log.Println("Start fsnotify example")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		for { // 循环监听接收
			select {
			case ev := <-watcher.Event:
				log.Println("event name: ", ev.Name)
				log.Println("event: ", ev)
			case err := <-watcher.Error:
				log.Println("error: ", err)
			}
		}
	}()

	err = watcher.Watch("/Users/pathbox") // 监听的文件目录
	if err != nil {
		log.Fatal(err)
	}

	<-done

	watcher.Close()
}

/*

原理解析： go 开一个goroutine不断接收 chan message，当目录中有文件发生写操作（增 删 改）时，原生系统会有通知信号，将这个通知信号转为 chan message信息，发送到chan中，则在接收chan的一方则可以接收到message了，相当于监听到了文件更改通知，如果没有通知chan会一直阻塞监听，但并不占CPU



For each event:

Name
IsCreate()
IsDelete()
IsModify()
IsRename()

*/
