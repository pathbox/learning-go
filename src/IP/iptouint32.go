package main

import (
	"fmt"
	"strconv"
	"strings"
)

func IPToUnit32(ip string) uint32 {
	bits := strings.Split(ip, ".")
	var item [4]int

	for i, v := range bits {
		item[i], _ = strconv.Atoi(v)
	}

	var result uint32
	for i, va := range item {
		ii := 24 - 8*i
		result += uint32(va) << uint32(ii)
	}

	return result
}

func main() {
	ip := "11.33.2.3"
	r := IPToUnit32(ip)
	fmt.Println(r)
}
