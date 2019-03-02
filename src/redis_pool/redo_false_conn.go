package main
 
import (
    "errors"
    "fmt"
    "io"
    "strings"
    "time"
 
    "github.com/gomodule/redigo/redis"
)
 
var (
    RedisClient *redis.Pool
)
 
func init() {
    var (
        host string
        auth string
        db   int
    )
    host = "127.0.0.1:6379"
    auth = ""
    db = 0
    RedisClient = &redis.Pool{
        MaxIdle:     100,
        MaxActive:   4000,
        IdleTimeout: 180 * time.Second,
        Dial: func() (redis.Conn, error) {
            c, err := redis.Dial("tcp", host, redis.DialPassword(auth), redis.DialDatabase(db))
            if nil != err {
                return nil, err
            }
            return c, nil
        },
        TestOnBorrow: func(c redis.Conn, t time.Time) error {
            if time.Since(t) < time.Minute {
                return nil
            }
            _, err := c.Do("PING")
            return err
        },
    }
 
}
 
func main() {
    rd := RedisClient.Get()
    defer rd.Close()
 
    fmt.Println("please kill redis server")
    time.Sleep(5 * time.Second)
 
    fmt.Println("please start redis server")
    time.Sleep(5 * time.Second)
 
	resp, err := redis.String(redo("SET", "push_primay", "locked"))
    fmt.Println(resp, err)
}
 
func IsConnError(err error) bool {
    var needNewConn bool
 
    if err == nil {
        return false
    }
 
    if err == io.EOF {
        needNewConn = true
    }
    if strings.Contains(err.Error(), "use of closed network connection") {
        needNewConn = true
    }
    if strings.Contains(err.Error(), "connect: connection refused") {
        needNewConn = true
    }
    return needNewConn
}
 
// 在pool加入TestOnBorrow方法来去除扫描坏连接
func redo(command string, opt ...interface{}) (interface{}, error) {
    rd := RedisClient.Get()
    defer rd.Close()
 
    var conn redis.Conn
    var err error
    var maxretry = 3
    var needNewConn bool
 
    resp, err := rd.Do(command, opt...)
    needNewConn = IsConnError(err)
    if needNewConn == false {
        return resp, err
    } else {
        conn, err = RedisClient.Dial()
    }
 
    for index := 0; index < maxretry; index++ {
        if conn == nil && index+1 > maxretry {
            return resp, err
        }
        if conn == nil {
            conn, err = RedisClient.Dial()
        }
        if err != nil {
            continue
        }
 
        resp, err := conn.Do(command, opt...)
        needNewConn = IsConnError(err)
        if needNewConn == false {
            return resp, err
        } else {
            conn, err = RedisClient.Dial()
        }
    }
 
    conn.Close()
    return "", errors.New("redis error")
}
