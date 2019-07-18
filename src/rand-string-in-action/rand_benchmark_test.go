package main
import (
	"math/rand"
	"testing"
	"time"
	"github.com/nu7hatch/gouuid"
	"crypto/sha256"
	"encoding/hex"
)
// Implementations
func init() {
	rand.Seed(time.Now().UnixNano())
}
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
func RandStringBytesMask(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; {
		if idx := int(rand.Int63() & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i++
		}
	}
	return string(b)
}
func RandStringBytesMaskImpr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}
var src = rand.NewSource(time.Now().UnixNano())
func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func RandStringUUID() string {
	uuidStr1, _ := uuid.NewV4()
	hk := sha256.New()
		hk.Write([]byte(uuidStr1.String()))
		return hex.EncodeToString(hk.Sum(nil))
}
// Benchmark functions
const n = 64
func BenchmarkRunes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringRunes(n)
	}
}
func BenchmarkBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytes(n)
	}
}
func BenchmarkBytesRmndr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytesRmndr(n)
	}
}
func BenchmarkBytesMask(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytesMask(n)
	}
}
func BenchmarkBytesMaskImpr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytesMaskImpr(n)
	}
}
func BenchmarkBytesMaskImprSrc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytesMaskImprSrc(n)
	}
}

func BenchmarkRandStringUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringUUID()
	}
}
// 加入uuid + sha256 但速度不够快

// go test -bench=. rand_string_benchmark_test.go
/*
goos: darwin
goarch: amd64
BenchmarkRunes-4                  300000              6138 ns/op
BenchmarkBytes-4                  300000              5120 ns/op
BenchmarkBytesRmndr-4             500000              2391 ns/op
BenchmarkBytesMask-4              500000              3086 ns/op
BenchmarkBytesMaskImpr-4         2000000               525 ns/op
BenchmarkBytesMaskImprSrc-4      5000000               391 ns/op
BenchmarkRandStringUUID-4        1000000              1733 ns/op
*/