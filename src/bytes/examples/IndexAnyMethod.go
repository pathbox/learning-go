import "unicode/utf8"

func IndexAny(s []byte, chars string) int {
	if chars == "" {
		return -1 // Advoid scanning all of s
	}

	if len(s) > 8 { // len(s) > 8, use as.contains(c)
		if as, isASCII := makeASCIISet(chars); isASCII {
			for i, c := range s {
				if as.contains(c) {
					return i
				}
			}
			return -1
		}
	}

	var width int
	// double for, find first index
	for i := 0; i < len(s); i += width {
		r := rune(s[i])
		if r < utf8.RuneSelf {
			width = 1
		} else {
			r, width = utf8.DecodeRune(s[i:])
		}
		for _, ch := range chars {
			if r == ch {
				return i
			}
		}
	}
	return -1
}