package gootp

import (
	"fmt"
	"testing"
)

func TestHOTP(t *testing.T) {
	secret := []byte("it works")
	codes := []string{"589511",
		"577914",
		"161463",
		"679608",
		"798396",
		"296513",
		"960531",
		"015650",
		"706124",
		"632316",
		"084253",
		"365671",
		"208259",
		"928972",
		"578454",
		"273403",
		"064499",
		"373756",
		"732496",
		"200182",
		"601882",
		"622765",
		"068158",
		"398765",
		"468796",
		"151151",
		"543124",
		"483908",
		"823855",
		"179797",
		"972202",
		"100346"}
	hotp := NewHOTP(secret, 6)
	for index, code := range codes {
		counter := uint64(index)
		fmt.Printf("input code: %s  code string: %s\n", code, hotp.At(counter))
		if !hotp.Verify(code, counter) {
			t.Errorf("%d, %s", index, code)
			return
		}
	}
}
