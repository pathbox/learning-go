package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type configuration struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
}

func main() {
	file, _ := os.Open("conf.json")
	defer file.Close()

	decoder := json.NewDecoder(file) // decode json
	conf := &configuration{}
	err := decoder.Decode(conf) // json to struct
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(conf.Path)

}
