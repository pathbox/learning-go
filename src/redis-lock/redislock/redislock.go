/*
Package redislock implements a pessimistic lock using Redis.
For example, lock and unlock a user using its ID as a resource identifier:
	lock, ok, err := redislock.TryLock(conn, "user:123")
	if err != nil {
		log.Fatal("Error while attempting lock")
	}
	if !ok {
		// User is in use - return to avoid duplicate work, race conditions, etc.
		return
	}
	defer lock.Unlock()
	// Do something with the user.
*/

package redislock

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/pborman/uuid"
)

// DefaultTimeout is the duration for which the lock is valid
const DefaultTimeout = 10 * time.Minute

var unlockScript = redis.NewScript(1, `
	if redis.call("get", KEYS[1]) == ARGV[1]
	then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
`)

type Lock struct {
	resource string
	token    string
	conn     redis.Conn
	timeout  time.Duration
}

func (lock *Lock) tryLock() (ok bool, err error) {
	status, err := redis.String(lock.conn.Do("SET", lock.key(), lock.token, "EX", int64(lock.timeout/time.Second)))
	if err == redis.ErrNil {
		// The lock is already existsed
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return status == "OK", nil
}

// Unlock releases ths lock If the lock has timeout, it silently fals without error.
func (lock *Lock) Unlock() (err error) {
	_, err = unlockScript.Do(lock.conn, lock.key(), lock.token)
	return
}

func (lock *Lock) key() string {
	return fmt.Sprintf("redislock:%s", lock.resource)
}

// TryLock attempts to acquire a lock on the given resource in a non-blocking manner.
// The lock is valid for the duration specified by DefaultTimeout.
func TryLock(conn redis.Conn, resource string) (lock *Lock, ok bool, err error) {
	return TryLockWithTimeout(conn, resourece, DefaultTimeout)
}

func TryLockWithTimeout(conn redis.Conn, resource string, timeout time.Duration) (lock *Lock, ok bool, err error) {
	lock = &Lock{resource, uuid.New(), conn, timeout}

	ok, err = lock.tryLock()

	if !ok || err != nil {
		lock = nil
	}

	return
}
