package main

import (
	// "encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	json "github.com/json-iterator/go"
)

type Person struct {
	Name string
}

func tryParams(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var p Person

	err := json.Unmarshal(body, &p)
	log.Println(err)

	fmt.Println(p)

}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", tryParams)

	log.Println("Listening at :9090")
	log.Fatal(http.ListenAndServe(":9090", mux))
}
