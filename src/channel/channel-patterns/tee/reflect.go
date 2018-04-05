func fanOutReflect(ch <-chan interface{}, out []chan interface{}) {
	go func() {
		defer func() {
			for i := 0; i < len(out); i++ {
				close(out[i])
			}
		}()
		cases := make([]reflect.SelectCase, len(out))
		for i := range cases {
			cases[i].Dir = reflect.SelectSend
		}
		for v := range ch {
			v := v
			for i := range cases {
				cases[i].Chan = reflect.ValueOf(out[i])
				cases[i].Send = reflect.ValueOf(v)
			}
			for _ = range cases { // for each channel
				chosen, _, _ := reflect.Select(cases)
				cases[chosen].Chan = reflect.ValueOf(nil)
			}
		}
	}()
}