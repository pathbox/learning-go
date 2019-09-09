package main

import (
	"fmt"
	"strings"
)

// parsePhone 解析用户手机号国家码
func parsePhone(phoneInfo string) (countryCode, userPhone string) {
	if phoneInfo == "" {
		return "", ""
	}
	// 没有国家码，默认86
	countryCode, userPhone = "86", phoneInfo
	indexLeft, indexRight := strings.Index(phoneInfo, "("), strings.Index(phoneInfo, ")")
	if indexLeft >= 0 && indexRight >= 0 && indexRight > indexLeft+1 {
		countryCode, userPhone = phoneInfo[indexLeft+1:indexRight], phoneInfo[indexRight+1:]
	}

	return countryCode, strings.TrimSpace(userPhone)
}

func main() {
	phone := "(86)    18666666666"
	code, userPhone := parsePhone(phone)
	fmt.Printf("code:%s, phone:%s\n", code, userPhone)
}
