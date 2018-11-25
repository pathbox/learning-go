package apns

import (
	"container/list"
	"context"
	"errors"
	"syn conc"
	"time"

	"log"
)

// connection pool
type ConnPool struct {
	ctx context.Context
	dialFunc func(ctx context.Context)(*ApnsConn, error)
	poolSize int
	pool *list.List
	running bool
	cancel context.CancelFunc
	mutex sync.RWMutex
}

func NewConnPool(poolSize int,
	parentCtx context.Context,
	dialFunc func(ctx context.Context) (*ApnsConn, error)) (*ConnPool, error) {

	ctx,cancel := context.WithCancel(parentCtx)
	pool := &ConnPool{
		ctx:      ctx,
		cancel:cancel,
		poolSize: poolSize,
		dialFunc: dialFunc,
		running:  true,
		pool: list.New()}

	err := pool.enhancedPool(poolSize)
	if nil != err {
		return nil, err
	}
	return pool, nil
}

func (self *ConnPool) enhancedPool(size int)error {

	for i := 0; i < size; i++{
		j := 0
		var err error
		var conn *ApnsConn
		for ; j< 3; j++{
			conn,err = self.dialFunc(self.ctx)
			if err != nil || conn == nil {
				log.Printf("POOL_FACTORY|CREATE CONNECTION|INIT|FAIL|%s", err)
				continue
			} else {
				break
			}
		}

		if j >= 3 || conn == nil {
			return errors.New("POOL_FACTORY|CREATE CONNECTION|INIT|FAIL|%s" + err.Error())
		}
		self.pool.PushFront(conn)
	}

	return nil
}

func (self *ConnPool) Get() (*ApnsConn,error) {
	if !self.running {
		return nil,errors.New("POOL_FACTORY|POOL IS SHUTDOWN")
	}

	self.mutex.Lock()
	defer self.mutex.Unlock()

	var conn *ApnsConn
	//先从Idealpool中获取如果存在那么就直接使用
	for e := self.pool.Back(); nil != e; e = e.Prev() {
		conn = e.Value.(*ApnsConn)
		//要么是不存活都需要移除
		if conn.alive {
			self.pool.MoveToFront(e)
			break
		} else {
			//什么都不干
		}
	}
	//找到一个存活的链接
	if nil !=conn && conn.alive{
		return conn,nil
	}

	//如果没有找到合格的一个连接，那么主动队列尾部的
	e := self.pool.Back()
	conn = e.Value.(*ApnsConn)
	var err error
	//要么是不存活都需要移除
	if !conn.alive {
		//移除队列尾部并主动创建
		conn.Destroy()
		self.pool.Remove(e)

		conn, err = self.dialFunc(self.ctx)
		if nil == err && nil != conn {
			self.pool.PushFront(conn)
		}
	}

	return conn,err
}

func (self *ConnPool) Shutdown() {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.running = false
	self.cancel()
	for i := 0; i < 3; {
		//等待五秒中结束
		time.Sleep(5 * time.Second)
		if self.pool.Len() <= 0 {
			break
		}

		log.Printf("CONNECTION POOL|CLOSEING|WORK POOL SIZE|:%d\n", self.pool.Len())
		i++
	}

	var idleconn *ApnsConn
	//关闭掉空闲的client
	for e := self.pool.Front(); e != nil; e = e.Next() {
		idleconn = e.Value.(*ApnsConn)
		idleconn.Destroy()
		idleconn = nil
	}

	log.Println("CONNECTION_POOL|SHUTDOWN")
}