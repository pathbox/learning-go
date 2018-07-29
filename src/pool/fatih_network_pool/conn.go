package pool

import (
	"net"
	"sync"
)

type PoolConn struct {
	net.Conn //  当前正在使用的conn
	mu       sync.RWMutex
	c        *channelPool
	unusable bool
}

func (p *PoolConn) Close() error { // 就是把当前conn put back to pool
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.unusable {
		if p.Conn != nil {
			return p.Conn.Close()
		}
		return nil
	}
	return p.c.put(p.Conn)
}

// MarkUnusable() marks the connection not usable any more, to let the pool close it instead of returning it to pool.
func (p *PoolConn) MarkUnusable() {
	p.mu.Lock()
	p.unusable = true
	p.mu.Unlock()
}

func (c *channelPool) wrapConn(conn net.Conn) net.Conn {
	p := &PoolConn{c: c}
	p.Conn = conn
	return p
}
