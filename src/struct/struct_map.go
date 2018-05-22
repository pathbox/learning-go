package main

import "fmt"

type User struct {
	Name     string
	ColorMap map[string]string
}

func main() {
	user := &User{
		Name:     "Cary",
		ColorMap: make(map[string]string),
	}

	user.ColorMap["Red"] = "Red"
	fmt.Println("Before:", user)
	AddColor(user.ColorMap)
	fmt.Println("After:", user)
}

func AddColor(cMap map[string]string) {
	cMap["Green"] = "Green"
	cMap["Blue"] = "Blue"
}
