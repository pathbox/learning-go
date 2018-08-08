package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Person struct {
	Name           string
	Age            int
	ServerResponse bool
}

func index(w http.ResponseWriter, r *http.Request) {

	// var person Person
	// body, _ := ioutil.ReadAll(r.Body)
	// xml.Unmarshal(body, &person)
	// person.ServerResponse = true

	// responseXML, _ := xml.MarshalIndent(person, "", "  ")

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	m := make(map[string]interface{})

	fmt.Println("Body:", string(body))
	log.Println(json.Unmarshal(body, &m))
	fmt.Println(m)

	fmt.Fprintf(w, "Hello world")
}

func main() {
	fmt.Println("go")
	http.HandleFunc("/", index)
	http.ListenAndServe(":9090", nil)
}
