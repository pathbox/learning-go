// 将读取的值发送给每个输出channel， 异步模式可能会产生很多的goroutine
func fanOut(ch <-chan interface{}, out []chan interface{}, async bool) {
	go func() {
		defer func() {
			for i := 0; i < len(out); i++{
				close(out[i])
			}
		}()

		for v := range ch {
			for i := 0; i < len(out); i++{
				if async {
					go func() {
						out[i] <- v
					}()
				} else {
					out[i] <- v 
				}
			}
		}
	}()
}