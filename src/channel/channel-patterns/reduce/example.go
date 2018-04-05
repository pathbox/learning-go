
2
3
4
5
6
7
8
 // Like reduce in Ruby
func reduce(in <-chan interface{}, fn func(r, v interface{}) interface{}) interface{} {
	if in == nil {
		return nil
	}
	out := <-in
	for v := range in {
		out = fn(out, v)
	}
	return out
}
// 你可以用`reduce`实现`sum`、`max`、`min`等聚合操作