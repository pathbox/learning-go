package main

import (
	"log"
	"os"
	"time"
)

func main() {
	err := os.Chmod("files/empty_new.txt", 0777)
	if err != nil {
		log.Println(err)
	}

	err = os.Chown("files/empty_new.txt", os.Getuid(), os.Getegid())
	if err != nil {
		log.Panicln(err)
	}

	// Change timestamps
	twoDaysFromNow := time.Now().Add(48 * time.Hour)
	lastAccessTime := twoDaysFromNow
	lastModifyTime := twoDaysFromNow
	err = os.Chtimes("files/empty_new.txt", lastAccessTime, lastModifyTime)
	if err != nil {
		log.Println(err)
	}
}
