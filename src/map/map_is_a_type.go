package main

import "fmt"

type resMap map[string]interface{}

func main() {
	m := getMap()

	fmt.Println("map is: ", m)
}

func getMap() resMap {
	r := make(map[string]interface{})

	r["name"] = "Cary"
	r["age"] = 27
	return r
}
