package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	beep "github.com/faiface/beep"
	mp3 "github.com/faiface/beep/mp3"
	speaker "github.com/faiface/beep/speaker"
)

func downloadAndPlay(dir string, name string, u string, done chan bool) {
	fileName := filepath.Join(dir, name+".mp3")

	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Printf("failed to create file %s: %v\n", fileName, err)
		return
	}
	resp, err := http.Get(u)
	if err != nil {
		fmt.Printf("failed to downlaod mp3 %s: %v\n", u, err)
		return
	}
	io.Copy(f, resp.Body)
	resp.Body.Close()
	f.Close()

	play(fileName, done)
}

func play(fileName string, done chan bool) {
	defer func() {
		recover()
	}()

	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("failed to read mp3: %v\n", err)
		return
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		fmt.Printf("failed to decode mp3: %v\n", err)
		f.Close()
		os.Remove(fileName)
		return

	}
	defer f.Close()
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	selfDone := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		selfDone <- true
	})))
	<-selfDone
	time.Sleep(200 * time.Millisecond)

	done <- true // 说明play是在一个新的goroutinue中执行

}
