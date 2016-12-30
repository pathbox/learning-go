package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/garyburd/redigo/redis"
)

// 重写生成连接池方法
func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   100,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

// init pool
var pool = newPool()

func redisServer(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	// 从连接池里面获得一个连接
	c := pool.Get()
	defer c.Close()
	dbkey := "netgame:info"
	if ok, err := redis.Bool(c.Do("LPUSH", dbkey, "yangzetao")); ok {

	} else {
		log.Println(err)
	}
	msg := fmt.Sprintf("time: %s", time.Now().Sub(startTime))
	io.WriteString(w, msg+"\n\n")
}

func main() {
	// 利用cpu多核来处理http请求，这个没有用go默认就是单核处理http的，这个压测过了，请一定要相信我
	runtime.GOMAXPROCS(runtime.NumCPU())
	http.HandleFunc("/", redisServer)
	http.ListenAndServe(":9090", nil)
}
