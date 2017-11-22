package main

import (
	"fmt"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

type ColorGroup struct {
	ID     int
	Name   string
	Colors []string
}

func main() {

	resp, _ := http.Get("http://localhost:9009/get")

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	fmt.Println("data:", string(data))

	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	group := &ColorGroup{}

	json.Unmarshal(data, group)

	fmt.Println("result: ", group.ID, group.Name, group.Colors[0])
}

// val := []byte(`{"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"]}`)
// jsoniter.Get(val, "Colors", 0).ToString()
