package gootp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	"strconv"
)

type HOTP struct {
	secret []byte
	digits int
}

// Implement RFC4226
// At generate code
func (h HOTP) At(counter uint64) string {
	counterBytes := make([]byte, 8) // 8字节长度
	binary.BigEndian.PutUint64(counterBytes, counter)
	hash := hmac.New(sha1.New, h.secret)
	hash.Write(counterBytes)
	hs := hash.Sum(nil) // hmac 算法得到hs值 是一个字符串
	offset := hs[19] & 0x0f
	binCodeBytes := make([]byte, 4) // 4字节长度
	binCodeBytes[0] = hs[offset] & 0x7f
	binCodeBytes[1] = hs[offset+1] & 0xff
	binCodeBytes[2] = hs[offset+2] & 0xff
	binCodeBytes[3] = hs[offset+3] & 0xff
	binCode := binary.BigEndian.Uint32(binCodeBytes)
	mod := uint32(1)
	for i := 0; i < h.digits; i++ {
		mod *= 10
	}
	code := binCode % mod // 获取这个整数的后6位（可以根据需要取后8位）
	codeString := strconv.FormatUint(uint64(code), 10)
	if len(codeString) < h.digits { // 将数字转成字符串，不够6位前面补0
		paddingByteLength := h.digits - len(codeString)
		paddingBytes := make([]byte, paddingByteLength)
		for i := 0; i < paddingByteLength; i++ {
			paddingBytes[i] = '0'
		}
		codeString = string(paddingBytes) + codeString
	}
	return codeString
}

// NewHOTP generate new HOTP instance  secret is a secret string, as []byte pass to the func
func NewHOTP(secret []byte, digits int) (h *HOTP) {
	h = new(HOTP)
	h.secret = secret
	h.digits = digits
	return
}

// Verify verify OTP code
func (h HOTP) Verify(code string, counter uint64) bool {
	realCode := h.At(counter) // code是客户端传过来的code 6位数字字符串  realCode是服务端算法生成的code, counter是一个随机整数
	if realCode == code {
		return true
	}
	return false
}
