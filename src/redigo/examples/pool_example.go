pool := &redis.Pool{
	Dial: func()(redis.Conn, error) {
		c, err := redis.Dial("tcp", ":6379")
		if err != nil {
			return nil, err
		}
		// 做密码校验
		if _, err := c.Do("AUTH", password); err != nil {
			c.Close() // return to pool
			return nil, err
		}
		if _, err := c.Do("SELECT", db); err != nil {
			c.Close()
			return nil, err
		}
		_, err := c.Do("PING"); err != nil {
			return nil, err
		}

		return c, nil
	},
}

redisPool := redis.NewPool(func() (redis.Conn, error){
	return redis.Dial("tcp",":6379")
}, 10)

func NewPool(newFn func()(Conn, error), maxIdle int) *Pool {
	return &Pool{Dial: newFn, MaxIdle: maxIdle}
}

func (p *Pool) Get() Conn {
	pc, err := p.get(nil)
	if err != nil {
		return errorConn{err}
	}
	return &activeConn{p: p, pc: pc}
}