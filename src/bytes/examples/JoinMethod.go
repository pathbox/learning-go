func Join(s [][]byte, sep []byte) {
	if len(s) == 0 {
		return []byte{}
	}
	if len(s) == 1 {
		return append([]byte(nil), s[0]...)
	}
	n := len(sep) * (len(s) - 1) // count len(sep)
	for _, v := range s {
		n += len(v) // count len(v)
	}
	b := make([]byte, n)
	bp := copy(b, s[0])
	for _, v := range s[1:] {
		bp += copy(b[bp:], sep) // add sep
		bp += copy(b[bp:], v)   // add v
	}
}