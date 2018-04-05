func mapChan(in <-chan interface{}, fn func(interface{}) interface{}) <-chan interface{} {
	out := make(chan interface{})
	if in == nil {
		close(out)
		return out
	}
	go func() {
		defer close(out)
		for v := range in {
			out <- fn(v)
		}
	}()
	return out
}

// Like map in Ruby