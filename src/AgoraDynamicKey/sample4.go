package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

func GeneratePublicSharingKey(appID, appCertificate, channelName string, unixTs, randomInt, uid, expiredTs uint32) string {
	return generateDynamicKey(appID, appCertificate, channelName, unixTs, randomInt, uid, expiredTs, "APSS")
}

func GenerateRecordingKey(appID, appCertificate, channelName string, unixTs, randomInt, uid, expiredTs uint32) string {
	return generateDynamicKey(appID, appCertificate, channelName, unixTs, randomInt, uid, expiredTs, "ARS")
}

func GenerateMediaChannelKey(appID, appCertificate, channelName string, unixTs, randomInt, uid, expiredTs uint32) string {
	return generateDynamicKey(appID, appCertificate, channelName, unixTs, randomInt, uid, expiredTs, "ACS")
}

func generateDynamicKey(appID, appCertificate, channelName string, unixTs, randomInt, uid, expiredTs uint32, serviceType string) string {
	// version := "004"
	unixTsStr := fmt.Sprintf("%010d", unixTs)
	fmt.Println("unixTsStr", unixTsStr)
	randomIntStr := fmt.Sprintf("%08x", randomInt)
	uidStr := fmt.Sprintf("%010d", uid)
	expiredTsStr := fmt.Sprintf("%010d", expiredTs)
	signature := generateSignature(appID, appCertificate, channelName, unixTsStr, randomIntStr, uidStr, expiredTsStr, serviceType)
	fmt.Println("signature: ", signature)

	buffer := strings.Join([]string{signature, appID, unixTsStr, randomIntStr, expiredTsStr}, "")
	return buffer
}

func generateSignature(appID, appCertificate, channelName, unixTsStr, randomIntStr, uidStr, expiredTsStr, serviceType string) string {
	buffer := strings.Join([]string{serviceType, appID, unixTsStr, randomIntStr, channelName, uidStr, expiredTsStr}, "")
	fmt.Println("buffer: ", buffer)

	signature := hmac.New(sha1.New, []byte(appCertificate))
	fmt.Println("before", string(signature))
	signature.Write([]byte(buffer))
	result := signature.Sum(nil)
	return hex.EncodeToString(result)
}

func main() {
	appID := "970ca35de60c44645bbae8a215061b33"
	appCertificate := "5cfd2fd1755d40ecb72977518be15d3b"
	channelName := "7d72365eb983485397e3e3f9d460bdda"
	unixTs := uint32(1446455472)
	uid := uint32(2882341273)
	randomInt := uint32(58964981)
	expiredTs := uint32(1446455471)

	mediaChannelKey := GenerateMediaChannelKey(appID, appCertificate, channelName, unixTs, randomInt, uid, expiredTs)
	fmt.Println("MediaChannelKey:", mediaChannelKey)
	// 	recordingKey := GenerateRecordingKey(appID, appCertificate, channelName, unixTs, randomInt, uid, expiredTs)
	// 	fmt.Println("recordingKey:", recordingKey)
	// 	publicSharingKey := GeneratePublicSharingKey(appID, appCertificate, channelName, unixTs, randomInt, uid, expiredTs)
	// 	fmt.Println("GeneratePublicSharingKey:", publicSharingKey)
}

/*
MediaChannelKey: 004d0ec5ee3179c964fe7c0485c045541de6bff332b970ca35de60c44645bbae8a215061b3314464554720383bbf51446455471
recordingKey: 004e0c24ac56aae05229a6d9389860a1a0e25e56da8970ca35de60c44645bbae8a215061b3314464554720383bbf51446455471
GeneratePublicSharingKey: 004ec32c0d528e58ef90e8ff437a9706124137dc795970ca35de60c44645bbae8a215061b3314464554720383bbf51446455471
*/
