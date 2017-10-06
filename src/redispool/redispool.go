package redispool

import (
	"sync/atomic"
	"time"

	"github.com/fzzy/radix/redis"
)

// A simple connection pool. It will create a small pool of initial connections,
// and if more connections are needed they will be created on demand. If a
// connection is returned and the pool is full it will be closed.
type RedisPool struct {
	pwd      string
	network  string
	addr     string
	pool     chan *redis.Client
	qpsLimit int32
	qps      int32
}

const DEFAULT_QPS_LIMIT = 2000

func (p *RedisPool) getRedisClient() (*redis.Client, error) {
	client, err := redis.Dial(p.network, p.addr) // 创建redis connecttion
	if err != nil {
		return nil, err
	}

	if len(p.pwd) != 0 {
		if client, err = p.redisAuth(client); err != nil {
			return nil, err
		}
	}
	return client, nil // 返回client
}

func (p *RedisPool) redisAuth(client *redis.Client) (*redis.Client, error) {
	if _, err := client.Cmd("AUTH", p.pwd).Str(); err != nil {
		client.Close()
		return client, err
	}
	return client, nil
}

// Creates a new Pool whose connections are all created using
// redis.Dial(network, addr). The size indicates the maximum number of idle
// connections to have waiting to be used at any given moment
func NewPool(network, addr, password string, size int) (*RedisPool, error) {
	if size <= 0 {
		size = 1
	}
	rp := RedisPool{
		network:  network,
		addr:     addr,
		pool:     make(chan *redis.Client, size),
		pwd:      password,
		qpsLimit: int32(DEFAULT_QPS_LIMIT),
		qps:      int32(0),
	}

	clientList := make([]*redis.Client, 0, size) // len = 0, cap = size
	for i := 0; i < size; i++ {
		client, err := rp.getRedisClient()
		if err != nil { // 只要发生错误，就关闭所有client_list 中的redis client
			for _, client = range clientList {
				client.Close()
			}
			return nil, err
		}
		if client != nil {
			clientList = append(clientList, client) // 将client 放入clientList
		}
		// 将clientList 中的redis client，放入 rp.pool 可缓冲的chan中
		for i := range clientList {
			rp.pool <- clientList[i]
		}
	}
	return &rp, nil
}

// Calls NewPool, but if there is an error it return a pool of the same size but
// without any connections pre-initialized (can be used the same way, but if
// this happens there might be something wrong with the redis instance you're
// connecting to)
func NewOrEmptyPool(network, addr, password string, size int) *RedisPool {
	pool, err := NewPool(network, addr, password, size)
	if err != nil {
		if size <= 0 {
			size = 1
		}
		pool = &RedisPool{
			network:  network,
			addr:     addr,
			pool:     make(chan *redis.Client, size),
			pwd:      password,
			qpsLimit: int32(DEFAULT_QPS_LIMIT),
			qps:      int32(0),
		}
	}
	return pool
}

func (p *RedisPool) get() (client *redis.Client, err error) {
	select {
	case conn := <-p.pool:
		return conn, nil
	default: // 如果 p.pool 中没有conn了，都被使用了，就主动请求创建新的redis client. 这样会不会导致太多redis conn 产生呢？
		return p.getRedisClient()
	}
}

// Retrieves an available redis client. If there are none available it will
// create a new one on the fly
func (p *RedisPool) Get() (client *redis.Client, err error) {
	if p.qpsLimit > 0 {
		for p.qps > p.qpsLimit {
			time.Sleep(time.Millisecond * time.Duration(10))
		}
	}

	for i := 0; i < 3; i++ { // 尝试三次 成功则终止循环
		if client, err = p.get(); err != nil { // redis client 是从p.pool中获取或者新建一个redis conn
			break
		} else if pstate, perr := client.Cmd("PING").Str(); pstate == "PONG" && perr == nil {
			break
		}
	}
	if err == nil {
		atomic.AddInt32(&p.qps, 1)
	}
	return
}

// Returns a client back to the pool. If the pool is full the client is closed
// instead. If the client is already closed (due to connection failure or
// what-have-you) it should not be put back in the pool. The pool will create
// more connections as needed.
func (p *RedisPool) Put(conn *redis.Client) {
	select {
	case p.pool <- conn:
	default:
		conn.Close()
	}
	if p.qps > 0 {
		atomic.AddInt32(&p.qps, -1)
	}
}

// Removes and calls Close() on all the connections currently in the pool.
// Assuming there are no other connections waiting to be Put back this method
// effectively closes and cleans up the pool.
func (p *RedisPool) Empty() {
	var conn *redis.Client

	for {
		select {
		case conn = <-p.pool:
			conn.Close()
		default:
			return
		}
		if p.qps > 0 {
			atomic.AddInt32(&p.qps, -1)
		}
	}
}
