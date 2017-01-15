package main

import (
	"fmt"
)

type Name struct {
	first  string
	middle string
	last   string
}

func (n *Name) showFirstName() string {
	return n.first
}

func showFirst(first *string) string {
	return *first
}

func (n Name) showMiddleName() string {
	return n.middle
}

func showMiddle(middle string) string {
	return middle
}

func showLastName(n Name) string {
	return n.last
}

func showLastNamePointer(n *Name) string {
	return n.last
}

func main() {
	n := Name{first: "Joe", middle: "Hello", last: "Waterson"}
	fmt.Println("First Name: ", n.showFirstName())
	fmt.Println("Middle Name: ", n.showMiddleName())
	fmt.Println("Middle Name: ", showMiddle(n.middle))
	fmt.Println("First Name: ", showFirst(&n.first))

	fmt.Println("Last Name: ", showLastName(n))
	fmt.Println("Last Name: ", showLastNamePointer(&n))

}
