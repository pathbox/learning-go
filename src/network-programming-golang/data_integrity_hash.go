package main

import (
	"crypto/md5"
	"fmt"
)

func main() {
	hash := md5.New()
	bytes := []byte("hello sjfoisjdfoijsofjasfjipsadjfiosjdo\n")
	hash.Write(bytes)
	hashValue := hash.Sum(nil)
	// hashSize := hash.Size()

	// for n := 0; n < hashSize; n += 4 { // 进一步的进制处理?
	// 	var val uint32
	// 	val = uint32(hashValue[n])<<24 +
	// 		uint32(hashValue[n+1])<<16 +
	// 		uint32(hashValue[n+2])<<8 +
	// 		uint32(hashValue[n+3])

	// 	fmt.Printf("%x ", val)
	// }
	fmt.Printf("%x", hashValue)
}
