package main

import (
	"fmt"
	"sync"
	"time"
)

var MLock = &sync.RWMutex{}
var Cache = make(map[int][]interface{})

type LP struct {
	L string
	P int
}

func main() {

	go func() {
		var i int
		for {
			r := GetLCache(100, 0)
			fmt.Println(r)
			i++
			fmt.Println("GetLCache:", i)
		}
	}()

	go func() {
		for {
			MLock.RLock()
			// defer MLock.RUnlock() // It create dead lock 这会导致死锁
			for k, _ := range Cache {
				fmt.Println(k)
			}
			MLock.RUnlock()
		}
	}()

	go func() {
		// var i int
		for {
			lpAry := make([]interface{}, 512)
			lpAry[0] = LP{P: 1, L: "A"}
			MLock.Lock()
			Cache[100] = lpAry
			MLock.Unlock()
		}
	}()
	time.Sleep(10 * time.Minute)
}

func GetLACache(key int) ([]interface{}, bool) {
	MLock.RLock()
	defer MLock.RUnlock() // 要在map取key的这一步加读锁,同理在写map的key-value时候，因为是map不是并发安全的
	lpAry, ok := Cache[key]
	return lpAry, ok
}

func GetLPCache(key, bit int) (LP, bool) {
	lpAry, ok := GetLACache(key)
	if !ok {
		return LP{}, false
	}
	if lpAry[bit] == nil {
		return LP{}, false
	}
	// MLock.RLock()
	// defer MLock.RUnlock() 没有用，会报并发读写错误
	return lpAry[bit].(LP), true
}

func GetLCache(key, bit int) string {
	lp, ok := GetLPCache(key, bit)
	if ok {
		return lp.L
	}
	return ""
}

func GetPCache(key, bit int) (int, bool) {
	lp, ok := GetLPCache(key, bit)
	if ok {
		return lp.P, true
	}
	return 0, false
}
