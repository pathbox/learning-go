func or(chans ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		var once sync.Once
		for _, c := range chans {
			go func(c <-chan interface{}) {
				select {
				case <-c:
					once.Do(func() { close(out)})
				case <-out:
				}
			}(c)
		}
	}()
	return out
}