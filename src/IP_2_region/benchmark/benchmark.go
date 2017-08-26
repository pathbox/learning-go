package benchmark

import (
	"fmt"
	"github.com/mohong122/ip2region/binding/golang"
)

func MemorySearch() {
	region, err := ip2region.New(".././ip2region.db")
	defer region.Close()
	if err != nil {
		fmt.Println("============")
		fmt.Println(err)
		return
	}
	ip, err := region.BtreeSearch("123.95.223.18")
	fmt.Println(ip, err)
}
