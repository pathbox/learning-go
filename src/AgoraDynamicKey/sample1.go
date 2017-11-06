package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

func Generate(appID, appCertificate, channelName string, unixTs, randomInt uint32) string {
	return generateDynamicKey(appID, appCertificate, channelName, unixTs, randomInt)
}

func generateDynamicKey(appID, appCertificate, channelName string, unixTs, randomInt uint32) string {
	unixTsStr := fmt.Sprintf("%010d", unixTs)
	randomIntStr := fmt.Sprintf("%08x", randomInt)
	signature := generateSignature(appID, appCertificate, channelName, unixTsStr, randomIntStr)
	buffer := strings.Join([]string{signature, appID, unixTsStr, randomIntStr}, "")
	return buffer
}

func generateSignature(appID, appCertificate, channelName, unixTsStr, randomIntStr string) string {
	buffer := strings.Join([]string{appID, unixTsStr, randomIntStr, channelName}, "")
	signature := hmac.New(sha1.New, []byte(appCertificate))
	signature.Write([]byte(buffer))
	return hex.EncodeToString(signature.Sum(nil))
}

func main() {
	appID := "970ca35de60c44645bbae8a215061b33"
	appCertificate := "5cfd2fd1755d40ecb72977518be15d3b"
	channelName := "7d72365eb983485397e3e3f9d460bdda"
	unixTs := uint32(1446455472)
	// uid := uint32(2882341273)
	randomInt := uint32(58964981)
	// expiredTs := uint32(1446455471)
	result := Generate(appID, appCertificate, channelName, unixTs, randomInt)
	fmt.Println("result:", result)
}
