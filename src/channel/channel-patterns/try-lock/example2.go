type Mutex struct {
	ch chan struct{}
}

func NewMutex() *Mutex {
	mu := &Mutex{make(chan struct{}, 1)}// 初始化带1缓冲的chan
	mu.ch <- struct{}{}
	return mu
}

func (m *Mutex) Lock() {
	<-m.ch // 从ch读取，阻塞直到有写入为止
}

func (m *Mutex) Unlock() {
	select {
	case m.ch <- struct{}{}: // 写入ch
	default: 
		panic("unlock of unlocked mutex")
	}
}

func (m *Mutex) TryLock() bool {
	select {
	case <-m.ch:
		return true
	default:
	}
	return false
}

func (m *Mutex) IsLocked() bool {
	return len(m.ch) == 0 // ch中struct{}被取走，说明lock了
}

// 在某个时刻，只有一个goroutine 能Lock成功，其他goroutine需要等待Unlock后，才能Lock
// struct{}相当于要争抢的锁🔐
