package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rfyiamcool/go-ringtimer"
)

func multi() {
	tw, err := timewheel.NewTimeWheel(timewheel.SecondInterval, 60)
	if err != nil {
		panic(err.Error())
	}

	tw.Start()
	defer tw.Stop()

	var (
		incr  int32 = 0
		wg          = sync.WaitGroup{}
		start       = time.Now()

		round         = 60
		countPerRound = 20000
	)

	wg.Add(1)
	go func() {
		defer wg.Done()

		addStart := time.Now()
		for index := 1; index <= round; index++ {
			delay := index

			go func() {
				for idx := 0; idx < countPerRound; idx++ {
					// s := time.Now()
					tw.AfterFunc(time.Duration(delay)*time.Second, func() {
						atomic.AddInt32(&incr, 1)
						// d := time.Now().Sub(s)
						// fmt.Println("trigger cost", d.Seconds(), i)
					})
				}
			}()

		}
		fmt.Println("multi add time cost: ", time.Now().Sub(addStart))
	}()

	wg.Add(1)
	go func() {
		tc := 0
		for {
			fmt.Println("incr: ", incr, "cost: ", tc)
			if atomic.LoadInt32(&incr) == int32(round*countPerRound) {
				wg.Done()
				return
			}

			time.Sleep(1 * time.Second)
			tc++
		}
	}()

	wg.Wait()
	fmt.Println("finish cost: ", time.Now().Sub(start).Seconds())
}

func main() {
	multi()
}