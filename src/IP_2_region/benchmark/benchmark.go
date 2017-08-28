package benchmark

import (
	"fmt"
	"github.com/mohong122/ip2region/binding/golang"
)

var IP_SEARCH, _ = ip2region.New(".././data/ip2region.db")

func MemorySearch() {

	// IP_SEARCH.BinarySearch("123.95.223.18")
	// IP_SEARCH.BtreeSearch("123.95.223.90")
	IP_SEARCH.MemorySearch("123.95.223.90")
	fmt.Println(ip, err)
}

// var IP_SEARCH, _ = ip2region.New(".././data/ip2region.db") //这里使用全局变量 可以达到最好的效果
