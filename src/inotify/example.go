package main

import (
	inotify "github.com/goinbox/inotify"

	"fmt"
	"path/filepath"
)

func main() {
	path := "/tmp/a.log"

	watcher, _ := inotify.NewWatcher()
	watcher.AddWatch(path, inotify.IN_ALL_EVENTS)
	watcher.AddWatch(filepath.Dir(path), inotify.IN_ALL_EVENTS)

	for i := 0; i < 5; i++ {
		events, _ := watcher.ReadEvents()
		for _, event := range events {
			if watcher.IsUnreadEvent(event) {
				fmt.Println("It is a last remaining event")
			}
			showEvent(event)
		}
	}

	watcher.Free()
	fmt.Println("Bye")
}

func showEvent(event *inotify.Event) {
	fmt.Println(event)

	if event.InIgnored() {
		fmt.Println("inotify.IN_IGNORED")
	}
	if event.InAttrib() {
		fmt.Println("inotify.IN_ATTRIB")
	}

	if event.InModify() {
		fmt.Println("inotify.IN_MODIFY")
	}

	if event.InMoveSelf() {
		fmt.Println("inotify.IN_MOVE_SELF")
	}

	if event.InMovedFrom() {
		fmt.Println("inotify.IN_MOVED_FROM")
	}

	if event.InMovedTo() {
		fmt.Println("inotify.IN_MOVED_TO")
	}

	if event.InDeleteSelf() {
		fmt.Println("inotify.IN_DELETE_SELF")
	}

	if event.InDelete() {
		fmt.Println("inotify.IN_DELETE")
	}

	if event.InCreate() {
		fmt.Println("inotify.IN_CREATE")
	}
}
