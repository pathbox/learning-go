package basexxx

import "encoding/binary"

func BytesToBase(bs []byte, baseFunc func(ui uint64) string) string {
	l := len(bs)
	s := ""

	for i := 0; i < l; i += 8 {
		li := i + 8
		if li > l {
			li = l
		}

		b := bs[i:li]
		lb := len(b)
		if lb < 8 {
			b = append(make([]byte, 8-lb), b...)
		}
		ui := binary.BigEndian.Uint64(b)
		s += baseFunc(ui)
	}

	return s
}