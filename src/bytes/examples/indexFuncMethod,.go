import "unicode/utf8"

func indexFunc(s []byte, f func(r rune) bool, truth bool) int {
	start := 0
	for start < len(s) {
		wid := 1
		r := rune(s[start])
		if r >= utf8.RuneSelf {
			r, wid = utf8.DecodeRune(s[start:])
		}
		if f(r) == truth {
			return start
		}
		start += wid
	}
	return -1
}