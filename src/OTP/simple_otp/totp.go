package gootp

import "time"

// TOTP implementation
type TOTP struct {
	*HOTP
	x int64
}

func NewTOTP(secret []byte, digits, x int) (t *TOTP) {
	t = new(TOTP)
	t.HOTP = NewHOTP(secret, digits)
	t.x = int64(x)
	return t
}

func (t TOTP) At(timestamp int64) string {
	counter := uint64(timestamp / t.x) // counter 是对应时间戳组成
	return t.HOTP.At(counter)
}

// Now get current TOTP
func (t TOTP) Now() string {
	now := time.Now()
	timestamp := now.Unix()
	return t.At(timestamp)
}

// Verify verify OTP code
func (t TOTP) Verify(code string) bool {
	realCode := t.Now()
	if realCode == code {
		return true
	}
	return false
}
