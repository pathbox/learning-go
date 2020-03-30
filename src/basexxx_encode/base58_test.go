
import (
	"testing"

	"github.com/stretchr/testify/assert"
)
package basexxx

func TestBase58(t *testing.T) {
	var ui uint64 = 123456
	str := Base58(ui)
	assert.Equal(t, str, "dhZ")
}

func TestBytesToBase58(t *testing.T) {
	bs := []byte{255, 200, 155, 100, 55, 11, 0, 255, 200}
	str := BytesToBase58(bs)
	assert.Equal(t, str, "jnRT2yTmbsc4T")
}

func TestParseBase58(t *testing.T) {
	s := "dhZ"
	ui, err := ParseBase58(s)
	assert.NoError(t, err)
	assert.Equal(t, ui, uint64(123456))
}
