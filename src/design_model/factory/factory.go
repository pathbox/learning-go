package factory

import "fmt"

type AppleFactory struct{}

func (appleFactory AppleFactory) CreateFruit() Fruit {
	return &Apple{}
}

func Fruit interface {
	Eat()
}
func (apple Apple) Eat() {
	fmt.Println("eat apple")
}