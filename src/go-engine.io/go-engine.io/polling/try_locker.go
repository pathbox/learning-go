package polling

type Locker struct {
	locker chan struct{}
}

func NewLocker() *Locker {
	return &Locker{
		locker: make(chan struct{}, 1),
	}
}

func (l *Locker) Lock() {
	l.locker <- struct{}{}
}

func (l *Locker) TryLock() bool {
	select {
	case l.locker <- struct{}{}:
		return true
	default:
		return false

	}
}

func (l *Locker) Unlock() {
	<-l.locker
}

// 使用了 chan 的特性 模拟了锁的操作

// Unlock() 会阻塞 直到 l.locker <- struct{}{} 操作后，才会执行，并且继续执行下去
