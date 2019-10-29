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
	counter := uint64(timestamp / t.x) // counter = 时间戳 / 步长(过期时间) 是整除，所以在t.x的误差内， 理论上客户端App和服务端得到的counter的值是一样的，got it  nice way，这样算法的参数是一样的，就能进行校验比对了
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
