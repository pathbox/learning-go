func fanOut(ch <-chan interface{}, out []chan interface{}) {
	go func() {
		defer func() {
			for i := 0; i < len(out); i++ {
				close(out[i])
			}
		}()
		// roundrobin
		var i = 0
		var n = len(out)
		for v := range ch {
			v := v
			out[i] <- v
			i = (i + 1) % n
		}
	}()
}