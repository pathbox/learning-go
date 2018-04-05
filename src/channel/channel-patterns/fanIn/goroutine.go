// 每个goroutine起一个goroutine

func fanIn(chans ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		var wg sync.WaitGroup
		wg.Add(len(chans))

		for _, c := range chans{ // 每个chan 起一个goroutine
			go func(c <-chan interface{}) {
				for v := range c {
					out <- v 
				}
				wg.Done()
			}(c)
		}
		wg.Wait()
		close(out)
	}()
	return out
}