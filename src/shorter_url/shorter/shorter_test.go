package shorter

import (
	"testing"
)

func TestGetShortUrl(t *testing.T) {
	id := 75
	result := "bn"

	url := shorter.GetShortUrl(id)

	if url != result {
		t.Error("The shorter url is not result")
	}
}
