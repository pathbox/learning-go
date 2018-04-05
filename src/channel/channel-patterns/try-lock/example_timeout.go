
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
25
26
27
28
29
30
31
32
33
34
35
36
type Mutex struct {
	ch chan struct{}
}
func NewMutex() *Mutex {
	mu := &Mutex{make(chan struct{}, 1)}
	mu.ch <- struct{}{}
	return mu
}
func (m *Mutex) Lock() {
	<-m.ch
}
func (m *Mutex) Unlock() {
	select {
	case m.ch <- struct{}{}:
	default:
		panic("unlock of unlocked mutex")
	}
}

func (m *Mutex) TryLock(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	select {
	case <-m.ch:
		timer.Stop()
		return true 
	case <-time.After(timeout):
	}
	return false
}

func (m *Mutex) IsLocked() bool {
	return len(m.ch) == 0
}