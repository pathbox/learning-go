package main

import (
	"fmt"
	"hash/crc32"
	"hash/fnv"
	"io"
	"os"
)

func main() {
	hostname, err := os.Hostname() // 或者可以根据唯一的IP地址，来构造ID
	fmt.Println("hostname: ", hostname)
	if err != nil {
		fmt.Println(err)
	}

	// h := md5.New()
	h := fnv.New64a()
	io.WriteString(h, hostname)
	defaultID := int64(crc32.ChecksumIEEE(h.Sum(nil)) % 1024) // 表示集群实例数量最多 1024个
	fmt.Println("defaultID: ", defaultID)
}

// inspired by NSQ options.go NewOptions()
