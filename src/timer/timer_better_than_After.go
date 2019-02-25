package timer

import "time"

func UseNewTimer(in <-chan string) {
	defer wg.Done()

	idleDuration := 3 * time.Minute
	idleDelay := time.NewTimer(idleDuration)
	Running := true
	for Running {
		idleDelay.Reset(idleDuration)
		select {
		case _, ok := <-in:
			if !ok {
				return
			}

			// handlke something
		case <-idleDelay.C: // 这里不再使用time.After(t) 会因为 t的时间长导致object占用大量heap内存,从而引起GC等性能问题
			return // 三分钟超时，返回不再继续for Running
		}
	}
}
