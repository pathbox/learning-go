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
			cases[i].Chan = reflect.ValueOf(out[i])
		}
		for v := range ch {
			v := v
			for i := range cases {
				cases[i].Send = reflect.ValueOf(v)
			}
			_, _, _ = reflect.Select(cases)
		}
	}()
}