package main

import (
	"fmt"
	"time"

	"github.com/connectordb/golang-lru-timedcache"
)

func main() {
  tc, _ := golanglrutimedcache.NewTimedCache(20000, 2000, nil)

  var keyid int64 = 1
  tc.Set("stringkey", keyid, "value")

  time.Sleep(2 * time.Second)
  v ,ok := tc.GetByName("stringkey")
  fmt.Println(v, ok)

  time.Sleep(1 * time.Second)
	v1, ok1 := tc.GetByName("stringkey")
	fmt.Println(v1, ok1)
}
