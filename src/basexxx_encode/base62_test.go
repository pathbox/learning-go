package basexxx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase62(t *testing.T) {
	var ui uint64 = 123456
	str := Base62(ui)
	assert.Equal(t, str, "W7E")
}

func TestBytesToBase62(t *testing.T) {
	bs := []byte{255, 200, 155, 100, 55, 11, 0, 255, 200}
	str := BytesToBase62(bs)
	assert.Equal(t, str, "LxWs8cX6Xt93E")
}

func TestParseBase62(t *testing.T) {
	s := "W7E"
	ui, err := ParseBase62(s)
	assert.NoError(t, err)
	assert.Equal(t, ui, uint64(123456))
}
