package main

import (
	"fmt"
	"sort"
)

type Person struct {
	Name     string
	Age      int
	ShoeSize float32
}

type PeopleByShoesSize []Person

func (p PeopleByShoesSize) Len() int {
	return len(p)
}

func (p PeopleByShoesSize) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PeopleByShoesSize) Less(i, j int) bool {
	return (p[i].ShoeSize < p[j].ShoeSize)
}

func main() {
	people := []Person{
		{
			Name:     "Person1",
			Age:      25,
			ShoeSize: 8,
		},
		{
			Name:     "Person2",
			Age:      21,
			ShoeSize: 4,
		},
		{
			Name:     "Person3",
			Age:      15,
			ShoeSize: 9,
		},
		{
			Name:     "Person4",
			Age:      45,
			ShoeSize: 15,
		},
		{
			Name:     "Person5",
			Age:      25,
			ShoeSize: 8.5,
		},
	}
	fmt.Println(people)
	sort.Sort(PeopleByShoesSize(people))
	fmt.Println(people)
}
