package basexxx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase32(t *testing.T) {
	var ui uint64 = 123456
	str := Base32(ui)
	assert.Equal(t, str, "da1y")
}

func TestBytesToBase32(t *testing.T) {
	bs := []byte{255, 200, 155, 100, 55, 11, 0, 255, 200}
	str := BytesToBase32(bs)
	assert.Equal(t, str, "x91r5co5osy89ge")
}

func TestParseBase32(t *testing.T) {
	s := "da1y"
	ui, err := ParseBase32(s)
	assert.NoError(t, err)
	assert.Equal(t, ui, uint64(123456))
}
