package locker

// Mutext struct
type Mutex struct {
	lock chan struct{}
}

// 创建一个互斥锁
func NewMutex() *Mutex {
	return &Mutex{lock: make(chan struct{}, 1)}
}

// 锁操作
func (m *Mutex) Lock() {
	m.lock <- struct{}{}
}

// 解锁操作
func (m *Mutex) Unlock() {
	<-m.lock
}
