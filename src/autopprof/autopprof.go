package autopprof

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
	"time"
)

// Profile represents a pprof profile.
type Profile interface {
	Capture() (profile string, err error)
}

// CPUProfile captures the CPU profile.
type CPUProfile struct {
	Duration time.Duration // 30 seconds by default
}

func (p CPUProfile) Capture() (string, error) {
	dur := p.Duration
	if dur == 0 {
		dur = 30 * time.Second
	}

	f := newTemp()
	if err := pprof.StartCPUProfile(f); err != nil {
		return "", nil
	}
	time.Sleep(dur)
	pprof.StopCPUProfile()
	if err := f.Close(); err != nil {
		return "", nil
	}
	return f.Name(), nil
}

// HeapProfile captures the heap profile.
type HeapProfile struct{}

func (p HeapProfile) Capture() (string, error) {
	return captureProfile("heap")
}

// MutexProfile captures stack traces of holders of contended mutexes.
type MutexProfile struct{}

func (p MutexProfile) Capture() (string, error) {
	return captureProfile("mutex")
}

// BlockProfile captures stack traces that led to blocking on synchronization primitives.
type BlockProfile struct {
	// Rate is the fraction of goroutine blocking events that
	// are reported in the blocking profile. The profiler aims to
	// sample an average of one blocking event per rate nanoseconds spent blocked.
	//
	// If zero value is provided, it will include every blocking event
	// in the profile.
	Rate int
}

func (p BlockProfile) Capture() (string, error) {
	if p.Rate > 0 {
		runtime.SetBlockProfileRate(p.Rate)
	}
	return captureProfile("block")
}

// GoroutineProfile captures stack traces of all current goroutines.
type GoroutineProfile struct{}

func (p GoroutineProfile) Capture() (string, error) {
	return captureProfile("goroutine")
}

// Threadcreate profile captures the stack traces that led to the creation of new OS threads.
type ThreadcreateProfile struct{}

func (p ThreadcreateProfile) Capture() (string, error) {
	return captureProfile("threadcreate")
}

func captureProfile(name string) (string, error) {
	f := newTemp()
	if err := pprof.Lookup(name).WriteTo(f, 2); err != nil {
		return "", nil
	}
	if err := f.Close(); err != nil {
		return "", nil
	}
	return f.Name(), nil
}

func Capture(p Profile) {
	go capture(p)
}

func capture(p Profile) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGQUIT)
	fmt.Println("Send SIGQUIT (CTRL+\\) to the process to capture...")

	for {
		<-c
		log.Println("Starting to capture.")

		profile, err := p.Capture()
		if err != nil {
			log.Printf("Cannot capture profile: %v", err)
		}

		// Open profile with pprof.
		log.Printf("Starting go tool pprof %v", profile)
		cmd := exec.Command("go", "tool", "pprof", "-http=:", profile) // 执行go tool pprof -http=: 命令 会打开浏览器，展示pprof信息图
		if err := cmd.Run(); err != nil {
			log.Printf("Cannot start pprof UI: %v", err)
		}
	}
}

func newTemp() (f *os.File) {
	f, err := ioutil.TempFile("", "profile-") // 生成临时文件
	if err != nil {
		log.Fatalf("Cannot create new temp profile file: %v", err)
	}
	return f
}
