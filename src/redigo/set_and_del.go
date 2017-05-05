// "redis is a good way, don't use global map store huge key/value"

package main

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strconv"
)

var redisClient, _ = redis.Dial("tcp", "localhost:6379")

func main() {
	log.Println("here start test redis")

	http.HandleFunc("/", redisOpt)
	http.HandleFunc("/delete", redisDel)

	log.Fatal(http.ListenAndServe(":9090", nil))
}

func redisOpt(w http.ResponseWriter, r *http.Request) {
	str := "aaaaaaaaaaaaaaaaaaaaaaaa"
	val := "8ff98326-2187-4de2-924e-af5098921aba"
	for i := 0; i < 10000000; i++ {
		key := str + strconv.Itoa(i)
		redisClient.Do("SET", key, val)
	}

	log.Println("get result")
	// count, err := redisClient.Do("DBSIZE") // 获得key 的总数
	// if err != nil {
	//  log.Println(err)
	// }
	// log.Println("count: ", count)
	// w.Write([]byte(count.(string)))
	w.Write([]byte("done"))

}

func redisDel(w http.ResponseWriter, r *http.Request) {
	log.Println("start delete key")
	str := "aaaaaaaaaaaaaaaaaaaaaaaa"
	for i := 0; i < 10000000; i++ {
		key := str + strconv.Itoa(i)
		redisClient.Do("DEL", key)
	}
	w.Write([]byte("done"))
	log.Println("done delete")
}
