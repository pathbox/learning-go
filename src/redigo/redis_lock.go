package main

import (
	"fmt"
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Lock struct {
	resource string
	token    string
	conn     redis.Conn
	timeout  int
}

func (lock *Lock) tryLock() (bool, error) {
	_, err := redis.String(lock.conn.Do("SET", lock.key(), lock.token, "EX", int(lock.timeout), "NX")) // lock 就是 setNX 加过期时间
	if err == redis.ErrNil {
		// The lock was not successful, it already exists.
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil //成功获得锁 返回true
}

func (lock *Lock) Unlock() error { //unlock 就是del key
	_, err := lock.conn.Do("del", lock.key())
	return err
}

func (lock *Lock) key() string {
	return fmt.Sprintf("redislock:%s", lock.resource)
}

func (lock *Lock) AddTimeout(ex_time int64) (bool, error) {
	ttl_time, err := redis.Int64(lock.conn.Do("TTL", lock.key()))
	fmt.Println("ttl time", ttl_time)
	if err != nil {
		log.Fatal("redis get failed:", err)
	}

	if ttl_time > 0 {
		fmt.Println(11)
		_, err := redis.String(lock.conn.Do("SET", lock.key(), lock.token, "EX", int(ttl_time+ex_time)))
		if err == redis.ErrNil {
			return false, nil
		}
		if err != nil {
			return false, err
		}
	}
	return false, nil
}

func TryLock(conn redis.Conn, resource string, token string, DefaulTimeout int) (lock *Lock, ok bool, err error) {
	return TryLockWithTimeout(conn, resource, token, DefaulTimeout)
}

func TryLockWithTimeout(conn redis.Conn, resource string, token string, timeout int) (lock *Lock, ok bool, err error) {
	lock = &Lock{resource, token, conn, timeout}

	ok, err = lock.tryLock()

	if !ok || err != nil {
		lock = nil
	}

	return
}

func main() {
	fmt.Println("start")
	DefaultTimeout := 3
	conn, err := redis.Dial("tcp", "localhost:6379")

	lock, ok, err := TryLock(conn, "xiaoru.cc", "token", int(DefaultTimeout))
	if err != nil {
		log.Fatal("Error while attempting lock")
	}
	if !ok {
		log.Fatal("Lock")
	}
	lock.AddTimeout(100)

	time.Sleep(time.Duration(DefaultTimeout) * time.Second)
	fmt.Println("end")
	defer lock.Unlock()
}
