package main

import (
	"fmt"
	"github.com/json-iterator/go"
	"net/http"
)

type ColorGroup struct {
	ID     int
	Name   string
	Colors []string
}

func main() {

	http.HandleFunc("/get", getHandler)
	http.ListenAndServe(":9009", nil)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	group := &ColorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	str, _ := json.Marshal(group)

	fmt.Println(str)

	fmt.Fprintln(w, string(str))
}
