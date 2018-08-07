package main

import (
	"fmt"
	"os"

	"github.com/bcicen/jstream"
)

func main() {
	f, _ := os.Open("./input.json") // input.json is a json file has json data
	decoder := jstream.NewDecoder(f, 1)
	for mv := range decoder.Stream() {
		fmt.Println("%v-%v\n", mv.Value, mv.Depth)
	}
}
