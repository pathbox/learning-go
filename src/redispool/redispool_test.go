package redispool

import (
	"fmt"
	. "testing"

	"github.com/fzzy/radix/redis"
)

func TestPool(t *T) {
	pool, err := NewPool("tcp", "localhost:6379", "", 10)
	if err != nil {
		t.Fatal(err)
	}

	conns := make([]*redis.Client, 20)
	for i := range conns {
		if conns[i], err = pool.Get(); err != nil {
			t.Fatal(err)
		}
	}

	for i := range conns {
		conns[i].Cmd("hmset", "redispool", fmt.Sprintf("key%v", i), fmt.Sprintf("val%v", i))
		if i == 19 {
			conns[i].Cmd("expire", "redispool", 300)
		}
		pool.Put(conns[i])
	}

	pool.Empty()
	t.Log("Pass")
}
