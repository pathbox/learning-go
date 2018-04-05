func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil 
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func(){
		defer close(orDone)
		var cases []reflect.SeelctCase
		for _, c := range channels {
			cases = append(cases, reflect.SelectCase{
				Dir: reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}
		reflect.Select(cases)
	}()
	return orDone
}