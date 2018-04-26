package main

import (
	"os"
	"time"

	"github.com/schollz/progressbar"
)

func main() {
	bar := progressbar.New(1000)

	bar.Reset() //重置进度条

	for i := 0; i < 1000; i++ {
		bar.Add(1)
		//  code is here
		time.Sleep(5 * time.Millisecond)
	}

	bar.Reset() // 重置进度条
	bar.SetWriter(os.Stderr)
	for i := 0; i < 1000; i++ {
		bar.Add(1)
		//  code is here
		time.Sleep(5 * time.Millisecond)
	}
}
