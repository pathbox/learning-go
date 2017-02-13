package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Info struct {
	Database struct {
		Username  string
		Password  string
		Name      string
		Hostname  string
		Port      uint16
		Parameter string
	}
	Email struct {
		Username string
		Password string
		Hostname string
		Port     uint16
		From     string
	}
}

func ConfigPath() string {
	path := "./config.json"
	return path
}

func initInfo() *Info {
	file, err := ioutil.ReadFile(ConfigPath())
	if err != nil {
		panic(err)
	}

	info := new(Info) // &Info{}
	json.Unmarshal(file, info)
	return info
}

func main() {
	info := initInfo()
	fmt.Println(info)
	fmt.Println(info.Email.Username)
}
