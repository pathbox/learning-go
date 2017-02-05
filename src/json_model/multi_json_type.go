package main

import (
	"encoding/json"
	"fmt"
)

var jsonStr = []byte(`
  {
    "things": [
        {
            "name": "Alice",
            "age": 37
        },
        {
            "city": "Ipoh",
            "country": "Malaysia"
        },
        {
            "name": "Bob",
            "age": 36
        },
        {
            "city": "Northampton",
            "country": "England"
        }
    ]
  }`)

func main() {
	personsA, placesA := solutionA(jsonStr)
	fmt.Println(personsA, "--", placesA)
	fmt.Printf("%d %d\n", len(personsA), len(placesA))

	personsB, placesB := solutionB(jsonStr)
	fmt.Println(personsB, "--", placesB)
	fmt.Printf("%d %d\n", len(personsB), len(placesB))

	personsC, placesC := solutionC(jsonStr)
	fmt.Println(personsC, "--", placesC)
	fmt.Printf("%d %d\n", len(personsC), len(placesC))
}

type Person struct {
	Name string
	Age  int
}

type Place struct {
	City    string
	Country string
}

//Solution A
//Unmarshal into a map
//Type assert when we need it
func solutionA(jsonStr []byte) ([]Person, []Place) {
	persons := []Person{}
	places := []Place{}
	var data map[string][]map[string]interface{}

	err := json.Unmarshal(jsonStr, &data)
	if err != nil {
		fmt.Println(err)
		return persons, places
	}

	fmt.Println(data)
	for i := range data["things"] {
		item := data["things"][i]
		if item["name"] != nil {
			persons = addPerson(persons, item)
		} else {
			places = addPlace(places, item)
		}
	}
	return persons, places
}

func addPerson(persons []Person, item map[string]interface{}) []Person {
	name, _ := item["name"].(string)
	age, _ := item["age"].(int)
	person := Person{name, age}
	persons = append(persons, person)
	return persons
}

func addPlace(places []Place, item map[string]interface{}) []Place {
	city, _ := item["city"].(string)
	country, _ := item["city"].(string)
	place := Place{city, country}
	places = append(places, place)
	return places
}

//SolutionB

type Mixed struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	City    string `json:"city"`
	Country string `json:"country"`
}

func solutionB(jsonStr []byte) ([]Person, []Place) {
	persons := []Person{}
	places := []Place{}
	var data map[string][]Mixed
	err := json.Unmarshal(jsonStr, &data)
	if err != nil {
		fmt.Println(err)
		return persons, places
	}

	for i := range data["things"] {
		item := data["things"][i]
		if item.Name != "" {
			persons = append(persons, Person{item.Name, item.Age})
		} else {
			places = append(places, Place{item.City, item.Country})
		}
	}
	return persons, places
}

// SolutionC
func solutionC(jsonStr []byte) ([]Person, []Place) {
	people := []Person{}
	places := []Place{}
	var data map[string][]json.RawMessage
	err := json.Unmarshal(jsonStr, &data)
	if err != nil {
		fmt.Println(err)
		return people, places
	}
	for _, thing := range data["things"] {
		people = addPersonC(thing, people)
		places = addPlaceC(thing, places)
	}
	return people, places
}

func addPersonC(thing json.RawMessage, people []Person) []Person {
	person := Person{}
	if err := json.Unmarshal(thing, &person); err != nil {
		fmt.Println(err)
	} else {
		if person != *new(Person) {
			people = append(people, person)
		}
	}
	return people
}

func addPlaceC(thing json.RawMessage, places []Place) []Place {
	place := Place{}
	if err := json.Unmarshal(thing, &place); err != nil {
		fmt.Println(err)
	} else {
		if place != *new(Place) {
			places = append(places, place)
		}
	}
	return places
}
