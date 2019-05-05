package main

import (
	"fmt"

	sendcloud "github.com/smartwalle/sendcloud"
)

func main() {
	apiUser := "apiUser"
	apiKey := "apiKey"
	sendcloud.UpdateApiInfo(apiUser, apiKey)

	_, _, r := sendcloud.GetTemplateList("", 0, 0, 0, 0)
	fmt.Println(r)
}
